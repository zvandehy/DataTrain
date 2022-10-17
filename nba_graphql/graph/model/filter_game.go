package model

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

type GameFilter struct {
	TeamID          *int            `json:"teamID"`
	OpponentID      *int            `json:"opponentID"`
	OpponentMatch   *bool           `json:"opponentMatch"`
	PlayerID        *int            `json:"playerID"`
	GameID          *string         `json:"gameID"`
	Seasons         *[]SeasonOption `json:"seasons"`
	StartDate       *string         `json:"startDate"`
	EndDate         *string         `json:"endDate"`
	GameType        *GameType       `json:"gameType"`
	GameTypeMatch   *bool           `json:"gameTypeMatch"`
	HomeOrAway      *HomeOrAway     `json:"homeOrAway"`
	HomeOrAwayMatch *bool           `json:"homeOrAwayMatch"`
	StatFilters     *[]*StatFilter  `json:"statFilters"`
	LastX           *int            `json:"lastX"`
}

func (f GameFilter) String() string {
	return util.Print(f)
}

// Match returns true if the player game matches the given filter
func (gameFilter *GameFilter) MatchPlayerGame(g *PlayerGame) bool {
	// TODO: Would be more maintainable to use the same Match function for any SeasonOption filter
	if gameFilter.Seasons != nil && len(*gameFilter.Seasons) > 0 {
		found := false
		for _, season := range *gameFilter.Seasons {
			if g.Season == season {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if gameFilter.GameType != nil {
		if *gameFilter.GameType == "Playoffs" && !g.Playoffs {
			return false
		}
		if *gameFilter.GameType == "Regular Season" && g.Playoffs {
			return false
		}
	}
	if gameFilter.PlayerID != nil && *gameFilter.PlayerID != g.PlayerID {
		return false
	}
	if gameFilter.GameID != nil && *gameFilter.GameID != g.GameID {
		return false
	}
	if gameFilter.TeamID != nil && *gameFilter.TeamID != g.TeamID {
		return false
	}
	if gameFilter.OpponentID != nil && *gameFilter.OpponentID != g.OpponentID {
		return false
	}
	if gameFilter.HomeOrAway != nil && !strings.EqualFold(string(*gameFilter.HomeOrAway), g.HomeOrAway) {
		return false
	}
	date, err := time.Parse("2006-01-02", g.Date)
	if err != nil {
		logrus.Error("Error parsing game date: ", err)
		return false
	}
	if gameFilter.StartDate != nil {
		start, err := time.Parse("2006-01-02", *gameFilter.StartDate)
		if err != nil {
			logrus.Error("Error parsing filter startDate: ", err)
			return false
		}
		if start.After(date) {
			return false
		}
	}
	if gameFilter.EndDate != nil {
		end, err := time.Parse("2006-01-02", *gameFilter.EndDate)
		if err != nil {
			logrus.Error("Error parsing filter endDate: ", err)
			return false
		}
		// match the date only if the endfilter is after and not equal to the date
		if end.Before(date) || end.Equal(date) {
			return false
		}
	}
	if gameFilter.StatFilters != nil && len(*gameFilter.StatFilters) > 0 {
		for _, statFilter := range *gameFilter.StatFilters {
			if !statFilter.MatchGame(g) {
				return false
			}
		}
	}
	return true
}
