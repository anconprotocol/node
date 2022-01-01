package handler

import (
	"context"
	"encoding/json"
	"net/http"
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
	api := protocol.NewProtocolAPI(adapter, &anconCtx.Store, proofs)

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "dag", anconCtx)
		rq := c.Request.WithContext(ctx)

		// get query
		opts := gql.NewRequestOptions(rq)
		to := rq.PostForm.Get("to")
		from := rq.PostForm.Get("from")
		sig := rq.PostForm.Get("sig")
		result := api.Service.Call(to, from, []byte(sig), opts.Query)
		// if formatErrorFn := h.formatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		// 	formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		// 	for i, formattedError := range result.Errors {
		// 		formatted[i] = formatErrorFn(formattedError.OriginalError())
		// 	}
		// 	result.Errors = formatted
		// }

		acceptHeader := rq.Header.Get("Accept")
		_, raw := rq.URL.Query()["raw"]
		if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
			renderPlayground(w, rq)
			return
		}

		// use proper JSON Header
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		var buff []byte
		if h.pretty {
			w.WriteHeader(http.StatusOK)
			buff, _ = json.MarshalIndent(result, "", "\t")

			w.Write(buff)
		} else {
			w.WriteHeader(http.StatusOK)
			buff, _ = json.Marshal(result)

			w.Write(buff)
		}

		if h.resultCallbackFn != nil {
			h.resultCallbackFn(ctx, &params, result, buff)
		}

		server.ServeHTTP(c.Writer, rq)
	}
}
