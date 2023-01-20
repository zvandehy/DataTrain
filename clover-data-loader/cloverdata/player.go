package cloverdata

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

type Player struct {
	FirstName    string   `json:"first_name" bson:"first_name" db:"firstName"`
	LastName     string   `json:"last_name" bson:"last_name" db:"lastName"`
	Name         string   `json:"name" bson:"name" db:"name"`
	PlayerID     int      `json:"playerID" bson:"playerID" db:"playerID"`
	Seasons      []string `json:"seasons" bson:"seasons"` //todo db
	Position     string   `json:"position" bson:"position" db:"position"`
	CurrentTeam  string   `json:"currentTeam" bson:"teamABR" db:"teamABR"`
	TeamID       int      `json:"teamID" bson:"teamID" db:"teamID"`
	Height       string   `json:"height" bson:"height" db:"height"`
	HeightInches int      `json:"heightInches" bson:"heightInches" db:"heightInches"`
	Weight       int      `json:"weight" bson:"weight" db:"weight"`

	CreatedAt *time.Time `json:"CreatedAt" bson:"CreatedAt" db:"CreatedAt"`
	UpdatedAt *time.Time `json:"UpdatedAt" bson:"UpdatedAt" db:"UpdatedAt"`
	League    string     `json:"league" bson:"league" db:"league"`
}

func GetPlayer(leaguePlayer nba.CommonPlayer, defenseStats nba.PlayerDefenseStats) *Player {
	p := &Player{}
	splitName := strings.SplitN(leaguePlayer.PLAYER_NAME, " ", 2)
	p.FirstName = splitName[0]
	p.LastName = splitName[1]
	p.Name = leaguePlayer.PLAYER_NAME
	p.PlayerID = leaguePlayer.PLAYER_ID
	splitPosition := strings.Split(defenseStats.PLAYER_POSITION, "-")
	position := string(splitPosition[0][0])
	if len(splitPosition) > 1 {
		position = position + "-" + string(splitPosition[1][0])
	}
	p.Position = position
	p.CurrentTeam = leaguePlayer.TEAM_ABBREVIATION
	p.TeamID = leaguePlayer.TEAM_ID
	p.Height = leaguePlayer.PLAYER_HEIGHT
	p.HeightInches = int(leaguePlayer.PLAYER_HEIGHT_INCHES)
	weight, err := strconv.Atoi(leaguePlayer.PLAYER_WEIGHT)
	if err != nil {
		weight = 0
		logrus.Errorf("Error converting weight to int for player %s: %v", leaguePlayer.PLAYER_NAME, err)
	}
	p.Weight = weight
	p.League = leaguePlayer.League
	return p
}
