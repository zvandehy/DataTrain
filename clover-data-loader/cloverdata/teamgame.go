package cloverdata

import (
	"strings"
	"time"

	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

type TeamGame struct {
	GameID            string `db:"gameID"`
	TeamID            int    `db:"teamID"`
	OpponentID        int    `db:"opponentID"`
	HomeAway          string `db:"homeAway"`
	Date              string `db:"date"`
	Outcome           string `db:"outcome"`
	Margin            int    `db:"margin"`
	Playoffs          bool   `db:"playoffs"`
	Assists           int    `db:"assists"`
	Blocks            int    `db:"blocks"`
	DefensiveRebounds int    `db:"defensiveRebounds"`
	OffensiveRebounds int    `db:"offensiveRebounds"`

	FieldGoalPct           float64 `db:"fieldGoalPct"`
	FieldGoalsAttempted    int     `db:"fieldGoalsAttempted"`
	FieldGoalsMade         int     `db:"fieldGoalsMade"`
	FreeThrowPct           float64 `db:"freeThrowPct"`
	FreeThrowsAttempted    int     `db:"freeThrowsAttempted"`
	FreeThrowsMade         int     `db:"freeThrowsMade"`
	ThreePointPct          float64 `db:"threePointPct"`
	ThreePointersAttempted int     `db:"threePointersAttempted"`
	ThreePointersMade      int     `db:"threePointersMade"`
	Points                 int     `db:"points"`
	PersonalFouls          int     `db:"personalFouls"`
	Rebounds               int     `db:"rebounds"`
	Steals                 int     `db:"steals"`
	Turnovers              int     `db:"turnovers"`
	Season                 string  `db:"season"`

	// Advanced
	DefensiveReboundPct   float64 `db:"defensiveReboundPct"`
	OffensiveReboundPct   float64 `db:"offensiveReboundPct"`
	PlusMinusPerHundred   float64 `db:"plusMinusPerHundred"`
	Possessions           int     `db:"possessions"`
	EffectiveFieldGoalPct float64 `db:"effectiveFieldGoalPct"`

	CreatedAt *time.Time `db:"CreatedAt"`
	UpdatedAt *time.Time `db:"UpdatedAt"`
}

func GetTeamGame(leagueGame nba.TeamLeagueGame, advancedGame nba.TeamAdvancedGameLog) TeamGame {
	teamGame := TeamGame{}
	teamGame.GameID = leagueGame.GameID
	teamGame.TeamID = int(leagueGame.TeamID)
	matchupSplit := strings.Split(leagueGame.Matchup, " ")
	opponentAbbreviation := matchupSplit[len(matchupSplit)-1]
	opponentID := nba.TeamAbbreviationToTeamID[opponentAbbreviation]
	teamGame.OpponentID = int(opponentID)
	if strings.Contains(leagueGame.Matchup, "@") {
		teamGame.HomeAway = "AWAY"
	} else {
		teamGame.HomeAway = "HOME"
	}
	date, err := time.Parse("2006-01-02", leagueGame.GameDate)
	if err != nil {
		panic(err)
	}
	teamGame.Date = date.Format("2006-01-02")
	outcome := "PENDING"
	if leagueGame.WL == "W" {
		outcome = "WIN"
	}
	if leagueGame.WL == "L" {
		outcome = "LOSS"
	}
	teamGame.Outcome = outcome
	teamGame.Margin = int(leagueGame.PlusMinus)
	teamGame.Playoffs = leagueGame.Playoffs
	teamGame.Assists = int(leagueGame.Ast)
	teamGame.Blocks = int(leagueGame.Blk)
	teamGame.DefensiveRebounds = int(leagueGame.Dreb)
	teamGame.OffensiveRebounds = int(leagueGame.Oreb)
	teamGame.FieldGoalPct = leagueGame.FgPct
	teamGame.FieldGoalsAttempted = int(leagueGame.Fga)
	teamGame.FieldGoalsMade = int(leagueGame.Fgm)
	teamGame.FreeThrowPct = leagueGame.FtPct
	teamGame.FreeThrowsAttempted = int(leagueGame.Fta)
	teamGame.FreeThrowsMade = int(leagueGame.Ftm)
	teamGame.ThreePointPct = leagueGame.Fg3Pct
	teamGame.ThreePointersAttempted = int(leagueGame.Fg3A)
	teamGame.ThreePointersMade = int(leagueGame.Fg3M)
	teamGame.Points = int(leagueGame.Pts)
	teamGame.PersonalFouls = int(leagueGame.Pf)
	teamGame.Rebounds = int(leagueGame.Reb)
	teamGame.Steals = int(leagueGame.Stl)
	teamGame.Turnovers = int(leagueGame.Tov)
	teamGame.Season = SeasonIDToSeason(leagueGame.SeasonID)
	teamGame.Possessions = int(advancedGame.POSS)
	teamGame.DefensiveReboundPct = advancedGame.DREB_PCT
	teamGame.OffensiveReboundPct = advancedGame.OREB_PCT
	teamGame.PlusMinusPerHundred = (leagueGame.PlusMinus / float64(advancedGame.POSS)) * 100
	return teamGame
}
