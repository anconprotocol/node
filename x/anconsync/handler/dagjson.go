package handler

import (
	"bytes"
	"log"
	"strings"
	"time"

	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/buger/jsonparser"
	ecies "github.com/ecies/go"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/spf13/cast"

	"google.golang.org/protobuf/types/known/emptypb"
)

type DagJsonHandler struct {
	*sdk.AnconSyncContext
	Proof    *proofsignature.IavlProofService
	IPFSHost string
	RootKey  string
}

// @BasePath /v0
// DagJsonWrite godoc
// @Summary Stores JSON as dag-json
// @Schemes
// @Description Writes a dag-json block which syncs with IPFS. Returns a CID.
// @Tags dag-json
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/dagjson [post]
func (dagctx *DagJsonHandler) DagJsonWrite(c *gin.Context) {

	v, _ := c.GetRawData()

	from, _ := jsonparser.GetString(v, "from")

	if from == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing from").Error(),
		})
		return
	}
	signature, _ := jsonparser.GetString(v, "signature")

	if signature == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing signature").Error(),
		})
		return
	}

	doc, err := dagctx.Store.DataStore.Get(context.Background(), from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing did").Error(),
		})
		return
	}

	p := fmt.Sprintf("%s/%s/user", "/anconprotocol", dagctx.RootKey)

	temp, _ := jsonparser.GetUnsafeString(v, "data")
	ok, err := types.Authenticate(doc, []byte(temp), signature)
	if !ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid signature").Error(),
		})
		return
	}

	data, err := hexutil.Decode(temp)
	var buf bytes.Buffer
	isJSON := false
	if err != nil {
		isJSON = true
		err = json.Compact(&buf, []byte(temp))
		data = buf.Bytes()
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing payload data source").Error(),
		})
		return
	}

	encrypt, _ := jsonparser.GetString(v, "encrypt")
	hasEncrypt := encrypt == "true"
	var authorizedRecipients []string

	if hasEncrypt {
		recipients, _ := jsonparser.GetString(v, "authorizedRecipients")
		authorizedRecipients = strings.Split(recipients, ",")

		content, err := dagctx.Store.DataStore.Get(c.Request.Context(), authorizedRecipients[0])
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("authorized recipient not found").Error(),
			})
			return
		}

		pub, err := types.GetDidDocumentAuthentication((content))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing recipient public key").Error(),
			})
			return
		}

		data, err = ecies.Encrypt((*ecies.PublicKey)(pub), data)
		fmt.Println(data)

		if err != nil {
			log.Printf("failed to encrypt payload: %s", err)
			return
		}

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error while loading recipient public keys").Error(),
			})
			return
		}

	}

	path, _ := jsonparser.GetString(v, "path")

	if path == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}

	var n datamodel.Node
	if isJSON && !hasEncrypt {
		n, err = sdk.Decode(basicnode.Prototype.Any, string(data))
	} else {
		// TODO: fix
		n = basicnode.NewBytes(data)
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	cid := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)
	internalKey := fmt.Sprintf("%s/%s", p, cid)
	dagctx.Proof.Set([]byte(internalKey), data)
	commit, err := dagctx.Proof.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(commit, "root_hash")
	version, err := jsonparser.GetInt(commit, "version")
	lastHash := []byte(hash)
	blockNumber := cast.ToInt64(version)
	block := fluent.MustBuildMap(basicnode.Prototype.Map, 7, func(na fluent.MapAssembler) {
		addrrec, err := jsonparser.GetString((doc), "verificationMethod", "[0]", "ethereumAddress")
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("invalid did %v", err).Error(),
			})
			return
		}
	
		na.AssembleEntry("issuer").AssignString(addrrec)
		na.AssembleEntry("timestamp").AssignInt(time.Now().Unix())
		na.AssembleEntry("content").AssignLink(cid)
		na.AssembleEntry("commitHash").AssignString(string(lastHash))
		na.AssembleEntry("height").AssignInt(blockNumber)
		na.AssembleEntry("signature").AssignString(signature)
		na.AssembleEntry("key").AssignString(base64.StdEncoding.EncodeToString([]byte(internalKey)))
		na.AssembleEntry("parent").AssignString(p)
	})
	res := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, block)

	resp, _ := sdk.Encode(block)
	tx, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp))

	resp2, _ := sdk.Encode(n)
	m, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp2))

	// c1, _ := sdk.ParseCidLink(m)
	// c2, _ := sdk.ParseCidLink(tx)
	// impl.FetchBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, c1)
	// impl.FetchBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, c2)
	c.JSON(201, gin.H{
		"cid": res.String(),
		"ipfs": map[string]interface{}{
			"metadata": m,
			"tx":       tx,
		},
	})
}

// @BasePath /v0
// DagJsonRead godoc
// @Summary Reads JSON from a dag-json block
// @Schemes
// @Description Returns JSON
// @Tags dag-json
// @Accept json
// @Produce json
// @Success 200
// @Router /v0/dagjson/{cid}/{path} [get]
func (dagctx *DagJsonHandler) DagJsonRead(c *gin.Context) {
	lnk, err := sdk.ParseCidLink(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	p := fmt.Sprintf("%s/%s/user", "/anconprotocol", dagctx.RootKey)

	n, err := dagctx.Store.Load(ipld.LinkContext{
		LinkPath: ipld.ParsePath(p),
	}, lnk)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	data, err := sdk.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	c.JSON(200, json.RawMessage(data))
}
