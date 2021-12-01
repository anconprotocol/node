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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type MetadataTransferProofPacket struct {
	MetadataCid string
	ResultCid   string
	FromOwner   string
	ToOwner     string
	ToAddress   string
	TokenId     string
	Signature   string
}
type MetadataTransferPacket struct {
	MetadataCid string
	ResultCid   string
	FromOwner   string
	ToOwner     string
	ToAddress   string
	TokenId     string
}

func ProofTypeAbi() abi.Type {
	proofType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name:         "metadataCid",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "fromOwner",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "resultCid",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "toOwner",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "toAddress",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
	})

	return proofType
}

func ProofWithSignatureTypeAbi() abi.Type {
	proofType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name:         "metadataCid",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "fromOwner",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "resultCid",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "toOwner",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "toAddress",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
		{
			Name:         "signature",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
	})

	return proofType
}

// "requestWithProof(address toAddress, uint256 tokenId, MetadataTransferProofPacket memory proof)",
func ExecuteDagContractWithProofAbiMethod() abi.Method {

	// uint256Type, _ := abi.NewType("uint256", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	// bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	request := abi.NewMethod(
		"requestWithProof",
		"requestWithProof",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{abi.Argument{
			Name:    "toAddress",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "tokenId",
			Type:    stringType,
			Indexed: false,
		},
			abi.Argument{
				Name:    "proof",
				Type:    ProofTypeAbi(),
				Indexed: false,
			},
		},
		abi.Arguments{abi.Argument{
			Type: uintType,
		}},
	)

	return request
}
func ExecuteDagContractWithSignedProofAbiMethod() abi.Method {

	// uint256Type, _ := abi.NewType("uint256", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	// bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	request := abi.NewMethod(
		"requestWithProof",
		"requestWithProof",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{abi.Argument{
			Name:    "toAddress",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "tokenId",
			Type:    stringType,
			Indexed: false,
		},
			abi.Argument{
				Name:    "proof",
				Type:    ProofWithSignatureTypeAbi(),
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
) (hexutil.Bytes, error) {

	id := (tokenId)
	unsignedProofData := encodePacked(
		// Token Address
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
			[]byte(id)),
		))

	hash := crypto.Keccak256Hash(unsignedProofData)

	signature, err := crypto.Sign(hash.Bytes(), adapter.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("signing failed")
	}

	signedProofData := encodePacked(
		// Current metadata cid
		[]byte(metadataCid),
		// Current owner (opaque)
		[]byte(fromOwner), // New owner address
		[]byte(toOwner),
		// Updated metadata cid
		[]byte(resultCid),

		// Token Address
		[]byte(toAddress),
		// Token Id
		[]byte(id),
		// Signature
		(signature),
	)

	signedTxData := encodePacked(
		// Token Address
		[]byte(toAddress),
		// Token Id
		[]byte(id),
		// Proof
		signedProofData,
	)

	if err != nil {
		return nil, fmt.Errorf("packing for signature proof generation failed")
	}
	return signedTxData, nil
}

func (adapter *OnchainAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
