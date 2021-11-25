package cmd

import (
	"context"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	blocks "github.com/ipfs/go-block-format"
	gsync "github.com/ipfs/go-graphsync"
	graphsync "github.com/ipfs/go-graphsync/impl"
	gsmsg "github.com/ipfs/go-graphsync/message"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/multiformats/go-multiaddr"

	"github.com/ipld/go-ipld-prime/linking"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func NewRouter(ctx context.Context, gsynchost host.Host, lsys linking.LinkSystem, peerhost string) string {

	var pi *peer.AddrInfo
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ = peer.AddrInfoFromP2pAddr(addr)
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		gsynchost.Connect(ctx, *pi)
	}

	network := gsnet.NewFromLibp2pHost(gsynchost)

	// Add Ancon fsstore
	exchange := graphsync.New(ctx, network, lsys)

	finalResponseStatusChan := make(chan gsync.ResponseStatusCode, 1)
	exchange.RegisterCompletedResponseListener(func(p peer.ID, request gsync.RequestData, status gsync.ResponseStatusCode) {
		select {
		case finalResponseStatusChan <- status:
			fmt.Sprintf("%s", status)
		default:
		}
	})

	pi, _ = peer.AddrInfoFromP2pAddr(multiaddr.StringCast(peerhost))
	r := &net.Receiver{
		MessageReceived: make(chan net.ReceivedMessage),
	}

	network.SetDelegate(r)
	err := network.ConnectTo(ctx, pi.ID)
	if err != nil {
		panic(err)
	}

	defer gsynchost.Close()
	go func() {
		a := fmt.Sprintf("%s/p2p/%s", gsynchost.Addrs()[0].String(), gsynchost.ID().Pretty())
		fmt.Printf("Ancon Router Sync peer id is %s\n", a)

		var received gsmsg.GraphSyncMessage
		var receivedBlocks []blocks.Block
		for {
			var message net.ReceivedMessage

			sender := message.Sender
			received = message.Message
			fmt.Sprintf("%s %s", sender.String(), received)
			receivedBlocks = append(receivedBlocks, received.Blocks()...)
			receivedResponses := received.Responses()
			if len(receivedResponses) == 0 {
				continue
			}
			fmt.Sprintf("%s", receivedResponses[0].Status())
			if receivedResponses[0].Status() != gsync.PartialResponse {
				break
			}
		}
	}()
	return ""
}
