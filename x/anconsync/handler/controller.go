package handler

import (
	"github.com/99designs/keyring"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/ipfs/go-graphsync"
	"github.com/libp2p/go-libp2p-core/peer"
)

type AnconSyncContext struct {
	Store    anconsync.Storage
	Exchange graphsync.GraphExchange
	IPFSPeer *peer.AddrInfo
	Keyring  keyring.Keyring
}

func NewAnconSyncContext(s anconsync.Storage, exchange graphsync.GraphExchange, ipfspeer *peer.AddrInfo, keyring keyring.Keyring) *AnconSyncContext {
	return &AnconSyncContext{
		Store:    s,
		Exchange: exchange,
		IPFSPeer: ipfspeer,
		Keyring:  keyring,
	}
}
