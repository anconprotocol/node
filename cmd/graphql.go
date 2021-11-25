package cmd

import (
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/contract"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
)

func ReadFromStore(s Storage, hash string, path string) (string, error) {
	lnk, err := ParseCidLink(string(hash))
	if err != nil {
		return "", fmt.Errorf("parse anchor link error %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("read trusted anchor JSON  error %v", err)
	}

	node, err := s.Load(
		ipld.LinkContext{
			LinkPath: datamodel.ParsePath(path),
		},
		lnk,
	)
	if err != nil {
		return "", err
	}

	output, err := Encode(node)

	if err != nil {
		return "", err
	}

	return output, nil
}

func QueryGraphQL(s Storage) func(*gin.Context) {
	return func(c *gin.Context) {

		gqlschema := c.PostForm("schemacid")
		jsonPayload := c.PostForm("payloadcid")
		vars := c.PostForm("vars")

		// JSON Payload
		payload, err := ReadFromStore(s, jsonPayload, "")

		// GQL Schema
		schema, err := ReadFromStore(s, gqlschema, "")

		sch, err := contract.HandleSchema([]byte(schema))

		if err != nil {
			return nil, fmt.Errorf("validate payload error %v", err)
		}

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Error while executing data contract transaction. %v", err).Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"cid": cid,
		})
		return
	}
}
