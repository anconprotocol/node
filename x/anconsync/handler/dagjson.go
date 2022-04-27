package handler

import (
	"bytes"
	"strings"
	"time"

	"context"
	"encoding/json"
	"fmt"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-graphsync/ipldutil"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/status-im/go-waku/waku/v2/node"
	"github.com/status-im/go-waku/waku/v2/protocol"
)

type DagJsonHandler struct {
	*sdk.AnconSyncContext
	ProofHandler  *ProofHandler
	WakuPeer      *WakuHandler
	RootKey       string
	Moniker       string
	PreviousBlock datamodel.Link
	ContentTopic  protocol.ContentTopic
}
type Mutation struct {
	Path          string
	PreviousValue string
	NextValue     interface{}
	NextValueKind datamodel.Kind
}

func NewDagHandler(ctx *sdk.AnconSyncContext,
	proof *ProofHandler,
	wakuPeer *WakuHandler,
	rootKey string,
	moniker string) *DagJsonHandler {

	return &DagJsonHandler{
		AnconSyncContext: ctx,
		ProofHandler:     proof,
		WakuPeer:         wakuPeer,
		RootKey:          rootKey,
		Moniker:          moniker,
		ContentTopic:     protocol.NewContentTopic(moniker, 1, "dag", "json"),
	}

}

func (h *DagJsonHandler) Listen(ctx context.Context) {
	sub, err := h.WakuPeer.Subscribe(ctx, h.ContentTopic.String())

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	for value := range sub.C {

		if value.Message().ContentTopic == h.ContentTopic.String() {
			payload, err := node.DecodePayload(value.Message(), &node.KeyInfo{Kind: node.None})
			if err != nil {
				fmt.Println(err)
				return
			}
			// Decode payload
			block, err := ipldutil.DecodeNode(payload.Data)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Get event and cid properties
			node, err := block.LookupByString("event")
			if err != nil {
				fmt.Println(err)
				return
			}

			eventType := must.String(node)
			node, err = block.LookupByString("cid")
			if err != nil {
				fmt.Println(err)
				return
			}

			key := must.String(node)
			has, err := h.Store.DataStore.Has(ctx, key)

			if err != nil {
				fmt.Println(err)
				return
			}

			// If get, lookup and return block, otherwise put / store
			if eventType == "get" {
				if has {
					lnk, _ := sdk.ParseCidLink(key)
					h.Store.Load(ipld.LinkContext{}, lnk)
					h.WakuPeer.Publish(h.ContentTopic, block)
				}
			} else if eventType == "put" {
				// store the payload
				if !has {
					h.Store.Store(ipld.LinkContext{}, block)
				}
			}
		}

	}
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
// @Router /v1/dagjson [post]
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

	// p := fmt.Sprintf("%s/%s", types.GetUserPath(dagctx.Moniker), from)
	hexdata, _ := jsonparser.GetString(v, "data")

	temp, _ := jsonparser.GetUnsafeString(v, "data")
	data, err := hexutil.Decode(hexdata)
	data = []byte(hexdata)
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
	resolution, err := types.ResolveDIDDoc(from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did").Error(),
		})
		return
	}

	// if DID is valid, assume sigature is ok?
	ok, err := types.IsValidSignature(resolution, data, signature)
	if !ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did signature").Error(),
		})
		return
	}

	///	parentHash, _ := jsonparser.GetString(v, "parent")
	path, _ := jsonparser.GetString(v, "path")

	if path == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}

	digest := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)))
	var n datamodel.Node
	if isJSON {
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
	cid := dagctx.Store.Store(ipld.LinkContext{}, n)
	// get current

	lastHash, _ := dagctx.ProofHandler.proofs.GetCurrentVersion()

	dbl := &DagBlockResult{
		Issuer:        from,
		Timestamp:     time.Now().Unix(),
		ContentHash:   cid,
		Signature:     signature,
		Digest:        hexutil.Encode(digest),
		Network:       dagctx.Moniker,
		LastBlockHash: string(lastHash),
	}
	block := dagctx.ApplyDagBlock(dbl)
	key := dagctx.Store.LinkSystem.MustComputeLink(sdk.GetDagJSONLinkPrototype(), block)
	bz, err := block.AsBytes()
	err = dagctx.Store.DataStore.Put(c.Request.Context(), (key.Binary()), bz)

	// block, err = dagctx.Store.Load(ipld.LinkContext{
	// 	LinkPath: ipld.ParsePath("/"),
	// }, key)

	// fmt.Println(block, err)
	dagctx.ProofHandler.AddToPool(cid.String())
	// dagctx.PreviousBlock = res
	topic, err := jsonparser.GetString(v, "topic")
	//	dagctx.Store.DataStore.Put(c.Request.Context(), fmt.Sprintf("block:%d", blockNumber), []byte(res.String()))
	contentTopic, err := protocol.StringToContentTopic(topic)

	// block
	dagctx.WakuPeer.Publish(contentTopic, block)

	// metadata
	dagctx.WakuPeer.Publish(contentTopic, n)

	c.JSON(201, gin.H{
		"cid": key.String(),
	})
}

