package protocol

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/second-state/WasmEdge-go/wasmedge"
)

type ProtocolAPI struct {
	Namespace string
	Version   string
	Service   *ProtocolService
	Public    bool
}

type ProtocolService struct {
	Adapter *ethereum.OnchainAdapter
	Storage *sdk.Storage
	Proof   *proofsignature.IavlProofAPI
	Host    *Host
	wasm    *wasmedge.VM
}

func NewProtocolAPI(adapter *ethereum.OnchainAdapter, storage *sdk.Storage, proof *proofsignature.IavlProofAPI) *ProtocolAPI {

	host := NewEvmRelayHost(storage, proof, adapter)

	wasmedge.SetLogErrorLevel()

	/// Create configure
	var conf = wasmedge.NewConfigure(wasmedge.WASI)

	/// Create VM with configure
	var vm = wasmedge.NewVMWithConfig(conf)

	/// Init WASI
	var wasi = vm.GetImportObject(wasmedge.WASI)
	wasi.InitWasi(
		os.Args[1:],     /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping preopens
	)

	vm.RegisterImport(host.GetImports())

	return &ProtocolAPI{
		Namespace: "Protocol",
		Version:   "1.0",
		Service: &ProtocolService{
			Adapter: adapter,
			Storage: storage,
			Proof:   proof,
			Host:    host,
			wasm:    vm,
		},
		Public: true,
	}
}

func (s *ProtocolService) Call(to string, from string, data string) hexutil.Bytes {

	val := make(map[string]string, 2)

	payload, err := hexutil.Decode(data)

	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("fail unpack data").Error())))
	}
	// Execute graphql
	toClink, err := sdk.ParseCidLink(to)
	if err != nil {

	}

	dataNode, err := s.Storage.Load(linking.LinkContext{}, toClink)
	if err != nil {

	}
	dataDecoded, _ := sdk.Encode(dataNode)
	err = s.wasm.LoadWasmBuffer([]byte(dataDecoded))
	if err != nil {

	}
	//TODO Validate user signature
	s.wasm.Validate()
	s.wasm.Instantiate()
	res, err := s.wasm.ExecuteBindgen("execute", wasmedge.Bindgen_return_array, payload)
	if err != nil {

	}

	s.wasm.Cleanup()

	val["result"] = string(res.([]byte))
	jsonval, err := json.Marshal(val)
	if err != nil {
		return hexutil.Bytes(hexutil.Encode([]byte(fmt.Errorf("reverted, json marshal").Error())))
	}
	return jsonval
}
