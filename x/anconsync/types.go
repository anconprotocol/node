package anconsync

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"github.com/google/cel-go/common/types/ref"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-graphsync"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multihash"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// From: https://github.com/google/cel-go/blob/master/codelab/solution/codelab.go
// valueToJSON converts the CEL type to a protobuf JSON representation and
// marshals the result to a string.
func ValueToJSON(val ref.Val) string {
	v, err := val.ConvertToNative(reflect.TypeOf(&structpb.Value{}))
	if err != nil {
		glog.Exit(err)
	}
	marshaller := protojson.MarshalOptions{Indent: "    "}
	bytes, err := marshaller.Marshal(v.(proto.Message))
	if err != nil {
		glog.Exit(err)
	}
	return string(bytes)
}

func ReadFromStore(s Storage, hash string, path string) (string, error) {
	lnk, err := ParseCidLink(string(hash))
	if err != nil {
		return "", fmt.Errorf("parse anchor link error %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("read trusted anchor JSON  error %v", err)
	}
	//CURL fields:
	//grabar un schema cid
	//save with prop called schema in plain text
	//json base 64
	//variable json asplaintext
	//op (operations) string
	//resolve path
	//query
	node, err := s.Load(
		ipld.LinkContext{
			LinkPath: datamodel.ParsePath(path),
		},
		lnk,
	)
	if err != nil {
		return "", err
	}

	output, err := Encode(node)

	if err != nil {
		return "", err
	}

	return output, nil
}

// ParseCidLink parses a string cid multihash into a cidLink
func ParseCidLink(hash string) (cidlink.Link, error) {
	lnk, err := cid.Parse(hash)
	if err != nil {
		return cidlink.Link{}, status.Error(
			3,
			"Invalid CID Link",
		)
	}

	return cidlink.Link{Cid: lnk}, nil
}

type DagContractTrustedContext struct {
	Store    Storage
	Exchange graphsync.GraphExchange
}

func GetLinkPrototype() ipld.LinkPrototype {
	// tip: 0x0129 dag-json
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x71, // dag-cbor
		MhType:   0x12, // sha2-256
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}
}

// CreateCidLink takes a hash eg ethereum hash and converts it to cid multihash
func CreateCidLink(hash []byte) cidlink.Link {
	lchMh, err := multihash.Encode(hash, GetLinkPrototype().(cidlink.LinkPrototype).MhType)
	if err != nil {
		return cidlink.Link{}
	}
	lcCID := cid.NewCidV1(GetLinkPrototype().(cidlink.LinkPrototype).Codec, lchMh)
	lcLinkCID := cidlink.Link{Cid: lcCID}
	return lcLinkCID
}
