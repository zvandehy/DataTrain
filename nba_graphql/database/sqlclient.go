package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

type SQLClient struct {
	League  string
	Queries int
	Client  *sql.DB
	// PlayerSimilarity model.PlayerSnapshots
	// TeamSimilarity   model.TeamSnapshots
	// PlayerCache      map[string][]*model.Player
	// TeamCache        map[string][]*model.Team
}

func NewSQLClient(league string) (*SQLClient, error) {
	db, err := GetDatabase(league)
	if err != nil {
		return nil, err
	}
	return &SQLClient{
		League:  league,
		Client:  db,
		Queries: 0,
	}, nil
}

func GetDatabase(league string) (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", fmt.Sprintf("%s/%s?tls=true", dsn, strings.ToLower(league)))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, err
}

func (c *SQLClient) GetLeague() string {
	return c.League
}

func (c *SQLClient) GetPlayers(ctx context.Context, withGames bool, playerFilter *model.PlayerFilter) ([]*model.Player, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetPropositions(ctx context.Context, propositionFilter *model.PropositionFilter) ([]*model.Proposition, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetTeams(ctx context.Context, withGames bool, teamFilter *[]*model.TeamFilter) ([]*model.Team, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetSimilarPlayers(ctx context.Context, similarToPlayerID int, input *model.SimilarPlayerInput, endDate string) ([]*model.Player, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetSimilarTeams(ctx context.Context, similarToTeamID int, input *model.SimilarTeamInput, endDate string) ([]*model.Team, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetPropositionsByPlayerGame(ctx context.Context, game model.PlayerGame) ([]*model.Proposition, error) {
	panic("not implemented") // TODO: Implement
}
