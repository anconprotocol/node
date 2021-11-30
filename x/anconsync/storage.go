package anconsync

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/fsstore"
	"github.com/multiformats/go-multihash"
)

const (
	LINK_PROTO_VERSION = 1
)

type Storage struct {
	DataStore  fsstore.Store
	LinkSystem linking.LinkSystem
}

func NewStorage(folder string) Storage {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome := filepath.Join(userHomeDir, folder)
	store := fsstore.Store{}
	store.InitDefaults(DefaultNodeHome)
	lsys := cidlink.DefaultLinkSystem()
	//   you just need a function that conforms to the ipld.BlockWriteOpener interface.
	lsys.StorageWriteOpener = func(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		// change prefix
		buf := bytes.Buffer{}
		return &buf, func(lnk ipld.Link) error {
			key := strings.Join([]string{lnk.String(), lnkCtx.LinkPath.String()}, "/")
			wr, cb, err := store.PutStream(lnkCtx.Ctx)
			if err != nil {
				return fmt.Errorf("error while reading stream")
			}
			wr.Write(buf.Bytes())
			cb(key)
			return err
		}, nil
	}
	lsys.StorageReadOpener = func(lnkCtx ipld.LinkContext, link ipld.Link) (io.Reader, error) {
		key := strings.Join([]string{link.String(), lnkCtx.LinkPath.String()}, "/")
		reader, err := store.GetStream(lnkCtx.Ctx, key)
		if err != nil {
			return nil, fmt.Errorf("path not found")
		}
		return reader, nil
	}

	lsys.TrustedStorage = true

	return Storage{
		DataStore:  store,
		LinkSystem: lsys,
	}
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
