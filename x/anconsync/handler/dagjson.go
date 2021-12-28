package handler

import (
	"bytes"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"google.golang.org/protobuf/types/known/emptypb"

	"context"
	"encoding/json"
	"fmt"

	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

type DagJsonHandler struct {
	*sdk.AnconSyncContext
	Proof    *proofsignature.IavlProofService
	RootHash string
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

	temp, _ := jsonparser.GetUnsafeString(v, "data")
	//		temp = strings.ReplaceAll(temp, "\n", "")
	// temp = strings.ReplaceAll(temp, "\\","\"")
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(temp))
	data := buf.Bytes()
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing payload data source").Error(),
		})
		return
	}
	didCid, err := dagctx.Store.DataStore.Get(context.Background(), from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing did").Error(),
		})
		return
	}

	p := fmt.Sprintf("/anconprotocol/%s/%s/%s", dagctx.RootHash, "user", didCid)

	didDoc, err := types.GetDidDocument(string(didCid), &dagctx.Store)
	hash := crypto.Keccak256([]byte(data))
	sig := []byte(signature)
	ok, err := types.Authenticate(didDoc, hash, sig)
	if ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid signature").Error(),
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

	n, err := sdk.Decode(basicnode.Prototype.Any, string(data))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}
	cid := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)
	p = fmt.Sprintf("%s/%s", p, cid)
	dagctx.Proof.Set([]byte(p), data)
	dagctx.Proof.SaveVersion(&emptypb.Empty{})
	c.JSON(201, gin.H{
		"cid": cid,
	})
	impl.PushBlock(c.Request.Context(), dagctx.IPFSPeer, data, cid)
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
	lnk, err := cid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	p := fmt.Sprintf("%s/%s/user", "/anconprotocol", dagctx.RootHash)

	n, err := dagctx.Store.Load(ipld.LinkContext{
		LinkPath: ipld.ParsePath(p),
	}, cidlink.Link{Cid: lnk})

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
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
	c.JSON(200, data)
}
