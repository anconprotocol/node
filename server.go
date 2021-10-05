package main

import (
	"context"

	bridge "github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/aguaclara/bridge"
)

func main() {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// h1 := net.NewPeer(ctx, "/ip4/0.0.0.0/tcp/7779")

	bridge.NewAguaclara(
		ctx, "", "tcp://localhost:8899", "http://localhost:26657", "http://localhost:26657",

		1, "",
	)
	// cmd.NewRouter(ctx, h1)
	//  	run(ctx, h2, h1.Addrs()[0].String())
	// run(ctx, h2,
	// 	fmt.Sprintf("%s/p2p/%s", h1.Addrs()[0].String(), h1.ID().Pretty()))
}
