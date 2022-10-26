package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerFilter struct {
	Name                *string           `json:"name"`
	PlayerID            *int              `json:"playerID"`
	Seasons             *[]SeasonOption   `json:"seasons"`
	PositionStrict      *Position         `json:"positionStrict"`
	PositionStrictMatch *bool             `json:"positionStrictMatch"`
	PositionLoose       *Position         `json:"positionLoose"`
	PositionLooseMatch  *bool             `json:"positionLooseMatch"`
	TeamAbr             *string           `json:"teamABR"`
	TeamID              *int              `json:"teamID"`
	StartDate           *string           `json:"startDate"`
	EndDate             *string           `json:"endDate"`
	StatFilters         *[]*StatFilter    `json:"statFilters"`
	WithPropositions    *ProjectionFilter `json:"withPropositions"`
	WithGames           *GameFilter       `json:"withGames"`
}

func (input *PlayerFilter) MongoPipeline() mongo.Pipeline {
	orFilters := []bson.M{}
	// TODO: user could provide their own position rules
	positionFilter := bson.M{"position": bson.M{"$in": []string{"G", "F", "C", "G-F", "F-G", "F-C", "C-F"}}}
	if input.PositionStrict != nil {
		switch *input.PositionStrict {
		case PositionG:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"G"}}}
		case PositionF:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"F"}}}
		case PositionC:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"C"}}}
		case PositionFG:
			fallthrough
		case PositionGF:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"G-F", "F-G"}}}
		case PositionCF:
			fallthrough
		case PositionFC:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"F-C", "C-F"}}}
		}
		orFilters = append(orFilters, positionFilter)
	} else if input.PositionLoose != nil {
		switch *input.PositionLoose {
		case PositionG:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"G", "G-F", "F-G"}}}
		case PositionF:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"F", "F-G", "F-C", "G-F", "C-F"}}}
		case PositionC:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"C", "C-F", "F-C"}}}
		case PositionFG:
			fallthrough
		case PositionGF:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"G-F", "F-G", "G", "F"}}}
		case PositionCF:
			fallthrough
		case PositionFC:
			positionFilter = bson.M{"position": bson.M{"$in": []string{"F-C", "C-F", "F", "C"}}}
		}
		orFilters = append(orFilters, positionFilter)
	}

	if input.Name != nil {
		orFilters = append(orFilters, bson.M{"name": bson.M{"$regex": *input.Name, "$options": "i"}})
	}

	if input.PlayerID != nil {
		orFilters = append(orFilters, bson.M{"playerID": *input.PlayerID})
	}

	if input.TeamAbr != nil {
		orFilters = append(orFilters, bson.M{"$regex": *input.TeamAbr, "$options": "i"})
	}

	var andFilters bson.A
	if input.Seasons != nil {
		andFilters = append(andFilters, bson.M{"seasons": bson.M{"$in": *input.Seasons}})
	}

	if len(orFilters) > 0 {
		andFilters = append(andFilters, bson.M{"$or": orFilters})
	}

	var filter bson.M
	if len(andFilters) > 0 {
		if len(andFilters) == 1 {
			if _, ok := andFilters[0].(bson.M)["seasons"]; ok {
				andFilters = append(andFilters, bson.M{"exclude": false})
			}
		}
		filter = bson.M{
			"$and": andFilters,
		}
	} else {
		filter = bson.M{}
	}

	lookupGames := bson.M{
		"from":         "games",
		"localField":   "playerID",
		"foreignField": "playerID",
		"as":           "gamesCache",
	}

	pipeline := mongo.Pipeline{
		bson.D{primitive.E{Key: "$match", Value: filter}},
		// TODO: players' games are not filtered at all
		// if the player has many games from previous seasons that the query doesn't care about,
		// then this might be unnecessarily slow
		bson.D{primitive.E{Key: "$lookup", Value: lookupGames}},
	}
	if input.WithPropositions != nil {
		projectionFilter := input.WithPropositions.MongoFilter()
		projectionFilter["projections.date"] = projectionFilter["date"]
		projectionFilter["date"] = nil
		lookupProjections := bson.M{
			"from":         "projections",
			"localField":   "name",
			"foreignField": "playername",
			"as":           "projections",
		}
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$lookup", Value: lookupProjections}})
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$match", Value: projectionFilter}})
	}
	logrus.Warn(pipeline)
	return pipeline
}

