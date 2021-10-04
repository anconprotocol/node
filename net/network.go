package net

import (
	"context"
	"fmt"

	"time"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	gsync "github.com/ipfs/go-graphsync"
	gsmsg "github.com/ipfs/go-graphsync/message"
	blockstore "github.com/ipld/go-car/v2/blockstore"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	noise "github.com/libp2p/go-libp2p-noise"
)

func NewPeer(ctx context.Context, addr string) host.Host {
	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}

	gsynchost, err := libp2p.New(
		ctx,
		// Use the keypair we generated
		libp2p.Identity(priv),
		libp2p.Security(noise.ID, noise.New),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(addr),

		// support TLS connections
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,         // Lowwater
			400,         // HighWater,
			time.Minute, // GracePeriod
		)),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(newDHT),
		// Let this host use relays and advertise itself on relays if
		// it finds it is behind NAT. Use libp2p.Relay(options...) to
		// enable active relays and more.
		libp2p.EnableAutoRelay(),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}
	return gsynchost
}

//WriteCAR
func ReadCAR() ([]cid.Cid, blocks.Block, datamodel.Node, error) {
	//lsys := linkstore.NewStorageLinkSystemWithNewStorage(cidlink.DefaultLinkSystem())
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	selector := ssb.ExploreAll(ssb.Matcher()).Node()

	// car := carv1.NewSelectiveCar(context.Background(),
	// 	lsys.ReadStore,
	// 	[]carv1.Dag{{
	// 		Root:     root,
	// 		Selector: selector,
	// 	}})
	// file, err := os.ReadFile(filename)
	// if err != nil {
	// 	return err
	// }

	robs, _ := blockstore.OpenReadOnly("/home/dallant/Code/ancon-node/dagbridge-block-239-begin.car",
		blockstore.UseWholeCIDs(true),
	)

	roots, err := robs.Roots()

	res, _ := robs.Get(roots[0])

	return roots, res, selector, err
}

type ReceivedMessage struct {
	Message gsmsg.GraphSyncMessage
	Sender  peer.ID
}

// Receiver is an interface for receiving messages from the GraphSyncNetwork.
type Receiver struct {
	MessageReceived chan ReceivedMessage
}

func (r *Receiver) ReceiveMessage(
	ctx context.Context,
	sender peer.ID,
	incoming gsmsg.GraphSyncMessage) {

	select {
	case <-ctx.Done():
	case r.MessageReceived <- ReceivedMessage{incoming, sender}:
	}
}

func (r *Receiver) ReceiveError(_ peer.ID, err error) {
	fmt.Println("got receive err")
}

func (r *Receiver) Connected(p peer.ID) {
}

func (r *Receiver) Disconnected(p peer.ID) {
}

// VerifyHasErrors verifies that at least one error was sent over a channel
func VerifyHasErrors(ctx context.Context, errChan <-chan error) error {
	errCount := 0
	for {
		select {
		case e, ok := <-errChan:
			if ok {
				return nil
			} else {
				return e
			}
			errCount++
		case <-ctx.Done():
		}
	}
}

// VerifyHasErrors verifies that at least one error was sent over a channel
func PrintProgress(ctx context.Context, pgChan <-chan gsync.ResponseProgress) {
	errCount := 0
	for {
		select {
		case data, ok := <-pgChan:
			if ok {
				fmt.Sprintf("path: %s, last path: %s", data.Path.String(), data.LastBlock.Path.String())
			}
			errCount++
		case <-ctx.Done():
		}
	}
}
