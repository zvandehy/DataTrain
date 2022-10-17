package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// Propositions is the resolver for the propositions field.
func (r *projectionResolver) Propositions(ctx context.Context, obj *model.Projection, input *model.PropositionFilter) ([]*model.Proposition, error) {
	propositions := []*model.Proposition{}
	for _, p := range obj.Props {
		if p.Match(*input) {
			propositions = append(propositions, p)
		}
	}
	return propositions, nil
}

// Result is the resolver for the result field.
func (r *projectionResolver) Result(ctx context.Context, obj *model.Projection) (*model.PlayerGame, error) {
	panic(fmt.Errorf("(r *projectionResolver) Result not implemented"))
}

// Projection returns generated.ProjectionResolver implementation.
func (r *Resolver) Projection() generated.ProjectionResolver { return &projectionResolver{r} }

type projectionResolver struct{ *Resolver }
