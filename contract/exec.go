package contract

import (
	"encoding/json"

	ast "github.com/vektah/gqlparser/ast/v2"
)

func main() {}

//export New
func New(data []byte, schema *ast.Schema) (*DAGContract, error) {

	// Require
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		return nil, err
	}
	//	resolved := introspection.WrapSchema(schema)
	return &DAGContract{Schema: schema, JsonData: jsonData}, nil
}
