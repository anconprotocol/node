package handler

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/anconprotocol/contracts/adapters/ethereum/erc721/transfer"
	graphqlclient "github.com/anconprotocol/contracts/graphql/client"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"

	"github.com/anconprotocol/contracts/sdk/durin"

	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/gin-gonic/gin"
	dbm "github.com/tendermint/tm-db"
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

// Defining the dageth RPC handler
func RPCHandler(anconCtx sdk.AnconSyncContext, gqlAddress string) gin.HandlerFunc {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}

	proofs, err := proofsignature.NewIavlAPI(anconCtx.Store, anconCtx.Exchange, db, 2000, 0)

	if err != nil {
		panic(err)
	}

	gqlcli := graphqlclient.NewClient(http.DefaultClient, gqlAddress)
	durin := durin.NewDurinAPI(transfer.NewOnchainAdapter(anconCtx.PrivateKey), gqlcli)
	server := rpc.NewServer()

	err = server.RegisterName(proofs.Namespace, proofs.Service)
	if err != nil {
		panic(err)
	}
	err = server.RegisterName(durin.Namespace, durin.Service)

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
func JsonRPCHandler(anconCtx sdk.AnconSyncContext, gqlAddress string) gin.HandlerFunc {

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
