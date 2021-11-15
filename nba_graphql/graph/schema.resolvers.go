package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zvandehy/nba_graphql/graph/generated"
	"github.com/zvandehy/nba_graphql/graph/model"
)

func (r *queryResolver) Games(ctx context.Context, id *string) ([]*model.PlayerGame, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) LastNGames(ctx context.Context, id *string, n *int) ([]*model.PlayerGame, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TeamGames(ctx context.Context, id *string) ([]*model.TeamGame, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	return []*model.Player{{ID: "1", FirstName: "James", LastName: "LeBron"}}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
