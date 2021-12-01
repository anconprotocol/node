package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/spf13/cast"
)

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
func (dagctx *AnconSyncContext) DagCborWrite(c *gin.Context) {
	var v map[string]string

	c.BindJSON(&v)
	if v["path"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}
	if v["data"] == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing payload data source").Error(),
		})
		return
	}

	buff, _ := base64.StdEncoding.DecodeString(v["data"])

	n, err := anconsync.DecodeCBOR(basicnode.Prototype.Any, buff)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	cid := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(v["path"])}, n)
	c.JSON(201, gin.H{
		"cid": cid,
	})
	PushBlock(c.Request.Context(), dagctx, cid)
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
func (dagctx *AnconSyncContext) DagCborRead(c *gin.Context) {
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
	data, err := anconsync.EncodeCBOR(n)
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
