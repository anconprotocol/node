package handler

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/buger/jsonparser"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
)

type AvailableDid string

const (
	DidTypeWeb AvailableDid = "web"
	DidTypeKey AvailableDid = "key"
)

const (
	defaultPath  = "/.well-known/did.json"
	documentPath = "/did.json"
)

type Did struct {
	*sdk.AnconSyncContext
	Proof   *proofsignature.IavlProofService
	RootKey string
}

// BuildDidWeb ....
func (dagctx *Did) BuildDidWeb(vanityName string, pubkey []byte) (*did.Doc, error) {
	ti := time.Now()
	// did web
	base := append([]byte("did:web:ipfs:user:"), []byte(vanityName)...)
	// did web # id

	//Authentication method 2018
	didWebVer := did.NewVerificationMethodFromBytes(
		string(base),
		"Secp256k1VerificationKey2018",
		string(base),
		pubkey,
	)

	ver := []did.VerificationMethod{}
	ver = append(ver, *didWebVer)

	//	serv := []did.Service{{}, {}}

	// Secp256k1SignatureAuthentication2018
	auth := []did.Verification{{}}

	didWebAuthVerification := did.NewEmbeddedVerification(didWebVer, did.Authentication)

	auth = append(auth, *didWebAuthVerification)

	doc := did.BuildDoc(
		did.WithVerificationMethod(ver),
		///		did.WithService(serv),
		did.WithAuthentication(auth),
		did.WithCreatedTime(ti),
		did.WithUpdatedTime(ti),
	)
	doc.ID = string(base)
	return doc, nil
}

// BuildDidKey ....
func (dagctx *Did) BuildDidKey() (*did.Doc, error) {

	pubKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	ti := time.Now()
	multi := append([]byte(multicodec.Ed25519Pub.String()), pubKey...)
	code, _ := multibase.Encode(multibase.Base58BTC, multi)
	// did key
	base := append([]byte("did:key:"), code...)
	// did key # id
	id := append(base, []byte("#")...)

	didWebVer := did.NewVerificationMethodFromBytes(
		string(id),
		"Ed25519VerificationKey2018",
		string(base),
		[]byte(pubKey),
	)

	ver := []did.VerificationMethod{}
	ver = append(ver, *didWebVer)
	//	serv := []did.Service{{}, {}}

	// Secp256k1SignatureAuthentication2018
	auth := []did.Verification{{}}

	didWebAuthVerification := did.NewEmbeddedVerification(didWebVer, did.Authentication)

	auth = append(auth, *didWebAuthVerification)

	doc := did.BuildDoc(
		did.WithVerificationMethod(ver),
		did.WithAuthentication(auth),
		did.WithCreatedTime(ti),
		did.WithUpdatedTime(ti),
	)
	doc.ID = string(base)
	return doc, nil
}

func (dagctx *Did) ReadDidWebUrl(c *gin.Context) {
	did := c.Param("did")

	path := strings.Join([]string{"did:web:ipfs:user", did}, ":")

	value, err := dagctx.Store.DataStore.Get(c.Request.Context(), path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("did web not found %v", err),
		})
		return
	}

	lnk, err := sdk.ParseCidLink(string(value))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid hash %v", err),
		})
		return
	}

	n, err := dagctx.Store.Load(ipld.LinkContext{}, cidlink.Link{Cid: lnk.Cid})
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("block not found%v", err),
		})
		return
	}
	data, err := sdk.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed encoding %v", err),
		})
		return
	}
	c.JSON(200, data)

}
func (dagctx *Did) ReadDid(c *gin.Context) {
	did := c.Param("did")
	// address, _, err := dagctx.ParseDIDWeb(did, true)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": fmt.Errorf("%v", err),
	// 	})
	// 	return
	// }
	value, err := dagctx.Store.DataStore.Get(c.Request.Context(), did)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("did web not found %v", err),
		})
		return
	}

	lnk, err := sdk.ParseCidLink(string(value))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid hash %v", err),
		})
		return
	}

	n, err := dagctx.Store.Load(ipld.LinkContext{}, cidlink.Link{Cid: lnk.Cid})
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("block not found%v", err),
		})
		return
	}
	data, err := sdk.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed encoding %v", err),
		})
		return
	}
	c.JSON(200, data)
}

