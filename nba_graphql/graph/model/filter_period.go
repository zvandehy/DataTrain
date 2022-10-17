package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
	"go.mongodb.org/mongo-driver/bson"
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

func (p *Period) MongoFilter() bson.M {
	filter := map[string]interface{}{}
	if p.StartDate != nil {
		filter["date"] = map[string]interface{}{
			"$gte": *p.StartDate,
		}
	}
	if p.EndDate != nil {
		if _, ok := filter["date"]; !ok {
			filter["date"] = map[string]interface{}{}
		}
		filter["date"].(map[string]interface{})["$lte"] = *p.EndDate
	}
	if p.Seasons != nil && len(*p.Seasons) > 0 {
		filter["season"] = map[string]interface{}{
			"$in": *p.Seasons,
		}
	}
	return filter
}

func (p *Period) MatchGame(game *PlayerGame) bool {
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

func (p *Period) MatchProjection(projection *Projection) bool {
	date, err := time.Parse("2006-01-02", projection.Date)
	if err != nil {
		logrus.Error("Period.Match: invalid projection date")
		return false
	}
	if p.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *p.StartDate)
		if err != nil {
			logrus.Error("Period.Match: invalid start date")
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
		if date.After(endDate) {
			return false
		}
	}
	projectionSeason := SEASON_2022_23
	if date.After(util.SEASON_DATE(util.SEASON_START_2020_21)) && date.Before(util.SEASON_DATE(util.SEASON_END_2021_22)) {
		projectionSeason = SEASON_2021_22
	} else if date.After(util.SEASON_DATE(util.SEASON_START_2020_21)) && date.Before(util.SEASON_DATE(util.SEASON_END_2020_21)) {
		projectionSeason = SEASON_2020_21
	}

	if p.Seasons != nil && len(*p.Seasons) > 0 {
		inSeason := false
		for _, season := range *p.Seasons {
			if projectionSeason == season {
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
