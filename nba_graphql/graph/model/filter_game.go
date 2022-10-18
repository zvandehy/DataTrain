package model

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Outcome         *Outcome        `json:"outcome"`
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
	if gameFilter.Outcome != nil {
		outcome := strings.ToLower(g.Outcome)
		if *gameFilter.Outcome == OutcomeWin && outcome[0] != 'w' {
			return false
		}
		if *gameFilter.Outcome == OutcomeLoss && outcome[0] != 'l' {
			return false
		}
		if *gameFilter.Outcome == OutcomePending && outcome[0] != 'p' {
			return false
		}
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

func (f *GameFilter) MongoPipeline(playerPipeline mongo.Pipeline) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		//match games
		bson.D{primitive.E{Key: "$match", Value: f.MongoFilter()}},
		//find players for each matched game
		// TODO: what if a player has multiple games that match the filter?
		bson.D{primitive.E{Key: "$lookup", Value: bson.M{
			"from":         "players",
			"localField":   "playerID",
			"foreignField": "playerID",
			"as":           "player",
		}}},
		bson.D{primitive.E{Key: "$unwind", Value: "$player"}},
		bson.D{primitive.E{Key: "$group", Value: bson.M{
			"_id":              "$playerID",
			"playerID":         bson.M{"$first": "$player.playerID"},
			"name":             bson.M{"$first": "$player.name"},
			"first_name":       bson.M{"$first": "$player.first_name"},
			"last_name":        bson.M{"$first": "$player.last_name"},
			"position":         bson.M{"$first": "$player.position"},
			"seasons":          bson.M{"$first": "$player.seasons"},
			"teamABR":          bson.M{"$first": "$player.teamABR"},
			"height":           bson.M{"$first": "$player.height"},
			"weight":           bson.M{"$first": "$player.weight"},
			"lastModifiedDate": bson.M{"$first": "$player.lastModifiedDate"},
			"league":           bson.M{"$first": "$player.league"},
		}}},
	}
	pipeline = append(pipeline, playerPipeline...)
	return pipeline
}

func (f *GameFilter) MongoFilter() bson.M {
	matchGame := bson.M{}
	if f.Seasons != nil && len(*f.Seasons) > 0 {
		matchGame["season"] = bson.M{"$in": *f.Seasons}
	}
	if f.GameType != nil {
		if *f.GameType == GameTypePlayoffs {
			matchGame["playoffs"] = true
		}
		if *f.GameType == GameTypeRegularSeason {
			matchGame["playoffs"] = false
		}
	}
	if f.PlayerID != nil {
		matchGame["playerID"] = *f.PlayerID
	}
	if f.GameID != nil {
		matchGame["gameID"] = *f.GameID
	}
	if f.TeamID != nil {
		matchGame["teamID"] = *f.TeamID
	}
	if f.OpponentID != nil {
		matchGame["opponentID"] = *f.OpponentID
	}
	if f.HomeOrAway != nil {
		matchGame["home_or_away"] = *f.HomeOrAway
	}
	if f.StartDate != nil && f.EndDate == nil {
		matchGame["date"] = bson.M{"$gte": *f.StartDate}
	}
	if f.StartDate == nil && f.EndDate != nil {
		matchGame["date"] = bson.M{"$lt": *f.EndDate}
	}
	if f.StartDate != nil && f.EndDate != nil {
		if *f.StartDate == *f.EndDate {
			matchGame["date"] = *f.StartDate
		} else {
			matchGame["date"] = bson.M{"$gte": *f.StartDate, "$lt": *f.EndDate}
		}
	}
	if f.Outcome != nil {
		if *f.Outcome == OutcomeWin {
			matchGame["win_or_loss"] = bson.M{"$regex": "win", "$options": "i"}
		}
		if *f.Outcome == OutcomeLoss {
			matchGame["win_or_loss"] = bson.M{"$regex": "loss", "$options": "i"}
		}
		if *f.Outcome == OutcomePending {
			matchGame["win_or_loss"] = bson.M{"$regex": "pending", "$options": "i"}
		}
	}
	return matchGame
}
