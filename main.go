package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"strings"

	"github.com/0xPolygon/polygon-sdk/helper/keccak"
	"github.com/anconprotocol/node/docs"
	"github.com/anconprotocol/node/x/anconsync/handler"
	"github.com/anconprotocol/node/x/anconsync/handler/protocol/ethereum"
	"github.com/anconprotocol/node/x/anconsync/handler/types"
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/impl"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gsync "github.com/ipfs/go-graphsync"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/lucas-clemente/quic-go/http3"
	multiaddr "github.com/multiformats/go-multiaddr"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SubgraphConfig struct {
	CosmosPrimaryAddress string

	CosmosHeight    int
	EnableDagcosmos bool

	EvmAddress   string
	EvmChainId   string
	EnableDageth bool
}

// @title        Ancon Protocol Sync API v0.4.0
// @version      0.4.0
// @description  API

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.ancon.did.pa
// @BasePath  /v0
func main() {

	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/34075/p2p/16Uiu2HAmAihQ3QDNyNxfYiN8qAvaBZKzqosmudoG9KYFMhW2YDXd", "A remote peer to sync")
	addr := flag.String("addr", "/ip4/0.0.0.0/tcp/7702", "Host multiaddr")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")
	enableCors := flag.Bool("cors", false, "Allow CORS")
	allowOrigins := flag.String("origins", "*", "Must send a delimited string by commas")
	init := flag.Bool("init", false, "Initialize merkle tree storage")
	genKeys := flag.Bool("keys", false, "Generate keys")
	hostName := flag.String("hostname", "cerro-ancon", "Send custom host name")
	rootHash := flag.String("roothash", "", "root hash")
	rootkey := flag.String("rootkey", "", "root key")
	sync := flag.Bool("sync", false, "Syncronizes remote dag storage")
	seedPeers := flag.String("peers", "", "Array of peer addresses ")
	quic := flag.Bool("quic", false, "Enable QUIC")
	tlsKey := flag.String("tlscert", "", "TLS certificate")
	tlsCert := flag.String("tlskey", "", "TLS key")
	ipfshost := flag.String("ipfshost", "", "IPFS Host")

	subgraph := SubgraphConfig{}
	subgraph.EnableDageth = *flag.Bool("enable-dageth", false, "enable EVM subgraph")
	subgraph.EnableDagcosmos = *flag.Bool("enable-dagcosmos", false, "enable Cosmos subgraph")
	subgraph.CosmosPrimaryAddress = *flag.String("cosmos-primary", "", "primary")
	subgraph.EvmAddress = *flag.String("evm-node-address", "", "remote node address")
	subgraph.EvmChainId = *flag.String("evm-chain-id", "", "chain id")

	flag.Parse()
	if *genKeys == true {
		result, err := handler.GenerateKeys()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
		return
	}

	r := gin.Default()
	config := cors.DefaultConfig()

	if *enableCors == true {
		config.AllowOrigins = strings.Split(*allowOrigins, ",")
		r.Use(cors.New(config))
	}

	s := sdk.NewStorage(*dataFolder)

	ctx := context.Background()
	host := impl.NewPeer(ctx, *addr)

	exchange, ipfspeer := impl.NewRouter(ctx, host, s.LinkSystem, *peerAddr)

	fmt.Println(ipfspeer.ID)

	docs.SwaggerInfo.BasePath = "/v0"

	dagHandler := sdk.NewAnconSyncContext(s, exchange, ipfspeer)

	if *init == true {

		// Set your own keypair
		priv, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		var digest []byte

		keccak.Keccak256(digest, []byte(*hostName))
		signed, err := priv.Sign(rand.Reader, digest, nil)

		if err != nil {
			panic(err)
		}

		lnkCtx := ipld.LinkContext{
			LinkPath: ipld.ParsePath(types.ROOT_PATH),
		}

		n :=basicnode.NewString( base64.RawStdEncoding.EncodeToString(signed))

		link := dagHandler.Store.Store(lnkCtx, n) //Put(ctx, key, []byte(key))

		result, _, err := handler.InitGenesis(*hostName, link, priv)

		if err != nil {
			panic(err)
		}

		fmt.Println(result)

		return
	}

	proofHandler := handler.NewProofHandler(dagHandler)

	if *rootHash != "" && *sync == false {
		hash, err := proofHandler.VerifyGenesis(*rootkey)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("key validated with last current root hash: ", string(hash))

	}

	splitPeers := strings.Split(*seedPeers, ",")

	if *sync && len(splitPeers) > 0 {
		items := make([]peer.AddrInfo, len(splitPeers))
		for i, value := range splitPeers {

			multiCast := multiaddr.StringCast(value)
			currentAddrInfo, err := peer.AddrInfoFromP2pAddr(multiCast)
			if err != nil {
				fmt.Println(err)
				return
			}

			items[i] = *currentAddrInfo

		}
		rootKeyLink, err := sdk.ParseCidLink(*rootkey)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(rootKeyLink)

		go func() {
			impl.PushBlockWithExtData(ctx, exchange, &items[0], rootKeyLink, gsync.ExtensionData{}, impl.SelectAll)
		}()

		///	impl.PushBlockWithExtData(ctx, exchange, &items[0], rootKeyLink, gsync.ExtensionData{}, impl.SelectAll)
	}

	adapter := ethereum.NewOnchainAdapter("", "ropsten", 5)
	dagJsonHandler := handler.DagJsonHandler{
		AnconSyncContext: dagHandler,
		Proof:            proofHandler.GetProofService(),
		RootKey:          *rootkey,
		IPFSHost:         *ipfshost,
	}

	didHandler := handler.Did{
		AnconSyncContext: dagHandler,
		Proof:            proofHandler.GetProofService(),
		RootKey:          *rootkey,
		IPFSHost:         *ipfshost,
	}

	fileHandler := handler.FileHandler{
		RootKey:          *rootkey,
		AnconSyncContext: dagHandler,
	}
	g := handler.PlaygroundHandler(*dagHandler, adapter, proofHandler.GetProofAPI())

	// ticker := time.NewTicker(1500 * time.Millisecond)
	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case <-ticker.C:
	// 			block, hash, _ := proofHandler.Commit()
	// 			fmt.Printf("block at %d %s\r\n", block, hash)
	// 		}
	// 	}
	// }()

	// defer ticker.Stop()

	api := r.Group("/v0")
	{
		api.POST("/file", fileHandler.FileWrite)
		// api.POST("/code", fileHandler.UploadContract)
		api.GET("/graphql", g)
		api.POST("/graphql", g)
		api.GET("/file/:cid/*path", fileHandler.FileRead)
		api.GET("/dagjson/:cid/*path", dagJsonHandler.DagJsonRead)
		api.GET("/dag/:cid/*path", dagJsonHandler.DagJsonRead)
		api.POST("/dag", dagJsonHandler.DagJsonWrite)
		api.POST("/dagjson", dagJsonHandler.DagJsonWrite)
		api.PUT("/dag", dagJsonHandler.Update)
		api.PUT("/dagjson", dagJsonHandler.Update)
		// api.GET("/dagcbor/:cid/*path", dagCborHandler.DagCborRead)
		// api.POST("/dagcbor", dagCborHandler.DagCborWrite)
		api.POST("/did/key", didHandler.CreateDidKey)
		api.POST("/did/web", didHandler.CreateDidWeb)
		api.GET("/did/:did", didHandler.ReadDid)
		api.GET("/proof/:key", proofHandler.Read)
		api.GET("/proofs/lasthash", proofHandler.ReadCurrentRootHash)
		api.POST("/proofs/qr", proofHandler.ExtractQR)
	}
	// if subgraph.EnableDagcosmos {
	// 	ctx := context.WithValue(context.Background(), "dag", dagHandler)
	// 	indexer := cosmos.New(ctx, subgraph.CosmosPrimaryAddress, "/websocket")
	// 	r.GET("/indexer/cosmos/tip", indexer.TipEvent)
	// 	indexer.Subscribe(ctx, cosmos.NewBlock)
	// }
	r.GET("/user/:did/did.json", didHandler.ReadDidWebUrl)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// r.POST("/rpc", handler.EVMHandler(*dagHandler, proofHandler.GetProofAPI()))

	if *quic {
		http3.ListenAndServe(*apiAddr, *tlsCert, *tlsKey, r)
	} else {
		r.Run(*apiAddr)
	}

}
