package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// Players is the resolver for the players field.
func (r *queryResolver) Players(ctx context.Context, input *model.PlayerFilter) ([]*model.Player, error) {
	startTime := time.Now()
	if input == nil {
		input = &model.PlayerFilter{}
	}
	players, err := r.GetPlayers(ctx, true, input)
	if err != nil {
		logrus.Error(err)
		return []*model.Player{}, err
	}
	logrus.Info("Players query took ", time.Since(startTime))
	return players, nil
}

// Teams is the resolver for the teams field.
func (r *queryResolver) Teams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
	panic(fmt.Errorf("(r *queryResolver) Teams not implemented"))
}

// Games is the resolver for the games field.
func (r *queryResolver) Games(ctx context.Context, input model.GameFilter) ([]*model.PlayerGame, error) {
	games, err := r.GetPlayerGames(ctx, &input)
	if err != nil || len(games) == 0 {
		logrus.Error(err)
		return []*model.PlayerGame{}, err
	}
	return games, nil
}

// SeasonMatch is the resolver for the seasonMatch field.
func (r *gameFilterResolver) SeasonMatch(ctx context.Context, obj *model.GameFilter, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// PreviousSeasonMatch is the resolver for the previousSeasonMatch field.
func (r *gameFilterResolver) PreviousSeasonMatch(ctx context.Context, obj *model.GameFilter, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// MatchSeason is the resolver for the matchSeason field.
func (r *periodResolver) MatchSeason(ctx context.Context, obj *model.Period, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// MatchPreviousSeason is the resolver for the matchPreviousSeason field.
func (r *periodResolver) MatchPreviousSeason(ctx context.Context, obj *model.Period, data *bool) error {
	panic(fmt.Errorf("not implemented"))
}

// WithPropositions is the resolver for the withPropositions field.
func (r *playerFilterResolver) WithPropositions(ctx context.Context, obj *model.PlayerFilter, data *model.PropositionFilter) error {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// GameFilter returns generated.GameFilterResolver implementation.
func (r *Resolver) GameFilter() generated.GameFilterResolver { return &gameFilterResolver{r} }

// Period returns generated.PeriodResolver implementation.
func (r *Resolver) Period() generated.PeriodResolver { return &periodResolver{r} }

// PlayerFilter returns generated.PlayerFilterResolver implementation.
func (r *Resolver) PlayerFilter() generated.PlayerFilterResolver { return &playerFilterResolver{r} }

type queryResolver struct{ *Resolver }
type gameFilterResolver struct{ *Resolver }
type periodResolver struct{ *Resolver }
type playerFilterResolver struct{ *Resolver }
