package handler

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/anconprotocol/node/x/anconsync/handler/protocol"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"

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
