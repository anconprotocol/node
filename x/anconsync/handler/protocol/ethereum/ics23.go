package ethereum

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/umbracle/go-web3/abi"
)

type EnrollL2Account struct {
	did           []byte
	innerOpPrefix []byte
	key           []byte
	innerOpSuffix []byte
	prefix        []byte
}
type VerifyProof struct {
	value        []byte
	innerOpPrefix []byte
	key           []byte
	innerOpSuffix []byte
	prefix        []byte
}
type SubmitPacketWithProof struct {
	packet        []byte
	innerOpPrefix []byte
	key           []byte
	innerOpSuffix []byte
	prefix        []byte
}

func SubmitPacketWithProofAbi() *abi.Method {

	m, err := abi.NewMethod("submitPacketWithProof(string key, string packet,string prefix,string innerOpPrefix,string innerOpSuffix)")

	if err != nil {
		panic(err)
	}

	return m
}
func EnrollL2AccountAbi() *abi.Method {

	m, err := abi.NewMethod("enrollL2Account(string key, string did,string prefix,string innerOpPrefix,string innerOpSuffix)")

	if err != nil {
		panic(err)
	}

	return m
}

func VerifyProofAbi() *abi.Method {

	m, err := abi.NewMethod("verifyProof(string key, string value,string prefix,string innerOpPrefix,string innerOpSuffix)")

	if err != nil {
		panic(err)
	}

	return m
}

type OnchainAdapter struct {
	From                   string
	HostAddress            string
	DestinationHostAddress string
	VerifierAddress        string
	SubmitterAddress       string
	ChainName              string
	ChainID                int
}

func NewOnchainAdapter(from string, chainName string, chainID int) *OnchainAdapter {
	return &OnchainAdapter{
		From:      from,
		ChainName: chainName,
		ChainID:   chainID,
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

// func (adapter *OnchainAdapter) ApplyRequestWithProof(
// 	ctx context.Context,
// 	metadataCid string,
// 	resultCid string,
// 	fromOwner string,
// 	toOwner string,
// 	toAddress string,
// 	tokenId string,
// 	prefix string,
// ) ([]byte, string, error) {

// 	id := (tokenId)
// 	var proof []byte
// 	keccak.Keccak256(proof, encodePacked(
// 		// Current metadata cid
// 		[]byte(metadataCid),
// 		// Current owner (opaque)
// 		[]byte(fromOwner),
// 		// Updated metadata cid
// 		[]byte(resultCid),
// 		// New owner address
// 		[]byte(toOwner),
// 		// Token Address
// 		[]byte(toAddress),
// 		// Token Id
// 		[]byte(id),
// 		// Contract Prefix
// 		[]byte(prefix)))

// 	unsignedProofData := encodePacked(
// 		[]byte("\x19Ethereum Signed Message:\n32"),
// 		// Proof
// 		proof)

// 	var hash []byte
// 	keccak.Keccak256(hash, unsignedProofData)

// 	return nil, resultCid, nil
// }
func (adapter *OnchainAdapter) EnrollL2Account(
	updatedProof *EncodePackedExistenceProof,
	value []byte,
	data []byte,
) ([]byte, error) {

	packet := &EnrollL2Account{
		proof: encodePacked(
			updatedProof.Prefix,
			updatedProof.InnerOpPrefix,
			updatedProof.InnerOpSuffix,
			i32tob((uint32(updatedProof.InnerOpHashOp))),
		),
		key:   updatedProof.Key,
		value: value,
		data:  data,
	}

	signedProofData, err := EnrollL2AccountAbi().Inputs.Encode(packet)

	if err != nil {
		return nil, fmt.Errorf("packing for signature proof generation failed")
	}

	return signedProofData, nil
}

func (adapter *OnchainAdapter) ApplyRequestWithProof(
	updatedProof *EncodePackedExistenceProof,
	value []byte,
	data []byte,
) ([]byte, error) {

	packet := &SubmitPacketWithProof{
		ops: updatedProof.LeafOp,
		proof: encodePacked(
			updatedProof.Prefix,
			updatedProof.InnerOpPrefix,
			updatedProof.InnerOpSuffix,
			i32tob((uint32(updatedProof.InnerOpHashOp))),
		),
		key:   updatedProof.Key,
		value: value,
		data:  data,
	}

	signedProofData, err := SubmitPacketWithProofAbi().Inputs.Encode(packet)

	if err != nil {
		return nil, fmt.Errorf("packing for signature proof generation failed")
	}

	return signedProofData, nil
}

func (adapter *OnchainAdapter) GenerateVerificationProof(
	proof *EncodePackedExistenceProof,
	root []byte,
	value []byte,
) ([]byte, error) {

	packet := &VerifyProof{
		ops: proof.LeafOp,
		proof: encodePacked(
			proof.Prefix,
			proof.InnerOpPrefix,
			proof.InnerOpSuffix,
			i32tob((uint32(proof.InnerOpHashOp))),
		),
		root:  root,
		key:   proof.Key,
		value: value,
		data:  nil,
	}

	signedProofData, err := VerifyProofAbi().Inputs.Encode(packet)

	if err != nil {
		return nil, fmt.Errorf("packing for signature proof generation failed")
	}

	return signedProofData, nil
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
