package impl

import (
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-graphsync/ipldutil"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/0xPolygon/polygon-sdk/state"
	"github.com/0xPolygon/polygon-sdk/types"
	"github.com/anconprotocol/node/x/anconsync"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/spf13/cast"
)

func AddOnchainMetadataEvent() abi.Event {

	stringType, _ := abi.NewType("string", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)
	return abi.NewEvent(
		"AddOnchainMetadata",
		"AddOnchainMetadata",
		false,
		abi.Arguments{abi.Argument{
			Name:    "name",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "description",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "image",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "owner",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "parent",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "sources",
			Type:    bytesType,
			Indexed: false,
		}},
	)
}

func EncodeDagCborEvent() abi.Event {

	str, _ := abi.NewType("string", "", nil)
	return abi.NewEvent(
		"EncodeDagCbor",
		"EncodeDagCbor",
		false,
		abi.Arguments{abi.Argument{
			Name:    "path",
			Type:    str,
			Indexed: false,
		}, abi.Argument{
			Name:    "hexdata",
			Type:    str,
			Indexed: false,
		}},
	)
}

func StoreDagBlockDoneEvent() abi.Event {

	str, _ := abi.NewType("string", "", nil)
	return abi.NewEvent(
		"StoreDagBlockDone",
		"StoreDagBlockDone",
		false,
		abi.Arguments{abi.Argument{
			Name:    "path",
			Type:    str,
			Indexed: false,
		}, abi.Argument{
			Name:    "cid",
			Type:    str,
			Indexed: true,
		}},
	)
}

func encodeDagCborBlock(s anconsync.Storage, inputs abi.Arguments, data []byte, txHash types.Hash, blockHash types.Hash, chainID int64) (datamodel.Node, datamodel.Link, error) {

	props, err := inputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}

	///	path := props[0].(string)
	values := props[1].(string)
	bz := common.Hex2Bytes(values)

	n, _ := ipldutil.DecodeNode(bz)
	transactionBlock := fluent.MustBuildMap(basicnode.Prototype.Any, 3, func(ma fluent.MapAssembler) {
		// ma.AssembleEntry("dagblock").AssignLink(nodelink)
		ma.AssembleEntry("transactionHash").AssignBytes(txHash[:])
		ma.AssembleEntry("blockHash").AssignBytes(blockHash.Bytes())
		ma.AssembleEntry("chainId").AssignInt(chainID)
	})

	p := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   0x12, // sha2-256
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}

	lnk, err := s.LinkSystem.Store(ipld.LinkContext{
		LinkNode: transactionBlock,
	}, p, n)

	// lnk := p.BuildLink(data)

	if err != nil {
		return nil, nil, err
	}

	return n, lnk, nil
}

func encodeAnconMetadata(s anconsync.Storage, inputs abi.Arguments, data []byte, txHash types.Hash, blockHash types.Hash, chainID int64) (datamodel.Node, datamodel.Link, error) {

	props, err := inputs.Unpack(data)

	bz := props[5].([]byte)
	// bz := common.Hex2Bytes(values)

	var sources map[string]string

	js := hexutil.Bytes{}
	js.UnmarshalJSON(bz)

	err = json.Unmarshal(js, &sources)

	if err != nil {
		return nil, nil, err
	}

	n := fluent.MustBuildMap(basicnode.Prototype.Map, 10, func(na fluent.MapAssembler) {
		// TODO:
		na.AssembleEntry("name").AssignString(props[0].(string))
		na.AssembleEntry("description").AssignString(props[1].(string))
		na.AssembleEntry("image").AssignString(props[2].(string))

		na.AssembleEntry("owner").AssignString(props[3].(string))

		if props[4] != nil {
			na.AssembleEntry("parent").AssignString(props[4].(string))
		} else {
			na.AssembleEntry("parent").AssignNull()
		}

		// Sources
		if len(sources) > 0 {

			na.AssembleEntry("sources").CreateList(cast.ToInt64(len(sources)), func(la fluent.ListAssembler) {
				for _, v := range sources {
					lnk, err := anconsync.ParseCidLink((v))
					if err != nil {
						continue
					}
					la.AssembleValue().AssignLink(lnk)

				}
			})

		} else {
			na.AssembleEntry("sources").AssignNull()
		}

	})

	transactionBlock := fluent.MustBuildMap(basicnode.Prototype.Any, 3, func(ma fluent.MapAssembler) {
		ma.AssembleEntry("transactionHash").AssignBytes(txHash[:])
		ma.AssembleEntry("blockHash").AssignBytes(blockHash.Bytes())
		ma.AssembleEntry("chainId").AssignInt(chainID)
	})

	p := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   0x12, // sha2-256
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}

	lnk, err := s.LinkSystem.Store(ipld.LinkContext{
		LinkNode: transactionBlock,
	}, p, n)

	if err != nil {
		return nil, nil, err
	}

	return n, lnk, nil
}

