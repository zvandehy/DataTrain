package cloverdata

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

type PlayerGame struct {
	// LeagueGameFinder
	Assists                int        `db:"assists"`
	Blocks                 int        `db:"blocks"`
	Date                   *time.Time `db:"date"`
	DefensiveRebounds      int        `db:"defensiveRebounds"`
	FieldGoalPercentage    float64    `db:"fieldGoalPct"`
	FieldGoalsAttempted    int        `db:"fieldGoalsAttempted"`
	FieldGoalsMade         int        `db:"fieldGoalsMade"`
	FreeThrowsAttempted    int        `db:"freeThrowsAttempted"`
	FreeThrowsMade         int        `db:"freeThrowsMade"`
	FreeThrowsPercentage   float64    `db:"freeThrowPct"`
	GameID                 string     `db:"gameID"`
	HomeOrAway             string     `db:"homeAway"`
	Minutes                float64    `db:"minutes"`
	OffensiveRebounds      int        `db:"offensiveRebounds"`
	OpponentID             int        `db:"opponentID"`
	Outcome                string     `db:"outcome"`
	PersonalFouls          int        `db:"personalFouls"`
	PlayerID               int        `db:"playerID"`
	Playoffs               bool       `db:"playoffs"`
	Points                 int        `db:"points"`
	Rebounds               int        `db:"rebounds"`
	Season                 string     `db:"season"`
	Steals                 int        `db:"steals"`
	TeamID                 int        `db:"teamID"`
	ThreePointPercentage   float64    `db:"threePointPct"`
	ThreePointersAttempted int        `db:"threePointersAttempted"`
	ThreePointersMade      int        `db:"threePointersMade"`
	Turnovers              int        `db:"turnovers"`

	// PlayerBoxScoreAdvanced
	AssistPercentage             float64 `db:"assistPct"`
	Usage                        float64 `db:"usage"`
	TrueShootingPercentage       float64 `db:"trueShootingPct"`
	EffectiveFieldGoalPercentage float64 `db:"effectiveFieldGoalPct"`
	DefensiveReboundPercentage   float64 `db:"defensiveReboundPct"`
	OffensiveReboundPercentage   float64 `db:"offensiveReboundPct"`

	// PassingBoxScore
	PotentialAssists     int     `db:"potentialAssists"`
	Passes               int     `db:"passes"`
	AssistConversionRate float64 `db:"assistConversionRate"`

	//Misc
	PersonalFoulsDrawn int `db:"personalFoulsDrawn"`

	// TeamBoxScore
	Margin int `db:"margin"`

	CreatedAt *time.Time `db:"CreatedAt"`
	UpdatedAt *time.Time `db:"UpdatedAt"`
}

func TransformNBAPlayerGame(playerGameBucket PlayerGameBucket) *PlayerGame {
	pg := &PlayerGame{}
	pg.fillPlayerGameFromLeagueGame(playerGameBucket.LeagueGame)
	pg.fillPlayerGameFromBoxScoreAdvanced(playerGameBucket.BoxScoreAdvanced)
	pg.fillPlayerGameFromPassingStats(playerGameBucket.PassingStats)
	pg.fillPlayerGameFromMiscStats(playerGameBucket.MiscStats)
	pg.fillPlayerGameFromTeamBoxScore(playerGameBucket.TeamBoxScore)
	return pg
}

type PlayerGameBucket struct {
	LeagueGame       nba.PlayerLeagueGame
	BoxScoreAdvanced nba.PlayerAdvancedGameLog
	PassingStats     nba.PlayerPassingStats
	MiscStats        nba.PlayerMiscGameLog
	TeamBoxScore     nba.TeamTraditionalGameLog
}

func (pg *PlayerGame) fillPlayerGameFromBoxScoreAdvanced(boxScoreAdvanced nba.PlayerAdvancedGameLog) {
	pg.AssistPercentage = boxScoreAdvanced.AST_PCT
	pg.Usage = boxScoreAdvanced.USG_PCT
	pg.TrueShootingPercentage = boxScoreAdvanced.TS_PCT
	pg.EffectiveFieldGoalPercentage = boxScoreAdvanced.EFG_PCT
	pg.DefensiveReboundPercentage = boxScoreAdvanced.DREB_PCT
	pg.OffensiveReboundPercentage = boxScoreAdvanced.OREB_PCT
}

