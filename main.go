package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	cmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/cmd/tendermint/commands/debug"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"

	"path/filepath"
	"strings"

	tmconfig "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/anconprotocol/node/docs"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/anconprotocol/sdk"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lucas-clemente/quic-go/http3"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SubgraphConfig struct {
	CosmosPrimaryAddress string

	CosmosHeight    int
	EnableDagcosmos bool

	EvmAddress   string
	EvmChainId   string
	EnableDageth bool
}

// @title        Ancon Protocol Sync API v1.5.0
// @version      1.5.0
// @description  API

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.ancon.did.pa
// @BasePath  /v1
func main() {

	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/34075/p2p/16Uiu2HAmAihQ3QDNyNxfYiN8qAvaBZKzqosmudoG9KYFMhW2YDXd", "A remote peer to sync")
	wakuAddr := flag.String("wakuaddr", "0.0.0.0:8876", "Waku address")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")
	enableCors := flag.Bool("cors", false, "Allow CORS")
	allowOrigins := flag.String("origins", "*", "Must send a delimited string by commas")

	rootkey := flag.String("rootkey", "", "root key")
	moniker := flag.String("moniker", "anconprotocol", "DAG Store rootname")
	//	seedPeers := flag.String("peers", "", "Array of peer addresses ")
	quic := flag.Bool("quic", false, "Enable QUIC")
	tlsKey := flag.String("tlscert", "", "TLS certificate")
	tlsCert := flag.String("tlskey", "", "TLS key")
	privateKeyPath := flag.String("privatekeypath", "", "")

	flag.Parse()

	r := gin.Default()
	config := cors.DefaultConfig()

	if *enableCors == true {
		config.AllowOrigins = strings.Split(*allowOrigins, ",")
		r.Use(cors.New(config))
	}

	//////// ctx := context.Background()
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// os.OpenFile(,,privateKeyPath)

	folder := filepath.Join(userHomeDir, *dataFolder)
	db, err := dbm.NewGoLevelDB(handler.DBName, folder)
	if err != nil {
		panic(err)
	}

	s := sdk.NewStorage(db, 0, 1024)

	dagHandler := &sdk.AnconSyncContext{Store: s}

	docs.SwaggerInfo.BasePath = "/v1"
	app := sdk.NewAnconAppChain(&s)
	wakuHandler := handler.NewWakuHandler(dagHandler, *peerAddr, *wakuAddr, *privateKeyPath)
	wakuHandler.Start()
	proofHandler := handler.NewProofHandler(dagHandler, wakuHandler, *moniker, *privateKeyPath)

	tm := "tcp://0.0.0.0:26657"
	// nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	// privvalKey := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	nodeFunc := func(cfg *tmconfig.Config, logger tmlog.Logger) (*node.Node, error) {
		nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
		if err != nil {
			return nil, fmt.Errorf("failed to load or gen node key %s: %w", cfg.NodeKeyFile(), err)
		}
		sec, err := time.ParseDuration("14s")

		cfg.Consensus.TimeoutCommit = sec
		return node.NewNode(cfg,
			privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
			nodeKey,
			proxy.NewLocalClientCreator(app),
			node.DefaultGenesisDocProviderFunc(cfg),
			func(d *node.DBContext) (dbm.DB, error) {
				return db, nil
			},
			node.DefaultMetricsProvider(cfg.Instrumentation),
			logger,
		)
	}

	client, err := rpcclient.New(tm, "/websocket")
	dagJsonHandler := handler.NewDagHandler(
		dagHandler,
		proofHandler,
		client,
		wakuHandler,
		*rootkey,
		*moniker,
	)

	didHandler := handler.NewDidHandler(
		dagHandler,
		wakuHandler,
		*rootkey,
		*moniker,
	)

	fileHandler := handler.FileHandler{
		RootKey:          *rootkey,
		AnconSyncContext: dagHandler,
		Moniker:          *moniker,
	}
	//	g := handler.PlaygroundHandler(*dagHandler, adapter, proofHandler.GetProofAPI())

	api := r.Group("/v1")
	{
		// api.GET("/graphql", g)
		// api.POST("/graphql", g)
		api.GET("/file/:cid/*path", fileHandler.FileRead)
		api.GET("/dag/:cid", dagJsonHandler.DagJsonRead)
		api.POST("/dag", dagJsonHandler.DagJsonWrite)
		api.PUT("/dag", dagJsonHandler.Update)
		api.POST("/did", didHandler.CreateDid)
		api.POST("/did/web", didHandler.CreateDid)
		api.GET("/did/:did", didHandler.ReadDid)
		api.GET("/proof/:key", dagJsonHandler.Read)
		api.GET("/proof", dagJsonHandler.Read)
		api.GET("/proofs/lasthash", dagJsonHandler.ReadCurrentRootHash)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	go func() {
		if *quic {
			http3.ListenAndServe(*apiAddr, *tlsCert, *tlsKey, r)
		} else {
			r.Run(*apiAddr)
		}
	}()

	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GenValidatorCmd,
		cmd.InitFilesCmd,
		cmd.ProbeUpnpCmd,
		cmd.LightCmd,
		cmd.ReplayCmd,
		cmd.ReplayConsoleCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ResetStateCmd,
		cmd.ShowValidatorCmd,
		cmd.TestnetFilesCmd,
		cmd.ShowNodeIDCmd,
		cmd.GenNodeKeyCmd,
		cmd.VersionCmd,
		cmd.RollbackStateCmd,
		debug.DebugCmd,
		cli.NewCompletionCmd(rootCmd, true),
	)

	// NOTE:
	// Users wishing to:
	//	* Use an external signer for their validators
	//	* Supply an in-proc abci app
	//	* Supply a genesis doc file from another source
	//	* Provide their own DB implementation
	// can copy this file and use something other than the
	// DefaultNewNode function
	// nodeFunc := nm.DefaultNewNode

	// Create & start node
	rootCmd.AddCommand(cmd.NewRunNodeCmd(nodeFunc))

	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", *dataFolder)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

}
