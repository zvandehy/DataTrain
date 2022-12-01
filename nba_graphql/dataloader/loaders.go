package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

type LoadersKey string

const loadersKey LoadersKey = "dataloaders"

const waitTime = 10 * time.Millisecond
const maxBatch = 50

type Loaders struct {
	PlayerByID          PlayerLoaderID
	PlayerWithGamesByID PlayerLoaderID
	TeamByID            TeamLoaderID
	TeamByAbr           TeamLoaderABR
	// PlayerByFilter           PlayerLoader
	PlayerGameByFilter PlayerGameLoader
	// TeamGameByPlayerGame     TeamGameLoader
	// OpponentGameByPlayerGame TeamGameLoader
	SimilarPlayerLoader SimilarPlayerLoader
}

func Middleware(conn database.BasketballRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			PlayerByID: *NewPlayerLoaderID(
				PlayerLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch:    fetchPlayerByID(conn, r.Context(), false),
				},
			),
			PlayerWithGamesByID: *NewPlayerLoaderID(
				PlayerLoaderIDConfig{
					MaxBatch: maxBatch,
					Wait:     waitTime,
					Fetch:    fetchPlayerByID(conn, r.Context(), true),
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
			SimilarPlayerLoader: *NewSimilarPlayerLoader(
				SimilarPlayerLoaderConfig{
					MaxBatch: maxBatch * 2,
					Wait:     waitTime * 2,
					Fetch:    fetchSimilarPlayerIDs(conn, r.Context()),
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
			PlayerGameByFilter: *NewPlayerGameLoader(
				PlayerGameLoaderConfig{
					MaxBatch: maxBatch * 10,
					Wait:     waitTime,
					Fetch:    fetchPlayerGamesByFilter(conn, r.Context()),
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

func similarPlayerInputKey(input model.SimilarPlayerInput) string {
	poolKey := ""
	if input.PlayerPoolFilter != nil {
		poolKey = input.PlayerPoolFilter.Key()
	}
	return fmt.Sprintf("%v-%s", input.StatsOfInterest, poolKey)
}

func similarPlayerQueryKey(query model.SimilarPlayerQuery) string {
	return fmt.Sprintf("%s-%s", query.EndDate.Format("2006-01-02"), similarPlayerInputKey(query.SimilarPlayerInput))
}

func fetchSimilarPlayerIDs(db database.BasketballRepository, ctx context.Context) func(similarPlayerQueries []*model.SimilarPlayerQuery) ([][]model.Player, []error) {
	return func(similarPlayerQueries []*model.SimilarPlayerQuery) ([][]model.Player, []error) {
		logrus.Warn("FETCHING STANDARDIZED STATS")
		uniqueQueries := make(map[string][]model.SimilarPlayerQuery)
		queryToSimilarPlayerIDs := make(map[string][]int)
		for _, query := range similarPlayerQueries {
			key := similarPlayerQueryKey(*query)
			if _, ok := uniqueQueries[key]; !ok {
				uniqueQueries[key] = []model.SimilarPlayerQuery{*query}
				logrus.Warnf("add unique query: %s", key)
			} else {
				uniqueQueries[key] = append(uniqueQueries[key], *query)
			}
		}
		for _, queries := range uniqueQueries {
			playerZScores, err := db.GetStandardizedPlayerStats(ctx, queries[0], lo.Map(queries, func(query model.SimilarPlayerQuery, _ int) int {
				return query.ToPlayerID
			})...)
			if err != nil {
				return nil, []error{err}
			}
			for _, query := range queries {
				mostSimilarIDs := database.FindMostSimilarPlayerIDs(query.ToPlayerID, query.SimilarPlayerInput.Limit, playerZScores)
				if len(mostSimilarIDs) == 0 {
					return nil, []error{fmt.Errorf("no similar players found")}
				}
				queryToSimilarPlayerIDs[fmt.Sprintf("%d-%d-%s", query.ToPlayerID, query.SimilarPlayerInput.Limit, similarPlayerQueryKey(query))] = mostSimilarIDs
			}
		}

		// get all similar players
		allSimilarPlayerIDs := []*model.PlayerFilter{}
		for _, similarPlayerIDs := range queryToSimilarPlayerIDs {
			allSimilarPlayerIDs = append(allSimilarPlayerIDs, lo.Map(similarPlayerIDs, func(id int, _ int) *model.PlayerFilter {
				return &model.PlayerFilter{
					PlayerID: &id,
				}
			})...)
		}
		allPlayers, err := db.GetPlayers(ctx, true, allSimilarPlayerIDs...)
		if err != nil {
			return nil, []error{err}
		}
		similarPlayers := make([][]model.Player, len(similarPlayerQueries))
		for i, query := range similarPlayerQueries {
			similarIDs := queryToSimilarPlayerIDs[fmt.Sprintf("%d-%d-%s", query.ToPlayerID, query.SimilarPlayerInput.Limit, similarPlayerQueryKey(*query))]
			// fmt.Printf("%d ==> %v\n", query.ToPlayerID, similarIDs)
			similarPlayers[i] = lo.FilterMap(allPlayers, func(player *model.Player, _ int) (model.Player, bool) {
				return *player, lo.Contains(similarIDs, player.PlayerID)
			})

		}
		return similarPlayers, nil
	}
}

func wrapPlayerIDs(ids []int) []*model.PlayerFilter {
	wrapped := make([]*model.PlayerFilter, len(ids))
	for i := range ids {
		wrapped[i] = &model.PlayerFilter{PlayerID: &(ids[i])}
	}
	return wrapped
}

func fetchPlayerByID(db database.BasketballRepository, ctx context.Context, loadGames bool) func(playerIDs []int) ([]*model.Player, []error) {
	return func(playerIDs []int) ([]*model.Player, []error) {
		idToPlayerMap := make(map[int]*model.Player, len(playerIDs))
		players, err := db.GetPlayers(ctx, loadGames, wrapPlayerIDs(playerIDs)...)
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

func checkFiltersForSameSeason(filters []model.GameFilter) bool {
	expect := *filters[0].Seasons
	for _, filter := range filters {
		if len(*filter.Seasons) != len(expect) {
			return false
		}
		for i, season := range *filter.Seasons {
			if season != expect[i] {
				return false
			}
		}
	}
	return true
}

func fetchPlayerGamesByFilter(db database.BasketballRepository, ctx context.Context) func(keys []model.GameFilter) ([][]*model.PlayerGame, []error) {
	// panic("fetch playergamesbyfilter not implemented")
	return func(keys []model.GameFilter) ([][]*model.PlayerGame, []error) {
		if !checkFiltersForSameSeason(keys) {
			return nil, []error{fmt.Errorf("filters must be for the same season")}
		}
		playerIDToGamesMap := make(map[int][]*model.PlayerGame)

		games, err := db.GetPlayerGames(ctx, keys...)
		if err != nil {
			return nil, []error{err}
		}
		for _, game := range games {
			playerIDToGamesMap[game.PlayerID] = append(playerIDToGamesMap[game.PlayerID], game)
		}
		playerGames := make([][]*model.PlayerGame, len(keys))
		for i, filter := range keys {
			playerGames[i] = playerIDToGamesMap[*filter.PlayerID]
		}
		return playerGames, nil
	}

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
