package handler

import (
	"bytes"
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
	Proof *proofsignature.IavlProofService

	RootKey string
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

	didDoc, _ := types.GetDidDocument(string(doc))
	temp, _ := jsonparser.GetUnsafeString(v, "data")
	ok, err := types.Authenticate(didDoc, []byte(temp), signature)
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

	path, _ := jsonparser.GetString(v, "path")

	if path == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}

	var n datamodel.Node
	if isJSON {
		n, err = sdk.Decode(basicnode.Prototype.Any, string(data))
	} else {
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
		lnk, _ := sdk.ParseCidLink((from))
		na.AssembleEntry("issuer").AssignLink(lnk)
		na.AssembleEntry("timestamp").AssignInt(time.Now().Unix())
		na.AssembleEntry("content").AssignLink(cid)
		na.AssembleEntry("commitHash").AssignString(string(lastHash))
		na.AssembleEntry("height").AssignInt(blockNumber)
		na.AssembleEntry("signature").AssignString(signature)
		na.AssembleEntry("key").AssignString(base64.StdEncoding.EncodeToString([]byte(internalKey)))
		na.AssembleEntry("parent").AssignString(p)

	})
	res := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, block)

	c.JSON(201, gin.H{
		"cid": res,
	})
	pin, _ := jsonparser.GetString(v, "pin")

	if pin == "true" {
		impl.PushBlock(c.Request.Context(), dagctx.IPFSPeer, data, cid)
	}
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
