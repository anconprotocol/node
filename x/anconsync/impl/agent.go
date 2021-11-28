package impl

import (
	"context"

	gsync "github.com/ipfs/go-graphsync"
	graphsync "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"

	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func NewAgent(ctx context.Context, gsynchost host.Host, router string) gsync.GraphExchange {

	// The last step to get fully up and running would be to connect to
	// bootstrap peers (or any other peers). We leave this commented as
	// this is an example and the peer will die as soon as it finishes, so
	// it is unnecessary to put strain on the network.

	// This connects to public bootstrappers
	// for _, addr := range dht.DefaultBootstrapPeers {
	//	pi, _ := peer.AddrInfoFromP2pAddr(multiaddr.StringCast(bootstrap))
	// We ignore errors as some bootstrap peers may be down
	// and that is fine.
	// This connects to public bootstrappers
	var pi *peer.AddrInfo
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ = peer.AddrInfoFromP2pAddr(addr)
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		gsynchost.Connect(ctx, *pi)
	}

	network := gsnet.NewFromLibp2pHost(gsynchost)

	//add carv1
	var exchange gsync.GraphExchange
	exchange = graphsync.New(ctx, network, cidlink.DefaultLinkSystem())
	
	// finalResponseStatusChan := make(chan gsync.ResponseStatusCode, 1)
	// exchange.RegisterCompletedResponseListener(func(p peer.ID, request gsync.RequestData, status gsync.ResponseStatusCode) {
	// 	select {
	// 	case finalResponseStatusChan <- status:
	// 		fmt.Sprintf("%s", status)
	// 	default:
	// 	}
	// })
	// c, _ := cid.Parse("bafyreigiumx5ficjmdwdgpsxddfeyx2vh6cbod5s454pqeaosue33w2fpq")
	// link := cidlink.Link{Cid: c}
	// ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	// selector := ssb.ExploreAll(ssb.Matcher()).Node()

	// r := &net.Receiver{
	// 	MessageReceived: make(chan net.ReceivedMessage),
	// }

	// pi, _ = peer.AddrInfoFromP2pAddr(multiaddr.StringCast(router))
	// network.SetDelegate(r)
	// err := network.ConnectTo(ctx, pi.ID)
	// if err != nil {
	// 	panic(err)
	// }
	// pgChan, errChan := exchange.Request(ctx, pi.ID, link, selector)
	// net.VerifyHasErrors(ctx, errChan)
	// net.PrintProgress(ctx, pgChan)
	// defer gsynchost.Close()

	// a := fmt.Sprintf("%s/p2p/%s", gsynchost.Addrs()[0].String(), gsynchost.ID().Pretty())
	// fmt.Printf("Hello World, my hosts ID is %s\n", a)

	// var received gsmsg.GraphSyncMessage
	// var receivedBlocks []blocks.Block
	// for {
	// 	var message net.ReceivedMessage

	// 	sender := message.Sender
	// 	received = message.Message
	// 	fmt.Sprintf("%s %s", sender.String(), received)
	// 	receivedBlocks = append(receivedBlocks, received.Blocks()...)
	// 	receivedResponses := received.Responses()
	// 	if len(receivedResponses) == 0 {
	// 		continue
	// 	}
	// 	fmt.Sprintf("%s", receivedResponses[0].Status())
	// 	if receivedResponses[0].Status() != gsync.PartialResponse {
	// 		break
	// 	}
	// }

	return exchange

}
