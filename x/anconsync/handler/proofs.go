package handler

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/bigqueue"
	"github.com/cosmos/iavl"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ipld/go-ipld-prime/datamodel"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
	"github.com/status-im/go-waku/waku/v2/protocol"

	"github.com/anconprotocol/sdk"
	dbm "github.com/tendermint/tm-db"
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
	RootKey                string
	ContentTopic           protocol.ContentTopic
	Moniker                string
	privateKey             *ecdsa.PrivateKey
	pendingTransactionPool *bigqueue.FileQueue
}

func NewProofHandler(ctx *sdk.AnconSyncContext, wakuPeer *WakuHandler, moniker string, privateKeyPath string) *ProofHandler {
	return &ProofHandler{AnconSyncContext: ctx,
		WakuPeer: wakuPeer,
	}

}

type PoolItem struct {
	Block DagBlockResult
	Cid   string
}

func QueueItemBuilder() interface{} {
	return &PoolItem{}
}

func (h *ProofHandler) AddToPool(item []byte) (int64, error) {

	// Add an item to the queue
	return h.pendingTransactionPool.Enqueue(item)
}

func (h *ProofHandler) HandleIncomingProofRequests() {
	go func() {
		for _ = range time.Tick(time.Second * 12) {
			// Commit tx batch
			h.Store.Commit()

			// get latest
			fmt.Println("new block: ", h.Store.GetTreeVersion())
		}
	}()

	// Notify

}

type DagBlockResult struct {
	Path          string         `json:"path"`
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

func (dagctx *ProofHandler) Apply(n datamodel.Node, height int64, hash string, key string) datamodel.Node {
	prog := traversal.Progress{
		Cfg: &traversal.Config{
			LinkSystem:                     dagctx.Store.LinkSystem,
			LinkTargetNodePrototypeChooser: basicnode.Chooser,
		},
	}
	current, _ := prog.FocusedTransform(
		n,
		datamodel.ParsePath("height"),
		func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			nb := basicnode.Prototype.Any.NewBuilder()
			nb.AssignInt(int64(height))
			return nb.Build(), nil
		}, false)

	block, _ := prog.FocusedTransform(
		current,
		datamodel.ParsePath("commitHash"),
		func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			nb := basicnode.Prototype.Any.NewBuilder()
			nb.AssignString(hash)
			return nb.Build(), nil
		}, false)

	dagblock, _ := prog.FocusedTransform(
		block,
		datamodel.ParsePath("key"),
		func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			nb := basicnode.Prototype.Any.NewBuilder()
			nb.AssignString(key)
			return nb.Build(), nil
		}, false)
	return dagblock
}

func (h *ProofHandler) VerifyGenesis(moniker string, key string) ([]byte, error) {

	version := 0
	s, err := iavl.NewMutableTree(h.db, int(2000))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iavl s")
	}
	if _, err = s.LoadVersion(int64(version)); err != nil {
		return nil, errors.Wrapf(err, "unable to load version %d", version)
	}
	key = fmt.Sprintf("%s%s", moniker, key)

	_, v, err := s.GetWithProof([]byte(key))
	if err != nil && v != nil {
		return nil, errors.Wrap(err, "Unable to get with proof")
	}

	bz := s.Hash()
	err = v.Verify(bz)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get rawkey")
	}

	return bz, nil
}

func InitGenesis(s *sdk.Storage, hostName string, moniker string, cidlink datamodel.Link, priv *ecdsa.PrivateKey) (string, string, error) {

	value := fmt.Sprintf(
		`{
		data: "%s",
		signature: "%s",
		}`,
		hostName, cidlink.String(),
	)

	s.Put([]byte(moniker), (cidlink.String()), []byte(value))

	s.Commit()

	proofData, err := s.GetWithProof([]byte(cidlink.String()))
	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to get with proof")
	}

	hash := s.GetTreeHash()
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

	return message, cidlink.String(), nil
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