func (input PlayerFilter) FilterPlayerStats(players []*Player, toPlayer *Player) []*Player {
	// TODO: Filter other PlayerFilter fields...
	if input.StatFilters != nil && len(*input.StatFilters) > 0 {
		filteredPlayers := []*Player{}
		for _, player := range players {
			matches := true
			if toPlayer != nil {
				if input.PositionStrictMatch != nil && *input.PositionStrictMatch {
					switch toPlayer.Position {
					case "G":
						if player.Position != "G" {
							matches = false
							break
						}
					case "F":
						if player.Position != "F" {
							matches = false
							break
						}
					case "C":
						if player.Position != "C" {
							matches = false
							break
						}
					case "F-G":
						fallthrough
					case "G-F":
						if player.Position != "F-G" && player.Position != "G-F" {
							matches = false
							break
						}
					case "F-C":
						fallthrough
					case "C-F":
						if player.Position != "F-C" && player.Position != "C-F" {
							matches = false
							break
						}
					}
				} else if input.PositionLooseMatch != nil && *input.PositionLooseMatch {
					switch toPlayer.Position {
					case "G":
						if player.Position != "G" && player.Position != "F-G" && player.Position != "G-F" {
							matches = false
							break
						}
					case "F":
						if player.Position != "F" && player.Position != "F-G" && player.Position != "F-C" && player.Position != "C-F" && player.Position != "G-F" {
							matches = false
							break
						}
					case "C":
						if player.Position != "C" && player.Position != "F-C" && player.Position != "C-F" {
							matches = false
							break
						}
					case "F-G":
						fallthrough
					case "G-F":
						if player.Position != "G" && player.Position != "F-G" && player.Position != "G-F" && player.Position != "F" {
							matches = false
							break
						}
					case "F-C":
						fallthrough
					case "C-F":
						if player.Position != "C" && player.Position != "F-C" && player.Position != "C-F" && player.Position != "F" {
							matches = false
							break
						}
					}
				}
			}
			for _, statFilter := range *input.StatFilters {
				if !statFilter.MatchPlayer(player) {
					matches = false
					break
				}
			}
			if matches {
				filteredPlayers = append(filteredPlayers, player)
			}
		}
		return filteredPlayers
	}
	return players
}

func (f PlayerFilter) String() string {
	return util.Print(f)
}

func (f PlayerFilter) Key() string {
	name := ""
	if f.Name != nil {
		name = *f.Name
	}
	playerID := ""
	if f.PlayerID != nil {
		playerID = strconv.Itoa(*f.PlayerID)
	}
	seasons := fmt.Sprintf("%v", f.Seasons)
	positionStrict := ""
	if f.PositionStrict != nil {
		positionStrict = string(*f.PositionStrict)
	}
	positionLoose := ""
	if f.PositionLoose != nil {
		positionLoose = string(*f.PositionLoose)
	}
	teamAbr := ""
	if f.TeamAbr != nil {
		teamAbr = *f.TeamAbr
	}
	teamID := ""
	if f.TeamID != nil {
		teamID = strconv.Itoa(*f.TeamID)
	}
	start := ""
	if f.StartDate != nil {
		start = *f.StartDate
	}
	end := ""
	if f.EndDate != nil {
		end = *f.EndDate
	}
	// statFilters := ""
	// if f.StatFilters != nil {
	// 	for _, statFilter := range *f.StatFilters {
	// 		statFilters += fmt.Sprintf("per:%v,stat:%v,mode:%v,op:%v,val:%v", statFilter.Period, statFilter.Stat, statFilter.Mode, statFilter.Operator, statFilter.Value)
	// 	}
	// }
	return fmt.Sprintf("Name:%v,PlayerID:%v,Seasons:%v,PositionStrict:%v,PositionLoose:%v,TeamAbr:%v,TeamID:%v,StartDate:%v,EndDate:%v",
		name,
		playerID,
		seasons,
		positionStrict,
		positionLoose,
		teamAbr,
		teamID,
		start,
		end,
	)
}

func containsSeason(seasons *[]SeasonOption, season SeasonOption) bool {
	for _, s := range *seasons {
		if s == season {
			return true
		}
	}
	return false
}

func (f PlayerFilter) GetEarliestSeasonStartDate() (*time.Time, error) {
	startDate := util.SEASON_START_2020_21
	if f.Seasons != nil {
		if containsSeason(f.Seasons, SEASON_2020_21) {
			startDate = util.SEASON_START_2020_21
		} else if containsSeason(f.Seasons, SEASON_2021_22) {
			startDate = util.SEASON_START_2021_22
		} else if containsSeason(f.Seasons, SEASON_2022_23) {
			startDate = util.SEASON_START_2022_23
		}
	}
	start, err := time.Parse(util.DATE_FORMAT, startDate)
	if err != nil {
		logrus.Errorf("Error parsing game date %v", startDate)
		return nil, fmt.Errorf("error parsing game date %v", startDate)
	}
	return &start, nil
}
