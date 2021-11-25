package contract

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
)

type DAGContract struct {
	Schema   graphql.ExecutableSchema
	JsonData interface{}
}

func (d *DAGContract) ApplyExec(op, query string) *graphql.Response {
	exec := executor.New(d.Schema)
	ctx := context.Background()
	rc, err := exec.CreateOperationContext(ctx, &graphql.RawParams{
		Query:         query,
		OperationName: op,
	})
	if err != nil {
		return exec.DispatchError(ctx, err)
	}

	resp, ctx2 := exec.DispatchOperation(ctx, rc)

	return resp(ctx2)
}
