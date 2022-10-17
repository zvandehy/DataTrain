package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// Players is the resolver for the players field.
func (r *queryResolver) Players(ctx context.Context, input model.PlayerFilter) ([]*model.Player, error) {
	players, err := r.Db.GetPlayers(ctx, &input)
	if err != nil {
		return []*model.Player{}, err
	}
	return players, nil
}

// Teams is the resolver for the teams field.
func (r *queryResolver) Teams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
	panic(fmt.Errorf("(r *queryResolver) Teams not implemented"))
}

// Projections is the resolver for the projections field.
func (r *queryResolver) Projections(ctx context.Context, input model.ProjectionFilter) ([]*model.Projection, error) {
	projections, err := r.Db.GetProjections(ctx, input)
	if err != nil {
		logrus.Errorf("Error getting projections: %v", err)
		return []*model.Projection{}, err
	}
	return projections, nil
}

// PositionStrictMatch is the resolver for the positionStrictMatch field.
func (r *playerFilterResolver) PositionStrictMatch(ctx context.Context, obj *model.PlayerFilter, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// PositionLooseMatch is the resolver for the positionLooseMatch field.
func (r *playerFilterResolver) PositionLooseMatch(ctx context.Context, obj *model.PlayerFilter, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// PlayerFilter returns generated.PlayerFilterResolver implementation.
func (r *Resolver) PlayerFilter() generated.PlayerFilterResolver { return &playerFilterResolver{r} }

type queryResolver struct{ *Resolver }
type playerFilterResolver struct{ *Resolver }
