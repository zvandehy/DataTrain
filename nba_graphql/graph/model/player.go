package model

import (
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
	similarity "github.com/zvandehy/DataTrain/nba_graphql/math"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// If key is found, use the value.
var PlayerNames map[string]string = map[string]string{
	"Nah'Shon Hyland":     "Bones Hyland",
	"Ty Harris":           "Tyasha Harris",
	"Naz Hillman":         "Naz Hillmon",
	"Olivia Nelson-Odada": "Olivia Nelson-Ododa",
}

type Player struct {
	FirstName   string   `json:"first_name" bson:"first_name"`
	LastName    string   `json:"last_name" bson:"last_name"`
	Name        string   `json:"name" bson:"name"`
	PlayerID    int      `json:"playerID" bson:"playerID"`
	Seasons     []string `json:"seasons" bson:"seasons"`
	Position    string   `json:"position" bson:"position"`
	CurrentTeam string   `json:"currentTeam" bson:"teamABR"`
	Height      string   `json:"height" bson:"height"`
	Weight      int      `json:"weight" bson:"weight"`
	// When retrieving a player, also retrieve all of the games they've played within the minimum start date and maximum end date.
	GamesCache []*PlayerGame `json:"gamesCache" bson:"gamesCache"`
}

func NewPlayerAverage(games []*PlayerGame, player *Player) PlayerAverage {
	average := PlayerAverage{}
	average.GamesPlayed = float64(len(games))
	average.Player = *player
	average.Height = float64(player.HeightInInches())
	average.Weight = float64(player.Weight)

	for _, game := range games {
		min, err := ConvertMinutesToFloat(game.Minutes)
		if err != nil {
			min = 0
		}
		average.Assists += float64(game.Assists)
		average.Blocks += float64(game.Blocks)
		average.DefensiveRebounds += float64(game.DefensiveRebounds)
		average.FieldGoalsAttempted += float64(game.FieldGoalsAttempted)
		average.FieldGoalsMade += float64(game.FieldGoalsMade)
		average.FreeThrowsAttempted += float64(game.FreeThrowsAttempted)
		average.FreeThrowsMade += float64(game.FreeThrowsMade)
		average.Minutes += min
		average.OffensiveRebounds += float64(game.OffensiveRebounds)
		average.PersonalFoulsDrawn += float64(game.PersonalFoulsDrawn)
		average.PersonalFouls += float64(game.PersonalFouls)
		average.Points += float64(game.Points)
		average.Rebounds += float64(game.Rebounds)
		average.Steals += float64(game.Steals)
		average.ThreePointersAttempted += float64(game.ThreePointersAttempted)
		average.ThreePointersMade += float64(game.ThreePointersMade)
		average.Turnovers += float64(game.Turnovers)
	}

	average.Assists /= float64(len(games))
	average.Blocks /= float64(len(games))
	average.DefensiveRebounds /= float64(len(games))
	average.FieldGoalsAttempted /= float64(len(games))
	average.FieldGoalsMade /= float64(len(games))
	average.FreeThrowsAttempted /= float64(len(games))
	average.FreeThrowsMade /= float64(len(games))
	average.Minutes /= float64(len(games))
	average.OffensiveRebounds /= float64(len(games))
	average.PersonalFoulsDrawn /= float64(len(games))
	average.PersonalFouls /= float64(len(games))
	average.Points /= float64(len(games))
	average.Rebounds /= float64(len(games))
	average.Steals /= float64(len(games))
	average.ThreePointersAttempted /= float64(len(games))
	average.ThreePointersMade /= float64(len(games))
	average.Turnovers /= float64(len(games))

	return average
}

func (p *Player) HeightInInches() int {
	var feet, inches int
	fmt.Sscanf(p.Height, "%d-%d", &feet, &inches)
	return feet*12 + inches
}

func (p Player) String() string {
	return util.Print(p)
}

func (p PlayerGame) String() string {
	return util.Print(p)
}

// TODO: To add a new stat to player similarity, add it to all of: ... // TODO: look into this, see if there is a more maintainable way to do this
type PlayerAverage struct {
	AllMinutes             []string `json:"all_minutes" bson:"all_minutes"`
	Assists                float64  `json:"assists" bson:"assists"`
	Blocks                 float64  `json:"blocks" bson:"blocks"`
	DefensiveRebounds      float64  `json:"defensive_rebounds" bson:"defensive_rebounds"`
	FieldGoalsAttempted    float64  `json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade         float64  `json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted    float64  `json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade         float64  `json:"free_throws_made" bson:"free_throws_made"`
	GamesPlayed            float64  `json:"games_played" bson:"games_played"`
	Height                 float64  `json:"height" bson:"height"`
	Minutes                float64  `json:"minutes" bson:"minutes"`
	OffensiveRebounds      float64  `json:"offensive_rebounds" bson:"offensive_rebounds"`
	Player                 Player   `json:"player" bson:"player"`
	PersonalFoulsDrawn     float64  `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls          float64  `json:"personal_fouls" bson:"personal_fouls"`
	Points                 float64  `json:"points" bson:"points"`
	Rebounds               float64  `json:"rebounds" bson:"rebounds"`
	Steals                 float64  `json:"steals" bson:"steals"`
	ThreePointersAttempted float64  `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade      float64  `json:"three_pointers_made" bson:"three_pointers_made"`
	Turnovers              float64  `json:"turnovers" bson:"turnovers"`
	Weight                 float64  `json:"weight" bson:"weight"`
}

func (p *PlayerAverage) Normalize(stats ...StatOfInterest) PlayerAverage {
	normalized := *p
	for _, stat := range stats {
		switch stat.Stat {
		case Points:
			normalized.Points = stat.ZScore(*p)
		case Rebounds:
			normalized.Rebounds = stat.ZScore(*p)
		case Assists:
			normalized.Assists = stat.ZScore(*p)
		case Steals:
			normalized.Steals = stat.ZScore(*p)
		case Blocks:
			normalized.Blocks = stat.ZScore(*p)
		case FieldGoalsAttempted:
			normalized.FieldGoalsAttempted = stat.ZScore(*p)
		case FieldGoalsMade:
			normalized.FieldGoalsMade = stat.ZScore(*p)
		case ThreePointersAttempted:
			normalized.ThreePointersAttempted = stat.ZScore(*p)
		case ThreePointersMade:
			normalized.ThreePointersMade = stat.ZScore(*p)
		case FreeThrowsAttempted:
			normalized.FreeThrowsAttempted = stat.ZScore(*p)
		case FreeThrowsMade:
			normalized.FreeThrowsMade = stat.ZScore(*p)
		case OffensiveRebounds:
			normalized.OffensiveRebounds = stat.ZScore(*p)
		case DefensiveRebounds:
			normalized.DefensiveRebounds = stat.ZScore(*p)
		case PersonalFouls:
			normalized.PersonalFouls = stat.ZScore(*p)
		case PersonalFoulsDrawn:
			normalized.PersonalFoulsDrawn = stat.ZScore(*p)
		case Turnovers:
			normalized.Turnovers = stat.ZScore(*p)
		case GamesPlayed:
			normalized.GamesPlayed = stat.ZScore(*p)
		case Minutes:
			normalized.Minutes = stat.ZScore(*p)
		case Height:
			normalized.Height = stat.ZScore(*p)
		case Weight:
			normalized.Weight = stat.ZScore(*p)
		}
	}
	return normalized
}

type PlayerDiff PlayerAverage

// type PlayerDiff struct {
// 	PlayerAverage
// 	Distance float64
// }

func (p *PlayerAverage) AverageMinutes() (float64, error) {
	var minutes float64
	for _, str := range p.AllMinutes {
		f, err := ConvertMinutesToFloat(str)
		if err != nil {
			return 0, nil
		}
		minutes += f
	}
	return minutes, nil
}

func (p *PlayerAverage) Difference(fromPlayer PlayerAverage) PlayerDiff {
	d := PlayerDiff{
		Assists:                similarity.RoundFloat(fromPlayer.Assists-p.Assists, 2),
		Blocks:                 similarity.RoundFloat(fromPlayer.Blocks-p.Blocks, 2),
		DefensiveRebounds:      similarity.RoundFloat(fromPlayer.DefensiveRebounds-p.DefensiveRebounds, 2),
		FieldGoalsAttempted:    similarity.RoundFloat(fromPlayer.FieldGoalsAttempted-p.FieldGoalsAttempted, 2),
		FieldGoalsMade:         similarity.RoundFloat(fromPlayer.FieldGoalsMade-p.FieldGoalsMade, 2),
		FreeThrowsAttempted:    similarity.RoundFloat(fromPlayer.FreeThrowsAttempted-p.FreeThrowsAttempted, 2),
		FreeThrowsMade:         similarity.RoundFloat(fromPlayer.FreeThrowsMade-p.FreeThrowsMade, 2),
		GamesPlayed:            similarity.RoundFloat(fromPlayer.GamesPlayed-p.GamesPlayed, 2),
		Height:                 similarity.RoundFloat(fromPlayer.Height-p.Height, 2),
		Minutes:                similarity.RoundFloat(fromPlayer.Minutes-p.Minutes, 2),
		OffensiveRebounds:      similarity.RoundFloat(fromPlayer.OffensiveRebounds-p.OffensiveRebounds, 2),
		Points:                 similarity.RoundFloat(fromPlayer.Points-p.Points, 2),
		Rebounds:               similarity.RoundFloat(fromPlayer.Rebounds-p.Rebounds, 2),
		Steals:                 similarity.RoundFloat(fromPlayer.Steals-p.Steals, 2),
		ThreePointersAttempted: similarity.RoundFloat(fromPlayer.ThreePointersAttempted-p.ThreePointersAttempted, 2),
		ThreePointersMade:      similarity.RoundFloat(fromPlayer.ThreePointersMade-p.ThreePointersMade, 2),
		Turnovers:              similarity.RoundFloat(fromPlayer.Turnovers-p.Turnovers, 2),
		Weight:                 similarity.RoundFloat(fromPlayer.Weight-p.Weight, 2),
		PersonalFoulsDrawn:     similarity.RoundFloat(fromPlayer.PersonalFoulsDrawn-p.PersonalFoulsDrawn, 2),
		PersonalFouls:          similarity.RoundFloat(fromPlayer.PersonalFouls-p.PersonalFouls, 2),
		Player:                 p.Player}
	return d
}

// EuclideanDistance calculates the euclidean distance of a PlayerDiff object that stores the difference between two players' averages.
func EuclideanDistance(diff PlayerDiff, statsOfInterest []Stat) float64 {
	// TODO: could add user-inputed weights for different stats
	sum := 0.0
	for _, stat := range statsOfInterest {
		switch stat {
		// math.Pow(diff.Assists, 2) +
		case Assists:
			sum += math.Pow(diff.Assists, 2)
		// math.Pow(diff.Blocks, 2) +
		case Blocks:
			sum += math.Pow(diff.Blocks, 2)
		// math.Pow(diff.DefensiveRebounds, 2) +
		case DefensiveRebounds:
			sum += math.Pow(diff.DefensiveRebounds, 2)
		case FieldGoalsAttempted:
			sum += math.Pow(diff.FieldGoalsAttempted, 2)
		case FieldGoalsMade:
			sum += math.Pow(diff.FieldGoalsMade, 2)
		case FreeThrowsAttempted:
			sum += math.Pow(diff.FreeThrowsAttempted, 2)
		// math.Pow(diff.FreeThrowsMade, 2) +
		case FreeThrowsMade:
			sum += math.Pow(diff.FreeThrowsMade, 2)
		case GamesPlayed:
			sum += math.Pow(diff.GamesPlayed, 2)
		case Height:
			sum += math.Pow(diff.Height, 2)
		case Minutes:
			sum += math.Pow(diff.Minutes, 2)
		case OffensiveRebounds:
			sum += math.Pow(diff.OffensiveRebounds, 2)
		case PersonalFoulsDrawn:
			sum += math.Pow(diff.PersonalFoulsDrawn, 2)
		case PersonalFouls:
			sum += math.Pow(diff.PersonalFouls, 2)
		case Points:
			sum += math.Pow(diff.Points, 2)
		case Rebounds:
			sum += math.Pow(diff.Rebounds, 2)
		case Steals:
			sum += math.Pow(diff.Steals, 2)
		case ThreePointersAttempted:
			sum += math.Pow(diff.ThreePointersAttempted, 2)
		case ThreePointersMade:
			sum += math.Pow(diff.ThreePointersMade, 2)
		case Turnovers:
			sum += math.Pow(diff.Turnovers, 2)
		case Weight:
			sum += math.Pow(diff.Weight, 2)
		default:
			sum += 0.0
		}
	}
	return math.Sqrt(sum)

}

func (p *PlayerAverage) Score(stat Stat) float64 {
	switch NewStat(string(stat)) {
	case Points:
		return float64(p.Points)
	case Assists:
		return float64(p.Assists)
	case Rebounds:
		return float64(p.Rebounds)
	case Steals:
		return float64(p.Steals)
	case Blocks:
		return float64(p.Blocks)
	case ThreePointersMade:
		return float64(p.ThreePointersMade)
	case ThreePointersAttempted:
		return float64(p.ThreePointersAttempted)
	case FreeThrowsMade:
		return float64(p.FreeThrowsMade)
	case FreeThrowsAttempted:
		return float64(p.FreeThrowsAttempted)
	case FieldGoalsMade:
		return float64(p.FieldGoalsMade)
	case FieldGoalsAttempted:
		return float64(p.FieldGoalsAttempted)
	case Minutes:
		return float64(p.Minutes)
	case OffensiveRebounds:
		return float64(p.OffensiveRebounds)
	case DefensiveRebounds:
		return float64(p.DefensiveRebounds)
	case Turnovers:
		return float64(p.Turnovers)
	case PersonalFouls:
		return float64(p.PersonalFouls)
	case PersonalFoulsDrawn:
		return float64(p.PersonalFoulsDrawn)
	case PointsReboundsAssists:
		return float64(p.Points + p.Rebounds + p.Assists)
	case PointsRebounds:
		return float64(p.Points + p.Rebounds)
	case PointsAssists:
		return float64(p.Points + p.Assists)
	case ReboundsAssists:
		return float64(p.Rebounds + p.Assists)
	case BlocksSteals:
		return float64(p.Blocks + p.Steals)
	case FantasyScore:
		return float64(p.Points) + float64(p.Rebounds)*1.2 + float64(p.Assists)*1.5 + float64(p.Steals)*3 + float64(p.Blocks)*3 - float64(p.Turnovers)
	case Weight:
		return float64(p.Weight)
	case Height:
		return float64(p.Height)
	case GamesPlayed:
		return float64(p.GamesPlayed)
	default:
		logrus.Warnf("Unknown stat: %s", stat)
		return 0
	}
}

func (p *PlayerAverage) AverageStats() *AverageStats {
	return &AverageStats{
		Points:                 similarity.RoundFloat(p.Points, 2),
		Assists:                similarity.RoundFloat(p.Assists, 2),
		Rebounds:               similarity.RoundFloat(p.Rebounds, 2),
		Steals:                 similarity.RoundFloat(p.Steals, 2),
		Blocks:                 similarity.RoundFloat(p.Blocks, 2),
		ThreePointersMade:      similarity.RoundFloat(p.ThreePointersMade, 2),
		ThreePointersAttempted: similarity.RoundFloat(p.ThreePointersAttempted, 2),
		FreeThrowsMade:         similarity.RoundFloat(p.FreeThrowsMade, 2),
		FreeThrowsAttempted:    similarity.RoundFloat(p.FreeThrowsAttempted, 2),
		FieldGoalsMade:         similarity.RoundFloat(p.FieldGoalsMade, 2),
		FieldGoalsAttempted:    similarity.RoundFloat(p.FieldGoalsAttempted, 2),
		Minutes:                similarity.RoundFloat(p.Minutes, 2),
		OffensiveRebounds:      similarity.RoundFloat(p.OffensiveRebounds, 2),
		DefensiveRebounds:      similarity.RoundFloat(p.DefensiveRebounds, 2),
		Turnovers:              similarity.RoundFloat(p.Turnovers, 2),
		PersonalFouls:          similarity.RoundFloat(p.PersonalFouls, 2),
		PersonalFoulsDrawn:     similarity.RoundFloat(p.PersonalFoulsDrawn, 2),
		Weight:                 similarity.RoundFloat(p.Weight, 2),
		Height:                 similarity.RoundFloat(p.Height, 2),
		GamesPlayed:            similarity.RoundFloat(p.GamesPlayed, 2),
	}
}

func (startValue *AverageStats) PercentChange(finalValue *AverageStats) *AverageStats {
	return &AverageStats{
		Points:                 similarity.RoundFloat(((finalValue.Points-startValue.Points)/startValue.Points)*100, 2),
		Assists:                similarity.RoundFloat(((finalValue.Assists-startValue.Assists)/startValue.Assists)*100, 2),
		Rebounds:               similarity.RoundFloat(((finalValue.Rebounds-startValue.Rebounds)/startValue.Rebounds)*100, 2),
		Steals:                 similarity.RoundFloat(((finalValue.Steals-startValue.Steals)/startValue.Steals)*100, 2),
		Blocks:                 similarity.RoundFloat(((finalValue.Blocks-startValue.Blocks)/startValue.Blocks)*100, 2),
		ThreePointersMade:      similarity.RoundFloat(((finalValue.ThreePointersMade-startValue.ThreePointersMade)/startValue.ThreePointersMade)*100, 2),
		ThreePointersAttempted: similarity.RoundFloat(((finalValue.ThreePointersAttempted-startValue.ThreePointersAttempted)/startValue.ThreePointersAttempted)*100, 2),
		FreeThrowsMade:         similarity.RoundFloat(((finalValue.FreeThrowsMade-startValue.FreeThrowsMade)/startValue.FreeThrowsMade)*100, 2),
		FreeThrowsAttempted:    similarity.RoundFloat(((finalValue.FreeThrowsAttempted-startValue.FreeThrowsAttempted)/startValue.FreeThrowsAttempted)*100, 2),
		FieldGoalsMade:         similarity.RoundFloat(((finalValue.FieldGoalsMade-startValue.FieldGoalsMade)/startValue.FieldGoalsMade)*100, 2),
		FieldGoalsAttempted:    similarity.RoundFloat(((finalValue.FieldGoalsAttempted-startValue.FieldGoalsAttempted)/startValue.FieldGoalsAttempted)*100, 2),
		Minutes:                similarity.RoundFloat(((finalValue.Minutes-startValue.Minutes)/startValue.Minutes)*100, 2),
		OffensiveRebounds:      similarity.RoundFloat(((finalValue.OffensiveRebounds-startValue.OffensiveRebounds)/startValue.OffensiveRebounds)*100, 2),
		DefensiveRebounds:      similarity.RoundFloat(((finalValue.DefensiveRebounds-startValue.DefensiveRebounds)/startValue.DefensiveRebounds)*100, 2),
		Turnovers:              similarity.RoundFloat(((finalValue.Turnovers-startValue.Turnovers)/startValue.Turnovers)*100, 2),
		PersonalFouls:          similarity.RoundFloat(((finalValue.PersonalFouls-startValue.PersonalFouls)/startValue.PersonalFouls)*100, 2),
		PersonalFoulsDrawn:     similarity.RoundFloat(((finalValue.PersonalFoulsDrawn-startValue.PersonalFoulsDrawn)/startValue.PersonalFoulsDrawn)*100, 2),
		Weight:                 similarity.RoundFloat(((finalValue.Weight-startValue.Weight)/startValue.Weight)*100, 2),
		Height:                 similarity.RoundFloat(((finalValue.Height-startValue.Height)/startValue.Height)*100, 2),
		GamesPlayed:            similarity.RoundFloat(((finalValue.GamesPlayed-startValue.GamesPlayed)/startValue.GamesPlayed)*100, 2),
	}
}

type TeamAverage struct {
	WinsAndLosses      []string `json:"wins_and_losses" bson:"wins_and_losses"`
	GamesWon           float64  `json:"games_won" bson:"games_won"`
	GamesLost          float64  `json:"games_lost" bson:"games_lost"`
	Points             float64  `json:"points" bson:"points"`
	OppPoints          float64  `json:"opponent_points" bson:"opponent_points"`
	Assists            float64  `json:"assists" bson:"assists"`
	OppAssists         float64  `json:"opponent_assists" bson:"opponent_assists"`
	Rebounds           float64  `json:"rebounds" bson:"rebounds"`
	OppRebounds        float64  `json:"opponent_rebounds" bson:"opponent_rebounds"`
	Steals             float64  `json:"steals" bson:"steals"`
	Blocks             float64  `json:"blocks" bson:"blocks"`
	Turnovers          float64  `json:"turnovers" bson:"turnovers"`
	ThreePointersMade  float64  `json:"three_pointers_made" bson:"three_pointers_made"`
	PersonalFouls      float64  `json:"personal_fouls" bson:"personal_fouls"`
	PersonalFoulsDrawn float64  `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	Team               Team     `json:"team" bson:"team"`
}