type DagBlock struct {
	Issuer      string         `json:"issuer"`
	Timestamp   int64          `json:"timestamp"`
	ContentHash datamodel.Link `json:"content_hash"`
	Signature   string         `json:"signature"`
	Digest      string         `json:"digest"`
	Network     string         `json:"network"`
}

func (dagctx *DagJsonHandler) ApplyDagBlock(args *DagBlockResult) datamodel.Node {

	block := fluent.MustBuildMap(basicnode.Prototype.Map, 13, func(na fluent.MapAssembler) {
		na.AssembleEntry("issuer").AssignString(args.Issuer)
		na.AssembleEntry("timestamp").AssignInt(args.Timestamp)
		na.AssembleEntry("contentHash").AssignLink(args.ContentHash)
		na.AssembleEntry("signature").AssignString(args.Signature)
		na.AssembleEntry("digest").AssignString(args.Digest)
		na.AssembleEntry("network").AssignString(args.Network)

		if args.LastBlockHash != "" {
			lnk, _ := sdk.ParseCidLink(args.LastBlockHash)
			na.AssembleEntry("lastBlockHash").AssignLink(lnk)
		}
	})

	return block
}

func (dagctx *DagJsonHandler) ApplyFocusedTransform(node datamodel.Node, mutations []Mutation) (datamodel.Node, error) {
	var current datamodel.Node
	var err error
	current = node

	prog := traversal.Progress{
		Cfg: &traversal.Config{
			LinkSystem:                     dagctx.Store.LinkSystem,
			LinkTargetNodePrototypeChooser: basicnode.Chooser,
		},
	}
	for _, m := range mutations {
		current, err = prog.FocusedTransform(
			current,
			datamodel.ParsePath(m.Path),
			func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {

				// update
				if prev != nil && !prev.IsAbsent() && progress.Path.String() == m.Path && must.String(prev) == (m.PreviousValue) {
					nb := prev.Prototype().NewBuilder()
					switch prev.Kind() {
					case datamodel.Kind_Float:
						nb.AssignFloat(m.NextValue.(float64))
					case datamodel.Kind_Bytes:
						nb.AssignBytes(m.NextValue.([]byte))
					case datamodel.Kind_Int:
						nb.AssignInt(m.NextValue.(int64))
					case datamodel.Kind_Link:
						lnk, err := sdk.ParseCidLink(m.NextValue.(string))
						if err != nil {
							return nil, err
						}
						nb.AssignLink(lnk)
					case datamodel.Kind_Bool:
						nb.AssignBool(m.NextValue.(bool))
					default:
						nb.AssignString(m.NextValue.(string))
					}
					return nb.Build(), nil
				} else if progress.Path.String() == m.Path && m.PreviousValue == "" {
					// previous is absent, add

					nb := basicnode.Prototype.Any.NewBuilder()
					switch m.NextValueKind {
					case datamodel.Kind_Float:
						nb.AssignFloat(m.NextValue.(float64))
					case datamodel.Kind_Bytes:
						nb.AssignBytes(m.NextValue.([]byte))
					case datamodel.Kind_Int:
						nb.AssignInt(m.NextValue.(int64))
					case datamodel.Kind_Link:
						lnk, err := sdk.ParseCidLink(m.NextValue.(string))
						if err != nil {
							return nil, err
						}
						nb.AssignLink(lnk)
					case datamodel.Kind_Bool:
						nb.AssignBool(m.NextValue.(bool))
					default:
						nb.AssignString(m.NextValue.(string))
					}

					return nb.Build(), nil

				} else if progress.Path.String() == m.Path && m.NextValue == "" {
					// next is absent, remove
				}
				return nil, fmt.Errorf("%s not found", m.Path)
			}, false)

		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// @BasePath /v0
// Update godoc
// @Summary Stores JSON as dag-json
// @Schemes
// @Description updates a dag-json block which syncs with IPFS. Returns a CID.
// @Tags dag-json
// @Accept json
// @Produce json
// @Success 200 {string} cid
// @Router /v1/dagjson [put]
func (dagctx *DagJsonHandler) Update(c *gin.Context) {

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

	// doc, err := dagctx.Store.DataStore.Get(context.Background(), from)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": fmt.Errorf("missing did").Error(),
	// 	})
	// 	return
	// }

	temp, _ := jsonparser.GetUnsafeString(v, "data")

	data, err := hexutil.Decode(temp)
	var buf bytes.Buffer
	if err != nil {
		err = json.Compact(&buf, []byte(temp))
		data = buf.Bytes()
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing payload data source").Error(),
		})
		return
	}

	resolution, err := types.ResolveDIDDoc(from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did").Error(),
		})
		return
	}

	// if DID is valid, assume sigature is ok?
	ok, err := types.IsValidSignature(resolution, data, signature)
	if err != nil || !ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did signature").Error(),
		})
		return
	}

	digest := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)))
	var items []map[string]interface{}
	json.Unmarshal(data, &items)

	mutations := make([]Mutation, len(items))
	json.Unmarshal(data, &mutations)
	nodecid, _ := jsonparser.GetString(v, "cid")

	if nodecid == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing cid").Error(),
		})
		return
	}

	topic, err := jsonparser.GetString(v, "topic")
	currentCid, err := sdk.ParseCidLink(nodecid)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid cid").Error(),
		})
		return
	}
	current, err := dagctx.Store.Load(ipld.LinkContext{}, currentCid)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing cid").Error(),
		})
		return
	}
	n, err := dagctx.ApplyFocusedTransform(current, mutations)

	cid := dagctx.Store.Store(ipld.LinkContext{}, n)

	lastHash, _ := dagctx.ProofHandler.proofs.GetCurrentVersion()

	dbl := &DagBlockResult{
		Issuer:        from,
		Timestamp:     time.Now().Unix(),
		ContentHash:   cid,
		Signature:     signature,
		Digest:        hexutil.Encode(digest),
		Network:       dagctx.Moniker,
		LastBlockHash: string(lastHash),
	}
	block := dagctx.ApplyDagBlock(dbl)
	key := dagctx.Store.LinkSystem.MustComputeLink(sdk.GetDagJSONLinkPrototype(), block)
	bz, err := block.AsBytes()
	err = dagctx.Store.DataStore.Put(c.Request.Context(), (key.Binary()), bz)
	// dagctx.PreviousBlock = res
	//	dagctx.Store.DataStore.Put(c.Request.Context(), fmt.Sprintf("block:%d", blockNumber), []byte(res.String()))
	contentTopic, err := protocol.StringToContentTopic(topic)

	// block
	dagctx.WakuPeer.Publish(contentTopic, block)

	// metadata
	dagctx.WakuPeer.Publish(contentTopic, n)

	c.JSON(201, gin.H{
		"cid": key.String(),
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
// @Router /v1/dagjson/{cid}/{path} [get]
func (dagctx *DagJsonHandler) DagJsonRead(c *gin.Context) {
	lnk, err := sdk.ParseCidLink(c.Param("cid"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err),
		})
		return
	}
	p := types.GetUserPath(dagctx.Moniker)

	path := c.Param("path")

	var n datamodel.Node
	if path != "" {
		var traversalPath ipld.Path
		if c.Query("namespace") != "" {
			traversalPath = ipld.ParsePath(c.Query("namespace"))
		} else {
			traversalPath = ipld.ParsePath(p)
		}
		prog := traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     dagctx.Store.LinkSystem,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
			// Path: traversalPath,
		}

		n, err = dagctx.Store.Load(ipld.LinkContext{
			LinkPath: traversalPath,
		}, lnk)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		if path != "/" {
			path = strings.TrimPrefix(path, "/")

			node, err := prog.Get(n, datamodel.ParsePath(path))

			if err != nil {
				child, _ := n.LookupByString(path)
				if child.Kind() == datamodel.Kind_Link {
					// l, _ := child.AsLink()
					err = prog.Focus(child, datamodel.ParsePath("/"), func(p traversal.Progress, n datamodel.Node) error {
						l, _ := n.AsLink()
						node, err = dagctx.Store.Load(ipld.LinkContext{
							LinkPath: traversalPath,
						}, l)

						return err
					})
				}
			}
			if err != nil {
				c.JSON(400, gin.H{
					"error": fmt.Errorf("%v", err),
				})
				return
			}

			trasEnc, err := sdk.Encode(node)
			c.JSON(200, json.RawMessage(trasEnc))
			if err != nil {
				c.JSON(400, gin.H{
					"error": fmt.Errorf("%s", err.Error()),
				})
				return
			}
			return
		}

	}

	data, err := sdk.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%s", err.Error()),
		})
		return
	}

	var contentData string

	datanode, err := n.LookupByString("contentHash")
	if err == nil && c.Query("compact") != "true" {
		lnkNode, _ := datanode.AsLink()
		// fetch
		contentHashNode, _ := dagctx.Store.Load(ipld.LinkContext{
			LinkPath: ipld.ParsePath("/"),
		}, lnkNode)
		contentData, err = sdk.Encode(contentHashNode)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%s", err.Error()),
			})
			return
		}
		response, err := jsonparser.Set([]byte(data), []byte(contentData), "content")
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%s", err.Error()),
			})
			return
		}
		c.JSON(200, json.RawMessage(response))
	} else {
		c.JSON(200, json.RawMessage(data))
	}
}

