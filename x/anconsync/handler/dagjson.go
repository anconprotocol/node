package handler

import (
	"bytes"
	"log"
	"strings"
	"time"

	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/buger/jsonparser"
	ecies "github.com/ecies/go"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/spf13/cast"

	"google.golang.org/protobuf/types/known/emptypb"
)

type DagJsonHandler struct {
	*sdk.AnconSyncContext
	Proof         *proofsignature.IavlProofService
	IPFSHost      string
	RootKey       string
	Moniker       string
	PreviousBlock datamodel.Link
}
type Mutation struct {
	Path          string
	PreviousValue string
	NextValue     interface{}
	NextValueKind datamodel.Kind
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

	p := fmt.Sprintf("%s/%s", types.GetUserPath(dagctx.Moniker), from)
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
	ok, err := types.Authenticate(doc, data, signature)
	if !ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid signature").Error(),
		})
		return
	}

	encrypt, _ := jsonparser.GetString(v, "encrypt")
	hasEncrypt := encrypt == "true"
	var authorizedRecipients []string

	if hasEncrypt {
		recipients, _ := jsonparser.GetString(v, "authorizedRecipients")
		authorizedRecipients = strings.Split(recipients, ",")

		content, err := dagctx.Store.DataStore.Get(c.Request.Context(), authorizedRecipients[0])
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("authorized recipient not found").Error(),
			})
			return
		}

		pub, err := types.GetDidDocumentAuthentication((content))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing recipient public key").Error(),
			})
			return
		}

		data, err = ecies.Encrypt((*ecies.PublicKey)(pub), data)
		fmt.Println(data)

		if err != nil {
			log.Printf("failed to encrypt payload: %s", err)
			return
		}

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error while loading recipient public keys").Error(),
			})
			return
		}

	}

	parentHash, _ := jsonparser.GetString(v, "parent")
	path, _ := jsonparser.GetString(v, "path")

	if path == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing path").Error(),
		})
		return
	}

	digest := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)))
	var n datamodel.Node
	if isJSON && !hasEncrypt {
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
	internalKey := fmt.Sprintf("%s/%s", p, cid)
	dagctx.Proof.Set([]byte(internalKey), data)
	commit, err := dagctx.Proof.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(commit, "root_hash")
	version, err := jsonparser.GetInt(commit, "version")
	lastHash := []byte(hash)
	blockNumber := cast.ToInt64(version)
	addrrec, err := jsonparser.GetString((doc), "verificationMethod", "[0]", "ethereumAddress")
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did %v", err).Error(),
		})
		return
	}

	l := ""
	if dagctx.PreviousBlock != nil {
		l = dagctx.PreviousBlock.String()
	}
	block := dagctx.Apply(&DagBlockResult{
		Issuer:        addrrec,
		Timestamp:     time.Now().Unix(),
		Content:       n,
		ContentHash:   cid,
		CommitHash:    string(lastHash),
		Height:        blockNumber,
		Signature:     signature,
		Digest:        hexutil.Encode(digest),
		Network:       dagctx.Moniker,
		Key:           base64.StdEncoding.EncodeToString([]byte(internalKey)),
		RootKey:       base64.StdEncoding.EncodeToString([]byte(p)),
		LastBlockHash: l,
		ParentHash:    parentHash,
	})
	res := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(types.GetUserPath(dagctx.Moniker))}, block)
	dagctx.PreviousBlock = res
	topic, err := jsonparser.GetString(v, "topic")

	if topic != "" {
		topic = topic + ":" + addrrec
		dagctx.Store.DataStore.Put(c.Request.Context(), topic, []byte(res.String()))
	}

	dagctx.Store.DataStore.Put(c.Request.Context(), fmt.Sprintf("block:%d", blockNumber), []byte(res.String()))
	resp, err := sdk.Encode(block)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("bad encoding %s", err.Error()).Error(),
		})
		return
	}
	tx, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp))
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("bad encoding %s", err.Error()).Error(),
		})
		return
	}
	resp2, _ := sdk.Encode(n)
	m, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp2))

	c.JSON(201, gin.H{
		"cid": res.String(),
		"ipfs": map[string]interface{}{
			"metadata": m,
			"tx":       tx,
		},
	})
}

