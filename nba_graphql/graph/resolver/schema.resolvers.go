package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// Players is the resolver for the players field.
func (r *queryResolver) Players(ctx context.Context, input *model.PlayerFilter) ([]*model.Player, error) {
	if input == nil {
		input = &model.PlayerFilter{}
	}
	players, err := r.GetPlayers(ctx, true, input)
	if err != nil {
		return []*model.Player{}, err
	}
	return players, nil
}

// Teams is the resolver for the teams field.
func (r *queryResolver) Teams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
	panic(fmt.Errorf("(r *queryResolver) Teams not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
