package metadata

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/adapters/ethereum/erc721/transfer"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func ProofTypeAbi() abi.Type {
	proofType, _ := abi.NewType("MetadataTransferProofPacket", "", []abi.ArgumentMarshaling{
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
	proofType, _ := abi.NewType("MetadataTransferProofPacket", "", []abi.ArgumentMarshaling{
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

	uint256Type, _ := abi.NewType("uint256", "", nil)
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
			Type:    uint256Type,
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

	uint256Type, _ := abi.NewType("uint256", "", nil)
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
			Type:    uint256Type,
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

type EthereumAdapter struct {
}

func (adapter *EthereumAdapter) ApplyRequestWithProof(
	metadataCid string,
	resultCid string,
	fromOwner string,
	toOwner string,
	toAddress string,
	tokenId uint64,
) (hexutil.Bytes, error) {

	pk, has := os.LookupEnv("ETHEREUM_ADAPTER_KEY")
	if !has {
		return nil, fmt.Errorf("environment key ETHEREUM_ADAPTER_KEY not found")
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid ETHEREUM_ADAPTER_KEY")
	}

	proofData, err :=
		ExecuteDagContractWithProofAbiMethod().Inputs.Pack(
			toAddress,
			tokenId,
			&transfer.MetadataTransferPacket{
				MetadataCid: metadataCid,
				ResultCid:   resultCid,
				FromOwner:   fromOwner,
				ToOwner:     toOwner,
				ToAddress:   toAddress,
				TokenId:     tokenId,
			},
		)
	if err != nil {
		return nil, fmt.Errorf("packing for proof generation failed")
	}

	hash := crypto.Keccak256Hash(proofData)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("signing failed")
	}

	signedPacket := &transfer.MetadataTransferProofPacket{
		MetadataCid: metadataCid,
		ResultCid:   resultCid,
		FromOwner:   fromOwner,
		ToOwner:     toOwner,
		ToAddress:   toAddress,
		TokenId:     tokenId,
		Signature:   hexutil.Encode(signature),
	}

	encoded, err := ExecuteDagContractWithProofAbiMethod().Inputs.Pack(toAddress, tokenId, signedPacket)
	if err != nil {
		return nil, fmt.Errorf("packing for signature proof generation failed")
	}
	return encoded, nil
}

func (adapter *EthereumAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