func EncodeDagJsonEvent() abi.Event {

	str, _ := abi.NewType("string", "", nil)
	return abi.NewEvent(
		"EncodeDagJson",
		"EncodeDagJson",
		false,
		abi.Arguments{abi.Argument{
			Name:    "path",
			Type:    str,
			Indexed: false,
		}, abi.Argument{
			Name:    "hexdata",
			Type:    str,
			Indexed: false,
		}},
	)
}
func encodeDagJsonBlock(s anconsync.Storage, inputs abi.Arguments, data []byte, txHash, blockHash types.Hash, chainID int64) (datamodel.Node, datamodel.Link, error) {

	props, err := inputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}

	///	path := props[0].(string)
	values := props[1].(string)
	bz := common.Hex2Bytes(values)

	js := hexutil.Bytes{}
	js.UnmarshalJSON(bz)

	n, _ := anconsync.Decode(basicnode.Prototype.Any, string(js))

	// var nodelink datamodel.Link

	transactionBlock := fluent.MustBuildMap(basicnode.Prototype.Any, 3, func(ma fluent.MapAssembler) {
		// ma.AssembleEntry("dagblock").AssignLink(nodelink)
		ma.AssembleEntry("transactionHash").AssignBytes(txHash[:])
		ma.AssembleEntry("blockHash").AssignBytes(blockHash.Bytes())
		ma.AssembleEntry("chainId").AssignInt(chainID)
	})

	p := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x12, // sha2-256
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}

	lnk, err := s.LinkSystem.Store(ipld.LinkContext{
		LinkNode: transactionBlock,
	}, p, n)

	if err != nil {
		return nil, nil, err
	}

	return n, lnk, nil
}

func PostTxProcessing(s anconsync.Storage, t *state.Transition) error {
	for _, log := range t.Txn().Logs() {
		for _, topic := range log.Topics {

			if len(log.Topics) == 0 {
				continue
			}
			blockHash := t.GetBlockHash(t.GetTxContext().Number)
			txHash := blockHash

			var node datamodel.Node
			var lnk datamodel.Link
			var err error
			switch {
			case common.Hash(topic) == AddOnchainMetadataEvent().ID:
				node, lnk, err = encodeAnconMetadata(s, AddOnchainMetadataEvent().Inputs, log.Data, txHash, blockHash, t.GetTxContext().ChainID)

			case common.Hash(topic) == EncodeDagJsonEvent().ID:
				node, lnk, err = encodeDagJsonBlock(s, EncodeDagJsonEvent().Inputs, log.Data, txHash, blockHash, t.GetTxContext().ChainID)

			case common.Hash(topic) == EncodeDagCborEvent().ID:
				node, lnk, err = encodeDagCborBlock(s, EncodeDagCborEvent().Inputs, log.Data, txHash, blockHash, t.GetTxContext().ChainID)

			default:
				return fmt.Errorf("failed to decode")
			}
			// if !ContractAllowed(log.Address) {
			// 	// Check the contract whitelist to prevent accidental native call.
			// 	continue
			// }
			if err != nil {
				return err
			}

			fmt.Println(lnk.String())
			fmt.Println(node)
			StoreDagBlockDoneEvent().Inputs.Pack()
			t.EmitLog(log.Address, log.Topics, []byte(lnk.String()))

			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func GetHooks(s anconsync.Storage) func(t *state.Transition) {
	return func(t *state.Transition) {
		PostTxProcessing(s, t)
	}
}
