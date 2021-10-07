package main

import (
	"context"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/bridge"
)

func main() {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// h1 := net.NewPeer(ctx, "/ip4/0.0.0.0/tcp/7779")

	aguaclara, _ := bridge.NewAguaclara(
		ctx,
		"anconprotocol_9000-1",
		"tcp://localhost:8899",
		"http://ancon.did.pa:26657",
		"http://ancon.did.pa:26657",
		27,
		"6D58C14836E7A951D06684FA2AB515835B3C9EB068DBD1ADF8EA58E6F5FD5294",
	)
	err := aguaclara.Proxy.ListenAndServe()

	fmt.Errorf("%s", err)
	// cmd.NewRouter(ctx, h1)
	//  	run(ctx, h2, h1.Addrs()[0].String())
	// run(ctx, h2,
	// 	fmt.Sprintf("%s/p2p/%s", h1.Addrs()[0].String(), h1.ID().Pretty()))
}
