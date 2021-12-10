module github.com/anconprotocol/node

go 1.16

require (
	github.com/anconprotocol/contracts v0.1.1
	github.com/anconprotocol/sdk v0.1.0
	github.com/ipfs/go-cid v0.1.0
	github.com/ipld/go-ipld-prime v0.14.0
)

require (
	github.com/99designs/gqlgen v0.14.0
	github.com/buger/jsonparser v0.0.0-20181115193947-bf1c66bbce23
	github.com/ethereum/go-ethereum v1.10.13
	github.com/gin-gonic/gin v1.7.4
	github.com/hyperledger/aries-framework-go v0.1.7
	github.com/multiformats/go-multibase v0.0.3
	github.com/multiformats/go-multicodec v0.3.0
	github.com/spf13/cast v1.4.1
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.3.3
	github.com/swaggo/swag v1.7.6
	github.com/tendermint/tendermint v0.35.0
	github.com/tendermint/tm-db v0.6.6
)

replace github.com/libp2p/go-libp2p => github.com/libp2p/go-libp2p v0.14.1

replace github.com/libp2p/go-libp2p-core v0.10.0 => github.com/libp2p/go-libp2p-core v0.9.0

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
