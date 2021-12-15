package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/anconprotocol/node/docs"
	"github.com/anconprotocol/node/subgraphs/cosmos"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/gin-gonic/gin"

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

// @title        Ancon Protocol Sync API v0.4.0
// @version      0.4.0
// @description  API

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.ancon.did.pa
// @BasePath  /v0
func main() {

	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/34075/p2p/16Uiu2HAmAihQ3QDNyNxfYiN8qAvaBZKzqosmudoG9KYFMhW2YDXd", "A remote peer to sync")
	addr := flag.String("addr", "/ip4/0.0.0.0/tcp/7702", "Host multiaddr")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")

	subgraph := SubgraphConfig{}
	subgraph.EnableDageth = *flag.Bool("enable-dageth", false, "enable EVM subgraph")
	subgraph.EnableDagcosmos = *flag.Bool("enable-dagcosmos", false, "enable Cosmos subgraph")
	subgraph.CosmosPrimaryAddress = *flag.String("cosmos-primary", "", "primary")
	subgraph.EvmAddress = *flag.String("evm-node-address", "", "remote node address")
	subgraph.EvmChainId = *flag.String("evm-chain-id", "", "chain idd")

	flag.Parse()

	gqlAddress := fmt.Sprintf("%s/v0/query", apiAddr)
	s := sdk.NewStorage(*dataFolder)

	ctx := context.Background()
	host := impl.NewPeer(ctx, *addr)

	exchange, ipfspeer := impl.NewRouter(ctx, host, s.LinkSystem, *peerAddr)
	fmt.Println(ipfspeer.ID)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/v0"

	dagHandler := sdk.NewAnconSyncContext(s, exchange, ipfspeer)
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
	proofHandler := handler.ProofHandler{
		AnconSyncContext: dagHandler,
	}

	api := r.Group("/v0")
	{
		api.POST("/file", fileHandler.FileWrite)
		api.POST("/query", handler.GraphqlHandler(s))
		api.GET("/query", handler.PlaygroundHandler(s))
		api.GET("/file/:cid/*path", fileHandler.FileRead)
		api.GET("/dagjson/:cid/*path", dagJsonHandler.DagJsonRead)
		api.GET("/dagcbor/:cid/*path", dagCborHandler.DagCborRead)
		api.POST("/dagjson", dagJsonHandler.DagJsonWrite)
		api.POST("/dagcbor", dagCborHandler.DagCborWrite)
		api.POST("/did/key", didHandler.CreateDidKey)
		api.POST("/did/web", didHandler.CreateDidWeb)
		api.GET("/did/:did", didHandler.ReadDid)
		api.GET("/proofs/get/:key", proofHandler.Read)
		api.POST("/proofs", proofHandler.Create)
		api.GET("/proofs/lasthash", proofHandler.ReadCurrentRootHash)
	}
	if subgraph.EnableDagcosmos {
		ctx := context.WithValue(context.Background(), "dag", dagHandler)
		indexer := cosmos.New(ctx, subgraph.CosmosPrimaryAddress, "/websocket")
		r.GET("/indexer/cosmos/tip", indexer.TipEvent)
		indexer.Subscribe(ctx, cosmos.NewBlock)
	}
	r.GET("/user/:did/did.json", didHandler.ReadDidWebUrl)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/rpc", handler.RPCHandler(*dagHandler, gqlAddress))
	r.Run(*apiAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
