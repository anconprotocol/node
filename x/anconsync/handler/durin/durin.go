package durin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anconprotocol/contracts/adapters/ethereum/erc721/transfer"
	graphqlclient "github.com/anconprotocol/contracts/graphql/client"
	"github.com/anconprotocol/contracts/hexutil"
)

type DurinAPI struct {
	Namespace string
	Version   string
	Service   *DurinService
	Public    bool
}

type DurinService struct {
	Adapter   *transfer.OnchainAdapter
	GqlClient *graphqlclient.Client
}

func NewDurinAPI(evm transfer.OnchainAdapter, gqlClient *graphqlclient.Client) *DurinAPI {
	return &DurinAPI{
		Namespace: "durin",
		Version:   "1.0",
		Service: &DurinService{
			Adapter:   &evm,
			GqlClient: gqlClient,
		},
		Public: true,
	}
}

func msgHandler(ctx *DurinService, to string, name string, args map[string]string) (hexutil.Bytes, string, error) {
	switch name {
	default:
		tokenId := args["tokenId"]
		input := graphqlclient.MetadataTransactionInput{
			Path:     "/",
			Cid:      args["metadataCid"],
			Owner:    args["fromOwner"],
			NewOwner: args["toOwner"],
		}
		// Send graphql mutation for IPLD DAG computing
		res, err := ctx.GqlClient.TransferOwnership(context.Background(), input)
		if err != nil {
			return nil, "", fmt.Errorf("transfer ownership reverted")
		}
		metadataCid := args["metadataCid"]
		newCid := res.Metadata.Cid
		newOwner := args["toOwner"]
		fromOwner := args["fromOwner"]
		prefix := args["prefix"]

		// Apply signature to create proof
		txdata, resultCid, err := ctx.Adapter.ApplyRequestWithProof(context.Background(),
			metadataCid,
			newCid,
			fromOwner,
			newOwner,
			to,
			tokenId,
			prefix)
		if err != nil {
			return nil, "", fmt.Errorf("request with proof raw tx failed")
		}
		return txdata, resultCid, nil
	}
}

func (s *DurinService) Call(to string, from string, data json.RawMessage, abis json.RawMessage) hexutil.Bytes {

	p := []byte(data)
	var values map[string]string
	val := make(map[string]string, 2)
	err := json.Unmarshal(p, &values)
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("fail unpack data").Error())))
	}
	// Execute graphql
	txdata, resultCid, err := msgHandler(s, to, "", values)
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("reverted").Error())))
	}

	val["txdata"] = txdata.String()
	val["resultCid"] = resultCid
	jsonval, err := json.Marshal(val)
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("reverted, json marshal").Error())))
	}
	return jsonval
}
