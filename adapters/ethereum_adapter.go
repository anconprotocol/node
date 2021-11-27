package adapters

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func ExecuteDagContractWithProofAbiMethod() abi.Method {
	//addressType, _ := abi.NewType("address", "", nil)
	uint8Type, _ := abi.NewType("uint8", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	// SchemaCid        string
	// DataSourceCid    string
	// Variables        map[string]interface{}
	// ContractMutation string

	metadata := abi.NewMethod(
		"executeDagContract",
		"executeDagContract",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{abi.Argument{
			Name:    "v",
			Type:    uint8Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "r",
			Type:    bytes32Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "s",
			Type:    bytes32Type,
			Indexed: false,
		}, abi.Argument{
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
		},
		abi.Arguments{abi.Argument{
			Type: uintType,
		}},
	)

	return metadata
}

func ExecuteDagContractAbiMethod() abi.Method {
	//addressType, _ := abi.NewType("address", "", nil)
	uint8Type, _ := abi.NewType("uint8", "", nil)
	uintType, _ := abi.NewType("uint", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	// SchemaCid        string
	// DataSourceCid    string
	// Variables        map[string]interface{}
	// ContractMutation string

	metadata := abi.NewMethod(
		"executeDagContract",
		"executeDagContract",
		abi.Function,
		"nonpayable",
		false,
		false,
		abi.Arguments{abi.Argument{
			Name:    "v",
			Type:    uint8Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "r",
			Type:    bytes32Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "s",
			Type:    bytes32Type,
			Indexed: false,
		}, abi.Argument{
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
			}, abi.Argument{
				Name:    "result",
				Type:    stringType,
				Indexed: false,
			},
		},
		abi.Arguments{abi.Argument{
			Type: uintType,
		}},
	)

	return metadata
}

type EthereumAdapter struct {
	ChainID        *big.Int
	AdapterAddress common.Address
}

func (adapter *EthereumAdapter) ExecuteDagContract(
	schemaCid string,
	dataSourceCid string,
	variables string,
	contractMutation string,
	result string,
) (*DagTransaction, error) {

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		return nil, err
	}
	data, err := ExecuteDagContractAbiMethod().Inputs.Pack(schemaCid, dataSourceCid, variables, contractMutation, result)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	return &DagTransaction{
		SchemaCid:        schemaCid,
		DataSourceCid:    dataSourceCid,
		Variables:        variables,
		ContractMutation: contractMutation,
		Result:           result,
		Signature:        hexutil.Encode(signature),
	}, nil
}

func (adapter *EthereumAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
