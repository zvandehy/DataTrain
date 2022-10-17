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
func (r *teamResolver) Players(ctx context.Context, obj *model.Team) ([]*model.Player, error) {
	panic(fmt.Errorf("(r *teamResolver) Players not implemented"))
}

// Opponent is the resolver for the opponent field.
func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
	panic(fmt.Errorf("(r *teamResolver) Players not implemented"))
}

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// TeamGame returns generated.TeamGameResolver implementation.
func (r *Resolver) TeamGame() generated.TeamGameResolver { return &teamGameResolver{r} }

type teamResolver struct{ *Resolver }
type teamGameResolver struct{ *Resolver }
