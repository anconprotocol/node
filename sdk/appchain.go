package sdk

import (
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)

var _ abcitypes.Application = (*AnconAppChain)(nil)

type AnconAppChain struct {
	abcitypes.Application

	StorageManager Storage
	Store          types.CommitMultiStore
}

func NewAnconAppChain(key string, db dbm.DB) *AnconAppChain {
	store := store.NewCommitMultiStore(db)

	return &AnconAppChain{
		Store:          store,
		StorageManager: NewStorage(key, store, db),
	}
}

func (app *AnconAppChain) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *AnconAppChain) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {

	return abcitypes.ResponseInfo{
		// Data:             "",
		Version:          version.ABCIVersion,
		AppVersion:       1,
		LastBlockHeight:  app.Store.LastCommitID().Version,
		LastBlockAppHash: app.Store.LastCommitID().Hash,
	}
}

func (app *AnconAppChain) isValid(data []byte) uint32 {
	return 0
}

func (app *AnconAppChain) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	// node := basicnode.NewBytes(req.Tx)
	// issuer, _ := node.LookupByString("issuer")
	// cid, _ := node.LookupByString("contentHash")
	// sig, _ := node.LookupByString("signature")
	// path, _ := node.LookupByString("path")

	// events := []abcitypes.Event{
	// 	{
	// 		Type: "dagblock",
	// 		Attributes: []abcitypes.EventAttribute{
	// 			{Key: []byte("issuer"), Value: []byte(must.String(issuer)), Index: true},
	// 			{Key: []byte("contentHash"), Value: []byte(must.String(cid)), Index: true},
	// 			{Key: []byte("signature"), Value: []byte(must.String(sig)), Index: true},
	// 			{Key: []byte("path"), Value: []byte(must.String(path)), Index: true},
	// 		},
	// 	},
	// }
	return abcitypes.ResponseDeliverTx{Code: abcitypes.CodeTypeOK}
}

func (app *AnconAppChain) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Validate block exists and is signed
	return abcitypes.ResponseCheckTx{Code: 0}
}

func (app *AnconAppChain) Commit() abcitypes.ResponseCommit {

	res := app.Store.Commit()

	return abcitypes.ResponseCommit{Data: res.Hash}
}

func (app *AnconAppChain) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {

	st := app.StorageManager.dataStore
	s := st.GetCommitStore(STORE_KEY_TYPE)
	iavlstore := s.(*iavl.Store)

	queryableStore := store.Queryable(iavlstore)

	return queryableStore.Query(req)
}
func (AnconAppChain) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *AnconAppChain) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	// no op
	return abcitypes.ResponseBeginBlock{}
}

func (app *AnconAppChain) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *AnconAppChain) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (AnconAppChain) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (AnconAppChain) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (app *AnconAppChain) ApplySnapshotChunk(req abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
