package main

import (
	"context"
	"fmt"

	graphqlclient "github.com/anconprotocol/node/contracts/graphql/client"
	"github.com/anconprotocol/node/x/anconsync/handler/durin"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Contract struct {
	durin *durin.DurinService
}

func NewContract(durin *durin.DurinService) *Contract {
	return &Contract{durin: durin}
}

// export execute
func (c *Contract) execute(args ...interface{}) (hexutil.Bytes, string, error) {
	tokenId := args[0].(string)
	input := graphqlclient.MetadataTransactionInput{
		Path:     "/",
		Cid:      args[1].(string),
		Owner:    args[2].(string),
		NewOwner: args[3].(string),
	}
	// Send graphql mutation for IPLD DAG computing
	res, err := c.durin.GqlClient.TransferOwnership(context.Background(), input)
	if err != nil {
		return nil, "", fmt.Errorf("transfer ownership reverted")
	}
	toAddress := args[4].(string)
	metadataCid := args[2].(string)
	newCid := res.Metadata.Cid
	newOwner := args[3].(string)
	fromOwner := args[2].(string)
	prefix := args[5].(string)

	// Apply signature to create proof
	return c.durin.Adapter.ApplyRequestWithProof(context.Background(),
		metadataCid,
		newCid,
		fromOwner,
		newOwner,
		toAddress,
		tokenId,
		prefix)
}
