package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

type LoadersKey string

const loadersKey LoadersKey = "dataloaders"

const waitTime = 150 * time.Millisecond
const maxBatch = 50

type Loaders struct {
	PlayerByID PlayerLoaderID
	TeamByID   TeamLoaderID
	TeamByAbr  TeamLoaderABR
	// PlayerByFilter           PlayerLoader
	// PlayerGameByFilter       PlayerGameLoader
	// TeamGameByPlayerGame     TeamGameLoader
	// OpponentGameByPlayerGame TeamGameLoader
}

func Middleware(conn database.BasketballRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			PlayerByID: *NewPlayerLoaderID(
				PlayerLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch:    fetchPlayerByID(conn, r.Context()),
				},
			),
			TeamByID: *NewTeamLoaderID(
				TeamLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch:    fetchTeamByID(conn, r.Context()),
				},
			),
			TeamByAbr: *NewTeamLoaderABR(
				TeamLoaderABRConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch:    fetchTeamByAbbreviation(conn, r.Context()),
				},
			),
			// TeamGameByPlayerGame: *NewTeamGameLoader(
			// 	TeamGameLoaderConfig{
			// 		MaxBatch: maxBatch,
			// 		Wait:     waitTime,
			// 		Fetch:    fetchTeamGameByPlayerGame(conn, r.Context()),
			// 	},
			// ),
			// OpponentGameByPlayerGame: *NewTeamGameLoader(
			// 	TeamGameLoaderConfig{
			// 		MaxBatch: maxBatch,
			// 		Wait:     waitTime,
			// 		Fetch:    fetchOpponentGameByPlayerGame,
			// 	},
			// ),
			// PlayerGameByFilter: *NewPlayerGameLoader(
			// 	PlayerGameLoaderConfig{
			// 		MaxBatch: maxBatch * 10,
			// 		Wait:     waitTime,
			// 		Fetch:    fetchPlayerGamesByFilter,
			// 	},
			// ),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func wrapPlayerNames(names []string) []*model.PlayerFilter {
	wrapped := make([]*model.PlayerFilter, len(names))
	for i := range names {
		wrapped[i] = &model.PlayerFilter{Name: &(names[i])}
	}
	return wrapped
}

func fetchPlayerByName(db database.BasketballRepository, ctx context.Context) func(playerNames []string) ([]*model.Player, []error) {
	return func(playerNames []string) ([]*model.Player, []error) {
		nameToPlayerMap := make(map[string]*model.Player, len(playerNames))
		players, err := db.GetPlayers(ctx, false, wrapPlayerNames(playerNames)...)
		if err != nil {
			return nil, []error{err}
		}
		for _, player := range players {
			nameToPlayerMap[player.Name] = player
		}
		players = make([]*model.Player, len(playerNames))
		for i, name := range playerNames {
			players[i] = nameToPlayerMap[name]
		}
		return players, nil
	}
}

func wrapPlayerIDs(ids []int) []*model.PlayerFilter {
	wrapped := make([]*model.PlayerFilter, len(ids))
	for i := range ids {
		wrapped[i] = &model.PlayerFilter{PlayerID: &(ids[i])}
	}
	return wrapped
}

func fetchPlayerByID(db database.BasketballRepository, ctx context.Context) func(playerIDs []int) ([]*model.Player, []error) {
	return func(playerIDs []int) ([]*model.Player, []error) {
		idToPlayerMap := make(map[int]*model.Player, len(playerIDs))
		players, err := db.GetPlayers(ctx, false, wrapPlayerIDs(playerIDs)...)
		if err != nil {
			logrus.Errorf("error fetching players by id: %v", err)
			return nil, []error{err}
		}
		if len(players) == 0 {
			return nil, []error{fmt.Errorf("no players found")}
		}
		for _, player := range players {
			idToPlayerMap[player.PlayerID] = player
		}
		players = make([]*model.Player, len(playerIDs))
		for i, id := range playerIDs {
			players[i] = idToPlayerMap[id]
		}
		return players, nil
	}
}

func wrapTeamIDs(ids []int) []*model.TeamFilter {
	wrapped := make([]*model.TeamFilter, len(ids))
	for i := range ids {
		wrapped[i] = &model.TeamFilter{TeamID: &(ids[i])}
	}
	return wrapped
}

func fetchTeamByID(db database.BasketballRepository, ctx context.Context) func(teamIDs []int) ([]*model.Team, []error) {
	return func(teamIDs []int) ([]*model.Team, []error) {
		idToTeamMap := make(map[int]*model.Team, len(teamIDs))
		teams, err := db.GetTeams(ctx, false, wrapTeamIDs(teamIDs)...)
		if err != nil {
			logrus.Errorf("error fetching teams by id: %v", err)
			return nil, []error{err}
		}
		if len(teams) == 0 {
			return nil, []error{fmt.Errorf("no teams found")}
		}
		for _, team := range teams {
			idToTeamMap[team.TeamID] = team
		}
		teams = make([]*model.Team, len(teamIDs))
		for i, id := range teamIDs {
			teams[i] = idToTeamMap[id]
		}
		return teams, nil
	}
}

