package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"strings"

	"github.com/0xPolygon/polygon-sdk/helper/keccak"
	"github.com/anconprotocol/node/docs"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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
	init := flag.Bool("init", false, "Initialize merkle tree storage")
	genKeys := flag.Bool("keys", false, "Generate keys")
	hostName := flag.String("hostname", "cerro-ancon", "Send custom host name")
	rootHash := flag.String("roothash", "", "root hash")
	rootkey := flag.String("rootkey", "", "root key")
	moniker := flag.String("moniker", "anconprotocol", "DAG Store rootname")
	//	seedPeers := flag.String("peers", "", "Array of peer addresses ")
	quic := flag.Bool("quic", false, "Enable QUIC")
	tlsKey := flag.String("tlscert", "", "TLS certificate")
	tlsCert := flag.String("tlskey", "", "TLS key")
	privateKeyPath := flag.String("privatekeypath", "", "")

	flag.Parse()
	if *genKeys == true {
		result, err := handler.GenerateKeys()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
		return
	}

	r := gin.Default()
	config := cors.DefaultConfig()

	if *enableCors == true {
		config.AllowOrigins = strings.Split(*allowOrigins, ",")
		r.Use(cors.New(config))
	}

	ctx := context.Background()

	s := sdk.NewStorage(*dataFolder)
	dagHandler := &sdk.AnconSyncContext{Store: s}

	docs.SwaggerInfo.BasePath = "/v1"

	if *init == true {

		// Set your own keypair
		priv, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		var digest []byte

		keccak.Keccak256(digest, []byte(*hostName))
		signed, err := priv.Sign(rand.Reader, digest, nil)

		if err != nil {
			panic(err)
		}

		lnkCtx := ipld.LinkContext{
			LinkPath: ipld.ParsePath(types.GetNetworkPath(*moniker)),
		}

		n := basicnode.NewString(base64.RawStdEncoding.EncodeToString(signed))

		link := dagHandler.Store.Store(lnkCtx, n) //Put(ctx, key, []byte(key))

		result, _, err := handler.InitGenesis(*hostName, *moniker, link, priv)

		if err != nil {
			panic(err)
		}

		fmt.Println(result)

		return
	}
	wakuHandler := handler.NewWakuHandler(dagHandler, *peerAddr, *wakuAddr, *privateKeyPath)
	wakuHandler.Start()
	proofHandler := handler.NewProofHandler(dagHandler, wakuHandler, *moniker, *privateKeyPath)

	if *rootHash != "" {
		hash, err := proofHandler.VerifyGenesis(*rootkey, *moniker)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("key validated with last current root hash: ", string(hash))

	}

	dagJsonHandler := handler.NewDagHandler(
		dagHandler,
		proofHandler,
		wakuHandler,
		*rootkey,
		*moniker,
	)

	didHandler := handler.NewDidHandler(
		dagHandler,
		proofHandler.GetProofService(),
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
		api.GET("/dag/:cid/*path", dagJsonHandler.DagJsonRead)
		api.POST("/dag", dagJsonHandler.DagJsonWrite)
		api.PUT("/dag", dagJsonHandler.Update)
		api.POST("/did", didHandler.CreateDid)
		api.POST("/did/web", didHandler.CreateDid)
		api.GET("/did/:did", didHandler.ReadDid)
		api.GET("/proof/:key", proofHandler.Read)
		api.GET("/proof", proofHandler.Read)
		api.GET("/proofs/lasthash", proofHandler.ReadCurrentRootHash)
		api.GET("/topics", dagJsonHandler.UserDag)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	dagJsonHandler.ListenAndSync(ctx)
	proofHandler.Listen(ctx)
	proofHandler.HandleIncomingProofRequests()

	if *quic {
		http3.ListenAndServe(*apiAddr, *tlsCert, *tlsKey, r)
	} else {
		r.Run(*apiAddr)
	}

}
