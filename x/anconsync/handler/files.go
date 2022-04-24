package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
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
	RootKey string
	Moniker string
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
// @Router /v1/file [post]
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
	p := types.GetUserPath(dagctx.Moniker)
	lnk := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)

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
// @Router /v1/file/{cid}/{path} [get]
func (dagctx *FileHandler) FileRead(c *gin.Context) {
	lnk, err := cid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error. %v", err).Error(),
		})
		return
	}
	p := types.GetUserPath(dagctx.Moniker)

	n, err := dagctx.Store.Load(ipld.LinkContext{
		LinkPath: ipld.ParsePath(p),
	}, cidlink.Link{Cid: lnk})

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
// @Router /v1/code [post]
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

	hash := fmt.Sprintf(`%s%s`, code, from)
	ok, err := types.Authenticate(didCid, []byte(hash), string(sig))

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

	p := types.GetUserPath(dagctx.Moniker)
	lnk := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, n)
	c.JSON(201, gin.H{
		"address": lnk.String(),
	})

	// impl.PushBlock(c.Request.Context(), "https://ipfs.xdv.digital", []byte(js), cid)
}
