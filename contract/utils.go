package contract

import (
	"encoding/json"

	"github.com/vektah/gqlparser/ast"
)

func ToDeployment(schema ast.Schema) ([]byte, error) {
	bz, err := json.Marshal(schema)

	if err != nil {
		return nil, err
	}

	return bz, nil
}

func HandleSchema(schema []byte) (*ast.Schema, error) {
	var schemaObj ast.Schema
	err := json.Unmarshal(schema, &schemaObj)

	return &schemaObj, err
}
