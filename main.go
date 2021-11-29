package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	gqlgenh "github.com/99designs/gqlgen/graphql/handler"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/docs"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/codegen/graph"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/codegen/graph/generated"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/handler"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/impl"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Defining the Graphql handler
func graphqlHandler(s anconsync.Storage) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := gqlgenh.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &handler.DagContractTrustedContext{
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
		ctx := context.WithValue(c.Request.Context(), "dag", &handler.DagContractTrustedContext{
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

// @host      ancon.did.pa/api
// @BasePath  /v0
func main() {
	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWAGyXSBPPo7Zq16WoCe6BtHDRQpFXPg9VCDQ1EPXcHWMw", "A remote peer to sync")
	addr := flag.String("addr", "/ip4/0.0.0.0/tcp/7702", "Host multiaddr")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")

	flag.Parse()

	s := anconsync.NewStorage(*dataFolder)
	ctx := context.Background()
	host := net.NewPeer(ctx, *addr)
	// peerhost := "/ip4/192.168.50.138/tcp/7702/p2p/12D3KooWA7vRHFLC8buiEP2xYwUN5kdCgzEtQRozMtnCPDi4n4HM"
	// "/ip4/190.34.226.207/tcp/29557/p2p/12D3KooWGd9mLtWx7WGEd9mnWPbCsr1tFCxtEi7RkgsJYxAZmZgi"

	exchange, ipfspeer := impl.NewRouter(ctx, host, s, *peerAddr)
	fmt.Println(ipfspeer.ID)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/v0"

	dagHandler := handler.NewDagContractContext(s, exchange, ipfspeer)
	api := r.Group("/v0")
	{
		api.POST("/file", dagHandler.FileWrite)
		api.POST("/query", graphqlHandler(s))
		api.GET("/query", playgroundHandler(s))
		api.GET("/file/:cid", dagHandler.FileRead)
		api.GET("/dagjson/:cid/*path", dagHandler.DagJsonRead)
		api.GET("/dagcbor/:cid/*path", dagHandler.DagCborRead)
		api.POST("/dagjson", dagHandler.DagJsonWrite)
		api.POST("/dagcbor", dagHandler.DagCborWrite)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(*apiAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
