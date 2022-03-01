package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

type LoadersKey string

const loadersKey LoadersKey = "dataloaders"

const waitTime = 100 * time.Millisecond
const maxBatch = 50

type Loaders struct {
	TeamByAbr          TeamLoaderABR
	TeamByID           TeamLoaderID
	PlayerByID         PlayerLoaderID
	PlayerByFilter     PlayerLoader
	PlayerGameByFilter PlayerGameLoader
}

func Middleware(conn *database.NBADatabaseClient, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			TeamByAbr: *NewTeamLoaderABR(
				TeamLoaderABRConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []string) ([]*model.Team, []error) {
						teams := make([]*model.Team, len(keys))
						teamsByAbr := make(map[string]*model.Team, len(keys))
						errs := make([]error, len(keys))
						cur, err := conn.GetTeamsByAbr(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							team := &model.Team{}
							err := cur.Decode(&team)
							if err != nil {
								return nil, []error{err}
							}
							teamsByAbr[team.Abbreviation] = team
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, team := range keys {
							teams[i] = teamsByAbr[team]
						}
						return teams, errs
					},
				},
			),
			TeamByID: *NewTeamLoaderID(
				TeamLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []int) ([]*model.Team, []error) {
						teams := make([]*model.Team, len(keys))
						teamsById := make(map[int]*model.Team, len(keys))
						errs := make([]error, len(keys))
						cur, err := conn.GetTeamsById(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							team := &model.Team{}
							err := cur.Decode(&team)
							if err != nil {
								return nil, []error{err}
							}
							teamsById[team.TeamID] = team
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, team := range keys {
							teams[i] = teamsById[team]
						}
						return teams, errs
					},
				},
			),
			//TODO: Replace this with PlayerByFilter
			PlayerByID: *NewPlayerLoaderID(
				PlayerLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []int) ([]*model.Player, []error) {
						players := make([]*model.Player, len(keys))
						playersById := make(map[int]*model.Player, len(keys))
						errs := make([]error, len(keys))
						filters := make([]model.PlayerFilter, len(keys))
						for i := range keys {
							filters = append(filters, model.PlayerFilter{PlayerID: &keys[i]})
						}
						cur, err := conn.GetPlayers(r.Context(), filters)
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							player := &model.Player{}
							err := cur.Decode(&player)
							if err != nil {
								return nil, []error{err}
							}
							playersById[player.PlayerID] = player
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, player := range keys {
							players[i] = playersById[player]
						}
						return players, errs
					},
				},
			),
			//TODO: fix firstName to be full name
			PlayerByFilter: *NewPlayerLoader(
				PlayerLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []model.PlayerFilter) ([]*model.Player, []error) {
						players := make([]*model.Player, len(keys))
						playersByName := make(map[string]*model.Player, len(keys))
						errs := make([]error, len(keys))
						cur, err := conn.GetPlayers(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							player := &model.Player{}
							err := cur.Decode(&player)
							if err != nil {
								return nil, []error{err}
							}
							name := player.FirstName + " " + player.LastName
							playersByName[name] = player
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, player := range keys {
							players[i] = playersByName[*player.Name]
						}
						return players, errs
					},
				},
			),
			PlayerGameByFilter: *NewPlayerGameLoader(
				PlayerGameLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []model.GameFilter) ([][]*model.PlayerGame, []error) {
						games := make([][]*model.PlayerGame, len(keys))
						errs := make([]error, len(keys))
						gamesByPlayerID := make(map[int][]*model.PlayerGame, len(keys))
						cur, err := conn.GetPlayerGames(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							game := &model.PlayerGame{}
							err := cur.Decode(&game)
							if err != nil {
								return nil, []error{err}
							}
							gamesByPlayerID[game.PlayerID] = append(gamesByPlayerID[game.PlayerID], game)
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, filter := range keys {
							games[i] = gamesByPlayerID[*filter.PlayerID]
						}
						return games, errs
					},
				},
			),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
