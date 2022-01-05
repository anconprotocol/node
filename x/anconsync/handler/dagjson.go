package handler

import (
	"bytes"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"context"
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
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

type DagJsonHandler struct {
	*sdk.AnconSyncContext
	Proof   *proofsignature.IavlProofService
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
	// TODO:
	// {

	// 	metadata: object,
	// 	ts: datetime,
	// 	did: user,
	// 	proofLink: "ipfs://babaaaf"
	// }

	// DAG Store
	// Set merkle tree
	// Commit
	// Get Proof
	// cid = Focused Transform (proofLink)

	// return cid, ...
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
	doc, err := dagctx.Store.DataStore.Get(context.Background(), from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing did").Error(),
		})
		return
	}

	p := fmt.Sprintf("%s/%s/user", "/anconprotocol", dagctx.RootKey)

	didDoc, _ := types.GetDidDocument(string(doc))
	sig, _ := hexutil.Decode(signature)
	ok, err := types.Authenticate(didDoc, data, sig)
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
	internalKey := fmt.Sprintf("%s/%s", p, cid)
	dagctx.Proof.Set([]byte(internalKey), data)
	commithash, _ := dagctx.Proof.SaveVersion(&emptypb.Empty{})
	proof, err := dagctx.Proof.GetCommitmentProof([]byte(internalKey))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("read Error %v", err).Error(),
		})
		return
	}
	existPayload, _, _, err := jsonparser.Get(proof, "proof", "proofs", "[0]", "Proof", "exist")

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("parse Error %v", err).Error(),
		})
		return
	}
	proofnode, err := sdk.Decode(basicnode.Prototype.Any, string(existPayload))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}
	block := fluent.MustBuildMap(basicnode.Prototype.Map, 7, func(na fluent.MapAssembler) {
		lnk, _ := sdk.ParseCidLink((from))
		na.AssembleEntry("issuer").AssignLink(lnk)
		na.AssembleEntry("timestamp").AssignInt(time.Now().Unix())
		na.AssembleEntry("content").AssignLink(cid)
		na.AssembleEntry("commitHash").AssignString(string(commithash))
		na.AssembleEntry("signature").AssignString(signature)
		na.AssembleEntry("proof").AssignNode(proofnode)
		na.AssembleEntry("key").AssignString(p)

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
