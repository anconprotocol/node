package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/adapters/ethereum/erc721/transfer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cast"
)

type DurinAPI struct {
	Namespace string
	Version   string
	Service   *DurinService
	Public    bool
}

type DurinService struct {
	Adapter   transfer.EthereumAdapter
	GqlClient *Client
}

func NewDurinAPI(evm transfer.EthereumAdapter, gqlClient *Client) *DurinAPI {
	return &DurinAPI{
		Namespace: "durin",
		Version:   "1.0",
		Service: &DurinService{
			Adapter:   evm,
			GqlClient: gqlClient,
		},
		Public: true,
	}
}

func (s *DurinService) msgHandler(to string, name string, args map[string]interface{}) (hexutil.Bytes, error) {
	switch name {
	default:
		tokenId := cast.ToUint64(args["tokenId"])
		input := MetadataTransactionInput{
			Path:     "/",
			Cid:      args["metadataCid"].(string),
			Owner:    args["fromOwner"].(string),
			NewOwner: args["toOwner"].(string),
		}
		// Send graphql mutation for IPLD DAG computing
		res, err := s.GqlClient.TransferOwnership(context.Background(), input)
		if err != nil {
			return nil, fmt.Errorf("transfer ownership reverted")
		}

		// Apply signature to create proof
		txdata, err := s.Adapter.ApplyRequestWithProof(context.Background(), input.Cid, res.Metadata.Cid, input.Owner, input.NewOwner, to, tokenId)
		if err != nil {
			return nil, fmt.Errorf("request with proof raw tx failed")
		}
		return txdata, nil
	}
}

func (s *DurinService) Call(to string, from string, data json.RawMessage, abis json.RawMessage) hexutil.Bytes {

	p := []byte(data)
	var values map[string]interface{}
	err := json.Unmarshal(p, &values)
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("fail unpack data").Error())))
	}
	// Execute graphql
	txdata, err := s.msgHandler(to, "", values)

	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("reverted").Error())))
	}
	return txdata
}
