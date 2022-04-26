package types

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/buger/jsonparser"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
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

func ResolveDIDDoc(did string) ([]byte, error) {
	return Read((did))
}

func resolveDID(id string) ([]byte, error) {
	// u, err := url.ParseRequestURI(fmt.Sprint("https://dev.uniresolver.io/1.0/identifiers/", url.PathEscape(id)))
	resp, err := http.Get(fmt.Sprintf("https://dev.uniresolver.io/1.0/identifiers/%s", id))
	if err != nil {
		return nil, fmt.Errorf("HTTP create get request failed: %w", err)
	}

	var gotBody []byte

	gotBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode == 200 {
		return gotBody, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("DID does not exist for request")
	}
	defer resp.Body.Close()

	return nil, fmt.Errorf("unsupported response from DID resolver [%v] header [%s] body [%s]",
		resp.StatusCode, resp.Header.Get("Content-type"), gotBody)
}

// Read implements didresolver.DidMethod.Read interface (https://w3c-ccg.github.io/did-resolution/#resolving-input)
func Read(didID string, _ ...vdrapi.DIDMethodOption) ([]byte, error) {

	data, err := resolveDID(didID)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, vdrapi.ErrNotFound
	}

	// documentResolution, err := did.ParseDocumentResolution(data)
	// if err != nil {
	// 	if !errors.Is(err, did.ErrDIDDocumentNotExist) {
	// 		return nil, err
	// 	}
	// } else {
	// 	return documentResolution, nil
	// }

	return data, nil
}

func GetDidVerificationMethod(data []byte) (*did.VerificationMethod, bool) {
	didDoc, _ := did.ParseDocument(data)
	jsonWebKey := didDoc.VerificationMethods()
	id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	return did.LookupPublicKey(id, didDoc)
}

func IsValidSignature(didoc []byte, data []byte, signature string) (bool, error) {
	vtype, err := jsonparser.GetString(didoc, "didDocument", "verificationMethod", "[0]", "type")
	if err != nil {
		return false, fmt.Errorf("invalid did, missing ethereumAddress")
	}

	if vtype == "Secp256k1VerificationKey2018" || vtype == "EcdsaSecp256k1RecoveryMethod2020" {
		return true, nil
		// return Authenticate(didoc, data, (signature))
	} else if vtype == "Ed25519VerificationKey2018" {
		// return ed25519.Verify(verificationMethod.Value, data, ([]byte(signature))), nil
		return true, nil
	}
	return false, nil
}

func GetDidDocumentAuthentication(data []byte) (*ecdsa.PublicKey, error) {
	didDoc, _ := did.ParseDocument(data)
	jsonWebKey := didDoc.VerificationMethods()
	id := jsonWebKey[did.Authentication][0].VerificationMethod.ID
	pub, _ := did.LookupPublicKey(id, didDoc)

	return crypto.ParsePublicKey(pub.Value)

}
func Authenticate(diddoc []byte, data []byte, sig string) (bool, error) {
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
