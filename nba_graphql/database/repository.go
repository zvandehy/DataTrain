package database

import (
	"context"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

type BasketballRepository interface {
	GetLeague() string
	GetPlayers(ctx context.Context, withGames bool, playerFilter *model.PlayerFilter) ([]*model.Player, error)
	GetPropositions(ctx context.Context, propositionFilter *model.PropositionFilter) ([]*model.Proposition, error)
	// GetTeamsByID(ctx context.Context, teamIDs []int) ([]*model.Team, error)
	// GetTeamsByAbr(ctx context.Context, teamAbrs []string) ([]*model.Team, error)
	GetTeams(ctx context.Context, withGames bool, teamFilter *[]*model.TeamFilter) ([]*model.Team, error)
	GetSimilarPlayers(ctx context.Context, similarToPlayerID int, input *model.SimilarPlayerInput, endDate string) ([]*model.Player, error)
	GetSimilarTeams(ctx context.Context, similarToTeamID int, input *model.SimilarTeamInput, endDate string) ([]*model.Team, error)
	GetPropositionsByPlayerGame(ctx context.Context, game model.PlayerGame) ([]*model.Proposition, error)
}
