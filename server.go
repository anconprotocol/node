package main

import (
	"context"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/cmd"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
)

func main() {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// h1 := net.NewPeer(ctx, "/ip4/0.0.0.0/tcp/7779")
	h2 := net.NewPeer(ctx, "/ip4/0.0.0.0/tcp/7777")

	router := "/ip4/192.168.50.138/tcp/7779/p2p/12D3KooWLNeo1sqTtMsrReTurqLTQ7fdjGwPaXsEBMbgnTgBJEbt"
	cmd.NewEdge(ctx, h2, router)
	// cmd.NewRouter(ctx, h1)
	//  	run(ctx, h2, h1.Addrs()[0].String())
	// run(ctx, h2,
	// 	fmt.Sprintf("%s/p2p/%s", h1.Addrs()[0].String(), h1.ID().Pretty()))
}
