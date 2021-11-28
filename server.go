package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cast"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/handler"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/impl"
)

func main() {
	peerAddr := flag.String("peeraddr", "/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWAGyXSBPPo7Zq16WoCe6BtHDRQpFXPg9VCDQ1EPXcHWMw", "A remote peer to sync")
	addr := flag.String("addr", "/ip4/0.0.0.0/tcp/7702", "Host multiaddr")
	apiAddr := flag.String("apiaddr", "0.0.0.0:7788", "API address")
	dataFolder := flag.String("data", ".ancon", "Data directory")

	flag.Parse()

	s := anconsync.NewStorage(*dataFolder)
	ctx := context.Background()
	host := net.NewPeer(ctx, *addr)
	// peerhost := "/ip4/192.168.50.138/tcp/7702/p2p/12D3KooWA7vRHFLC8buiEP2xYwUN5kdCgzEtQRozMtnCPDi4n4HM"
	// "/ip4/190.34.226.207/tcp/29557/p2p/12D3KooWGd9mLtWx7WGEd9mnWPbCsr1tFCxtEi7RkgsJYxAZmZgi"

	exchange, ipfspeer := impl.NewRouter(ctx, host, s, *peerAddr)
	fmt.Println(ipfspeer.ID)
	r := gin.Default()
	r.POST("/file", func(c *gin.Context) {
		w, fn, err := s.DataStore.PutStream(c.Request.Context())
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error while getting stream. %v", err).Error(),
			})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error in form file %v", err).Error(),
			})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("cannot open file. %v", err).Error(),
			})
			return
		}
		defer src.Close()
		// var bz []byte

		_, err = io.Copy(w, src)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("failed reading file. %v", err).Error(),
			})
			return
		}

		var bz []byte
		bz, _ = json.Marshal(file.Header)

		lnk := anconsync.CreateCidLink(bz)

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("cid error. %v", err).Error(),
			})
			return
		}
		fn(lnk.String())
		c.JSON(201, gin.H{
			"cid": lnk.String(),
		})
	})
	r.POST("/dagcontract", handler.QueryGraphQL(s))
	r.POST("/graphqli", handler.QueryGraphQL(s))
	r.GET("/file/:cid", func(c *gin.Context) {
		lnk, err := cid.Parse(c.Param("cid"))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("cid error. %v", err).Error(),
			})
			return
		}
		reader, err := s.DataStore.GetStream(c.Request.Context(), lnk.String())

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error while getting stream. %v", err).Error(),
			})
			return
		}

		contentLength := cast.ToInt64(-1)
		contentType := c.ContentType()

		extraHeaders := map[string]string{
			//  "Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	r.GET("/dagjson/:cid/*path", handler.DagJsonRead(s, exchange, ipfspeer))
	r.GET("/dagcbor/:cid/*path", handler.DagCborRead(s))
	r.POST("/dagjson", handler.DagJsonWrite(s, exchange, ipfspeer))
	r.POST("/dagcbor", handler.DagCborWrite(s, exchange, ipfspeer))
	r.Run(*apiAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
