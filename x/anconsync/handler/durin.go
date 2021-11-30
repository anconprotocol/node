package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/adapters/ethereum/erc721/transfer"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type DurinService struct {
	Adapter   transfer.EthereumAdapter
	GqlClient *Client
}

func NewDurinService(evm transfer.EthereumAdapter, gqlClient *Client) *DurinService {
	return &DurinService{
		Adapter:   evm,
		GqlClient: gqlClient,
	}
}

func (s *DurinService) MsgHandler(to string, name string, args map[string]interface{}) (hexutil.Bytes, error) {
	switch name {
	default:
		tokenId := args["tokenId"].(uint64)
		input := MetadataTransactionInput{
			Path:     args["path"].(string),
			Cid:      args["cid"].(string),
			Owner:    args["fromOwner"].(string),
			NewOwner: args["toOwner"].(string),
		}
		// Send graphql mutation for IPLD DAG computing
		res, err := s.GqlClient.TransferOwnership(context.Background(), input)
		if err != nil {
			return nil, fmt.Errorf("transfer ownership reverted")
		}

		// Apply signature to create proof
		txdata, err := s.Adapter.ApplyRequestWithProof(input.Cid, res.Metadata.Cid, input.Owner, input.NewOwner, to, tokenId)
		if err != nil {
			return nil, fmt.Errorf("request with proof raw tx failed")
		}
		return txdata, nil
	}
}

func (s *DurinService) DurinCall(params ...interface{}) hexutil.Bytes {
	to := params[0].(string)
	data, err := (hexutil.Decode(params[1].(string)))
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("fail reading data").Error())))
	}
	abis := params[2].(json.RawMessage)

	// Get the function selector
	selector := string(data)[:10]

	iface, err := abi.JSON(bytes.NewReader(abis))

	fn, err := iface.MethodById([]byte(selector))
	var values map[string]interface{}
	err = iface.UnpackIntoMap(values, fn.Name, data)

	// Execute graphql
	txdata, err := s.MsgHandler(to, fn.Name, values)

	return txdata
}
