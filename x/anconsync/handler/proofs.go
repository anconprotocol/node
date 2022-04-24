package handler

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/0xPolygon/polygon-sdk/helper/keccak"
	"github.com/buger/jsonparser"
	"github.com/cosmos/iavl"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-graphsync/ipldutil"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/joncrlsn/dque"
	"github.com/status-im/go-waku/waku/v2/node"
	"github.com/status-im/go-waku/waku/v2/protocol"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	dbm "github.com/tendermint/tm-db"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
)

type Commit struct {
	LastHash []byte
	Height   int64
}
type ProofHandler struct {
	*sdk.AnconSyncContext
	db                     dbm.DB
	WakuPeer               *WakuHandler
	LastCommit             *Commit
	api                    proofsignature.IavlProofAPI
	proofs                 proofsignature.IavlProofService
	RootKey                string
	ContentTopic           protocol.ContentTopic
	Moniker                string
	privateKey             *ecdsa.PrivateKey
	rwLock                 sync.RWMutex
	pendingTransactionPool *dque.DQue
}

func (h *ProofHandler) Commit() (int64, string, error) {

	// p := fmt.Sprintf("%s/%s/hash", "/anconprotocol", h.RootKey)

	// parent, err := h.proofs.Hash(&emptypb.Empty{})
	// parentHash, _ := jsonparser.GetString(parent, "lastHash")
	// h.proofs.Set([]byte(p), []byte(parentHash))

	v, err := h.proofs.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(v, "root_hash")
	version, err := jsonparser.GetInt(v, "version")
	lastHash := []byte(hash)
	blockNumber := cast.ToInt64(version)

	h.LastCommit = &Commit{Height: blockNumber, LastHash: lastHash}
	return cast.ToInt64(version), hash, err
}
func (h *ProofHandler) GetProofService() *proofsignature.IavlProofService {
	return &h.proofs

}

func (h *ProofHandler) GetProofAPI() *proofsignature.IavlProofAPI {
	return &h.api

}
func NewProofHandler(ctx *sdk.AnconSyncContext, wakuPeer *WakuHandler, moniker string, privateKeyPath string) *ProofHandler {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.GenerateOrReadPrivateKey(privateKeyPath)
	if err != nil {
		//		panic(err)
	}

	// os.OpenFile(,,privateKeyPath)

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}
	proofs, _ := proofsignature.NewIavlAPI(ctx.Store, ctx.Exchange, db, 2000, 0)
	contentTopic := protocol.NewContentTopic(moniker, 1, "proofs", "json")

	queueName := "memqueue"
	queueDir := folder
	segmentSize := 50

	// Create a new queue with segment size of 50
	pendingTransactionPool, err := dque.NewOrOpen(queueName, queueDir, segmentSize, QueueItemBuilder)

	if err != nil {
		panic(err)
	}
	return &ProofHandler{AnconSyncContext: ctx,
		WakuPeer: wakuPeer,
		db:       db,
		proofs:   *proofs.Service, api: *proofs,
		rwLock:                 sync.RWMutex{},
		privateKey:             privateKey,
		ContentTopic:           contentTopic,
		Moniker:                moniker,
		pendingTransactionPool: pendingTransactionPool,
	}

}

type PoolItem struct {
	Block DagBlockResult
	Cid   string
}

func QueueItemBuilder() interface{} {
	return &PoolItem{}
}

func (h *ProofHandler) AddToPool(item *PoolItem) error {

	// Add an item to the queue
	return h.pendingTransactionPool.Enqueue(item)
}

