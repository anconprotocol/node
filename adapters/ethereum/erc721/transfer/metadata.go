package transfer

import (
	"bytes"
	"context"
	"fmt"

	"github.com/99designs/keyring"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	Keyring   keyring.Keyring
	ChainName string
	ChainID   int
}

func NewOnchainAdapter(keyring keyring.Keyring) OnchainAdapter {

	return OnchainAdapter{
		Keyring:   keyring,
		ChainName: "Ethereum",
		ChainID:   5,
	}
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
	// dag := ctx.Value("dag").(*handler.AnconSyncContext)

	// pk, has := os.LookupEnv("ETHEREUM_ADAPTER_KEY")
	key, err := adapter.Keyring.Get("ethereum")

	if err != nil {
		return nil, fmt.Errorf("environment key ETHEREUM_ADAPTER_KEY not found")
	}

	privateKey, err := crypto.HexToECDSA(string(key.Data))
	if err != nil {
		return nil, fmt.Errorf("invalid ETHEREUM_ADAPTER_KEY")
	}

	proofData, err :=
		ExecuteDagContractWithProofAbiMethod().Inputs.Pack(
			toAddress,
			tokenId,
			MetadataTransferPacket{
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

	signedPacket := &MetadataTransferProofPacket{
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

func (adapter *OnchainAdapter) GetTransaction(ctx context.Context, signedEthereumTx []byte) (types.Transaction, error) {
	tx := types.Transaction{}
	buff := bytes.NewReader(signedEthereumTx)
	streamInstance := rlp.NewStream(buff, 100000)
	err := tx.DecodeRLP(streamInstance)
	return tx, err
}
