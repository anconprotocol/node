package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
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

func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v\n", result.Errors)
	}
	return result
}

func QueryGraphQL(s Storage) func(*gin.Context) {
	return func(c *gin.Context) {

		// gqlschema := c.PostForm("schemacid")
		jsonPayload := c.PostForm("payloadcid")
		query := c.PostForm("query")

		// JSON Payload
		payload, err := ReadFromStore(s, jsonPayload, "")

		if err != nil {
			return nil, fmt.Errorf("validate payload error %v", err)
		}

		res := ExecuteQuery("", graphql.Schema{})

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

func FilterUser(data []map[string]interface{}, args map[string]interface{}) map[string]interface{} {
	for _, user := range data {
		for k, v := range args {
			if user[k] != v {
				goto nextuser
			}
			return user
		}

	nextuser:
	}
	return nil
}

func ImportJSONData(content []byte) (*graphql.Schema, error) {

	var data []map[string]interface{}

	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)
	for _, item := range data {
		for k := range item {
			fields[k] = &graphql.Field{
				Type: graphql.String,
			}
			args[k] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
		}
	}

	var userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "User",
			Fields: fields,
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: userType,
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return FilterUser(data, p.Args), nil
					},
				},
			},
		})

	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)

	return &schema, err
}
