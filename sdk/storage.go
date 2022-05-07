package sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	ibc "github.com/cosmos/ibc-go/v2/modules/core/23-commitment/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/cache"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multihash"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	LINK_PROTO_VERSION = 1
)

type Storage struct {
	dataStore  types.CommitMultiStore
	LinkSystem linking.LinkSystem
	RootHash   cidlink.Link
}

var STORE_KEY = "anconprotocol"
var STORE_KEY_TYPE = types.NewKVStoreKey(STORE_KEY)

func NewStorage(key string, store types.CommitMultiStore, db dbm.DB) Storage {
	store.MountStoreWithDB(STORE_KEY_TYPE, types.StoreTypeIAVL, db)

	lsys := cidlink.DefaultLinkSystem()
	s := Storage{
		dataStore:  store,
		LinkSystem: lsys,
	}
	//   you just need a function that conforms to the ipld.BlockWriteOpener interface.
	lsys.StorageWriteOpener = func(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {

		// change prefix
		buf := bytes.Buffer{}
		return &buf, func(lnk ipld.Link) error {
			kvstore := store.GetCommitStore(STORE_KEY_TYPE)
			path := []byte(lnkCtx.LinkPath.String())

			kvs := prefix.NewStore(kvstore.(types.CommitKVStore), path)

			kvs.Set([]byte(lnk.String()), buf.Bytes())

			return nil

		}, nil
	}
	lsys.StorageReadOpener = func(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {

		path := []byte(lnkCtx.LinkPath.String())
		kvstore := store.GetCommitStore(STORE_KEY_TYPE)
		kvs := prefix.NewStore(kvstore.(types.CommitKVStore), path)
		value := kvs.Get([]byte(lnk.String()))
		return bytes.NewReader(value), nil
	}

	lsys.TrustedStorage = true
	s.LinkSystem = lsys
	return s
}

func (s *Storage) Get(path []byte, id string) ([]byte, error) {

	store := s.dataStore.GetCommitStore(STORE_KEY_TYPE)
	kvs := prefix.NewStore(store.(types.KVStore), path)
	value := kvs.Get([]byte(id))
	return value, nil
}

func (s *Storage) Remove(path []byte, id string) error {
	store := s.dataStore.GetCommitStore(STORE_KEY_TYPE)
	kvs := prefix.NewStore(store.(types.KVStore), path)
	kvs.Delete([]byte(id))
	return nil
}

func (s *Storage) Put(path []byte, id string, data []byte) (err error) {
	store := s.dataStore.GetCommitStore(STORE_KEY_TYPE)
	kvs := prefix.NewStore(store.(types.KVStore), path)

	kvs.Set([]byte(id), data)

	return nil
}

func (s *Storage) Iterate(path []byte, start, end []byte) (dbm.Iterator, error) {
	store := s.dataStore.GetCommitStore(STORE_KEY_TYPE)
	kvs := prefix.NewStore(store.(types.KVStore), path)

	return kvs.Iterator(start, end), nil
}

func (s *Storage) Has(path []byte, id []byte) (bool, error) {
	store := s.dataStore.GetCommitStore(STORE_KEY_TYPE)
	kvs := prefix.NewStore(store.(types.KVStore), path)

	return kvs.Has(id), nil
}

func (s *Storage) GetTreeHash() []byte {
	return s.dataStore.LastCommitID().Hash
}

func (s *Storage) GetTreeVersion() int64 {
	return s.dataStore.LastCommitID().Version
}

/*
CreateMembershipProof will produce a CommitmentProof that the given key (and queries value) exists in the iavl tree.
If the key doesn't exist in the tree, this will return an error.
*/
// func createMembershipProof(tree *iavl.MutableTree, key []byte, exist *ics23.ExistenceProof) (*ics23.CommitmentProof, error) {
// 	// exist, err := createExistenceProof(tree, key)
// 	proof := &ics23.CommitmentProof{
// 		Proof: &ics23.CommitmentProof_Exist{
// 			Exist: exist,
// 		},
// 	}
// 	return proof, nil
// 	// return ics23.CombineProofs([]*ics23.CommitmentProof{proof})
// }

// GetWithProof returns a result containing the IAVL tree version and value for
// a given key based on the current state (version) of the tree including a
// verifiable Merkle proof.
func (s *Storage) GetWithProof(key []byte, height int64) (json.RawMessage, error) {

	// create cache manager to unwrap
	mngr := cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
	mngr.GetStoreCache(STORE_KEY_TYPE, s.dataStore.GetCommitKVStore(STORE_KEY_TYPE))
	iavlstore := mngr.Unwrap(STORE_KEY_TYPE).(*iavl.Store)

	queryableStore := store.Queryable(iavlstore)

	result := make(map[string]interface{})
	res := queryableStore.Query(abci.RequestQuery{
		Data:   []byte(key),
		Path:   ("/key"),
		Height: height,
		Prove:  true,
	})
	mp, err := ibc.ConvertProofs(res.ProofOps)
	if err != nil {
		return nil, err
	}
	result["proof"] = mp
	result["value"] = res.Value

	hexres, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return hexres, nil
}

// GetCommitmentProof returns a result containing the IAVL tree version and value for
// a given key based on the current state (version) of the tree including a
// verifiable existing or not existing Commitment proof.
func (s *Storage) GetCommitmentProof(key []byte, version int64) (json.RawMessage, error) {

	// create cache manager to unwrap
	mngr := cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
	mngr.GetStoreCache(STORE_KEY_TYPE, s.dataStore.GetCommitKVStore(STORE_KEY_TYPE))
	iavlstore := mngr.Unwrap(STORE_KEY_TYPE).(*iavl.Store)

	queryableStore := store.Queryable(iavlstore)

	res := queryableStore.Query(abci.RequestQuery{
		Data:   []byte(key),
		Path:   ("/key"),
		Height: version,
		Prove:  true,
	})
	mp, err := ibc.ConvertProofs(res.ProofOps)
	if err != nil {
		return nil, err
	}

	hexres, err := json.Marshal(mp.Proofs)
	if err != nil {
		return nil, err
	}

	return hexres, nil
}

// eth-block	ipld	0x90	permanent	Ethereum Header (RLP)
// eth-block-list	ipld	0x91	permanent	Ethereum Header List (RLP)
// eth-tx-trie	ipld	0x92	permanent	Ethereum Transaction Trie (Eth-Trie)
// eth-tx	ipld	0x93	permanent	Ethereum Transaction (MarshalBinary)
// eth-tx-receipt-trie	ipld	0x94	permanent	Ethereum Transaction Receipt Trie (Eth-Trie)
// eth-tx-receipt	ipld	0x95	permanent	Ethereum Transaction Receipt (MarshalBinary)
// eth-state-trie	ipld	0x96	permanent	Ethereum State Trie (Eth-Secure-Trie)
// eth-account-snapshot	ipld	0x97	permanent	Eth	ereum Account Snapshot (RLP)
// eth-storage-trie	ipld	0x98	permanent	Ethereum Contract Storage Trie (Eth-Secure-Trie)
// eth-receipt-log-trie	ipld	0x99	draft	Ethereum Transaction Receipt Log Trie (Eth-Trie)
// eth-reciept-log	ipld	0x9a	draft	Ethereum Transaction Receipt Log (RLP)
var (
	DagEthCodecs map[string]uint64 = make(map[string]uint64)
)

func init() {
	DagEthCodecs["eth-block"] = 0x90
}

func GetDagEthereumLinkPrototype(codec string) ipld.LinkPrototype {
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  LINK_PROTO_VERSION,
		Codec:    DagEthCodecs[codec],
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: 32,                 // sha2-256 hash has a 32-byte sum.
	}}
}

func GetDagCBORLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  LINK_PROTO_VERSION,
		Codec:    cid.DagCBOR,        // dag-cbor
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: 32,                 // sha2-256 hash has a 32-byte sum.
	}}
}

func GetDagJSONLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  LINK_PROTO_VERSION,
		Codec:    0x0129,             // dag-json
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: 32,                 // sha2-256 hash has a 32-byte sum.
	}}
}

func GetDagJOSELinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  LINK_PROTO_VERSION,
		Codec:    cid.DagJOSE,        // dag-json
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: 32,                 // sha2-256 hash has a 32-byte sum.
	}}
}

func GetRawLinkPrototype() ipld.LinkPrototype {
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  LINK_PROTO_VERSION,
		Codec:    0x55,               // dag-json
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: 32,                 // sha2-256 hash has a 32-byte sum.
	}}
}

// Store node as  dag-json
func (k *Storage) Store(linkCtx ipld.LinkContext, node datamodel.Node) datamodel.Link {
	return k.LinkSystem.MustStore(linkCtx, GetDagJSONLinkPrototype(), node)
}

// Load node from  dag-json
func (k *Storage) Load(linkCtx ipld.LinkContext, link datamodel.Link) (datamodel.Node, error) {
	np := basicnode.Prototype.Any
	node, err := k.LinkSystem.Load(linkCtx, link, np)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// Store node as  dag-cbor
func (k *Storage) StoreDagCBOR(linkCtx ipld.LinkContext, node datamodel.Node) datamodel.Link {
	return k.LinkSystem.MustStore(linkCtx, GetDagCBORLinkPrototype(), node)
}

// Store node as  raw
func (k *Storage) StoreRaw(linkCtx ipld.LinkContext, node datamodel.Node) datamodel.Link {
	return k.LinkSystem.MustStore(linkCtx, GetRawLinkPrototype(), node)
}

// Store node as  dag-eth
func (k *Storage) StoreDagEth(linkCtx ipld.LinkContext, node datamodel.Node, codecFormat string) datamodel.Link {
	return k.LinkSystem.MustStore(linkCtx, GetDagEthereumLinkPrototype(codecFormat), node)
}

func Encode(n datamodel.Node) (string, error) {
	var sb strings.Builder
	err := dagjson.Encode(n, &sb)
	return sb.String(), err
}

func Decode(proto datamodel.NodePrototype, src string) (datamodel.Node, error) {
	nb := proto.NewBuilder()
	err := dagjson.Decode(nb, strings.NewReader(src))
	return nb.Build(), err
}

func EncodeCBOR(n datamodel.Node) ([]byte, error) {
	var sb bytes.Buffer
	err := dagcbor.Encode(n, &sb)
	return sb.Bytes(), err
}

func DecodeCBOR(proto datamodel.NodePrototype, src []byte) (datamodel.Node, error) {
	nb := proto.NewBuilder()
	err := dagcbor.Decode(nb, bytes.NewReader(src))
	return nb.Build(), err
}
