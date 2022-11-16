package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Projection struct {
	PlayerName  string         `json:"playername" bson:"playername"`
	OpponentAbr string         `json:"opponent" bson:"opponent"` //TODO: Refactor to OpponentAbr
	Props       []*Proposition `json:"propositions" bson:"propositions"`
	StartTime   string         `json:"startTime" bson:"startTime"`
	Date        string         `json:"date" bson:"date"`

	Player   Player `json:"playerCache" bson:"playerCache"`
	Opponent Team   `json:"opponentTeam" bson:"opponentTeam"`
}

func (p *Projection) UnmarshalBSON(data []byte) error {
	type Alias Projection
	bson.Unmarshal(data, (*Alias)(p))
	for i, prop := range p.Props {
		// prop.ProjectionRef = p
		p.Props[i] = prop
	}
	// p.PlayerMatrix = v.PlayerMatrix
	return nil
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
func ParsePrizePick(prop PrizePicksData, itemIDToNameMap map[string]string, projections []*Projection) ([]*Projection, error) {
	var playerName string
	var statType Stat

	// get playername and statType from id to name mapping
	if val, ok := itemIDToNameMap[prop.Relationships.Player.Data.ID]; ok {
		playerName = strings.TrimSpace(val)
	} else {
		return nil, errors.New("could not find player name")
	}
	if val, ok := itemIDToNameMap[prop.Relationships.StatType.Data.ID]; ok {
		var err error
		statType, err = NewStat(val)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("could not find stat type")
	}

	if playerName == "" {
		err := fmt.Errorf("error retrieving prizepick player name")
		logrus.Error(err)
		return nil, err
	}
	if statType == "" {
		err := fmt.Errorf("error retrieving prizepick stat type for player %s | %+v", playerName, prop)
		logrus.Error(err)
		return nil, err
	}

	target, err := strconv.ParseFloat(prop.Attributes.Line_score, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve prizepicks target")
	}

	dateSlice := strings.SplitN(prop.Attributes.Start_time, "T", 2)
	date := dateSlice[0]
	now := time.Now()

	// if player is already in list of projections, just add a proposition (for this prop type) to their projection
	for i, projection := range projections {
		if projection.PlayerName == playerName {
			if prop.Attributes.Is_promo {
				logrus.Warn("skipping promo prizepick")
				continue
			}
			projections[i].Props = append(projections[i].Props, &Proposition{Sportsbook: "PrizePicks", Target: target, Type: statType, LastModified: &now})
			return projections, nil
		}
	}

	// otherwise, create a new player projection with this proposition (prop type) in it
	projections = append(projections, &Projection{PlayerName: playerName, OpponentAbr: prop.Attributes.Description, Date: date, StartTime: prop.Attributes.Start_time, Props: []*Proposition{{Sportsbook: "PrizePicks", Target: target, Type: statType, LastModified: &now}}})
	return projections, nil
}

type UnderdogFantasy struct {
	Players        []UnderdogPlayer        `json:"players"`
	Appearances    []UnderdogAppearance    `json:"appearances"`
	OverUnderLines []UnderdogOverUnderLine `json:"over_under_lines"`
	Games          []UnderdogGame          `json:"games"`
}

type UnderdogPlayer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PlayerID  string `json:"id"`
	Sport     string `json:"sport_id"`
	TeamID    string `json:"team_id"`
}

type UnderdogAppearance struct {
	PlayerID     string `json:"player_id"`
	AppearanceID string `json:"id"`
}

type UnderdogOverUnderLine struct {
	OverUnder struct {
		Appearance struct {
			AppearanceID string `json:"appearance_id"`
			Category     string `json:"stat"`
		} `json:"appearance_stat"`
	} `json:"over_under"`

	Target string `json:"stat_value"`
}

type UnderdogGame struct {
	AwayTeamID string `json:"away_team_id"`
	HomeTeamID string `json:"home_team_id"`
	Title      string `json:"title"`
	StartTime  string `json:"scheduled_at"`
}

