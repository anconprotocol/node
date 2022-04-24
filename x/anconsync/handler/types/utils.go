package types

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/buger/jsonparser"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

func GetNetworkPath(moniker string) string {
	return fmt.Sprintf("%s", moniker)
}

func GetGraphPath(moniker string) string {
	return fmt.Sprintf("%s/graphs", moniker)
}

func GetUserPath(moniker string) string {
	return fmt.Sprintf("%s/users", moniker)
}

func GetDidDocument(data string) (*did.Doc, error) {
	bz := []byte(data)
	return did.ParseDocument(bz)
}

func GetDidDocumentAuthentication(data []byte) (*ecdsa.PublicKey, error) {
	didDoc, _ := did.ParseDocument(data)
	jsonWebKey := didDoc.VerificationMethods()
	id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	pub, _ := did.LookupPublicKey(id, didDoc)

	return crypto.ParsePublicKey(pub.Value)

}
func 			Authenticate(diddoc []byte, data []byte, sig string) (bool, error) {
	// jsonWebKey := didDoc.VerificationMethods()
	// id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	// pub, _ := did.LookupPublicKey(id, didDoc)
	// addrrec := pub.ID

	addrrec, err := jsonparser.GetString((diddoc), "verificationMethod", "[0]", "ethereumAddress")
	if err != nil {
		return false, fmt.Errorf("invalid did, missing ethereumAddress")
	}
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256([]byte(msg))
	signature := hexutil.MustDecode(sig)
	signature[64] -= 27

	rec, err := crypto.RecoverPubkey(signature, hash)
	addr := crypto.PubKeyToAddress(rec)

	fmt.Println(addr)

	return addrrec == addr.String(), err

}

func RecoverKey(data string, sig string) ([]byte, string, error) {

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256([]byte(msg))
	signature := hexutil.MustDecode(sig)
	signature[64] -= 27

	rec, err := crypto.RecoverPubkey(signature, hash)
	addr := crypto.PubKeyToAddress(rec)

	bz, err := crypto.Ecrecover(hash, signature)
	fmt.Println(addr)

	return bz, addr.String(), err

}
