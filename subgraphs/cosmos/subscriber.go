package cosmos

import (
	"context"
	"fmt"

	"github.com/anconprotocol/sdk"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type SubscriptionType string

const (
	Tx                  SubscriptionType = "tm.event='Tx'"
	ValidatorSetUpdates SubscriptionType = "tm.event='ValidatorSetUpdates'"
	NewBlock            SubscriptionType = "tm.event='NewBlock'"
	BankModule          SubscriptionType = "message.module='bank'"
)

type CosmosIndexer struct {
	AnconSyncContext *sdk.AnconSyncContext
	Client           *rpcclient.WSClient
	LastLinkNode     datamodel.Node
	LastLink         datamodel.Link
}

func (i *CosmosIndexer) Subscribe(ctx context.Context, subscriptionType SubscriptionType) {
	i.Client.Subscribe(ctx, string(subscriptionType))
}

// @BasePath /v0
// DagJsonWrite godoc
// @Summary Stores JSON as dag-json
// @Schemes
// @Description Writes a dag-json block which syncs with IPFS. Returns a CID.
// @Tags dag-json
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/tip [post]
func (i *CosmosIndexer) TipEvent(c *gin.Context) {
	c.JSON(201, gin.H{
		"cid": i.LastLink,
	})
}

func New(ctx context.Context, tmRPCAddr, tmEndpoint string) *CosmosIndexer {
	tmWsClient, err := rpcclient.NewWS(tmRPCAddr, tmEndpoint)
	dag := ctx.Value("dag").(*sdk.AnconSyncContext)
	if err != nil {
		panic(err)
	}
	err = tmWsClient.Start()
	if err != nil {
		panic(err)
	}
	i := &CosmosIndexer{Client: tmWsClient, AnconSyncContext: dag}
	go func() {
		for {
			select {
			case <-i.Client.Quit():
				return
			case res, ok := <-i.Client.ResponsesCh:
				if ok {
					block, err := sdk.Decode(basicnode.Prototype.Any, string(res.Result))
					if err != nil {
						fmt.Errorf("invalid json %v", err)
					}

					i.LastLink = i.AnconSyncContext.Store.Store(ipld.LinkContext{
						LinkNode: i.LastLinkNode,
					}, block)
					i.LastLinkNode = block
					// PushBlock(c.Request.Context(), dagctx, cid)

				}
			}
		}
	}()
	return i
}
