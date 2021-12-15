package handler

import (
	"context"

	gqlgenh "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/anconprotocol/contracts/graphql/server/graph"
	"github.com/anconprotocol/contracts/graphql/server/graph/generated"
	"github.com/anconprotocol/sdk"
	"github.com/gin-gonic/gin"
)

// Defining the Graphql handler
func GraphqlHandler(s sdk.Storage) gin.HandlerFunc {
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
func PlaygroundHandler(s sdk.Storage) gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", &sdk.AnconSyncContext{
			Store: s,
		})
		rq := c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, rq)
	}
}
