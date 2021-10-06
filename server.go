package main

import (
	"context"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/bridge"
)

func main() {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// h1 := net.NewPeer(ctx, "/ip4/0.0.0.0/tcp/7779")

	aguaclara, _ := bridge.NewAguaclara(
		ctx, "", "tcp://localhost:8899", "http://localhost:26657", "http://localhost:26657",
		819, "7CF4887E64FBA3EAA9A73D8D25B3DB80CB543800EA6F4D0F3D88D75E480049E8",

	aguaclara.Proxy.ListenAndServe()
	// cmd.NewRouter(ctx, h1)
	//  	run(ctx, h2, h1.Addrs()[0].String())
	// run(ctx, h2,
	// 	fmt.Sprintf("%s/p2p/%s", h1.Addrs()[0].String(), h1.ID().Pretty()))
}