func ParseUnderdogProjection(json UnderdogFantasy, sport string) ([]*Projection, error) {
	now := time.Now()
	var projections []*Projection
	for _, player := range json.Players {
		if strings.EqualFold(player.Sport, sport) {
			continue
		}
		playername := fmt.Sprintf("%s %s", player.FirstName, player.LastName)
		playername = strings.TrimSpace(playername)
		playername = strings.Replace(playername, "  ", " ", -1)
		if val, ok := PlayerNames[playername]; ok {
			playername = val
		}
		opponent, game, err := getOpponent(json.Games, player.TeamID)
		if err != nil {
			return nil, err
		}
		opponent = getAbbreviation(opponent)
		startTime, err := time.Parse("2006-01-02T15:04:05Z", game.StartTime)
		startTime = startTime.Add(time.Hour * -6)
		start := startTime.Format("2006-01-02T15:04:05Z")
		dateSlice := strings.SplitN(start, "T", 2)
		date := dateSlice[0]

		appearanceID, err := getAppearanceIDForPlayer(player.PlayerID, json.Appearances)
		if err != nil {
			return nil, err
		}
		var propositions []*Proposition
		for _, overUnder := range json.OverUnderLines {
			if overUnder.OverUnder.Appearance.AppearanceID == appearanceID {
				target, err := strconv.ParseFloat(overUnder.Target, 64)
				if err != nil {
					return nil, fmt.Errorf("couldn't get target for %v / %v", playername, overUnder.OverUnder.Appearance.Category)
				}
				category, ok := categories[overUnder.OverUnder.Appearance.Category]
				if !ok {
					category = overUnder.OverUnder.Appearance.Category
				}
				stat, err := NewStat(category)
				if err != nil {
					return nil, err
				}
				proposition := Proposition{
					Sportsbook:   "UnderdogFantasy",
					LastModified: &now,
					Target:       target,
					Type:         stat,
				}
				propositions = append(propositions, &proposition)
			}

		}
		projections = append(projections, &Projection{
			PlayerName:  playername,
			OpponentAbr: opponent,
			StartTime:   start,
			Date:        date,
			Props:       propositions,
		})
	}
	return projections, nil
}

var abbreviations = map[string]string{
	"PHX":  "PHO",
	"LA":   "LAS",
	"WSH":  "WAS",
	"NY":   "NYL",
	"CONN": "CON",
	"LV":   "LVA",
}

func getAbbreviation(opp string) string {
	if abbr, ok := abbreviations[opp]; ok {
		return abbr
	}
	return opp
}

var categories map[string]string = map[string]string{
	"pts_rebs_asts":     "Pts+Rebs+Asts",
	"points":            "Points",
	"assists":           "Assists",
	"three_points_made": "3-PT Made",
	"rebounds":          "Rebounds",
}

func getAppearanceIDForPlayer(playerID string, appearances []UnderdogAppearance) (string, error) {
	for _, appearance := range appearances {
		if appearance.PlayerID == playerID {
			return appearance.AppearanceID, nil
		}
	}
	return "", fmt.Errorf("Couldn't find appearance for playerID: '%s'", playerID)
}

func getOpponent(underdogGames []UnderdogGame, teamID string) (string, *UnderdogGame, error) {
	for _, game := range underdogGames {
		if game.AwayTeamID == teamID {
			return strings.SplitN(game.Title, " @ ", 2)[1], &game, nil
		}
		if game.HomeTeamID == teamID {
			return strings.SplitN(game.Title, " @ ", 2)[0], &game, nil
		}
	}
	return "", nil, fmt.Errorf("Couldn't find opponent for team: '%s'", teamID)
}

func GetBestProjection(projections []*Projection) *Projection {
	maxTargets := 0
	var bestProjections []*Projection
	for _, projection := range projections {
		if len(projection.Props) > maxTargets {
			maxTargets = len(projection.Props)
			bestProjections = []*Projection{projection}
		} else if len(projection.Props) == maxTargets {
			bestProjections = append(bestProjections, projection)
		}
	}
	return bestProjections[len(bestProjections)-1]
}
