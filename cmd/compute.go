package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/qri-io/jsonschema"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type DagCompute struct {
	// dagContract string
	Storage Storage
}

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

func (d DagCompute) ReadFromStore(hash string, path string) (string, error) {
	lnk, err := ParseCidLink(string(hash))
	if err != nil {
		return "", fmt.Errorf("parse anchor link error %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("read trusted anchor JSON  error %v", err)
	}

	ss := d.Storage
	node, err := ss.Load(
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

func MockGetDataContract() string {
	return `[payload.lastName, payload.firstName, inputs.say, did]`
}

func MockSetOutputSignature() bool {
	return true
}

func FromJSON(j string) interface{} {
	var val interface{}
	json.Unmarshal([]byte(j), &val)
	return val
}

func GetDataContractGlobals(jsonArgs, did, payload string) map[string]interface{} {
	return map[string]interface{}{
		"inputs":  FromJSON(jsonArgs),
		"did":     did,
		"payload": FromJSON(payload),
	}

}

func GetDataContractEnvironment() *cel.Env {
	env, _ := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("inputs", decls.Dyn),
			decls.NewVar("did", decls.Dyn),
			decls.NewVar("payload", decls.Dyn),
		),
	)

	return env
}

func (d DagCompute) ExecuteDataContractTransaction(jsonSchemaCid string, inputPayload string, jsonArgs string, did string, contractAddress string, contractId string) (datamodel.Link, error) {
	//inputs: all cids available in the anchoring list
	// var anch []byte

	// if k.HasAnchor(ctx, msg.Did, msg.InputCid) {
	// 	anch = k.GetAnchor(ctx, msg.Did, msg.InputCid)
	// } else {
	// 	return nil, fmt.Errorf("invalid anchor link")

	// }

	jschem := jsonschema.Schema{}

	offSchema, err := d.ReadFromStore(jsonSchemaCid, "")

	err = jschem.UnmarshalJSON([]byte(offSchema))

	payl, _ := d.ReadFromStore(inputPayload, "")

	_, err = jschem.ValidateBytes(context.Background(), []byte(payl))

	if err != nil {
		return nil, fmt.Errorf("validate payload error %v", err)
	}

	dataContr := MockGetDataContract()

	if err != nil {
		return nil, fmt.Errorf("get contract error", dataContr)
	}

	ast, issues := GetDataContractEnvironment().Compile(string(dataContr))

	if len(issues.Errors()) > 0 {

		return nil, fmt.Errorf("reverted %v", issues.Errors())
	}

	prog, err := GetDataContractEnvironment().Program(ast)
	if err != nil {
		return nil, fmt.Errorf("env error: %v", err)
	}
	out, _, err := prog.Eval(
		GetDataContractGlobals(jsonArgs, did, payl),
	)
	if err != nil {
		return nil, fmt.Errorf("eval node error", err)
	}

	node, err := Decode(basicnode.Prototype.Any, ValueToJSON(out))
	if err != nil {
		return nil, fmt.Errorf("Decode error", err)
	}

	clink := d.Storage.Store(ipld.LinkContext{}, node)

	return clink, nil
}