func wrapTeamAbbreviations(abbreviations []string) []*model.TeamFilter {
	wrapped := make([]*model.TeamFilter, len(abbreviations))
	for i := range abbreviations {
		wrapped[i] = &model.TeamFilter{Abbreviation: &(abbreviations[i])}
	}
	return wrapped
}

func fetchTeamByAbbreviation(db database.BasketballRepository, ctx context.Context) func(teamAbbreviations []string) ([]*model.Team, []error) {
	return func(teamAbbreviations []string) ([]*model.Team, []error) {
		abrToTeamMap := make(map[string]*model.Team, len(teamAbbreviations))
		teams, err := db.GetTeams(ctx, false, wrapTeamAbbreviations(teamAbbreviations)...)
		if err != nil {
			logrus.Errorf("error fetching teams by id: %v", err)
			return nil, []error{err}
		}
		if len(teams) == 0 {
			return nil, []error{fmt.Errorf("no teams found")}
		}
		for _, team := range teams {
			abrToTeamMap[team.Abbreviation] = team
		}
		teams = make([]*model.Team, len(teamAbbreviations))
		for i, id := range teamAbbreviations {
			teams[i] = abrToTeamMap[id]
		}
		return teams, nil
	}
}

func fetchTeamGameByPlayerGame(conn database.BasketballRepository, ctx context.Context) func(keys []model.PlayerGame) ([]*model.TeamGame, []error) {
	return func(keys []model.PlayerGame) ([]*model.TeamGame, []error) {
		panic("fetch teamgamebyplayergame not implemented")
		// opponentIDs := make([]int, len(keys))
		// gameIDs := make([]string, len(keys))
		// for i, key := range keys {
		// 	opponentIDs[i] = key.OpponentID
		// 	gameIDs[i] = key.GameID
		// }
		// teamGames := make([]*model.TeamGame, len(keys))
		// teamGamesByPlayerGame := make(map[string]*model.TeamGame, len(keys))
		// errs := make([]error, len(keys))
		// cur, err := conn.Collection("teamgames").Find(r.Context(), bson.M{"gameID": bson.M{"$in": gameIDs}, "opponent": bson.M{"$in": opponentIDs}})
		// if err != nil {
		// 	return nil, []error{err}
		// }
		// defer cur.Close(r.Context())

		// for cur.Next(r.Context()) {
		// 	game := &model.TeamGame{}
		// 	err := cur.Decode(&game)
		// 	if err != nil {
		// 		return nil, []error{err}
		// 	}
		// 	mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
		// 	teamGamesByPlayerGame[mapKey] = game
		// }
		// if err := cur.Err(); err != nil {
		// 	return nil, []error{err}
		// }
		// for i, game := range keys {
		// 	mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
		// 	teamGames[i] = teamGamesByPlayerGame[mapKey]
		// }
		// return teamGames, errs
	}
}

func fetchOpponentGameByPlayerGame(keys []model.PlayerGame) ([]*model.TeamGame, []error) {
	panic("fetch opponentgamebyplayergame not implemented")
	// 	teamIDs := make([]int, len(keys))
	// 	gameIDs := make([]string, len(keys))
	// 	for i, key := range keys {
	// 		teamIDs[i] = key.TeamID
	// 		gameIDs[i] = key.GameID
	// 	}
	// 	teamGames := make([]*model.TeamGame, len(keys))
	// 	teamGamesByPlayerGame := make(map[string]*model.TeamGame, len(keys))
	// 	errs := make([]error, len(keys))
	// 	cur, err := conn.Collection("teamgames").Find(r.Context(), bson.M{"gameID": bson.M{"$in": gameIDs}, "opponent": bson.M{"$in": teamIDs}})
	// 	if err != nil {
	// 		return nil, []error{err}
	// 	}
	// 	defer cur.Close(r.Context())

	// 	for cur.Next(r.Context()) {
	// 		game := &model.TeamGame{}
	// 		err := cur.Decode(&game)
	// 		if err != nil {
	// 			return nil, []error{err}
	// 		}
	// 		mapKey := fmt.Sprintf("%s:%d", game.GameID, game.TeamID)
	// 		teamGamesByPlayerGame[mapKey] = game
	// 	}
	// 	if err := cur.Err(); err != nil {
	// 		return nil, []error{err}
	// 	}
	// 	for i, game := range keys {
	// 		mapKey := fmt.Sprintf("%s:%d", game.GameID, game.OpponentID)
	// 		teamGames[i] = teamGamesByPlayerGame[mapKey]
	// 	}
	// 	return teamGames, errs
}

