package adapters

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func ExecuteDagContractWithProofAbiMethod() abi.Method {
	//addressType, _ := abi.NewType("address", "", nil)
	// uint8Type, _ := abi.NewType("uint8", "", nil)
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
			Name:    "schemaCid",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "dataSourceCid",
			Type:    stringType,
			Indexed: false,
		},
			abi.Argument{
				Name:    "variables",
				Type:    stringType,
				Indexed: false,
			}, abi.Argument{
				Name:    "contractMutation",
				Type:    stringType,
				Indexed: false,
			},
			abi.Argument{
				Name:    "result",
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
			Name:    "schemaCid",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "dataSourceCid",
			Type:    stringType,
			Indexed: false,
		},
			abi.Argument{
				Name:    "variables",
				Type:    stringType,
				Indexed: false,
			}, abi.Argument{
				Name:    "contractMutation",
				Type:    stringType,
				Indexed: false,
			},
			abi.Argument{
				Name:    "result",
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
	ChainID        *big.Int
	AdapterAddress common.Address
}

func (adapter *EthereumAdapter) ExecuteDagContract(
	metadatadCid string,
	resultCid string,
	fromOwner string,
	toOwner string,
) (*DagTransaction, error) {

	pk, has := os.LookupEnv("ETHEREUM_ADAPTER_KEY")
	if !has {
		return nil, fmt.Errorf("environment key ETHEREUM_ADAPTER_KEY not found")
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid ETHEREUM_ADAPTER_KEY")
	}

	data, err := ExecuteDagContractWithProofAbiMethod().Inputs.Pack(metadatadCid, fromOwner, resultCid, toOwner)
	if err != nil {
		return nil, fmt.Errorf("packing for proof generation failed")
	}

	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("signing failed")
	}

	return &DagTransaction{
		MetadataCid: metadatadCid,
		ResultCid: resultCid,
		FromOwner: fromOwner,
		ToOwner: toOwner,
		Signature:     hexutil.Encode(signature),
	}, nil
}

func (adapter *EthereumAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