func (pg *PlayerGame) fillPlayerGameFromMiscStats(miscStats nba.PlayerMiscGameLog) {
	pg.PersonalFoulsDrawn = int(miscStats.PFD)
}

func (pg *PlayerGame) fillPlayerGameFromTeamBoxScore(boxScore nba.TeamTraditionalGameLog) {
	pg.Margin = int(boxScore.PLUS_MINUS)
}

func (pg *PlayerGame) fillPlayerGameFromPassingStats(passingStats nba.PlayerPassingStats) {
	pg.Assists = int(passingStats.AST)
	pg.PotentialAssists = int(passingStats.POTENTIAL_AST)
	pg.Passes = int(passingStats.PASSES_MADE)
	if passingStats.POTENTIAL_AST > 0 {
		pg.AssistConversionRate = passingStats.AST / passingStats.POTENTIAL_AST
	} else {
		pg.AssistConversionRate = 0
	}
}

func (pg *PlayerGame) fillPlayerGameFromLeagueGame(leagueGameFinderResults nba.PlayerLeagueGame) {
	pg.Assists = int(leagueGameFinderResults.Ast)
	pg.Blocks = int(leagueGameFinderResults.Blk)
	pg.DefensiveRebounds = int(leagueGameFinderResults.Dreb)
	pg.FieldGoalPercentage = leagueGameFinderResults.FgPct
	pg.FieldGoalsAttempted = int(leagueGameFinderResults.Fga)
	pg.FieldGoalsMade = int(leagueGameFinderResults.Fgm)
	pg.FreeThrowsAttempted = int(leagueGameFinderResults.Fta)
	pg.FreeThrowsMade = int(leagueGameFinderResults.Ftm)
	pg.FreeThrowsPercentage = leagueGameFinderResults.FtPct
	pg.GameID = leagueGameFinderResults.GameID
	pg.HomeOrAway = HomeOrAway(leagueGameFinderResults)
	pg.Minutes = leagueGameFinderResults.Min
	pg.OffensiveRebounds = int(leagueGameFinderResults.Oreb)
	pg.OpponentID = OpponentID(leagueGameFinderResults)
	pg.Outcome = GameOutcome(leagueGameFinderResults)
	pg.PersonalFouls = int(leagueGameFinderResults.Pf)
	pg.PlayerID = int(leagueGameFinderResults.PlayerID)
	pg.Playoffs = leagueGameFinderResults.Playoffs
	pg.Points = int(leagueGameFinderResults.Pts)
	pg.Rebounds = int(leagueGameFinderResults.Reb)
	pg.Season = SeasonIDToSeason(leagueGameFinderResults.SeasonID)
	pg.Steals = int(leagueGameFinderResults.Stl)
	pg.TeamID = int(leagueGameFinderResults.TeamID)
	pg.ThreePointPercentage = leagueGameFinderResults.Fg3Pct
	pg.ThreePointersAttempted = int(leagueGameFinderResults.Fg3A)
	pg.ThreePointersMade = int(leagueGameFinderResults.Fg3M)
	pg.Turnovers = int(leagueGameFinderResults.Tov)

	date, err := time.Parse("2006-01-02", leagueGameFinderResults.GameDate)
	if err != nil {
		logrus.Errorf("couldn't parse date: %v", err)
	}
	pg.Date = &date
}

func SeasonIDToSeason(seasonID string) string {
	season, err := strconv.Atoi(seasonID[1:])
	if err != nil {
		logrus.Errorf("couldn't convert seasonID to int: %v", err)
	}
	return fmt.Sprintf("%s-%d", seasonID[1:], season+1)
}

func HomeOrAway(game nba.PlayerLeagueGame) string {
	if strings.Contains(game.Matchup, "@") {
		return "away"
	}
	return "home"
}

func GameOutcome(game nba.PlayerLeagueGame) string {
	if game.WL == "W" {
		return "WIN"
	}
	if game.WL == "L" {
		return "LOSS"
	}
	return "PENDING"
}

func OpponentID(game nba.PlayerLeagueGame) int {
	return int(nba.TeamAbbreviationToTeamID[game.Matchup[len(game.Matchup)-3:]])
}
