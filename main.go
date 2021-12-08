package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	gqlgenh "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/anconprotocol/node/adapters/ethereum/erc721/transfer"
	graphqlclient "github.com/anconprotocol/node/contracts/graphql/client"
	"github.com/anconprotocol/node/contracts/graphql/server/graph"
	"github.com/anconprotocol/node/contracts/graphql/server/graph/generated"
	"github.com/anconprotocol/node/docs"
	dagcosmos "github.com/anconprotocol/node/subgraphs/cosmos"
	"github.com/anconprotocol/node/x/anconsync"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/anconprotocol/node/x/anconsync/handler/durin"
	"github.com/anconprotocol/node/x/anconsync/impl"
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
func dagethRPCHandler(anconCtx handler.AnconSyncContext) gin.HandlerFunc {

	gqlcli := graphqlclient.NewClient(http.DefaultClient, "http://localhost:7788/v0/query")
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

// // Defining the dagcosmos RPC handler
// func dagcosmosRPCHandler(anconCtx handler.AnconSyncContext, cfg SubgraphConfig) gin.HandlerFunc {

// 	ctx := context.WithValue(context.Background(), "dag", anconCtx)
// 	p := dagcosmos.New(ctx, cfg.CosmosMoniker, cfg.CosmosPrimaryAddress, cfg.CosmosWitnessAddress, cfg.CosmosProxyAddress, cfg.CosmosAppHash, int(cfg.CosmosHeight))

// 	// 1) Register regular routes.
// 	r := proxy.RPCRoutes(p.Client)
// 	rpcserver.RegisterRPCFuncs(http.DefaultServeMux, r, p.Logger)

// 	// 2) Allow websocket connections.
// 	wmLogger := p.Logger.With("protocol", "websocket")
// 	wm := rpcserver.NewWebsocketManager(r,
// 		rpcserver.OnDisconnect(func(remoteAddr string) {
// 			err := p.Client.UnsubscribeAll(context.Background(), remoteAddr)
// 			if err != nil && err != tmpubsub.ErrSubscriptionNotFound {
// 				wmLogger.Error("Failed to unsubscribe addr from events", "addr", remoteAddr, "err", err)
// 			}
// 		}),
// 		rpcserver.ReadLimit(p.Config.MaxBodyBytes),
// 	)
// 	return func(c *gin.Context) {
// 		rq := c.Request.WithContext(ctx)

// 		wm.SetLogger(wmLogger)
// 		wm.WebsocketHandler(c.Writer, rq)
// 	}
// }

// Defining the JSON RPC handler
func jsonRPCHandler(anconCtx handler.AnconSyncContext) gin.HandlerFunc {

	gqlcli := graphqlclient.NewClient(http.DefaultClient, "http://localhost:7788/v0/query")
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
func graphqlHandler(s anconsync.Storage) gin.HandlerFunc {
	h := gqlgenh.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &handler.AnconSyncContext{
			Store: s,
		})
		rq := c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, rq)
	}
}

// Defining the Playground handler
func playgroundHandler(s anconsync.Storage) gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &handler.AnconSyncContext{
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
	subgraph.CosmosAppHash = *flag.String("cosmos-app-hash", "", "app hash")
	subgraph.CosmosHeight = *flag.Int("cosmos-height", 1, "height")
	subgraph.CosmosPrimaryAddress = *flag.String("cosmos-primary", "", "primary")
	subgraph.EvmAddress = *flag.String("evm-node-address", "", "remote node address")
	subgraph.EvmChainId = *flag.String("evm-chain-id", "", "chain idd")
	subgraph.CosmosMoniker = *flag.String("cosmos-moniker", "my-graph", "cosmos-moniker")
	moniker := flag.String("moniker", "my-graph", "moniker")
	flag.Parse()

	s := anconsync.NewStorage(*dataFolder)

	if *init {
		s.InitGenesis([]byte(*moniker))
		return
	} else {
		root := os.Getenv("ROOTHASH")
		subgraph.CosmosMoniker = os.Getenv("COSMOS_MONIKER")
		subgraph.CosmosAppHash = os.Getenv("COSMOS_APP_HASH")
		subgraph.CosmosHeight = cast.ToInt(os.Getenv("COSMOS_HEIGHT"))
		subgraph.CosmosPrimaryAddress = os.Getenv("COSMOS_PRIMARY_ADDRESS")
		subgraph.CosmosProxyAddress = os.Getenv("COSMOS_PROXY_ADDRESS")
		subgraph.EnableDagcosmos = cast.ToBool(os.Getenv("ENABLE_DAGCOSMOS"))

		s.LoadGenesis(root)
	}
	ctx := context.Background()
	host := impl.NewPeer(ctx, *addr)

	exchange, ipfspeer := impl.NewRouter(ctx, host, s, *peerAddr)
	fmt.Println(ipfspeer.ID)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/v0"

	dagHandler := handler.NewAnconSyncContext(s, exchange, ipfspeer, privateKey)
	api := r.Group("/v0")
	{
		api.POST("/file", dagHandler.FileWrite)
		api.POST("/query", graphqlHandler(s))
		api.GET("/query", playgroundHandler(s))
		api.GET("/file/:cid/*path", dagHandler.FileRead)
		api.GET("/dagjson/:cid/*path", dagHandler.DagJsonRead)
		api.GET("/dagcbor/:cid/*path", dagHandler.DagCborRead)
		api.POST("/dagjson", dagHandler.DagJsonWrite)
		api.POST("/dagcbor", dagHandler.DagCborWrite)
		api.POST("/did/key", dagHandler.CreateDidKey)
		api.POST("/did/web", dagHandler.CreateDidWeb)
		api.GET("/did/:did", dagHandler.ReadDid)
	}
	if subgraph.EnableDagcosmos {

		ctx := context.WithValue(context.Background(), "dag", dagHandler)
		indexer := dagcosmos.New(ctx, subgraph.CosmosPrimaryAddress, "/websocket")
		r.GET("/indexer/cosmos/tip", indexer.TipEvent)
		indexer.Subscribe(ctx, dagcosmos.NewBlock)

	}
	r.GET("/user/:did/did.json", dagHandler.ReadDidWebUrl)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/rpc", jsonRPCHandler(*dagHandler))
	r.Run(*apiAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