func fetchPlayerGamesByFilter(keys []model.GameFilter) ([][]*model.PlayerGame, []error) {
	panic("fetch playergamesbyfilter not implemented")
	// 	games := make([][]*model.PlayerGame, len(keys))
	// 	errs := make([]error, len(keys))
	// 	gamesByPlayerID := make(map[int][]*model.PlayerGame, len(keys))
	// 	cur, err := conn.GetPlayerGames(r.Context(), keys)
	// 	if err != nil {
	// 		return nil, []error{err}
	// 	}
	// 	defer cur.Close(r.Context())

	// 	for cur.Next(r.Context()) {
	// 		game := &model.PlayerGame{}
	// 		err := cur.Decode(&game)
	// 		if err != nil {
	// 			return nil, []error{err}
	// 		}
	// 		gamesByPlayerID[game.PlayerID] = append(gamesByPlayerID[game.PlayerID], game)
	// 	}
	// 	if err := cur.Err(); err != nil {
	// 		return nil, []error{err}
	// 	}
	// 	for i, filter := range keys {
	// 		games[i] = gamesByPlayerID[*filter.PlayerID]
	// 	}
	// 	return games, errs
}

// func fetchPlayerAverage(keys []model.GameFilter) ([]*model.PlayerAverage, []error) {
// 	averages := make([]*model.PlayerAverage, len(keys))
// 	errs := make([]error, len(keys))
// 	createFilterWithoutIDs := make([]model.GameFilter, len(keys))
// 	for i, filter := range keys {
// 		newFilter := filter
// 		newFilter.OpponentID = nil
// 		newFilter.GameID = nil
// 		newFilter.PlayerID = nil
// 		createFilterWithoutIDs[i] = newFilter
// 	}
// 	playerAverages, err := conn.GetAverages(r.Context(), createFilterWithoutIDs)
// 	if err != nil || len(*playerAverages) == 0 {
// 		return nil, []error{err}
// 	}
// 	for i, filter := range keys {
// 		findPlayerAverage := (*playerAverages)[0]
// 		for _, p := range *playerAverages {
// 			if p.Player.PlayerID == *filter.PlayerID {
// 				findPlayerAverage = p
// 			}
// 		}
// 		averages[i] = &findPlayerAverage
// 	}
// 	return averages, errs
// }

// func fetchSimilarPlayer(keys []model.SimilarPlayerInput) ([][]*model.Player, []error) {
// 	similarPlayers := make([][]*model.Player, len(keys))
// 	errs := make([]error, len(keys))
// 	createFilterWithoutIDs := make([]model.GameFilter, len(keys))
// 	for i, filter := range keys {
// 		newFilter := *filter.GameFilter
// 		newFilter.OpponentID = nil
// 		newFilter.GameID = nil
// 		newFilter.PlayerID = nil
// 		createFilterWithoutIDs[i] = newFilter
// 	}
// 	playerAverages, err := conn.GetAverages(r.Context(), createFilterWithoutIDs)
// 	if err != nil || len(*playerAverages) == 0 {
// 		return nil, []error{err}
// 	}
// 	for i, filter := range keys {
// 		targetPlayer := (*playerAverages)[0]
// 		for _, p := range *playerAverages {
// 			if p.Player.PlayerID == *filter.GameFilter.PlayerID {
// 				targetPlayer = p
// 			}
// 		}
// 		similarPlayers[i] = util.SimilarPlayers(*playerAverages, targetPlayer)
// 	}
// 	return similarPlayers, errs
// }

// func fetchSimilarTeam(keys []model.SimilarTeamInput) ([][]*model.Team, []error) {
// 	similarTeams := make([][]*model.Team, len(keys))
// 	errs := make([]error, len(keys))
// 	createFilterWithoutIDs := make([]model.GameFilter, len(keys))
// 	for i, filter := range keys {
// 		newFilter := *filter.GameFilter
// 		newFilter.OpponentID = nil
// 		newFilter.GameID = nil
// 		newFilter.PlayerID = nil
// 		createFilterWithoutIDs[i] = newFilter

// 	}
// 	teamAverages, err := conn.GetTeamAverages(r.Context(), createFilterWithoutIDs)
// 	if err != nil || len(*teamAverages) == 0 {
// 		return nil, []error{err}
// 	}
// 	for i, filter := range keys {
// 		targetTeam := (*teamAverages)[0]
// 		for _, t := range *teamAverages {
// 			if t.Team.TeamID == *filter.GameFilter.TeamID {
// 				targetTeam = t
// 			}
// 		}
// 		similarTeams[i] = util.SimilarTeams(*teamAverages, targetTeam)
// 	}
// 	return similarTeams, errs
// }
