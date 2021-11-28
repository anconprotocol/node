package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/net"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-graphsync"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/spf13/cast"
)

func DagCborWrite(s anconsync.Storage, exchange graphsync.GraphExchange, pi *peer.AddrInfo) func(*gin.Context) {
	return func(c *gin.Context) {

		var v map[string]string

		c.BindJSON(&v)
		if v["path"] == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing path").Error(),
			})
			return
		}
		if v["data"] == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("missing payload data source").Error(),
			})
			return
		}

		buff, _ := base64.StdEncoding.DecodeString(v["data"])

		n, err := anconsync.DecodeCBOR(basicnode.Prototype.Any, buff)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		cid := s.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(v["path"])}, n)
		c.JSON(201, gin.H{
			"cid": cid,
		})
		net.PushBlock(c.Request.Context(), exchange, pi.ID, cid)
	}
}

func DagCborRead(s anconsync.Storage) func(*gin.Context) {
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
		data, err := anconsync.EncodeCBOR(n)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		reader := bytes.NewReader(data)
		contentLength := cast.ToInt64(reader.Len())
		contentType := "application/cbor"

		extraHeaders := map[string]string{
			//  "Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	}
}
