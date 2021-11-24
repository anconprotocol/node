module github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync

go 1.16

require (
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/ipfs/go-block-format v0.0.3
	github.com/ipfs/go-cid v0.1.0
	github.com/ipfs/go-graphsync v0.9.3
	github.com/ipfs/go-ipfs-blockstore v1.0.4 // indirect
	github.com/ipld/go-car/v2 v2.0.2
	github.com/ipld/go-ipld-prime v0.14.0
	github.com/libp2p/go-libp2p v0.14.4
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.8.6
	github.com/libp2p/go-libp2p-kad-dht v0.13.1
	github.com/libp2p/go-libp2p-noise v0.2.0
	github.com/multiformats/go-multiaddr v0.3.3
	github.com/prometheus/procfs v0.7.0 // indirect
	github.com/proofzero/go-ipld-linkstore v1.0.0
	github.com/smartystreets/assertions v1.1.1 // indirect
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.5
)

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/snappy v0.0.4 // indirect
	github.com/multiformats/go-multihash v0.1.0
	github.com/spf13/cast v1.4.1
	github.com/tharsis/ethermint v0.7.1
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
)

replace github.com/libp2p/go-libp2p => github.com/libp2p/go-libp2p v0.14.1

replace github.com/libp2p/go-libp2p-core v0.10.0 => github.com/libp2p/go-libp2p-core v0.9.0

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