func (h *ProofHandler) Listen(ctx context.Context) {
	sub, err := h.WakuPeer.Subscribe(ctx, h.ContentTopic.String())

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	go func() {
		var i int
		for {

			if i == 5 {
				// Commit tx batch
				_, err = h.proofs.SaveVersion(&emptypb.Empty{})
				i = 0
			}
			i++
			time.Sleep(5000)
			has, err := h.pendingTransactionPool.Peek()
			if has == nil {
				return
			}
			iface, err := h.pendingTransactionPool.Dequeue()
			if err != nil {
				return
			}
			// Assert type of the response to an Item pointer so we can work with it
			item, ok := iface.(*PoolItem)
			if !ok {
				log.Fatal("Dequeued object is not an Item pointer")
			}

			// Load dag block
			lnk, _ := sdk.ParseCidLink(item.Cid)
			keypath := protocol.NewContentTopic(h.Moniker, 1, "block", lnk.String())
			k := []byte(keypath.String())

			// commit block
			block, err := h.Store.Load(ipld.LinkContext{}, lnk)
			bz, _ := block.AsBytes()
			commit, _ := h.proofs.Set(k, bz)
			hash, err := jsonparser.GetString(commit, "updated")

			// get latest
			_, latestBlockNumber := h.proofs.GetCurrentVersion()

			proofblock := h.Apply(&DagBlockResult{
				Issuer:        item.Block.Issuer,
				Timestamp:     item.Block.Timestamp,
				ContentHash:   item.Block.ContentHash,
				Signature:     item.Block.Signature,
				Digest:        item.Block.Digest,
				Network:       item.Block.Network,
				CommitHash:    string(hash),
				Height:        int64(latestBlockNumber),
				Key:           base64.StdEncoding.EncodeToString([]byte(k)),
				LastBlockHash: item.Block.LastBlockHash,
				// ParentHash:    currentCid.String(),
			})
			res := h.Store.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(types.GetUserPath(h.Moniker))}, proofblock)

			n, err := h.Store.Load(ipld.LinkContext{}, res)
			h.WakuPeer.Publish(h.ContentTopic, n)

		}

		// Notify

	}()

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
					block, _ := h.Store.Load(ipld.LinkContext{}, lnk)
					h.WakuPeer.Publish(h.ContentTopic, block)
				}
			}
		}

	}
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

func (dagctx *ProofHandler) Apply(args *DagBlockResult) datamodel.Node {

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

func (h *ProofHandler) VerifyGenesis(moniker string, key string) ([]byte, error) {

	version := 0
	tree, err := iavl.NewMutableTree(h.db, int(2000))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return nil, errors.Wrapf(err, "unable to load version %d", version)
	}
	key = fmt.Sprintf("%s%s", moniker, key)

	_, v, err := tree.GetWithProof([]byte(key))
	if err != nil && v != nil {
		return nil, errors.Wrap(err, "Unable to get with proof")
	}

	bz := tree.Hash()
	err = v.Verify(bz)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get rawkey")
	}

	return bz, nil
}

func InitGenesis(hostName string, moniker string, cidlink datamodel.Link, priv *ecdsa.PrivateKey) (string, string, error) {
	userHomeDir, err := os.UserHomeDir()
	version := 0
	if err != nil {
		panic(err)
	}

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}

	tree, err := iavl.NewMutableTree(db, int(2000))
	if err != nil {
		return " ", " ", errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return " ", " ", errors.Wrapf(err, "unable to load version %d", version)
	}

	// cidlink := sdk.CreateCidLink(signed[0:32])

	key := fmt.Sprintf("%s%s", moniker, cidlink.String())
	value := fmt.Sprintf(
		`{
		data: "%s",
		signature: "%s",
		}`,
		hostName, cidlink.String(),
	)

	tree.Set([]byte(key), []byte(value))

	_, _, err = tree.SaveVersion()

	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to commit")
	}

	_, proof, err := tree.GetWithProof([]byte(key))
	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to get with proof")
	}
	var proofData []byte
	proofData, err = proof.ToProto().Marshal()
	// err = proto.Unmarshal(proofData, )
	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to marshal")
	}

	hash := tree.Hash()
	rawKey, err := crypto.MarshalPrivateKey(priv)

	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to get rawkey")
	}
	stringRawKey := hexutil.Encode(rawKey)

	message := fmt.Sprintf(
		`Ancon protocol node initialize with: 
		*Sep256k1 private key: %s
		*Genesis value: %s
		*Genesis key: %s
		*Proof: %s
		*Last header hash: %s
		`, stringRawKey, hex.EncodeToString([]byte(value)), cidlink.String(), hex.EncodeToString(proofData), hex.EncodeToString(hash),
	)

	return message, key, nil
}

