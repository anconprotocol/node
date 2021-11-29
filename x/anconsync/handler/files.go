package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cast"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
)

// @BasePath /v0
// FileWrite godoc
// @Summary Stores files
// @Schemes
// @Description Writes a raw block which syncs with IPFS. Returns a CID.
// @Tags file
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/file [post]
func (dagctx *DagContractTrustedContext) FileWrite(c *gin.Context) {
	w, fn, err := dagctx.Store.DataStore.PutStream(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error while getting stream. %v", err).Error(),
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error in form file %v", err).Error(),
		})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cannot open file. %v", err).Error(),
		})
		return
	}
	defer src.Close()
	// var bz []byte

	_, err = io.Copy(w, src)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed reading file. %v", err).Error(),
		})
		return
	}

	var bz []byte
	bz, _ = json.Marshal(file.Header)

	lnk := anconsync.CreateCidLink(bz)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error. %v", err).Error(),
		})
		return
	}
	fn(lnk.String())
	c.JSON(201, gin.H{
		"cid": lnk.String(),
	})
}

// @BasePath /v0
// FileRead godoc
// @Summary Reads JSON from a dag-json block
// @Schemes
// @Description Returns JSON
// @Tags file
// @Accept json
// @Produce json
// @Success 200
// @Router /v0/file/{cid}/{path} [get]
func (dagctx *DagContractTrustedContext) FileRead(c *gin.Context) {
	lnk, err := cid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error. %v", err).Error(),
		})
		return
	}
	reader, err := dagctx.Store.DataStore.GetStream(c.Request.Context(), lnk.String())

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error while getting stream. %v", err).Error(),
		})
		return
	}

	contentLength := cast.ToInt64(-1)
	contentType := c.ContentType()

	extraHeaders := map[string]string{
		//  "Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
