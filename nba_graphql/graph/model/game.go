package model

import (
	"strconv"
)

type PlayerGame struct {
	AssistPercentage             float64      `json:"assist_percentage" bson:"assist_percentage"`
	Assists                      int          `json:"assists" bson:"assists"`
	Date                         string       `json:"date" bson:"date"`
	DefensiveReboundPercentage   float64      `json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	DefensiveRebounds            int          `json:"defensive_rebounds" bson:"defensive_rebounds"`
	EffectiveFieldGoalPercentage float64      `json:"effective_field_goal_percentage" bson:"effective_field_goal_percentage"`
	FieldGoalPercentage          float64      `json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted          int          `json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade               int          `json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted          int          `json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade               int          `json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage         float64      `json:"free_throws_percentage" bson:"free_throws_percentage"`
	GameID                       string       `json:"gameID" bson:"gameID"`
	HomeOrAway                   string       `json:"home_or_away" bson:"home_or_away"`
	Margin                       int          `json:"margin" bson:"margin"`
	Minutes                      string       `json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   float64      `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            int          `json:"offensive_rebounds" bson:"offensive_rebounds"`
	TeamID                       int          `json:"team" bson:"team"`
	OpponentID                   int          `json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           int          `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                int          `json:"personal_fouls" bson:"personal_fouls"`
	Points                       int          `json:"points" bson:"points"`
	PlayerID                     int          `json:"playerID" bson:"playerID"`
	PlayerRef                    *Player      `json:"playerRef" bson:"playerRef"`
	Playoffs                     bool         `json:"playoffs" bson:"playoffs"`
	Season                       SeasonOption `json:"season" bson:"season"`
	ThreePointPercentage         float64      `json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       int          `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            int          `json:"three_pointers_made" bson:"three_pointers_made"`
	Rebounds                     int          `json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       float64      `json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    int          `json:"turnovers" bson:"turnovers"`
	Blocks                       int          `json:"blocks" bson:"blocks"`
	Steals                       int          `json:"steals" bson:"steals"`
	Usage                        float64      `json:"usage" bson:"usage"`
	WinOrLoss                    string       `json:"win_or_loss" bson:"win_or_loss"`
}

func (g *PlayerGame) Score(stat Stat) float64 {
	switch stat {
	case Points:
		return float64(g.Points)
	case Assists:
		return float64(g.Assists)
	case Rebounds:
		return float64(g.Rebounds)
	case Steals:
		return float64(g.Steals)
	case Blocks:
		return float64(g.Blocks)
	case ThreePointersMade:
		return float64(g.ThreePointersMade)
	case ThreePointersAttempted:
		return float64(g.ThreePointersAttempted)
	case ThreePointPercentage:
		return g.ThreePointPercentage
	case FreeThrowsMade:
		return float64(g.FreeThrowsMade)
	case FreeThrowsAttempted:
		return float64(g.FreeThrowsAttempted)
	case FreeThrowsPercentage:
		return g.FreeThrowsPercentage
	case FieldGoalsMade:
		return float64(g.FieldGoalsMade)
	case FieldGoalsAttempted:
		return float64(g.FieldGoalsAttempted)
	case FieldGoalPercentage:
		return g.FieldGoalPercentage
	case EffectiveFieldGoalPercentage:
		return g.EffectiveFieldGoalPercentage
	case TrueShootingPercentage:
		return g.TrueShootingPercentage
	case Minutes:
		min, err := ConvertMinutesToFloat(g.Minutes)
		if err != nil {
			return 0
		}
		return float64(min)
	case OffensiveRebounds:
		return float64(g.OffensiveRebounds)
	case DefensiveRebounds:
		return float64(g.DefensiveRebounds)
	case AssistPercentage:
		return g.AssistPercentage
	case OffensiveReboundPercentage:
		return g.OffensiveReboundPercentage
	case DefensiveReboundPercentage:
		return g.DefensiveReboundPercentage
	case Usage:
		return g.Usage
	case Turnovers:
		return float64(g.Turnovers)
	case PersonalFouls:
		return float64(g.PersonalFouls)
	case PersonalFoulsDrawn:
		return float64(g.PersonalFoulsDrawn)
	case PointsReboundsAssists:
		return float64(g.Points + g.Rebounds + g.Assists)
	case PointsRebounds:
		return float64(g.Points + g.Rebounds)
	case PointsAssists:
		return float64(g.Points + g.Assists)
	case ReboundsAssists:
		return float64(g.Rebounds + g.Assists)
	case BlocksSteals:
		return float64(g.Blocks + g.Steals)
	case FantasyScore:
		return float64(g.Points) + float64(g.Rebounds)*1.2 + float64(g.Assists)*1.5 + float64(g.Steals)*3 + float64(g.Blocks)*3 - float64(g.Turnovers)
	case GamesPlayed:
		if g.Score(Minutes) > 0 {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ConvertMinutesToFloat(minutes string) (float64, error) {
	//convert "mm:ss" to minutes
	min, err := strconv.ParseFloat(minutes[:len(minutes)-3], 64)
	if err != nil {
		return 0, err
	}
	sec, err := strconv.ParseFloat(minutes[len(minutes)-2:], 64)
	if err != nil {
		return 0, err
	}
	return min + sec/60, nil
}
