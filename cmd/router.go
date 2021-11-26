package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	gsync "github.com/ipfs/go-graphsync"
	graphsync "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/multiformats/go-multiaddr"

	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func NewRouter(ctx context.Context, gsynchost host.Host, s Storage, peerhost string) (gsync.GraphExchange, *peer.AddrInfo) {

	var pi *peer.AddrInfo
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ = peer.AddrInfoFromP2pAddr(addr)
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		gsynchost.Connect(ctx, *pi)
	}

	network := gsnet.NewFromLibp2pHost(gsynchost)

	// Add Ancon fsstore
	exchange := graphsync.New(ctx, network, s.LinkSystem)

	// var receivedResponseData []byte
	// var receivedRequestData []byte

	exchange.RegisterIncomingResponseHook(
		func(p peer.ID, responseData gsync.ResponseData, hookActions gsync.IncomingResponseHookActions) {
			fmt.Println(responseData.Status().String(), responseData.RequestID())
		})

	exchange.RegisterIncomingRequestHook(func(p peer.ID, requestData gsync.RequestData, hookActions gsync.IncomingRequestHookActions) {
		// var has bool
		// receivedRequestData, has = requestData.Extension(td.extensionName)
		// if !has {
		// 	hookActions.TerminateWithError(errors.New("Missing extension"))
		// } else {
		// 	hookActions.SendExtensionData(td.extensionResponse)
		// }
		hookActions.ValidateRequest()

		has, _ := s.DataStore.Has(ctx, requestData.Root().String())
		if !has {
			hookActions.TerminateWithError(errors.New("not found"))
			net.FetchBlock(ctx, exchange, p, cidlink.Link{Cid: requestData.Root()})
		}
		hookActions.UseLinkTargetNodePrototypeChooser(basicnode.Chooser)
		fmt.Println(requestData.Root(), requestData.ID(), requestData.IsCancel())
	})
	finalResponseStatusChan := make(chan gsync.ResponseStatusCode, 1)
	exchange.RegisterCompletedResponseListener(func(p peer.ID, request gsync.RequestData, status gsync.ResponseStatusCode) {
		select {
		case finalResponseStatusChan <- status:
		default:
		}
	})
	pi, _ = peer.AddrInfoFromP2pAddr(multiaddr.StringCast(peerhost))
	// err := network.ConnectTo(ctx, pi.ID)
	// if err != nil {
	// 	panic(err)
	// }

	a := fmt.Sprintf("%s/p2p/%s", gsynchost.Addrs()[0].String(), gsynchost.ID().Pretty())

	//	defer gsynchost.Close()

	go func() {
		fmt.Printf("Ancon Router Sync peer id is %s\n", a)
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
	}()
	return exchange, pi
}