// @BasePath /v0
// UserDag godoc
// @Summary Reads JSON from a dag-json block
// @Schemes
// @Description Returns JSON
// @Tags dag-json
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/dag/topics/ [get]
func (dagctx *DagJsonHandler) UserDag(c *gin.Context) {
	user := c.Query("from")
	topic := c.Query("topic")
	concat := fmt.Sprintf("%s:%s", topic, user)

	concatLink, err := dagctx.Store.DataStore.Get(c.Request.Context(), concat)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err.Error()),
		})
	}

	lnk, err := sdk.ParseCidLink(string(concatLink))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err.Error()),
		})
		return
	}

	var n datamodel.Node
	p := types.GetUserPath(dagctx.Moniker)

	n, err = dagctx.Store.Load(ipld.LinkContext{LinkPath: ipld.ParsePath(p)}, lnk)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%v", err.Error()),
		})
		return
	}

	data, err := sdk.Encode(n)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("%s", err.Error()),
		})
		return
	}

	var contentData string

	datanode, err := n.LookupByString("contentHash")
	if err == nil && c.Query("compact") != "true" {
		lnkNode, _ := datanode.AsLink()
		// fetch
		contentHashNode, _ := dagctx.Store.Load(ipld.LinkContext{
			LinkPath: ipld.ParsePath("/"),
		}, lnkNode)
		contentData, err = sdk.Encode(contentHashNode)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%s", err.Error()),
			})
			return
		}
		response, err := jsonparser.Set([]byte(data), []byte(contentData), "content")
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%s", err.Error()),
			})
			return
		}
		c.JSON(200, json.RawMessage(response))
	} else {
		c.JSON(200, json.RawMessage(data))
	}
}
