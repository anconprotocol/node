package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/spf13/cast"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/cmd"
	"github.com/multiformats/go-multihash"
)

func GetLinkPrototype() ipld.LinkPrototype {
	// tip: 0x0129 dag-json
	return cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x71, // dag-cbor
		MhType:   0x12, // sha2-256
		MhLength: 32,   // sha2-256 hash has a 32-byte sum.
	}}
}

// CreateCidLink takes a hash eg ethereum hash and converts it to cid multihash
func CreateCidLink(hash []byte) cidlink.Link {
	lchMh, err := multihash.Encode(hash, GetLinkPrototype().(cidlink.LinkPrototype).MhType)
	if err != nil {
		return cidlink.Link{}
	}
	lcCID := cid.NewCidV1(GetLinkPrototype().(cidlink.LinkPrototype).Codec, lchMh)
	lcLinkCID := cidlink.Link{Cid: lcCID}
	return lcLinkCID
}

func main() {
	s := cmd.NewStorage(".ancon")
	r := gin.Default()
	r.POST("/file", func(c *gin.Context) {

		w, fn, err := s.DataStore.PutStream(c.Request.Context())
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Error while getting stream. %v", err).Error(),
			})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Error in form file %v", err).Error(),
			})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Cannot open file. %v", err).Error(),
			})
			return
		}
		defer src.Close()
		// var bz []byte

		_, err = io.Copy(w, src)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Failed reading file. %v", err).Error(),
			})
			return
		}

		var bz []byte
		bz, _ = json.Marshal(file.Header)

		lnk := CreateCidLink(bz)

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Cid error. %v", err).Error(),
			})
			return
		}
		fn(lnk.String())
		c.JSON(201, gin.H{
			"cid": lnk.String(),
		})
	})
	r.GET("/file/:cid", func(c *gin.Context) {
		lnk, err := cid.Parse(c.Param("cid"))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Cid error. %v", err).Error(),
			})
			return
		}
		reader, err := s.DataStore.GetStream(c.Request.Context(), lnk.String())

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Error while getting stream. %v", err).Error(),
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
	r.GET("/dagjson/:cid/*path", func(c *gin.Context) {
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
		data, err := cmd.Encode(n)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		c.PureJSON(201, data)
	})
	r.GET("/dagcbor/:cid/*path", func(c *gin.Context) {
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
		data, err := cmd.EncodeCBOR(n)
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
	})

	r.POST("/dagjson", func(c *gin.Context) {

		buff, _ := base64.StdEncoding.DecodeString(c.PostForm("data"))

		n, err := cmd.Decode(basicnode.Prototype.Any, string(buff))

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Decode Error %v", err).Error(),
			})
			return
		}
		cid := s.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(c.PostForm("path"))}, n)
		c.JSON(201, gin.H{
			"cid": cid,
		})
	})
	r.POST("/dagcbor", func(c *gin.Context) {

		buff, _ := base64.StdEncoding.DecodeString(c.PostForm("data"))

		n, err := cmd.DecodeCBOR(basicnode.Prototype.Any, buff)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("%v", err),
			})
			return
		}
		cid := s.Store(ipld.LinkContext{LinkPath: ipld.ParsePath(c.PostForm("path"))}, n)
		c.JSON(201, gin.H{
			"cid": cid,
		})
	})
	r.Run("0.0.0.0:7788") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
