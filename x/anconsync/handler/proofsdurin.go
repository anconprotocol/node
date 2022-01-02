package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/anconprotocol/node/x/anconsync/handler/protocol"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"
	gql "github.com/graphql-go/handler"

	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/gin-gonic/gin"
)

var (
	dbName string = "proofs-db"
	dbPath string = ".ancon/db/proofs"
)

// Defining the dageth RPC handler
func SmartContractHandler(anconCtx sdk.AnconSyncContext,
	adapter *ethereum.OnchainAdapter, proofs *proofsignature.IavlProofAPI) gin.HandlerFunc {
	api := protocol.NewProtocolAPI(adapter, &anconCtx.Store, proofs)
	server := rpc.NewServer()

	// err = server.RegisterName(proofs.Namespace, proofs.Service)
	// if err != nil {
	// 	panic(err)
	// }
	err := server.RegisterName("ancon", api.Service)

	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", anconCtx)
		rq := c.Request.WithContext(ctx)
		server.ServeHTTP(c.Writer, rq)
	}
}

func PlaygroundHandler(anconCtx sdk.AnconSyncContext,
	adapter *ethereum.OnchainAdapter, proofs *proofsignature.IavlProofAPI) gin.HandlerFunc {
	//	api := protocol.NewProtocolAPI(adapter, &anconCtx.Store, proofs)

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", anconCtx)
		rq := c.Request.WithContext(ctx)

		// get query
		opts := gql.NewRequestOptions(rq)
		to := opts.Variables["to"]
		from := opts.Variables["from"]
		sig := opts.Variables["sig"]
		var result string
		if to != nil {
			result = fmt.Sprintf(to.(string), from.(string), []byte(sig.(string)), opts.Query)

		}
		// if formatErrorFn := h.]ormatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		// 	formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		// 	for i, formattedError z:= range result.Errors {
		// 		formatted[i] = formatErrorFn(formattedError.OriginalError())
		// 	}
		// 	result.Errors = formatted
		// }

		acceptHeader := rq.Header.Get("Accept")
		_, raw := rq.URL.Query()["raw"]
		if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
			RenderPlayground(c.Writer, rq)
			return
		}

		// use proper JSON Header
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(200, result)
		// if h.resultCallbackFn != nil {
		// 	h.resultCallbackFn(ctx, &params, result, buff)
		// }

	}
}
