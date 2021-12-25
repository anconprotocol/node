package types

import (
	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/sdk"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

func GetDidDocument(data string, s *sdk.Storage) (*did.Doc, error) {
	return did.ParseDocument([]byte(data))
}

func Authenticate(didDoc *did.Doc, hash []byte, sig []byte) (bool, error) {
	jsonWebKey := didDoc.VerificationMethods()
	id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	pub, _ := did.LookupPublicKey(id, didDoc)

	bz := crypto.Keccak256([]byte(hash))
	ok, err := crypto.Ecrecover(bz, sig)

	return hexutil.Encode(ok) == hexutil.Encode(pub.Value), err

}
