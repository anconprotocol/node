package handler

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/buger/jsonparser"
	"github.com/status-im/go-waku/waku/v2/protocol"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
)

type AvailableDid string

const (
	DidTypeEthr AvailableDid = "ethr"
	DidTypeWeb  AvailableDid = "web"
	DidTypeKey  AvailableDid = "key"
)

const (
	defaultPath  = "/.well-known/did.json"
	documentPath = "/did.json"
)

type DidHandler struct {
	*sdk.AnconSyncContext
	WakuPeer     *WakuHandler
	RootKey      string
	Moniker      string
	ContentTopic protocol.ContentTopic
}

func NewDidHandler(ctx *sdk.AnconSyncContext,
	wakuPeer *WakuHandler,
	moniker string) *DidHandler {

	return &DidHandler{
		AnconSyncContext: ctx,
		WakuPeer:         wakuPeer,
		Moniker:          moniker,
		ContentTopic:     protocol.NewContentTopic(moniker, 1, "did", "json"),
	}

}

// BuildDidWeb ....
func (dagctx *DidHandler) BuildDidWeb(vanityName string, pubkey []byte) (*did.Doc, error) {
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

// BuildDidWeb ....
func (dagctx *DidHandler) BuildEthrDid(name string, pubkey []byte) (*did.Doc, error) {
	ti := time.Now()
	// did web
	base := []byte(name)
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
func (dagctx *DidHandler) BuildDidKey() (*did.Doc, error) {

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

func (dagctx *DidHandler) ReadDidWebUrl(c *gin.Context) {
	did := c.Param("did")

	// path := strings.Join([]string{"did:web:ipfs:user", did}, ":")
	p := types.GetUserPath(dagctx.Moniker)

	value, err := dagctx.Store.Get([]byte(p), did)
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
func (dagctx *DidHandler) ReadDid(c *gin.Context) {
	did := c.Param("did")
	// address, _, err := dagctx.ParseDIDWeb(did, true)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": fmt.Errorf("%v", err),
	// 	})
	// 	return
	// }
	p := types.GetUserPath(dagctx.Moniker)

	lnk, err := sdk.ParseCidLink(did)
	if err != nil {

		value, err := dagctx.Store.Get(([]byte(p)), did)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("did web not found %v", err),
			})
			return
		}

		if strings.HasPrefix(did, "raw:") || strings.HasPrefix(did, "did:") {
			c.JSON(200, json.RawMessage(value))
			return
		}

		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid hash %v", err),
		})
		return
	}

	n, err := dagctx.Store.Load(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, cidlink.Link{Cid: lnk.Cid})
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
	c.JSON(200, json.RawMessage(data))
}

func (dagctx *DidHandler) CreateDid(c *gin.Context) {
	var v map[string]string

	c.BindJSON(&v)
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

	pub := v["pub"]
	ethrdid := v["ethrdid"]
	domainName := v["domainName"]
	ethrpub, addr, err := types.RecoverKey((v["message"]), (v["signature"]))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})
		return
	}
	var cid datamodel.Link
	var key string
	if pub == "" {
		cid, key, err = dagctx.AddDid(DidTypeKey, "", addr, []byte{})
	} else if ethrdid != "" && pub != "" {
		cid, key, err = dagctx.AddDid(DidTypeEthr, ethrdid, addr, ethrpub)
	} else if domainName != "" && pub != "" {
		cid, key, err = dagctx.AddDid(DidTypeWeb, domainName, addr, ethrpub)

	} else {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})

		return

	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed to create did").Error(),
		})
	}

	p := types.GetUserPath(dagctx.Moniker)

	block := fluent.MustBuildMap(basicnode.Prototype.Map, 8, func(na fluent.MapAssembler) {
		na.AssembleEntry("issuer").AssignString(addr)
		na.AssembleEntry("timestamp").AssignInt(time.Now().Unix())
		na.AssembleEntry("contentHash").AssignLink(cid)
		na.AssembleEntry("signature").AssignString(v["signature"])
		na.AssembleEntry("key").AssignString(base64.StdEncoding.EncodeToString([]byte(key)))
		na.AssembleEntry("parent").AssignString(p)
	})
	res := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, block)

	resp, _ := sdk.Encode(block)
	if ethrdid != "" {
		dagctx.Store.Put([]byte(p), "raw:"+ethrdid, []byte(resp))
	}
	dagctx.WakuPeer.Publish(dagctx.ContentTopic, block)

	c.JSON(201, gin.H{
		"cid": res.String(),
	})
}

func (dagctx *DidHandler) AddDid(didType AvailableDid, domainName string, addr string, pubbytes []byte) (ipld.Link, string, error) {

	var didDoc *did.Doc
	var err error
	//	ctx := context.Background()

	if didType == DidTypeWeb {

		// exists, err := dagctx.Store.Has(ctx, domainName)
		// if err != nil {
		// 	return nil, "", fmt.Errorf("invalid domain name: %v", domainName)
		// }
		// if exists {
		// 	return nil, "", fmt.Errorf("invalid domain name: %v", domainName)
		// }

		didDoc, err = dagctx.BuildDidWeb(domainName, pubbytes)
		if err != nil {
			return nil, "", err
		}

	} else if didType == DidTypeEthr {

		didDoc, err = dagctx.BuildEthrDid(domainName, pubbytes)
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
	patch, err := jsonparser.Set(bz, []byte(fmt.Sprintf(`"%s"`, addr)), "verificationMethod", "[0]", "ethereumAddress")
	if err != nil {
		return nil, "", fmt.Errorf("invalid patch")
	}

	n, err := sdk.Decode(basicnode.Prototype.Any, string(patch))
	p := types.GetUserPath(dagctx.Moniker)

	lnk := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)
	if err != nil {
		return nil, "", err
	}

	//	dagctx.WakuPeer.Publish(dagctx.ContentTopic, n)

	dagctx.Store.Put([]byte(p), didDoc.ID, []byte(lnk.String()))
	dagctx.Store.Put([]byte(p), domainName, patch)

	return lnk, "", nil
}

func (dagctx *DidHandler) ParseDIDWeb(id string, useHTTP bool) (string, string, error) {
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
