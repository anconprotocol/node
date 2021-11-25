package contract

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
)

func ToDeployment(schema graphql.Schema) ([]byte, error) {
	bz, err := json.Marshal(schema)

	if err != nil {
		return nil, err
	}

	return bz, nil
}

func HandleSchema(schema []byte) (*graphql.Schema, error) {
	var schemaObj graphql.Schema
	err := json.Unmarshal(schema, &schemaObj)

	return &schemaObj, err
}
