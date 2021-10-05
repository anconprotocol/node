package bridge

import (
	"context"
	"fmt"

	"os"
	"time"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/light/proxy"

	httpprovider "github.com/tendermint/tendermint/light/provider/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/jsonrpc/server"
	dbm "github.com/tendermint/tm-db"
	badger "github.com/tendermint/tm-db/badgerdb"

	"github.com/tendermint/tendermint/light/provider"
	dbs "github.com/tendermint/tendermint/light/store/db"

	"github.com/tendermint/tendermint/light"
)

type AguaclaraAdapter struct {
	Height      int
	AppHash     []byte
	LightClient *light.Client
	DB          dbm.DB
	Proxy       *proxy.Proxy
}

func (adapter *AguaclaraAdapter) ReceivedHTLAMessage(evt *coretypes.ResultEvent) ([]byte, error) {
	//	types.Conte
	// rpctypes.Context

	// Send to chain b

	fmt.Sprintf("%s", evt.Data)

	return nil, nil
}

func NewAguaclara(ctx context.Context,
	chainID string, proxyAddr string,
	primary string, witness string,
	height int,
	appHash string) (*AguaclaraAdapter, error) {

	db, err := badger.NewDB("anconnode", "/tmp/badger")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	hash := []byte(appHash)

	node, err := NewLightTendermint(ctx,
		chainID,
		primary, witness, height, hash, dbm.DB(db))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	// // rpc
	c := server.DefaultConfig()
	proxy, err := proxy.NewProxy(node, proxyAddr, primary, c, log.NewNopLogger())
	// proxyerr := proxy.ListenAndServe()

	q := "tm.event='Tx'" /// AND ethereum.recipient='hexAddress'"
	a := &AguaclaraAdapter{
		Height:      height,
		AppHash:     hash,
		DB:          db,
		Proxy:       proxy,
		LightClient: node,
	}
	outChan, _ := proxy.Client.Subscribe(ctx, "localhost", q)

	select {
	case msg := <-outChan:
		a.ReceivedHTLAMessage(&msg)
	default:

	}
	return a, nil
}

func NewLightTendermint(ctx context.Context, chainID string,
	primary string, witness string, height int, hash []byte, db dbm.DB) (*light.Client, error) {

	primaryNode, _ := httpprovider.New(chainID, primary)
	witnessNode, _ := httpprovider.New(chainID, witness)
	c, _ := light.NewClient(
		ctx,
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: int64(height),
			Hash:   hash,
		},
		primaryNode,
		[]provider.Provider{witnessNode},
		dbs.New(db, "ancon-node"),
	)
	//_, err := c.Update(ctx, time.Now())

	return c, nil
}
