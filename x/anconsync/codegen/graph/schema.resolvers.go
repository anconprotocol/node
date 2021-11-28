package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/codegen/graph/generated"
	"github.com/Electronic-Signatures-Industries/ancon-ipld-router-sync/x/anconsync/codegen/graph/model"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/traversal"
)

func (r *mutationResolver) Apply(ctx context.Context, tx model.MetadataTransactionInput) (*model.DagLink, error) {
	dag := ctx.Value("dag").(*anconsync.DagContractTrustedContext)

	lnk, err := anconsync.ParseCidLink(tx.Cid)
	if err != nil {
		return nil, err
	}
	rootNode, err := dag.Store.Load(ipld.LinkContext{}, lnk)
	if err != nil {
		return nil, err
	}

	// owner update
	n, err := traversal.FocusedTransform(
		rootNode,
		datamodel.ParsePath("owner"),
		func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			if progress.Path.String() == "owner" && must.String(prev) == tx.Owner {
				nb := prev.Prototype().NewBuilder()
				nb.AssignString(tx.NewOwner)
				return nb.Build(), nil
			}
			return nil, fmt.Errorf("Owner not found")
		}, false)

	if err != nil {
		return nil, err
	}

	// parent update
	n, _ = traversal.FocusedTransform(
		n,
		datamodel.ParsePath("parent"),
		func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			nb := prev.Prototype().NewBuilder()
			// set previous hash, not current
			l, _ := anconsync.ParseCidLink(tx.Cid)
			nb.AssignLink(l)
			return nb.Build(), nil
		}, false)

	link := dag.Store.Store(ipld.LinkContext{}, n)

	return &model.DagLink{
		Path: "/",
		Cid:  link.String(),
	}, nil
}

func (r *queryResolver) Metadata(ctx context.Context, cid string, path string) (*model.Ancon721Metadata, error) {
	dag := ctx.Value("dag").(*anconsync.DagContractTrustedContext)

	jsonmodel, err := anconsync.ReadFromStore(dag.Store, cid, path)
	if err != nil {
		return nil, err
	}
	var metadata model.Ancon721Metadata
	json.Unmarshal([]byte(jsonmodel), &metadata)
	return &metadata, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
