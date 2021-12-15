package handler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-ipld-prime"
	"github.com/second-state/WasmEdge-go/wasmedge"

	"github.com/anconprotocol/contracts/wasmvm"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
)

type WasmVMHandler struct {
	*sdk.AnconSyncContext
	proofs *proofsignature.IavlProofAPI
	vm     *wasmedge.VM
}

func NewWasmVMHandler(ctx *sdk.AnconSyncContext, proofs *proofsignature.IavlProofAPI) *WasmVMHandler {

	host := wasmvm.NewHost(ctx.Store, proofs)
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

	var type1 = wasmedge.NewFunctionType(
		[]wasmedge.ValType{

			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		}, []wasmedge.ValType{
			wasmedge.ValType_I32,
		})
	var type2 = wasmedge.NewFunctionType(
		[]wasmedge.ValType{

			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		}, []wasmedge.ValType{
			//		wasmedge.ValType_I32,
		})
	n := wasmedge.NewImportObject("env")
	fn1 := wasmedge.NewFunction(type2, host.WriteStore, nil, 0)
	n.AddFunction("write_store", fn1)

	fn2 := wasmedge.NewFunction(type1, host.ReadStore, nil, 0)
	n.AddFunction("read_store", fn2)

	fn3 := wasmedge.NewFunction(type2, host.ReadDagBlock, nil, 0)
	n.AddFunction("read_dag_block", fn3)

	fn4 := wasmedge.NewFunction(type1, host.WriteDagBlock, nil, 0)
	n.AddFunction("write_dag_block", fn4)
	vm.RegisterImport(n)
	// wasi.InitWasi(

	defer vm.Release()

	defer conf.Release()

	return &WasmVMHandler{AnconSyncContext: ctx, proofs: proofs, vm: vm}

}

// @BasePath /v0
// Execute godoc
// @Summary Executes a WASM smart contract
// @Schemes
// @Description ...
// @Tags proofs
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/contracts/execute [post]
func (dagctx *WasmVMHandler) Execute(c *gin.Context) {

	v, _ := c.GetRawData()

	cid, _ := jsonparser.GetString(v, "contractCid")

	if cid == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing cid").Error(),
		})
		return
	}

	tx, _ := jsonparser.GetString(v, "transaction")

	if cid == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing transaction").Error(),
		})
		return
	}

	var txIface map[string]interface{}
	err := json.Unmarshal([]byte(tx), txIface)

	lnk, err := sdk.ParseCidLink(cid)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("cid error %v", err).Error(),
		})
		return
	}
	contract, err := dagctx.Store.Load(ipld.LinkContext{}, lnk)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("loading error %v", err).Error(),
		})
		return
	}
	payload, err := sdk.EncodeCBOR(contract)
	var bz []byte
	cbornode.DecodeInto(payload, bz)
	vm := dagctx.vm
	vm.LoadWasmBuffer(bz)
	vm.Validate()
	vm.Instantiate()

	// example: `query { metadata(cid:"babfy",path:"/")  {image}}`
	res, err := vm.ExecuteBindgen("execute",
		wasmedge.Bindgen_return_array,
		txIface["data"].([]byte),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("lastHash Error %v", err).Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"result": string(res.([]byte)),
	})
}
