package handler

import (
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/ipfs/go-graphsync"
	"github.com/libp2p/go-libp2p-core/peer"
)

type DagContractTrustedContext struct {
	Store    anconsync.Storage
	Exchange graphsync.GraphExchange
	IPFSPeer *peer.AddrInfo
}

func NewDagContractContext(s anconsync.Storage, exchange graphsync.GraphExchange, ipfspeer *peer.AddrInfo) *DagContractTrustedContext {
	return &DagContractTrustedContext{
		Store:    s,
		Exchange: exchange,
		IPFSPeer: ipfspeer,
	}
}
