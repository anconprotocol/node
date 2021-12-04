package transfer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

func SignedProofAbiMethod() abi.Method {

	// uint256Type, _ := abi.NewType("uint256", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	request := abi.NewMethod(
		"transferURIWithProof",
		"transferURIWithProof",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{{
			Name:    "metadataCid",
			Type:    stringType,
			Indexed: false,
		},
			{
				Name:    "fromOwner",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "resultCid",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "toOwner",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "toAddress",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "tokenId",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "prefix",
				Type:    stringType,
				Indexed: false,
			},
			{
				Name:    "signature",
				Type:    bytesType,
				Indexed: false,
			},
		},
		abi.Arguments{abi.Argument{
			Type: uintType,
		}},
	)

	return request
}

type OnchainAdapter struct {
	PrivateKey *ecdsa.PrivateKey
	ChainName  string
	ChainID    int
}

func NewOnchainAdapter(pk *ecdsa.PrivateKey) OnchainAdapter {

	return OnchainAdapter{
		PrivateKey: pk,
		ChainName:  "Ethereum",
		ChainID:    5,
	}
}

// https://gist.github.com/miguelmota/bc4304bb21a8f4cc0a37a0f9347b8bbb
func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func encodeBytesString(v string) []byte {
	decoded, err := hex.DecodeString(v)
	if err != nil {
		panic(err)
	}
	return decoded
}

func encodeUint256(v string) []byte {
	bn := new(big.Int)
	bn.SetString(v, 10)
	return math.U256Bytes(bn)
}

func encodeUint256Array(arr []string) []byte {
	var res [][]byte
	for _, v := range arr {
		b := encodeUint256(v)
		res = append(res, b)
	}

	return bytes.Join(res, nil)
}
func (adapter *OnchainAdapter) ApplyRequestWithProof(
	ctx context.Context,
	metadataCid string,
	resultCid string,
	fromOwner string,
	toOwner string,
	toAddress string,
	tokenId string,
	prefix string,
) (hexutil.Bytes, string, error) {

	id := (tokenId)
	unsignedProofData := encodePacked(
		[]byte("\x19Ethereum Signed Message:\n32"),
		// Proof
		crypto.Keccak256(encodePacked(
			// Current metadata cid
			[]byte(metadataCid),
			// Current owner (opaque)
			[]byte(fromOwner),
			// Updated metadata cid
			[]byte(resultCid),
			// New owner address
			[]byte(toOwner),
			// Token Address
			[]byte(toAddress),
			// Token Id
			[]byte(id),
			// Contract Prefix
			[]byte(prefix)),
		))

	hash := crypto.Keccak256Hash(unsignedProofData)

	signature, err := crypto.Sign(hash.Bytes(), adapter.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("signing failed")
	}

	signedProofData, err := SignedProofAbiMethod().Inputs.Pack(
		metadataCid, fromOwner, resultCid, toOwner, toAddress, id, prefix, signature)

	if err != nil {
		return nil, "", fmt.Errorf("packing for signature proof generation failed")
	}

	return signedProofData, resultCid, nil
}
