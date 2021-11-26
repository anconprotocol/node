package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-graphsync"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

func DagJsonRead(s Storage, exchange graphsync.GraphExchange, pi *peer.AddrInfo) func(*gin.Context) {
	return func(c *gin.Context) {
		lnk, err := cid.Parse(c.Param("cid"))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		n, err := s.Load(ipld.LinkContext{LinkPath: ipld.ParsePath(c.Param("path"))}, cidlink.Link{Cid: lnk})

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		data, err := Encode(n)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		c.PureJSON(201, data)
	}
}
