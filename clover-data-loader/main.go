package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/cloverdata"
	"github.com/zvandehy/DataTrain/clover-data-loader/database"
	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		logrus.Fatalf("DSN is empty")
	}

	db, err := database.NewDB(dsn)
	if err != nil {
		logrus.Fatalf("couldn't connect to database: %v", err)
	}

	go func() {
		getAndUploadData(db)
		tick := time.Tick(4 * time.Hour)
		for range tick {
			getAndUploadData(db)
		}
	}()

	//TODO: Move Getting Sportsbook Data to this service

	router := chi.NewRouter()
	router.Get("/update", func(w http.ResponseWriter, r *http.Request) {
		getAndUploadData(db)
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getAndUploadData(db *database.Database) {
	from, to := GetPreviousDays(3)
	err := UploadDataFromRange(db, from, to)
	if err != nil {
		logrus.Errorf("couldn't update past 3 days of games: %v", err)
	}

	randomDay := GetRandomPreviousDay()
	err = UploadDataFromRange(db, randomDay, randomDay)
	if err != nil {
		logrus.Errorf("couldn't update random day (%s) games: %v", randomDay.Format("2006-01-02"), err)
	}
}

func FindAndUploadMissingPlayers(db *database.Database, playerIDs []int) error {
	missingIDs, err := db.GetMissingPlayerIDs(playerIDs)
	if err != nil {
		logrus.Errorf("couldn't get missing player IDs: %v", err)
	}
	fmt.Println(missingIDs)

	if len(missingIDs) > 0 {
		leagueBios, err := nba.NBAPlayers()
		if err != nil {
			logrus.Errorf("couldn't get player bio data: %v", err)
		}

		defense, err := nba.NBAPlayerDefense(nba.DefenseDashboardParams{})
		if err != nil {
			logrus.Errorf("couldn't get player defense data: %v", err)
		}

		var players []cloverdata.Player
		for _, playerID := range missingIDs {
			playerBio := nba.FindLeagueBioByID(leagueBios.ResultSets[0].RowSet, playerID)
			playerDefense := nba.FindDefenseByID(defense.ResultSets[0].RowSet, playerID)
			player := cloverdata.GetPlayer(playerBio, playerDefense)
			players = append(players, *player)
		}

		err = db.UploadPlayers(players)
		if err != nil {
			logrus.Errorf("couldn't upload missing players: %v", err)
		}
	}
	return nil
}

func GetPreviousDays(days int) (time.Time, time.Time) {
	today := time.Now()
	from := today.AddDate(0, 0, -days)
	to := today.AddDate(0, 0, -1)
	return from, to
}

func GetRandomPreviousDay() time.Time {
	today := time.Now()
	startOfSeason, _ := time.Parse("2006-01-02", "2022-10-18")
	rand.Seed(time.Now().UnixNano())
	// get random day between 1 day ago and Oct 18th, 2022
	from := today.AddDate(0, 0, -1).AddDate(0, 0, -rand.Intn(365))
	for from.Before(startOfSeason) {
		from = today.AddDate(0, 0, -1).AddDate(0, 0, -rand.Intn(365))
	}
	return from
}

func UploadDataFromRange(db *database.Database, from, to time.Time) error {
	playerGames, err := GetPlayerGamesByDateRange(from, to)
	if err != nil {
		return err
	}
	logrus.Infof("got %d player games from the NBA API", len(playerGames))

	uniquePlayers := make(map[int]bool)
	for _, playerGame := range playerGames {
		uniquePlayers[playerGame.PlayerID] = true
	}
	playerIDs := make([]int, 0, len(uniquePlayers))
	for playerID := range uniquePlayers {
		playerIDs = append(playerIDs, playerID)
	}
	err = FindAndUploadMissingPlayers(db, playerIDs)
	if err != nil {
		return err
	}

	err = db.UploadPlayerGames(playerGames)
	if err != nil {
		logrus.Errorf("couldn't upload player games: %v", err)
		return err
	}

	teamGames, err := GetTeamGamesByDateRange(from, to)
	if err != nil {
		return err
	}
	logrus.Infof("got %d team games from the NBA API", len(teamGames))

	err = db.UploadTeamGames(teamGames)
	if err != nil {
		logrus.Errorf("couldn't upload team games: %v", err)
		return err
	}

	for day := from; day.Before(to.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		err := db.UploadStandardizedPlayerComparisons(day.Format("2006-01-02"))
		if err != nil {
			logrus.Errorf("couldn't upload standardized player comparisons: %v", err)
			return err
		}
	}
	return nil
}

func GetTeamGamesByDateRange(from, to time.Time) ([]cloverdata.TeamGame, error) {
	gamesByDay := make(map[string]GameAggregation)
	for day := from; day.Before(to.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		gamesByDay[day.Format("2006-01-02")] = GameAggregation{}
	}

	teamGameData, err := nba.NBATeamGameFinder(nba.LeagueGameFinderParams{
		DateTo:   to.Format("01/02/2006"),
		DateFrom: from.Format("01/02/2006"),
	})
	if err != nil {
		logrus.Errorf("couldn't get team league game finder: %v", err)
	}
	time.Sleep(5 * time.Second)
	teamLeagueGames := teamGameData.ResultSets[0].RowSet
	for _, game := range teamLeagueGames {
		if game.WL == "" {
			continue
		}
		date, err := GetDate(game.GameDate)
		if err != nil {
			logrus.Errorf("couldn't parse team game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.TeamLeagueGames = append(agg.TeamLeagueGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	advancedTeamGameData, err := nba.NBATeamAdvancedGameLog(nba.GameLogParams{DateFrom: from.Format("01/02/2006"),
		DateTo: to.Format("01/02/2006")})
	if err != nil {
		logrus.Errorf("couldn't get NBATeamGameLog data: %v", err)
	}
	time.Sleep(5 * time.Second)
	advancedTeamGames := advancedTeamGameData.ResultSets[0].RowSet
	for _, game := range advancedTeamGames {
		if game.WL == "" {
			continue
		}
		date, err := GetDate(game.GAME_DATE)
		if err != nil {
			logrus.Errorf("couldn't parse advanced team game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.AdvancedTeamGames = append(agg.AdvancedTeamGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}
	today := time.Now()
	var games []cloverdata.TeamGame
	for _, agg := range gamesByDay {
		for _, teamGame := range agg.TeamLeagueGames {
			for _, advancedTeamGame := range agg.AdvancedTeamGames {
				if teamGame.TeamID == advancedTeamGame.TEAM_ID {
					game := cloverdata.GetTeamGame(teamGame, advancedTeamGame)
					game.CreatedAt = &today
					game.UpdatedAt = &today
					games = append(games, game)
					break
				}
			}
		}
	}
	return games, nil
}

func GetPlayerGamesByDateRange(from, to time.Time) ([]cloverdata.PlayerGame, error) {
	leagueParams := nba.LeagueGameFinderParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006"),
	}

	gamesByDay := make(map[string]GameAggregation)
	for day := from; day.Before(to.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		passingGameData, err := nba.PassingStats(nba.DashboardParams{
			DateFrom: day.Format("2006-01-02"),
			DateTo:   day.Format("2006-01-02")})
		if err != nil {
			logrus.Errorf("couldn't get PassingStats data: %v", err)
		}
		passingPlayerGames := passingGameData.ResultSets[0].RowSet
		gamesByDay[day.Format("2006-01-02")] = GameAggregation{
			PassingGames: passingPlayerGames,
		}
		time.Sleep(5 * time.Second)
	}

	leagueGameData, err := nba.NBAPlayerGameFinder(leagueParams)
	if err != nil {
		logrus.Errorf("couldn't get LeagueGameFinder data: %v", err)
	}
	time.Sleep(5 * time.Second)

	playerLeagueGames := leagueGameData.ResultSets[0].RowSet
	for _, game := range playerLeagueGames {
		if game.WL == "" {
			continue
		}
		date, err := GetDate(game.GameDate)
		if err != nil {
			logrus.Errorf("couldn't parse league game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.LeagueGames = append(agg.LeagueGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	advancedGameData, err := nba.NBAAdvancedPlayerGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Errorf("couldn't get NBAPlayerBoxScoreAdvanced data: %v", err)
	}
	time.Sleep(5 * time.Second)

	advancedPlayerGames := advancedGameData.ResultSets[0].RowSet
	for _, game := range advancedPlayerGames {
		date, err := GetDate(game.GAME_DATE)
		if err != nil {
			logrus.Errorf("couldn't parse advanced game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.AdvancedGames = append(agg.AdvancedGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	miscGameData, err := nba.NBAMiscPlayerGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Errorf("couldn't get NBAPlayerMisc data: %v", err)
	}
	time.Sleep(5 * time.Second)

	miscPlayerGames := miscGameData.ResultSets[0].RowSet
	for _, game := range miscPlayerGames {
		date, err := GetDate(game.GAME_DATE)
		if err != nil {
			logrus.Errorf("couldn't parse misc game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.MiscGames = append(agg.MiscGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	teamGameData, err := nba.NBATeamTraditionalGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Errorf("couldn't get NBATeamGameLog data: %v", err)
	}
	time.Sleep(5 * time.Second)
	teamGames := teamGameData.ResultSets[0].RowSet
	for _, game := range teamGames {
		if game.WL == "" {
			continue
		}
		date, err := GetDate(game.GAME_DATE)
		if err != nil {
			logrus.Errorf("couldn't parse (player) team game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.TeamGames = append(agg.TeamGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	today := time.Now()

	var playerGames []cloverdata.PlayerGame
	for day, agg := range gamesByDay {
		logrus.Infof("transforming games for %v\n", day)
		logrus.Infof("league games: %v\n", len(agg.LeagueGames))
		logrus.Infof("advanced games: %v\n", len(agg.AdvancedGames))
		logrus.Infof("misc games: %v\n", len(agg.MiscGames))
		logrus.Infof("passing games: %v\n", len(agg.PassingGames))
		logrus.Infof("team games: %v\n\n", len(agg.TeamGames))
		logrus.Infof("advanced team games %v\n\n", len(agg.AdvancedTeamGames))

		for _, leagueGame := range agg.LeagueGames {
			passingGame := MatchingPassingGame(leagueGame, agg.PassingGames)
			if passingGame.PLAYER_NAME == "" {
				logrus.Warnf("no passing game for %v %v on %v. Setting passing stats to zero value.\n", leagueGame.PlayerName, leagueGame.Matchup, leagueGame.GameDate)
			}
			playerGameBucket := cloverdata.PlayerGameBucket{
				LeagueGame:       leagueGame,
				BoxScoreAdvanced: MatchingAdvancedGame(leagueGame, agg.AdvancedGames),
				MiscStats:        MatchingMiscGame(leagueGame, agg.MiscGames),
				PassingStats:     passingGame,
				TeamBoxScore:     MatchingTeamGame(leagueGame, agg.TeamGames),
			}
			playerGame := cloverdata.TransformNBAPlayerGame(playerGameBucket)
			playerGame.CreatedAt = &today
			playerGame.UpdatedAt = &today
			playerGames = append(playerGames, *playerGame)
		}
	}
	return playerGames, nil
}

func MatchingAdvancedGame(leagueGame nba.PlayerLeagueGame, advancedGames []nba.PlayerAdvancedGameLog) nba.PlayerAdvancedGameLog {
	for _, advancedGame := range advancedGames {
		if int(leagueGame.PlayerID) == advancedGame.PLAYER_ID && leagueGame.GameID == advancedGame.GAME_ID && int(leagueGame.TeamID) == advancedGame.TEAM_ID {
			return advancedGame
		}
	}

	logrus.Warnf("Couldn't find matching advanced game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerAdvancedGameLog{}
}

func MatchingMiscGame(leagueGame nba.PlayerLeagueGame, miscGames []nba.PlayerMiscGameLog) nba.PlayerMiscGameLog {
	for _, miscGame := range miscGames {
		if leagueGame.PlayerID == miscGame.PLAYER_ID && leagueGame.GameID == miscGame.GAME_ID && leagueGame.TeamID == miscGame.TEAM_ID {
			return miscGame
		}
	}

	logrus.Warnf("Couldn't find matching misc game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerMiscGameLog{}
}

func MatchingPassingGame(leagueGame nba.PlayerLeagueGame, passingGames []nba.PlayerPassingStats) nba.PlayerPassingStats {
	for _, passingGame := range passingGames {
		if leagueGame.PlayerID == passingGame.PLAYER_ID && leagueGame.TeamID == passingGame.TEAM_ID {
			return passingGame
		}
	}

	logrus.Warnf("Couldn't find matching passing game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerPassingStats{}
}

func MatchingTeamGame(leagueGame nba.PlayerLeagueGame, teamGames []nba.TeamTraditionalGameLog) nba.TeamTraditionalGameLog {
	for _, teamGame := range teamGames {
		if leagueGame.TeamID == teamGame.TEAM_ID && leagueGame.GameID == teamGame.GAME_ID {
			return teamGame
		}
	}
	logrus.Warnf("Couldn't find matching team game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.TeamTraditionalGameLog{}
}

type GameAggregation struct {
	LeagueGames       []nba.PlayerLeagueGame
	AdvancedGames     []nba.PlayerAdvancedGameLog
	PassingGames      []nba.PlayerPassingStats
	MiscGames         []nba.PlayerMiscGameLog
	TeamGames         []nba.TeamTraditionalGameLog
	TeamLeagueGames   []nba.TeamLeagueGame
	AdvancedTeamGames []nba.TeamAdvancedGameLog
}

func GetDate(gameDate string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", gameDate)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("2006-01-02T15:04:05", gameDate)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("2006-01-02T15:04:05.000", gameDate)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("2006-01-02T15:04:05.000Z", gameDate)
	if err == nil {
		return date, nil
	}
	return time.Time{}, fmt.Errorf("couldn't parse date: %v", gameDate)
}
