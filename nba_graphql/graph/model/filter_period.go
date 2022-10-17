package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// TODO: Add GameType to Period
type Period struct {
	StartDate *string         `json:"startDate"`
	EndDate   *string         `json:"endDate"`
	Seasons   *[]SeasonOption `json:"seasons"`
	Limit     *int            `json:"limit"`
}

func (p Period) String() string {
	return util.Print(p)
}

func (p *Period) Match(game *PlayerGame) bool {
	if p.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *p.StartDate)
		if err != nil {
			logrus.Error("Period.Match: invalid start date")
			return false
		}
		date, err := time.Parse("2006-01-02", game.Date)
		if err != nil {
			logrus.Error("Period.Match: invalid game date")
			return false
		}
		if date.Before(startDate) {
			return false
		}
	}
	if p.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *p.EndDate)
		if err != nil {
			logrus.Error("Period.Match: invalid end date")
			return false
		}
		date, err := time.Parse("2006-01-02", game.Date)
		if err != nil {
			logrus.Error("Period.Match: invalid game date")
			return false
		}
		if date.After(endDate) {
			return false
		}
	}
	if p.Seasons != nil && len(*p.Seasons) > 0 {
		inSeason := false
		for _, season := range *p.Seasons {
			if game.Season == season {
				inSeason = true
				break
			}
		}
		if !inSeason {
			return false
		}
	}
	return true
}
