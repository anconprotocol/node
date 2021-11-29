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
	github.com/smartystreets/assertions v1.1.1 // indirect
)

require (
	github.com/99designs/gqlgen v0.14.0
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/buger/jsonparser v1.1.1
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/ethereum/go-ethereum v1.10.13
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gin-gonic/gin v1.7.0
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang/glog v1.0.0
	github.com/google/cel-go v0.9.0
	github.com/graphql-go/graphql v0.8.0
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/multiformats/go-multihash v0.1.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/spf13/cast v1.4.1
	github.com/vektah/gqlparser/v2 v2.2.0
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/libp2p/go-libp2p => github.com/libp2p/go-libp2p v0.14.1

replace github.com/libp2p/go-libp2p-core v0.10.0 => github.com/libp2p/go-libp2p-core v0.9.0

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
