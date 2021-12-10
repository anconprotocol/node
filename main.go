package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/rpc"

	gqlgenh "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/anconprotocol/contracts/adapters/ethereum/erc721/transfer"
	graphqlclient "github.com/anconprotocol/contracts/graphql/client"
	"github.com/anconprotocol/contracts/graphql/server/graph"
	"github.com/anconprotocol/contracts/graphql/server/graph/generated"
	"github.com/anconprotocol/node/docs"
	"github.com/anconprotocol/node/subgraphs/cosmos"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	dbm "github.com/tendermint/tm-db"

	"github.com/anconprotocol/contracts/sdk/durin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SubgraphConfig struct {
	CosmosMoniker        string
	CosmosAppHash        string
	CosmosProxyAddress   string
	CosmosPrimaryAddress string
	CosmosWitnessAddress string
	CosmosHeight         int
	EnableDagcosmos      bool

	EvmAddress   string
	EvmChainId   string
	EnableDageth bool
}

// Defining the dageth RPC handler
func dagethRPCHandler(anconCtx sdk.AnconSyncContext, gqlAddress string) gin.HandlerFunc {

	db := dbm.NewMemDB()

	proofs, _ := proofsignature.NewIavlAPI(anconCtx.Store, anconCtx.Exchange, db, 2000, 0)
	gqlcli := graphqlclient.NewClient(http.DefaultClient, gqlAddress)
	durin := durin.NewDurinAPI(transfer.NewOnchainAdapter(anconCtx.PrivateKey), gqlcli)
	server := rpc.NewServer()

	err := server.RegisterName(durin.Namespace, durin.Service)
	if err != nil {
		panic(err)
	}

	err = server.RegisterName(proofs.Namespace, proofs.Service)

	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", anconCtx)
		rq := c.Request.WithContext(ctx)
		server.ServeHTTP(c.Writer, rq)
	}
}

// Defining the JSON RPC handler
func jsonRPCHandler(anconCtx sdk.AnconSyncContext, gqlAddress string) gin.HandlerFunc {

	gqlcli := graphqlclient.NewClient(http.DefaultClient, gqlAddress)
	durin := durin.NewDurinAPI(transfer.NewOnchainAdapter(anconCtx.PrivateKey), gqlcli)
	server := rpc.NewServer()

	err := server.RegisterName(durin.Namespace, durin.Service)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", anconCtx)
		rq := c.Request.WithContext(ctx)
		server.ServeHTTP(c.Writer, rq)
	}
}

// Defining the Graphql handler
func graphqlHandler(s sdk.Storage) gin.HandlerFunc {
	h := gqlgenh.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &sdk.AnconSyncContext{
			Store: s,
		})
		rq := c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, rq)
	}
}

// Defining the Playground handler
func playgroundHandler(s sdk.Storage) gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &sdk.AnconSyncContext{
			Store: s,
		})
		rq := c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, rq)
	}
}

// @title        Ancon Protocol Sync API v0.4.0
// @version      0.4.0
// @description  API

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.ancon.did.pa
// @BasePath  /v0
func main() {
	pk, has := os.LookupEnv("ETHEREUM_ADAPTER_KEY")
	if !has {
		panic(fmt.Errorf("environment key ETHEREUM_ADAPTER_KEY not found"))
	}
	// ring, _ := keyring.Open(keyring.Config{
	// 	AllowedBackends: []keyring.BackendType{
	// 		keyring.FileBackend,
	// 	},
	// 	ServiceName: "signer",
	// })
	// key, err := ring.Get("ethereum")
	// if err == keyring.ErrKeyNotFound {
	// 	_ = ring.Set(keyring.Item{
	// 		Key:  "ethereum",
	// 		Data: []byte(pk),
	// 	})
	// 	key, err = ring.Get("ethereum")
	// }

	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(fmt.Errorf("invalid ETHEREUM_ADAPTER_KEY"))
	}

	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/34075/p2p/16Uiu2HAmAihQ3QDNyNxfYiN8qAvaBZKzqosmudoG9KYFMhW2YDXd", "A remote peer to sync")
	addr := flag.String("addr", "/ip4/0.0.0.0/tcp/7702", "Host multiaddr")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")

	subgraph := SubgraphConfig{}
	init := flag.Bool("init", false, "genesis")
	subgraph.EnableDageth = *flag.Bool("enable-dageth", false, "enable EVM subgraph")
	subgraph.EnableDagcosmos = *flag.Bool("enable-dagcosmos", false, "enable Cosmos subgraph")
	subgraph.CosmosPrimaryAddress = *flag.String("cosmos-primary", "", "primary")
	subgraph.EvmAddress = *flag.String("evm-node-address", "", "remote node address")
	subgraph.EvmChainId = *flag.String("evm-chain-id", "", "chain idd")
	moniker := flag.String("moniker", "my-graph", "moniker")
	flag.Parse()

	gqlAddress := fmt.Sprintf("%s/v0/query", apiAddr)
	s := sdk.NewStorage(*dataFolder)

	if *init {
		s.InitGenesis([]byte(*moniker))
		return
	} else {
		root := os.Getenv("ROOTHASH")
		subgraph.CosmosPrimaryAddress = os.Getenv("COSMOS_PRIMARY_ADDRESS")
		subgraph.EnableDagcosmos = cast.ToBool(os.Getenv("ENABLE_DAGCOSMOS"))

		s.LoadGenesis(root)
	}
	ctx := context.Background()
	host := impl.NewPeer(ctx, *addr)

	exchange, ipfspeer := impl.NewRouter(ctx, host, s.LinkSystem, *peerAddr)
	fmt.Println(ipfspeer.ID)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/v0"

	dagHandler := sdk.NewAnconSyncContext(s, exchange, ipfspeer, privateKey)
	didHandler := handler.Did{
		AnconSyncContext: dagHandler,
	}

	dagJsonHandler := handler.DagJsonHandler{
		AnconSyncContext: dagHandler,
	}
	dagCborHandler := handler.DagCborHandler{
		AnconSyncContext: dagHandler,
	}
	fileHandler := handler.FileHandler{
		AnconSyncContext: dagHandler,
	}
	api := r.Group("/v0")
	{
		api.POST("/file", fileHandler.FileWrite)
		api.POST("/query", graphqlHandler(s))
		api.GET("/query", playgroundHandler(s))
		api.GET("/file/:cid/*path", fileHandler.FileRead)
		api.GET("/dagjson/:cid/*path", dagJsonHandler.DagJsonRead)
		api.GET("/dagcbor/:cid/*path", dagCborHandler.DagCborRead)
		api.POST("/dagjson", dagJsonHandler.DagJsonWrite)
		api.POST("/dagcbor", dagCborHandler.DagCborWrite)
		api.POST("/did/key", didHandler.CreateDidKey)
		api.POST("/did/web", didHandler.CreateDidWeb)
		api.GET("/did/:did", didHandler.ReadDid)
	}
	if subgraph.EnableDagcosmos {
		ctx := context.WithValue(context.Background(), "dag", dagHandler)
		indexer := cosmos.New(ctx, subgraph.CosmosPrimaryAddress, "/websocket")
		r.GET("/indexer/cosmos/tip", indexer.TipEvent)
		indexer.Subscribe(ctx, cosmos.NewBlock)
	}
	r.GET("/user/:did/did.json", didHandler.ReadDidWebUrl)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/rpc", dagethRPCHandler(*dagHandler, gqlAddress))
	r.Run(*apiAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
