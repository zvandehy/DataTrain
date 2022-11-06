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

const waitTime = 150 * time.Millisecond
const maxBatch = 50

type Loaders struct {
	// PlayerByID               PlayerLoaderID
	// PlayerByFilter           PlayerLoader
	// PlayerGameByFilter       PlayerGameLoader
	// TeamGameByPlayerGame     TeamGameLoader
	// OpponentGameByPlayerGame TeamGameLoader
}

func Middleware(conn database.BasketballRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			// PlayerByFilter: *NewPlayerLoader(
			// 	PlayerLoaderConfig{
			// 		MaxBatch: maxBatch,
			// 		Wait:     waitTime,
			// 		Fetch:    fetchPlayer,
			// 	},
			// ),
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

func fetchPlayer(keys []model.PlayerFilter) ([]*model.Player, []error) {
	panic("fetch player not implemented")
	// players := make([]*model.Player, len(keys))
	// playersByName := make(map[string]*model.Player, len(keys))
	// errs := make([]error, len(keys))
	// cur, err := conn.GetPlayersCursor(r.Context(), keys)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer cur.Close(r.Context())

	// for cur.Next(r.Context()) {
	// 	player := &model.Player{}
	// 	err := cur.Decode(&player)
	// 	if err != nil {
	// 		return nil, []error{err}
	// 	}
	// 	playersByName[player.Name] = player
	// }
	// if err := cur.Err(); err != nil {
	// 	return nil, []error{err}
	// }
	// for i, player := range keys {
	// 	players[i] = playersByName[*player.Name]
	// }
	// return players, errs
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