func (dagctx *Did) CreateDidKey(c *gin.Context) {
	var v map[string]string

	c.BindJSON(&v)
	// if v["pub"] == "" {
	// 	c.JSON(400, gin.H{
	// 		"error": fmt.Errorf("missing pub").Error(),
	// 	})
	// 	return
	// }

	domainName := ""
	pub := []byte{}
	cid, proof, err := dagctx.AddDid(DidTypeKey, domainName, pub)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})
	}
	c.JSON(201, gin.H{
		"cid":   cid,
		"proof": proof,
	})
}

func (dagctx *Did) CreateDidWeb(c *gin.Context) {
	var v map[string]string

	c.BindJSON(&v)
	if v["domainName"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing domainName").Error(),
		})
		return
	}
	if v["pub"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing pub").Error(),
		})
		return
	}
	if v["signature"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing signature").Error(),
		})
		return
	}
	if v["message"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing message").Error(),
		})
		return
	}

	domainName := v["domainName"]
	pub, err := types.RecoverKey((v["message"]), (v["signature"]))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})
	}
	cid, key, err := dagctx.AddDid(DidTypeWeb, domainName, pub)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})
	}
	commit, err := dagctx.Proof.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(commit, "root_hash")
	version, err := jsonparser.GetInt(commit, "version")
	lastHash := []byte(hash)

	_, err = dagctx.Proof.GetCommitmentProof([]byte(key), version)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(201, gin.H{
		"cid":    cid,
		"height": version,
		"hash":   lastHash,
		"key":    base64.StdEncoding.EncodeToString([]byte(key)),
	})
}

func (dagctx *Did) AddDid(didType AvailableDid, domainName string, pubbytes []byte) (ipld.Link, string, error) {

	var didDoc *did.Doc
	var err error
	ctx := context.Background()

	if didType == DidTypeWeb {

		exists, err := dagctx.Store.DataStore.Has(ctx, domainName)
		if err != nil {
			return nil, "", fmt.Errorf("invalid domain name: %v", domainName)
		}
		if exists {
			return nil, "", fmt.Errorf("invalid domain name: %v", domainName)
		}

		didDoc, err = dagctx.BuildDidWeb(domainName, pubbytes)
		if err != nil {
			return nil, "", err
		}

	} else if didType == DidTypeKey {
		didDoc, err = dagctx.BuildDidKey()
		if err != nil {
			return nil, "", err
		}

	} else {
		return nil, "", fmt.Errorf("Must create a did")
	}
	bz, err := didDoc.JSONBytes()
	n, err := sdk.Decode(basicnode.Prototype.Any, string(bz))
	lnk := dagctx.Store.Store(ipld.LinkContext{}, n)
	if err != nil {
		return nil, "", err
	}

	dagctx.Store.DataStore.Put(ctx, didDoc.ID, []byte(lnk.String()))

	// proofs
	//	key := fmt.Sprintf("%s/%s/user", "/anconprotocol", dagctx.RootKey)
	internalKey := fmt.Sprintf("%s/%s/user/%s", "/anconprotocol", dagctx.RootKey, lnk.String())
	_, err = dagctx.Proof.Set([]byte(internalKey), []byte(didDoc.ID))
	if err != nil {
		return nil, "", fmt.Errorf("invalid key")
	}

	if err != nil {
		return nil, "", fmt.Errorf("invalid commit")
	}

	return lnk, internalKey, nil
}

func (dagctx *Did) ParseDIDWeb(id string, useHTTP bool) (string, string, error) {
	var address, host string

	parsedDID, err := did.Parse(id)
	if err != nil {
		return address, host, fmt.Errorf("invalid did, does not conform to generic did standard --> %w", err)
	}

	pathComponents := strings.Split(parsedDID.MethodSpecificID, ":")

	pathComponents[0], err = url.QueryUnescape(pathComponents[0])
	if err != nil {
		return address, host, fmt.Errorf("error parsing did:web did")
	}

	host = strings.Split(pathComponents[0], ":")[0]

	protocol := "https://"
	if useHTTP {
		protocol = "http://"
	}

	switch len(pathComponents) {
	case 1:
		address = protocol + pathComponents[0] + defaultPath
	default:
		address = protocol + strings.Join(pathComponents, "/") + documentPath
	}

	return address, host, nil
}
