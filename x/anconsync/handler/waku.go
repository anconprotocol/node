package handler

import (
	"time"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"

	"fmt"
	"net"

	"github.com/status-im/go-waku/waku/v2/node"
	"github.com/status-im/go-waku/waku/v2/protocol"
	"github.com/status-im/go-waku/waku/v2/protocol/filter"
	"github.com/status-im/go-waku/waku/v2/protocol/lightpush"
	"github.com/status-im/go-waku/waku/v2/protocol/pb"
	"github.com/status-im/go-waku/waku/v2/protocol/relay"
	"github.com/status-im/go-waku/waku/v2/protocol/store"
	"github.com/status-im/go-waku/waku/v2/utils"

	"github.com/anconprotocol/sdk"

	"context"
)

type WakuHandler struct {
	Node        *node.WakuNode
	PeerAddress multiaddr.Multiaddr
}

func NewWakuHandler(ctx *sdk.AnconSyncContext, peerAddr string, address string, privateKeyPath string) *WakuHandler {
	hostAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprint(address))

	privateKey, err := crypto.GenerateOrReadPrivateKey(privateKeyPath)
	if err != nil {
		// try directly
		privateKey, err = crypto.BytesToPrivateKey([]byte(privateKeyPath))
	}
	wakuNode, err := node.New(context.Background(),
		node.WithPrivateKey(privateKey),
		node.WithHostAddress(hostAddr),
		node.WithWakuRelay(),
		node.WithLightPush(),
		node.WithWakuStore(true, true),
		node.WithWakuFilter(true),
	)
	if err != nil {
		panic(err)
	}
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		panic(err)
	}
	return &WakuHandler{
		Node:        wakuNode,
		PeerAddress: addr,
	}
}

func (h *WakuHandler) Start() {

	h.Node.AddPeer(h.PeerAddress, string(store.StoreID_v20beta4))
	h.Node.AddPeer(h.PeerAddress, string(lightpush.LightPushID_v20beta1))
	h.Node.AddPeer(h.PeerAddress, string(relay.WakuRelayID_v200))
	h.Node.AddPeer(h.PeerAddress, string(filter.FilterID_v20beta1))

	if err := h.Node.Start(); err != nil {
		panic(err)
	}
	h.Node.DialPeer(context.Background(), h.PeerAddress.String())
	waiting := true
	for i := 0; waiting; i++ {
		waiting = !(h.Node.Relay().EnoughPeersToPublish() == true)
		time.Sleep(100)
	}
}

func (h *WakuHandler) Subscribe(context context.Context, contentTopic string) (*relay.Subscription, error) {
	rel := h.Node.Relay()
	return rel.SubscribeToTopic(context, contentTopic)
}

func (h *WakuHandler) Publish(contentTopic protocol.ContentTopic, msgContent datamodel.Node) error {

	var version uint32 = 0
	var timestamp int64 = utils.GetUnixEpoch()

	p := new(node.Payload)
	data, err := msgContent.AsBytes()
	p.Data = data
	if err != nil {
		return errors.Wrap(err, "Bad dag block")
	}

	p.Key = &node.KeyInfo{Kind: node.None}

	payload, err := p.Encode(version)
	if err != nil {
		return errors.Wrap(err, "Encode error")
	}

	msg := &pb.WakuMessage{
		// any
		Payload: payload,
		// v1
		Version: version,
		// monike/1/appName/json
		ContentTopic: contentTopic.String(),
		Timestamp:    timestamp,
	}

	_, err = h.Node.Relay().Publish(context.Background(), msg)
	if err != nil {
		return errors.Wrap(err, "Error sending a message: ")
	}

	return nil
}

// func (h *WakuHandler) readLoop(ctx context.Context) {
// 	sub, err := h.Node.Relay().Subscribe(ctx)
// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 		return
// 	}

// 	for value := range sub.C {
// 		payload, err := node.DecodePayload(value.Message(), &node.KeyInfo{Kind: node.None})
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		// store the payload

// 	}
// }
