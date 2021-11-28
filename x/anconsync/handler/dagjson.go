package handler

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-graphsync"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

func DagJsonWrite(s anconsync.Storage, exchange graphsync.GraphExchange, pi *peer.AddrInfo) func(*gin.Context) {
	return func(c *gin.Context) {

		v, _ := c.GetRawData()

		path, _ := jsonparser.GetString(v, "path")

		if path == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing path").Error(),
			})
			return
		}

		temp, _ := jsonparser.GetUnsafeString(v, "data")
		//		temp = strings.ReplaceAll(temp, "\n", "")
		// temp = strings.ReplaceAll(temp, "\\","\"")
		var buf bytes.Buffer
		err := json.Compact(&buf, []byte(temp))
		data := buf.Bytes()
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing payload data source").Error(),
			})
			return
		}

		n, err := anconsync.Decode(basicnode.Prototype.Any, string(data))

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("decode Error %v", err).Error(),
			})
			return
		}
		cid := s.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(path)}, n)
		c.JSON(201, gin.H{
			"cid": cid,
		})
		net.PushBlock(c.Request.Context(), exchange, pi.ID, cid)
	}
}
func DagJsonRead(s anconsync.Storage, exchange graphsync.GraphExchange, pi *peer.AddrInfo) func(*gin.Context) {
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
		data, err := anconsync.Encode(n)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		c.JSON(200, data)
	}
}
