package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/buger/jsonparser"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/raw"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
)

type FileHandler struct {
	*sdk.AnconSyncContext
}

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
func (dagctx *FileHandler) FileWrite(c *gin.Context) {
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

	var w bytes.Buffer
	_, err = io.Copy(&w, src)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("failed reading file. %v", err).Error(),
		})
		return
	}

	n, err := DecodeNode(w.Bytes())
	lnk := dagctx.Store.Store(ipld.LinkContext{
		LinkPath: ipld.ParsePath(strings.Join([]string{"/", file.Filename}, "/")),
	}, n)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error. %v", err).Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"cid": lnk.String(),
	})
	// impl.PushBlockWithExtData(ctx, exchange, &items[0], rootKeyLink, gsync.ExtensionData{}, impl.SelectAll)
	// impl.PushBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, lnk)
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
func (dagctx *FileHandler) FileRead(c *gin.Context) {
	lnk, err := cid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error. %v", err).Error(),
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
	bz, err := EncodeNode(n)

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

	reader := bytes.NewReader(bz)
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
func EncodeNode(node ipld.Node) ([]byte, error) {
	var buffer bytes.Buffer
	err := raw.Encode(node, &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeNode(encoded []byte) (ipld.Node, error) {
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := raw.Decode(nb, bytes.NewReader(encoded)); err != nil {
		return nil, err
	}
	return nb.Build(), nil
}

// @BasePath /v0
// UploadContract godoc
// @Summary Upload hybrid smartcontracts
// @Schemes
// @Description Execute library smartcontracts.
// @Tags file
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/code [post]
func (dagctx *FileHandler) UploadContract(c *gin.Context) {

	v, _ := c.GetRawData()

	code, err := jsonparser.GetString(v, "body", "code")
	contract, err := hexutil.Decode(code)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error in code %v", err).Error(),
		})
		return
	}

	from, err := jsonparser.GetString(v, "body", "from")
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error in from %v", err).Error(),
		})
		return
	}

	signature, err := jsonparser.GetString(v, "body", "signature")
	sig := hexutil.MustDecode(signature)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error in signature %v", err).Error(),
		})
		return
	}

	didCid, err := dagctx.Store.DataStore.Get(c.Request.Context(), from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("did not found %v", err),
		})
		return
	}

	didDoc, err := types.GetDidDocument(string(didCid), &dagctx.Store)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("error in did Document %v", err).Error(),
		})
		return
	}

	hash := fmt.Sprintf(`%s%s`, code, from)
	ok, err := types.Authenticate(didDoc, []byte(hash), sig)

	if !ok || err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("Error parsing signer %v", err),
		})
		return
	}

	js := fmt.Sprintf(`{"code": "%s"}`, hexutil.Encode(contract))
	n, err := sdk.Decode(basicnode.Prototype.Any, js)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	cid := dagctx.Store.Store(ipld.LinkContext{}, n)
	c.JSON(201, gin.H{
		"address": cid,
	})

	impl.PushBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, cid)
}
