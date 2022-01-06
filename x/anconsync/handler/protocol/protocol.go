package protocol

import (
	"context"
	"fmt"
	"os"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/anconprotocol/node/x/anconsync/handler/hexutil"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
	"github.com/buger/jsonparser"
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
		Namespace: "ancon",
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

func (s *ProtocolService) Call(to string, from string, sig []byte, data string) string {
	doc, err := s.Storage.DataStore.Get(context.Background(), from)
	if err != nil {
		return (hexutil.Encode([]byte(fmt.Errorf("invalid signature").Error())))
	}

	didDoc, err := types.GetDidDocument(string(doc))
	hash := crypto.Keccak256([]byte(data))
	ok, err := types.Authenticate(didDoc, hash, string(sig))
	if !ok {
		return (hexutil.Encode([]byte(fmt.Errorf("user must registered as a did").Error())))
	}
	has := s.wasm.GetFunctionTypeRegistered(to, "execute")
	if has == nil {
		toClink, err := sdk.ParseCidLink(to)
		if err != nil {
			return (hexutil.Encode([]byte(fmt.Errorf("invalid cid link").Error())))
		}
		dataNode, err := s.Storage.Load(linking.LinkContext{}, toClink)
		if err != nil {
			return (hexutil.Encode([]byte(fmt.Errorf("no contract cid found").Error())))
		}
		dataDecoded, _ := sdk.Encode(dataNode)
		code, _ := jsonparser.GetString([]byte(dataDecoded), "code")
		buf, err := hexutil.Decode(code)

		if err != nil {
			return (hexutil.Encode([]byte(fmt.Errorf("invalid wasm contract, must be base64 encoded").Error())))
		}
		err = s.wasm.RegisterWasmBuffer(to, buf)
		if err != nil {
			return (hexutil.Encode([]byte(fmt.Errorf("invalid wasm contract, error while loading").Error())))
		}
	}
	//TODO Validate user signature

	///s.wasm.Validate()
	///	s.wasm.Instantiate()
	res, err := s.wasm.ExecuteBindgenRegistered(to, "execute", wasmedge.Bindgen_return_array, []byte(data))
	if err != nil {
		return (hexutil.Encode([]byte(fmt.Errorf("reverted, json marshal").Error())))
	}

	defer s.wasm.Cleanup()

	return hexutil.Encode(res.([]byte))
}
