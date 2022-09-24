package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
	"go.mongodb.org/mongo-driver/bson"
)

type LoadersKey string

const loadersKey LoadersKey = "dataloaders"

const waitTime = 150 * time.Millisecond
const maxBatch = 50

type Loaders struct {
	TeamByAbr                TeamLoaderABR
	TeamByID                 TeamLoaderID
	PlayerByID               PlayerLoaderID
	PlayerByFilter           PlayerLoader
	PlayerGameByFilter       PlayerGameLoader
	TeamGameByPlayerGame     TeamGameLoader
	OpponentGameByPlayerGame TeamGameLoader
	SimilarPlayerLoader      SimilarPlayerLoader
	SimilarTeamLoader        SimilarTeamLoader
	PlayerInjuryLoader       PlayerInjuryLoader
	TeamInjuryLoader         TeamInjuryLoader
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
							if teamsByAbr[team] == nil {
								errs[i] = fmt.Errorf("failed to get team: %v", team)
								logrus.Errorf("failed to get team: %v", team)
							} else {
								teams[i] = teamsByAbr[team]
							}
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
							if teamsById[team] == nil {
								errs[i] = fmt.Errorf("failed to get team: %v", team)
							} else {
								teams[i] = teamsById[team]
							}
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
							playersByName[player.Name] = player
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
			PlayerInjuryLoader: *NewPlayerInjuryLoader(
				PlayerInjuryLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []int) ([][]*model.Injury, []error) {
						playerInjuries := make([][]*model.Injury, len(keys))
						injuriesByPlayerID := make(map[int][]*model.Injury, len(keys))
						errs := make([]error, len(keys))
						injuries, err := conn.GetPlayerInjuries(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}

						for _, injury := range injuries {
							injuriesByPlayerID[injury.PlayerID] = append(injuriesByPlayerID[injury.PlayerID], injury)
						}
						for i, playerID := range keys {
							playerInjuries[i] = injuriesByPlayerID[playerID]
						}
						return playerInjuries, errs
					},
				},
			),
			TeamInjuryLoader: *NewTeamInjuryLoader(
				TeamInjuryLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []int) ([][]*model.Injury, []error) {
						teamInjuries := make([][]*model.Injury, len(keys))
						injuriesByTeamID := make(map[int][]*model.Injury, len(keys))
						errs := make([]error, len(keys))
						injuries, err := conn.GetTeamInjuries(r.Context(), keys)
						if err != nil {
							return nil, []error{err}
						}

						for _, injury := range injuries {
							injuriesByTeamID[injury.Team.TeamID] = append(injuriesByTeamID[injury.Team.TeamID], injury)
						}
						for i, teamID := range keys {
							teamInjuries[i] = injuriesByTeamID[teamID]
						}
						return teamInjuries, errs
					},
				},
			),
			TeamGameByPlayerGame: *NewTeamGameLoader(
				TeamGameLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []model.PlayerGame) ([]*model.TeamGame, []error) {
						opponentIDs := make([]int, len(keys))
						gameIDs := make([]string, len(keys))
						for i, key := range keys {
							opponentIDs[i] = key.OpponentID
							gameIDs[i] = key.GameID
						}
						teamGames := make([]*model.TeamGame, len(keys))
						teamGamesByPlayerGame := make(map[string]*model.TeamGame, len(keys))
						errs := make([]error, len(keys))
						cur, err := conn.Collection("teamgames").Find(r.Context(), bson.M{"gameID": bson.M{"$in": gameIDs}, "opponent": bson.M{"$in": opponentIDs}})
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							game := &model.TeamGame{}
							err := cur.Decode(&game)
							if err != nil {
								return nil, []error{err}
							}
							mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
							teamGamesByPlayerGame[mapKey] = game
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, game := range keys {
							mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
							teamGames[i] = teamGamesByPlayerGame[mapKey]
						}
						return teamGames, errs
					},
				},
			),
			OpponentGameByPlayerGame: *NewTeamGameLoader(
				TeamGameLoaderConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch: func(keys []model.PlayerGame) ([]*model.TeamGame, []error) {
						teamIDs := make([]int, len(keys))
						gameIDs := make([]string, len(keys))
						for i, key := range keys {
							teamIDs[i] = key.TeamID
							gameIDs[i] = key.GameID
						}
						teamGames := make([]*model.TeamGame, len(keys))
						teamGamesByPlayerGame := make(map[string]*model.TeamGame, len(keys))
						errs := make([]error, len(keys))
						cur, err := conn.Collection("teamgames").Find(r.Context(), bson.M{"gameID": bson.M{"$in": gameIDs}, "opponent": bson.M{"$in": teamIDs}})
						if err != nil {
							return nil, []error{err}
						}
						defer cur.Close(r.Context())

						for cur.Next(r.Context()) {
							game := &model.TeamGame{}
							err := cur.Decode(&game)
							if err != nil {
								return nil, []error{err}
							}
							mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
							teamGamesByPlayerGame[mapKey] = game
						}
						if err := cur.Err(); err != nil {
							return nil, []error{err}
						}
						for i, game := range keys {
							mapKey := fmt.Sprintf("%s:%d", game.GameID, game.OpponentID)
							teamGames[i] = teamGamesByPlayerGame[mapKey]
						}
						return teamGames, errs
					},
				},
			),
			PlayerGameByFilter: *NewPlayerGameLoader(
				PlayerGameLoaderConfig{
					MaxBatch: maxBatch * 10,
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
			SimilarPlayerLoader: *NewSimilarPlayerLoader(
				SimilarPlayerLoaderConfig{
					MaxBatch: maxBatch * 2,
					Wait:     waitTime,
					Fetch: func(keys []model.GameFilter) ([][]*model.Player, []error) {
						similarPlayers := make([][]*model.Player, len(keys))
						errs := make([]error, len(keys))
						createFilterWithoutIDs := make([]model.GameFilter, len(keys))
						for i, filter := range keys {
							createFilterWithoutIDs[i] = model.GameFilter{Season: filter.Season} //TODO: add other filters if applicable
						}
						playerAverages, err := conn.GetAverages(r.Context(), createFilterWithoutIDs)
						if err != nil || len(*playerAverages) == 0 {
							return nil, []error{err}
						}
						for i, filter := range keys {
							targetPlayer := (*playerAverages)[0]
							for _, p := range *playerAverages {
								if p.Player.PlayerID == *filter.PlayerID {
									targetPlayer = p
								}
							}
							similarPlayers[i] = util.SimilarPlayers(*playerAverages, targetPlayer)
						}
						return similarPlayers, errs
					},
				},
			),
			SimilarTeamLoader: *NewSimilarTeamLoader(
				SimilarTeamLoaderConfig{
					MaxBatch: maxBatch * 2,
					Wait:     waitTime,
					Fetch: func(keys []model.GameFilter) ([][]*model.Team, []error) {
						similarTeams := make([][]*model.Team, len(keys))
						errs := make([]error, len(keys))
						createFilterWithoutIDs := make([]model.GameFilter, len(keys))
						for i, filter := range keys {
							createFilterWithoutIDs[i] = model.GameFilter{Season: filter.Season} //TODO: add other filters if applicable
						}
						teamAverages, err := conn.GetTeamAverages(r.Context(), createFilterWithoutIDs)
						if err != nil || len(*teamAverages) == 0 {
							return nil, []error{err}
						}
						for i, filter := range keys {
							targetTeam := (*teamAverages)[0]
							for _, t := range *teamAverages {
								if t.Team.TeamID == *filter.TeamID {
									targetTeam = t
								}
							}
							similarTeams[i] = util.SimilarTeams(*teamAverages, targetTeam)
						}
						return similarTeams, errs
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