type DagBlockResult struct {
	Issuer        string         `json:"issuer"`
	Timestamp     int64          `json:"timestamp"`
	Content       datamodel.Node `json:"content"`
	ContentHash   datamodel.Link `json:"content_hash"`
	CommitHash    string         `json:"commit_hash"`
	Height        int64
	Signature     string `json:"signature"`
	Digest        string `json:"digest"`
	Network       string `json:"network"`
	Key           string `json:"key"`
	RootKey       string
	LastBlockHash string
	ParentHash    string `json:"parent_hash"`
}

func (dagctx *DagJsonHandler) Apply(args *DagBlockResult) datamodel.Node {

	block := fluent.MustBuildMap(basicnode.Prototype.Map, 13, func(na fluent.MapAssembler) {
		na.AssembleEntry("issuer").AssignString(args.Issuer)
		na.AssembleEntry("timestamp").AssignInt(args.Timestamp)
		na.AssembleEntry("contentHash").AssignLink(args.ContentHash)
		//	na.AssembleEntry("content").AssignNode(args.Content)
		na.AssembleEntry("commitHash").AssignString(args.CommitHash)
		na.AssembleEntry("height").AssignInt(args.Height)
		na.AssembleEntry("signature").AssignString(args.Signature)
		na.AssembleEntry("digest").AssignString(args.Digest)
		na.AssembleEntry("network").AssignString(args.Network)
		na.AssembleEntry("key").AssignString(args.Key)
		na.AssembleEntry("rootKey").AssignString(args.RootKey)
		if args.LastBlockHash != "" {
			lnk, _ := sdk.ParseCidLink(args.LastBlockHash)
			na.AssembleEntry("lastBlockHash").AssignLink(lnk)
		}
		if args.ParentHash != "" {
			lnk, _ := sdk.ParseCidLink(args.ParentHash)
			na.AssembleEntry("parentHash").AssignLink(lnk)
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
// @Router /v0/dagjson [put]
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

	doc, err := dagctx.Store.DataStore.Get(context.Background(), from)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing did").Error(),
		})
		return
	}

	p := fmt.Sprintf("%s/%s", types.GetUserPath(dagctx.Moniker), from)

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
	content, _ := sdk.Encode(n) // signature must match data + dif

	ok, err := types.Authenticate(doc, []byte(content), signature)
	if !ok {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid signature").Error(),
		})
		return
	}

	cid := dagctx.Store.Store(ipld.LinkContext{}, n)
	internalKey := fmt.Sprintf("%s/%s", p, cid)
	dagctx.Proof.Set([]byte(internalKey), data)
	commit, err := dagctx.Proof.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(commit, "root_hash")
	version, err := jsonparser.GetInt(commit, "version")
	lastHash := []byte(hash)
	blockNumber := cast.ToInt64(version)
	addrrec, err := jsonparser.GetString((doc), "verificationMethod", "[0]", "ethereumAddress")
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("invalid did %v", err).Error(),
		})
		return
	}

	l := ""
	if dagctx.PreviousBlock != nil {
		l = dagctx.PreviousBlock.String()
	}

	block := dagctx.Apply(&DagBlockResult{
		Issuer:        addrrec,
		Timestamp:     time.Now().Unix(),
		Content:       n,
		ContentHash:   cid,
		CommitHash:    string(lastHash),
		Height:        blockNumber,
		Signature:     signature,
		Digest:        hexutil.Encode(digest),
		Network:       dagctx.Moniker,
		Key:           base64.StdEncoding.EncodeToString([]byte(internalKey)),
		RootKey:       base64.StdEncoding.EncodeToString([]byte(p)),
		LastBlockHash: l,
		ParentHash:    currentCid.String()				,
	})
	res := dagctx.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(types.GetUserPath(dagctx.Moniker))}, block)

	if topic != "" {
		topic = topic + ":" + addrrec

		dagctx.Store.DataStore.Put(c.Request.Context(), topic, []byte(res.String()))
	}

	dagctx.PreviousBlock = res
	dagctx.Store.DataStore.Put(c.Request.Context(), fmt.Sprintf("block:%d", blockNumber), []byte(res.String()))
	resp, _ := sdk.Encode(block)
	tx, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp))

	resp2, _ := sdk.Encode(n)
	m, err := impl.PushBlock(c.Request.Context(), dagctx.IPFSHost, []byte(resp2))

	c.JSON(200, gin.H{
		"cid": res.String(),
		"ipfs": map[string]interface{}{
			"metadata": m,
			"tx":       tx,
		},
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
// @Router /v0/dagjson/{cid}/{path} [get]
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
// @Router /v0/dag/topics/ [get]
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
