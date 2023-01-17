package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/cloverdata"
	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

func main() {
	today := time.Now()
	playerGames, err := GetPlayerGamesByDateRange(today.AddDate(0, 0, -3), today.AddDate(0, 0, -1))
	if err != nil {
		logrus.Fatalf("couldn't get player games: %v", err)
	}

	fmt.Printf("got %d player games from the NBA API", len(playerGames))

	logrus.Infof("%+v", playerGames[0])
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
			logrus.Fatalf("couldn't get PassingStats data: %v", err)
		}
		passingPlayerGames := passingGameData.ResultSets[0].RowSet
		gamesByDay[day.Format("2006-01-02")] = GameAggregation{
			PassingGames: passingPlayerGames,
		}
		time.Sleep(5 * time.Second)
	}

	leagueGameData, err := nba.NBAPlayerGameFinder(leagueParams)
	if err != nil {
		logrus.Fatalf("couldn't get LeagueGameFinder data: %v", err)
	}
	time.Sleep(5 * time.Second)

	playerLeagueGames := leagueGameData.ResultSets[0].RowSet
	for _, game := range playerLeagueGames {
		date, err := time.Parse("2006-01-02", game.GameDate)
		if err != nil {
			logrus.Fatalf("couldn't parse league game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.LeagueGames = append(agg.LeagueGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	advancedGameData, err := nba.NBAAdvancedPlayerGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Fatalf("couldn't get NBAPlayerBoxScoreAdvanced data: %v", err)
	}
	time.Sleep(5 * time.Second)

	advancedPlayerGames := advancedGameData.ResultSets[0].RowSet
	for _, game := range advancedPlayerGames {
		date, err := time.Parse("2006-01-02T15:04:05", game.GAME_DATE)
		if err != nil {
			logrus.Fatalf("couldn't parse advanced game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.AdvancedGames = append(agg.AdvancedGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	miscGameData, err := nba.NBAMiscPlayerGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Fatalf("couldn't get NBAPlayerMisc data: %v", err)
	}
	time.Sleep(5 * time.Second)

	miscPlayerGames := miscGameData.ResultSets[0].RowSet
	for _, game := range miscPlayerGames {
		date, err := time.Parse("2006-01-02T15:04:05", game.GAME_DATE)
		if err != nil {
			logrus.Fatalf("couldn't parse misc game date: %v", err)
		}
		agg := gamesByDay[date.Format("2006-01-02")]
		agg.MiscGames = append(agg.MiscGames, game)
		gamesByDay[date.Format("2006-01-02")] = agg
	}

	teamGameData, err := nba.NBATeamTraditionalGameLog(nba.GameLogParams{
		DateFrom: from.Format("01/02/2006"),
		DateTo:   to.Format("01/02/2006")})
	if err != nil {
		logrus.Fatalf("couldn't get NBATeamGameLog data: %v", err)
	}
	time.Sleep(5 * time.Second)

	teamGames := teamGameData.ResultSets[0].RowSet
	for _, game := range teamGames {
		date, err := time.Parse("2006-01-02T15:04:05", game.GAME_DATE)
		if err != nil {
			logrus.Fatalf("couldn't parse team game date: %v", err)
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

func MatchingAdvancedGame(leagueGame nba.LeagueGame, advancedGames []nba.PlayerAdvancedGameLog) nba.PlayerAdvancedGameLog {
	for _, advancedGame := range advancedGames {
		if int(leagueGame.PlayerID) == advancedGame.PLAYER_ID && leagueGame.GameID == advancedGame.GAME_ID && int(leagueGame.TeamID) == advancedGame.TEAM_ID {
			return advancedGame
		}
	}

	logrus.Warnf("Couldn't find matching advanced game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerAdvancedGameLog{}
}

func MatchingMiscGame(leagueGame nba.LeagueGame, miscGames []nba.PlayerMiscGameLog) nba.PlayerMiscGameLog {
	for _, miscGame := range miscGames {
		if leagueGame.PlayerID == miscGame.PLAYER_ID && leagueGame.GameID == miscGame.GAME_ID && leagueGame.TeamID == miscGame.TEAM_ID {
			return miscGame
		}
	}

	logrus.Warnf("Couldn't find matching misc game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerMiscGameLog{}
}

func MatchingPassingGame(leagueGame nba.LeagueGame, passingGames []nba.PlayerPassingStats) nba.PlayerPassingStats {
	for _, passingGame := range passingGames {
		if leagueGame.PlayerID == passingGame.PLAYER_ID && leagueGame.TeamID == passingGame.TEAM_ID {
			return passingGame
		}
	}

	logrus.Warnf("Couldn't find matching passing game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.PlayerPassingStats{}
}

func MatchingTeamGame(leagueGame nba.LeagueGame, teamGames []nba.TeamTraditionalGameLog) nba.TeamTraditionalGameLog {
	for _, teamGame := range teamGames {
		if leagueGame.TeamID == teamGame.TEAM_ID && leagueGame.GameID == teamGame.GAME_ID {
			return teamGame
		}
	}
	logrus.Warnf("Couldn't find matching team game for: %v %v %v %v %v %v", leagueGame.PlayerName, leagueGame.PlayerID, leagueGame.Matchup, leagueGame.GameID, leagueGame.GameDate, leagueGame.TeamID)
	return nba.TeamTraditionalGameLog{}
}

type GameAggregation struct {
	LeagueGames   []nba.LeagueGame
	AdvancedGames []nba.PlayerAdvancedGameLog
	PassingGames  []nba.PlayerPassingStats
	MiscGames     []nba.PlayerMiscGameLog
	TeamGames     []nba.TeamTraditionalGameLog
}
