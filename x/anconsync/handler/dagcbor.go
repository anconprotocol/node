package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DagCborHandler struct {
	*sdk.AnconSyncContext
	Proof    *proofsignature.IavlProofService
	RootHash string
}

// @BasePath /v0

// DagCborWrite godoc
// @Summary Stores CBOR as dag-json
// @Schemes
// @Description Writes a dag-cbor block which syncs with IPFS. Returns a CID.
// @Tags dag-cbor
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/dagcbor [post]
func (dagctx *DagCborHandler) DagCborWrite(c *gin.Context) {
	var v map[string]string

	c.BindJSON(&v)
	if v["path"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}
	if v["from"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing from").Error(),
		})
		return
	}
	if v["signature"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing signature").Error(),
		})
		return
	}
	if v["data"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing payload data source").Error(),
		})
		return
	}

	data, _ := base64.StdEncoding.DecodeString(v["data"])
	didCid, err := dagctx.Store.DataStore.Get(context.Background(), v["from"])
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing did").Error(),
		})
		return
	}

	didDoc, err := types.GetDidDocument(string(didCid), &dagctx.Store)
	hash := crypto.Keccak256([]byte(data))
	sig := []byte(v["signature"])
	ok, err := types.Authenticate(didDoc, hash, sig)
	if ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid signature").Error(),
		})
		return
	}

	p := fmt.Sprintf("%s/%s/%s", dagctx.RootHash, "user", didCid)

	n, err := sdk.DecodeCBOR(basicnode.Prototype.Any, data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	cid := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)
	dagctx.Proof.Set([]byte(p), data)
	dagctx.Proof.SaveVersion(&emptypb.Empty{})

	c.JSON(201, gin.H{
		"cid": cid,
	})

	impl.PushBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, cid)
}

// @BasePath /v0
// DagCborRead godoc
// @Summary Reads CBOR from a dag-cbor block
// @Schemes
// @Description Returns CBOR
// @Tags dag-cbor
// @Accept json
// @Produce json
// @Success 200
// @Router /v0/dagcbor/{cid}/{path} [get]
func (dagctx *DagCborHandler) DagCborRead(c *gin.Context) {
	lnk, err := cid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	n, err := dagctx.Store.Load(ipld.LinkContext{LinkPath: ipld.ParsePath(c.Param("path"))}, cidlink.Link{Cid: lnk})
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	data, err := sdk.EncodeCBOR(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	reader := bytes.NewReader(data)
	contentLength := cast.ToInt64(reader.Len())
	contentType := "application/cbor"

	extraHeaders := map[string]string{
		//  "Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
