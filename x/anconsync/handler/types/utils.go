package types

import (
	"bytes"
	"fmt"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

func GetDidDocument(data string) (*did.Doc, error) {
	return did.ParseDocument([]byte(data))
}

func Authenticate(didDoc *did.Doc, data []byte, sig []byte) (bool, error) {
	jsonWebKey := didDoc.VerificationMethods()
	id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	pub, _ := did.LookupPublicKey(id, didDoc)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256([]byte(msg))
	sig[64] -= 27

	rec, err := crypto.RecoverPubkey(sig, hash)
	addr := crypto.PubKeyToAddress(rec)

	bz, err := crypto.Ecrecover(hash, sig)
	fmt.Println(addr)

	k := bytes.Equal(bz, pub.Value)
	return k, err

}

func RecoverKey(data string, sig string) ([]byte, error) {

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256([]byte(msg))
	signature := hexutil.MustDecode(sig)
	// signature[64] -= 27

	rec, err := crypto.RecoverPubkey(signature, hash)
	addr := crypto.PubKeyToAddress(rec)

	bz, err := crypto.Ecrecover(hash, signature)
	fmt.Println(addr)

	return bz, err

}