func GenerateKeys() (string, error) {
	// Set your own keypair
	priv, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	rawKey, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return " ", errors.Wrap(err, "Unable to get private key")
	}
	pub := crypto.MarshalPublicKey(&priv.PublicKey)
	if err != nil {
		return " ", errors.Wrap(err, "Unable to get public key")
	}
	stringRawKey := hexutil.Encode(rawKey)
	publicKeyBase58 := base58.Encode(pub)
	pubhex := hexutil.Encode(pub)
	message := fmt.Sprintf(
		`Generated keys
		Sep256k1 private key (hex): %s
		Sep256k1 public key (hex): %s
		Sep256k1 public key (base58): %s`,
		stringRawKey, pubhex, publicKeyBase58,
	)

	return message, nil
}

// @BasePath /v0
// Verify godoc
// @Summary Verifies an ics23 proofs
// @Schemes
// @Description Verifies an ics23 proof
// @Tags proofs
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v1/proofs/verify [post]
func (dagctx *ProofHandler) ReadCurrentRootHash(c *gin.Context) {

	lastHash, err := dagctx.proofs.Hash(nil)
	sig := c.Query("sig")

	if sig == "true" {
		var digest []byte
		// priv, err := crypto.GenerateKey()
		keccak.Keccak256(digest, []byte(lastHash))
		signed, err := dagctx.PrivateKey.Sign(rand.Reader, digest, nil) //priv.Sign(rand.Reader, digest, nil)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Sig query Error %v", err).Error(),
			})
			return
		}
		c.JSON(201, gin.H{
			"lastHash":  lastHash,
			"signature": signed,
		})
		return
	}

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("lastHash Error %v", err).Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"lastHash": lastHash,
	})
}

// @BasePath /v0
// Create godoc
// @Summary Create
// @Schemes
// @Description Writes an ics23 proof
// @Tags proofs
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v1/proofs [post]
func (dagctx *ProofHandler) Create(c *gin.Context) {

	v, _ := c.GetRawData()

	key, _ := jsonparser.GetString(v, "key")

	if key == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing key").Error(),
		})
		return
	}

	value, _ := jsonparser.GetString(v, "value")

	if value == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing value").Error(),
		})
		return
	}

	message, err := dagctx.proofs.Set([]byte(key), []byte(value))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	c.JSON(201, message)
	//	impl.PushBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, cid)
}

// @BasePath /v0
// Read godoc
// @Summary Reads an existing proof
// @Schemes
// @Description Returns JSON
// @Tags proofs
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/proofs/get/{path} [get]
func (dagctx *ProofHandler) Read(c *gin.Context) {

	key, _ := c.Params.Get("key")

	if key == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing key").Error(),
		})
		return
	}
	height, _ := c.GetQuery("height")

	if height == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing height").Error(),
		})
		return
	}
	version := cast.ToInt64(height)
	internalKey, _ := base64.StdEncoding.DecodeString(key)
	data, err := dagctx.proofs.GetCommitmentProof([]byte(internalKey), version)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	exportAs, _ := c.GetQuery("export")
	if exportAs == "qr" {
		qrc, err := qrcode.New(string(data))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("decode Error %v", err).Error(),
			})
			return
		}

		bg := c.Query("bgcolor")
		if bg == "" {
			bg = "#ffffff"
		} else {
			bg = "#" + bg
		}
		fg := c.Query("fgcolor")
		if fg == "" {
			fg = "#000000"
		} else {
			fg = "#" + fg
		}
		buf := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		wr := gzip.NewWriter(buf)

		w := standard.NewWithWriter(wr,
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			standard.WithBgColorRGBHex(bg),
			standard.WithFgColorRGBHex(fg),
		)
		qrc.Save(w)
		w.Close()
		rdr, err := gzip.NewReader(buf)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error %v", err).Error(),
			})
			return
		}

		data, err := io.ReadAll(rdr)
		buf2.Write(data)
		defer rdr.Close()

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error %v", err).Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"qr": base64.StdEncoding.EncodeToString(buf2.Bytes()),
		})
	} else {
		c.JSON(200, data)
	}
}
