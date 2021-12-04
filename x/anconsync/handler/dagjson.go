package handler

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/anconprotocol/node/x/anconsync"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

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
func (dagctx *AnconSyncContext) DagJsonWrite(c *gin.Context) {

	v, _ := c.GetRawData()

	path, _ := jsonparser.GetString(v, "path")

	if path == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
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

	n, err := anconsync.Decode(basicnode.Prototype.Any, string(data))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}
	cid := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(path)}, n)
	c.JSON(201, gin.H{
		"cid": cid,
	})
	PushBlock(c.Request.Context(), dagctx, cid)
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
func (dagctx *AnconSyncContext) DagJsonRead(c *gin.Context) {
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
	data, err := anconsync.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	c.JSON(200, data)
}
