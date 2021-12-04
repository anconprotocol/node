package handler

import (
	"crypto/ecdsa"

	"github.com/anconprotocol/node/x/anconsync"
	"github.com/ipfs/go-graphsync"
	"github.com/libp2p/go-libp2p-core/peer"
)

type AnconSyncContext struct {
	Store      anconsync.Storage
	Exchange   graphsync.GraphExchange
	IPFSPeer   *peer.AddrInfo
	PrivateKey *ecdsa.PrivateKey
}

func NewAnconSyncContext(s anconsync.Storage, exchange graphsync.GraphExchange, ipfspeer *peer.AddrInfo, privateKey *ecdsa.PrivateKey) *AnconSyncContext {
	return &AnconSyncContext{
		Store:      s,
		Exchange:   exchange,
		IPFSPeer:   ipfspeer,
		PrivateKey: privateKey,
	}
}
