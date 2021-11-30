package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)
func ProofTypeAbi() abi.Type {
	proofType, _ := abi.NewType("DagContractRequestProof", "", []abi.ArgumentMarshaling{
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
			Name:         "toReceiverContractAddress",
			Type:         "string",
			InternalType: "",
			Components:   []abi.ArgumentMarshaling{},
			Indexed:      false,
		},
	})

	return proofType
}

// "requestWithProof(address toReceiverContractAddress, uint256 tokenId, DagContractRequestProof memory proof)",
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
			Name:    "toReceiverContractAddress",
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

func ExecuteDagContracAbiMethod() abi.Method {
	//addressType, _ := abi.NewType("address", "", nil)
	// uint8Type, _ := abi.NewType("uint8", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	// bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	request := abi.NewMethod(
		"request",
		"request",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{abi.Argument{
			Name:    "metadataCid",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "resultCid",
			Type:    stringType,
			Indexed: false,
		},
			abi.Argument{
				Name:    "fromOwner",
				Type:    stringType,
				Indexed: false,
			}, abi.Argument{
				Name:    "toOwner",
				Type:    stringType,
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

type Proof struct {
	metadataCid string
	resultCid string
	fromOwner string
	toOwner string
	toReceiverContractAddress string
}

func (adapter *EthereumAdapter) ApplyRequestWithProof(
	metadataCid string,
	resultCid string,
	fromOwner string,
	toOwner string,
	toAddress string,
	tokenId string,
) (*DagTransaction, error) {

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
			&Proof{
				metadataCid:               metadataCid,
				resultCid:                 resultCid,
				fromOwner:                 fromOwner,
				toOwner:                   toOwner,
				toReceiverContractAddress: toAddress,
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

	return &DagTransaction{
		MetadataCid: metadataCid,
		ResultCid:   resultCid,
		FromOwner:   fromOwner,
		ToOwner:     toOwner,
		Signature:   hexutil.Encode(signature),
	}, nil
}

func (adapter *EthereumAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
