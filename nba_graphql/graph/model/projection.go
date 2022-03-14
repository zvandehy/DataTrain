package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Projection struct {
	PlayerName  string    `json:"playername" bson:"playername"`
	Sportsbook  string    `json:"sportsbook" bson:"sportsbook"`
	OpponentAbr string    `json:"opponent" bson:"opponent"`
	Targets     []*Target `json:"targets" bson:"targets"`
	StartTime   string    `json:"startTime" bson:"startTime"`
	Date        string    `json:"date" bson:"date"`
}

type PrizePicks struct {
	Data     []PrizePicksData     `json:"data" bson:"data"`
	Included []PrizePicksIncluded `json:"included" bson:"included"`
}

type PrizePicksData struct {
	//unique identifier of the projection
	PrizePickID string `json:"id" bson:"id"`
	Attributes  struct {
		//line available to be played on the board (i.e. 45.5), associated w/ a player and stat type
		Line_score string `json:"line_score" bson:"line_score"`
		//
		Risk_exposure float64 `json:"risk_exposure" bson:"risk_exposure"`
		//time the projection was created at
		Created_at string `json:"created_at" bson:"created_at"`
		//last time the projection was updated at
		Updated_at string `json:"updated_at" bson:"updated_at"`
		//the opponent of the player in the event offered, varies by league (i.e. for NFL the opposing team abbreviation)
		Description string `json:"description" bson:"description"`
		//time the event is expected to start
		Start_time string `json:"start_time" bson:"start_time"`
		// event states (i.e. pre-game, in-progress, final)
		Status string `json:"status" bson:"status"`
		//whether or not a projection is promotional in nature
		Is_promo            bool   `json:"is_promo" bson:"is_promo"`
		Discount_percentage string `json:"discount_percentage" bson:"discount_percentage"`
		//time the event is expected to end
		End_time string `json:"end_time" bson:"end_time"`
		Position struct {
			Data string `json:"data" bson:"data"`
		} `json:"position" bson:"position"`
	} `json:"attributes" bson:"attributes"`
	Relationships struct {
		Player struct {
			Data struct {
				ID string `json:"id" bson:"id"`
			} `json:"data" bson:"data"`
		} `json:"new_player" bson:"new_player"`
		StatType struct {
			Data struct {
				ID string `json:"id" bson:"id"`
			} `json:"data" bson:"data"`
		} `json:"stat_type" bson:"stat_type"`
	} `json:"relationships" bson:"relationships"`
}

type PrizePicksIncluded struct {
	ID         string `json:"id" bson:"id"`
	Attributes struct {
		Name string `json:"name" bson:"name"`
	} `json:"attributes" bson:"attributes"`
}

//ParsePrizePick creates a Projection and adds it to the projections slice or adds a Target to an existing projection
func ParsePrizePick(prop PrizePicksData, included []PrizePicksIncluded, projections []*Projection) ([]*Projection, error) {
	var playerName string
	var statType string
	for _, p := range included {
		if p.ID == prop.Relationships.Player.Data.ID {
			playerName = p.Attributes.Name
		}
		if p.ID == prop.Relationships.StatType.Data.ID {
			statType = p.Attributes.Name
		}
		if statType != "" && playerName != "" {
			break
		}
	}
	if playerName == "" {
		return nil, fmt.Errorf("error retrieving prizepick player name")
	}
	if statType == "" {
		return nil, fmt.Errorf("error retrieving prizepick stat type")
	}
	target, err := strconv.ParseFloat(prop.Attributes.Line_score, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve prizepicks target")
	}

	dateSlice := strings.SplitN(prop.Attributes.Start_time, "T", 2)
	date := dateSlice[0]
	for i, projection := range projections {
		if projection.PlayerName == playerName {
			projections[i].Targets = append(projections[i].Targets, &Target{Target: target, Type: statType})
			return projections, nil
		}
	}

	projections = append(projections, &Projection{Sportsbook: "PrizePicks", PlayerName: playerName, OpponentAbr: prop.Attributes.Description, Date: date, StartTime: prop.Attributes.Start_time, Targets: []*Target{{Target: target, Type: statType}}})
	return projections, nil
}
