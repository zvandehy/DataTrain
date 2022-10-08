package model

import (
	"fmt"
	"math"
	"time"

	similarity "github.com/zvandehy/DataTrain/nba_graphql/math"
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
	Games []PlayerGame `json:"games" bson:"games"`
}

// TODO: Provide filters other than the dateRange
func (p *Player) AverageStats(startDate, endDate time.Time) PlayerAverage {
	average := PlayerAverage{}
	var filteredGames []PlayerGame
	for _, game := range p.Games {
		gameDate, err := time.Parse("2006-01-02", game.Date)
		if err != nil {
			continue
		}
		if gameDate.After(startDate) && gameDate.Before(endDate) {
			filteredGames = append(filteredGames, game)
		}
	}
	average.GamesPlayed = float64(len(filteredGames))
	average.Player = *p
	average.Height = float64(p.HeightInInches())
	average.Weight = float64(p.Weight)

	for _, game := range filteredGames {
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

	average.Assists /= float64(len(filteredGames))
	average.Blocks /= float64(len(filteredGames))
	average.DefensiveRebounds /= float64(len(filteredGames))
	average.FieldGoalsAttempted /= float64(len(filteredGames))
	average.FieldGoalsMade /= float64(len(filteredGames))
	average.FreeThrowsAttempted /= float64(len(filteredGames))
	average.FreeThrowsMade /= float64(len(filteredGames))
	average.Minutes /= float64(len(filteredGames))
	average.OffensiveRebounds /= float64(len(filteredGames))
	average.PersonalFoulsDrawn /= float64(len(filteredGames))
	average.PersonalFouls /= float64(len(filteredGames))
	average.Points /= float64(len(filteredGames))
	average.Rebounds /= float64(len(filteredGames))
	average.Steals /= float64(len(filteredGames))
	average.ThreePointersAttempted /= float64(len(filteredGames))
	average.ThreePointersMade /= float64(len(filteredGames))
	average.Turnovers /= float64(len(filteredGames))

	return average
}

func (p *Player) HeightInInches() int {
	var feet, inches int
	fmt.Sscanf(p.Height, "%d'%d", &feet, &inches)
	return feet*12 + inches
}

func (p Player) String() string {
	return Print(p)
}

func (p PlayerGame) String() string {
	return Print(p)
}

type PlayerOpponentMatchup struct {
	PlayerID   int    `json:"playerID" bson:"playerID"`
	OpponentID int    `json:"opponent" bson:"opponent"`
	Date       string `json:"date" bson:"date"`
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

type PlayerDiff struct {
	PlayerAverage
	Distance float64
}

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

func (p *PlayerAverage) ToAverageStats() *AverageStats {
	mins := fmt.Sprintf("%d", int(p.Minutes))
	secs := fmt.Sprintf("%02d", int((p.Minutes-float64(int(p.Minutes)))*60))
	return &AverageStats{
		Assists:                p.Assists,
		Blocks:                 p.Blocks,
		DefensiveRebounds:      p.DefensiveRebounds,
		FieldGoalsAttempted:    p.FieldGoalsAttempted,
		FieldGoalsMade:         p.FieldGoalsMade,
		FreeThrowsAttempted:    p.FreeThrowsAttempted,
		FreeThrowsMade:         p.FreeThrowsMade,
		Minutes:                fmt.Sprintf("%s:%s", mins, secs),
		OffensiveRebounds:      p.OffensiveRebounds,
		PersonalFoulsDrawn:     p.PersonalFoulsDrawn,
		PersonalFouls:          p.PersonalFouls,
		Points:                 p.Points,
		Rebounds:               p.Rebounds,
		Steals:                 p.Steals,
		ThreePointersAttempted: p.ThreePointersAttempted,
		ThreePointersMade:      p.ThreePointersMade,
		Turnovers:              p.Turnovers,
	}
}

func (p *PlayerAverage) Difference(fromPlayer PlayerAverage) PlayerDiff {
	d := PlayerDiff{
		PlayerAverage: PlayerAverage{Assists: similarity.RoundFloat(fromPlayer.Assists-p.Assists, 2),
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
			PersonalFouls:          similarity.RoundFloat(fromPlayer.PersonalFouls-p.PersonalFouls, 2)},
	}
	d.Distance = EuclideanDistance(d)
	return d
}

// EuclideanDistance is the distance between two PlayerDiff vectors. If a stat should not be used in the distance calculation, set it to 0.
func EuclideanDistance(diff PlayerDiff) float64 {
	// TODO: could add user-inputed weights for different stats
	sum :=
		math.Pow(diff.Assists, 2) +
			math.Pow(diff.Blocks, 2) +
			math.Pow(diff.DefensiveRebounds, 2) +
			math.Pow(diff.FieldGoalsAttempted, 2) +
			math.Pow(diff.FieldGoalsMade, 2) +
			math.Pow(diff.FreeThrowsAttempted, 2) +
			math.Pow(diff.FreeThrowsMade, 2) +
			math.Pow(diff.GamesPlayed, 2) +
			math.Pow(diff.Height, 2) +
			math.Pow(diff.Minutes, 2) +
			math.Pow(diff.OffensiveRebounds, 2) +
			math.Pow(diff.PersonalFoulsDrawn, 2) +
			math.Pow(diff.PersonalFouls, 2) +
			math.Pow(diff.Points, 2) +
			math.Pow(diff.Rebounds, 2) +
			math.Pow(diff.Steals, 2) +
			math.Pow(diff.ThreePointersAttempted, 2) +
			math.Pow(diff.ThreePointersMade, 2) +
			math.Pow(diff.Turnovers, 2) +
			math.Pow(diff.Weight, 2)
	return math.Sqrt(sum)

}

func (p *PlayerAverage) Score(stat Stat) float64 {
	switch stat {
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
		return 0
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
