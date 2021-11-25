package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jensneuse/graphql-go-tools/pkg/engine/plan"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/jensneuse/abstractlogger"
	"github.com/jensneuse/graphql-go-tools/pkg/graphql"
)

var schema graphql.Schema

const jsonDataFile = "data.json"

func ReadFromStore(s Storage, hash string, path string) (string, error) {
	lnk, err := ParseCidLink(string(hash))
	if err != nil {
		return "", fmt.Errorf("parse anchor link error %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("read trusted anchor JSON  error %v", err)
	}
	//CURL fields:
	//grabar un schema cid
	//save with prop called schema in plain text
	//json base 64
	//variable json asplaintext
	//op (operations) string
	//resolve path
	//query
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
		// if c.PostForm("path") == "" {
		// 	c.JSON(400, gin.H{
		// 		"error": fmt.Errorf("Missing path").Error(),
		// 	})
		// 	return
		// }
		// if c.PostForm("op") == "" {
		// 	c.JSON(400, gin.H{
		// 		"error": fmt.Errorf("Missing operation").Error(),
		// 	})
		// 	return
		// }
		if c.PostForm("query") == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Missing query").Error(),
			})
			return
		}
		if c.PostForm("schemacid") == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Missing schema cid").Error(),
			})
			return
		}
		if c.PostForm("payloadcid") == "" {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Missing payload data source cid").Error(),
			})
			return
		}
		// if c.PostForm("variables") == "" {
		// 	c.JSON(400, gin.H{
		// 		"error": fmt.Errorf("Missing variables").Error(),
		// 	})
		// 	return
		// }

		gqlschema := c.PostForm("schemacid")
		jsonPayload := c.PostForm("payloadcid")
		variables := c.PostForm("variables")
		query := c.PostForm("query")
		op := c.PostForm("op")

		// JSON Payload
		payload, err := ReadFromStore(s, jsonPayload, "")

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("No JSON payload found. %v", err).Error(),
			})
			return
		}
		// GraphQL Schema
		schemaGQL, err := ReadFromStore(s, gqlschema, "")
		var v map[string]string
		json.Unmarshal([]byte(schemaGQL), &v)

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("No GraphQL Schema found. %v", err).Error(),
			})
			return
		}
		schema, err := NewSchemaFrom([]byte(v["schema"]))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Schema generation failed %v", err).Error(),
			})
			return
		}

		engineConf := graphql.NewEngineV2Configuration(schema)

		engineConf.AddDataSource(plan.DataSourceConfiguration{
			Custom: []byte(payload),
		})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		engine, err := graphql.NewExecutionEngineV2(ctx, abstractlogger.Noop{}, engineConf)

		operation := graphql.Request{
			OperationName: op,
			Variables:     []byte(variables),
			Query:         query,
		}

		resultWriter := graphql.NewEngineResultWriter()
		execCtx, execCtxCancel := context.WithCancel(context.Background())
		defer execCtxCancel()
		err = engine.Execute(execCtx, &operation, &resultWriter)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Error while executing data contract transaction. %v", err).Error(),
			})
			return
		}

		// buff, _ := base64.StdEncoding.DecodeString(resultWriter.String())
		n, err := Decode(basicnode.Prototype.Any, resultWriter.String())
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
		return
	}
}

func NewSchemaFrom(schemaBytes []byte) (*graphql.Schema, error) {

	schemaReader := bytes.NewBuffer(schemaBytes)
	schema, err := graphql.NewSchemaFromReader(schemaReader)

	return schema, err
}
