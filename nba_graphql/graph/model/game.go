package model

import (
	"strconv"
	"strings"
)

type PlayerGame struct {
	AssistPercentage             float32 `json:"assist_percentage" bson:"assist_percentage"`
	Assists                      int     `json:"assists" bson:"assists"`
	Date                         string  `json:"date" bson:"date"`
	DefensiveReboundPercentage   float32 `json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	DefensiveRebounds            int     `json:"defensive_rebounds" bson:"defensive_rebounds"`
	EffectiveFieldGoalPercentage float32 `json:"effective_field_goal_percentage" bson:"effective_field_goal_percentage"`
	FieldGoalPercentage          float32 `json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted          int     `json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade               int     `json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted          int     `json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade               int     `json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage         float32 `json:"free_throws_percentage" bson:"free_throws_percentage"`
	GameID                       string  `json:"gameID" bson:"gameID"`
	HomeOrAway                   string  `json:"home_or_away" bson:"home_or_away"`
	Margin                       int     `json:"margin" bson:"margin"`
	Minutes                      string  `json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   float32 `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            int     `json:"offensive_rebounds" bson:"offensive_rebounds"`
	TeamID                       int     `json:"team" bson:"team"`
	OpponentID                   int     `json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           int     `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                int     `json:"personal_fouls" bson:"personal_fouls"`
	Points                       int     `json:"points" bson:"points"`
	PlayerID                     int     `json:"playerID" bson:"playerID"`
	Playoffs                     bool    `json:"playoffs" bson:"playoffs"`
	Season                       string  `json:"season" bson:"season"`
	ThreePointPercentage         float32 `json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       int     `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            int     `json:"three_pointers_made" bson:"three_pointers_made"`
	Rebounds                     int     `json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       float32 `json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    int     `json:"turnovers" bson:"turnovers"`
	Blocks                       int     `json:"blocks" bson:"blocks"`
	Steals                       int     `json:"steals" bson:"steals"`
	Usage                        float32 `json:"usage" bson:"usage"`
	WinOrLoss                    string  `json:"win_or_loss" bson:"win_or_loss"`
}

type Stat string

const (
	Points                       Stat = "points"
	Assists                      Stat = "assists"
	Rebounds                     Stat = "rebounds"
	Steals                       Stat = "steals"
	Blocks                       Stat = "blocks"
	ThreePointersMade            Stat = "three_pointers_made"
	ThreePointersAttempted       Stat = "three_pointers_attempted"
	ThreePointPercentage         Stat = "three_point_percentage"
	FreeThrowsMade               Stat = "free_throws_made"
	FreeThrowsAttempted          Stat = "free_throws_attempted"
	FreeThrowsPercentage         Stat = "free_throws_percentage"
	FieldGoalsMade               Stat = "field_goals_made"
	FieldGoalsAttempted          Stat = "field_goals_attempted"
	FieldGoalPercentage          Stat = "field_goal_percentage"
	EffectiveFieldGoalPercentage Stat = "effective_field_goal_percentage"
	TrueShootingPercentage       Stat = "true_shooting_percentage"
	Minutes                      Stat = "minutes"
	OffensiveRebounds            Stat = "offensive_rebounds"
	DefensiveRebounds            Stat = "defensive_rebounds"
	AssistPercentage             Stat = "assist_percentage"
	OffensiveReboundPercentage   Stat = "offensive_rebound_percentage"
	DefensiveReboundPercentage   Stat = "defensive_rebound_percentage"
	Usage                        Stat = "usage"
	Turnovers                    Stat = "turnovers"
	PersonalFouls                Stat = "personal_fouls"
	PersonalFoulsDrawn           Stat = "personal_fouls_drawn"
	PointsReboundsAssists        Stat = "points_rebounds_assists"
	PointsRebounds               Stat = "points_rebounds"
	PointsAssists                Stat = "points_assists"
	ReboundsAssists              Stat = "rebounds_assists"
	BlocksSteals                 Stat = "blocks_steals"
	FantasyScore                 Stat = "fantasy_score"
	GamesPlayed                  Stat = "games_played"
	Height                       Stat = "height"
	Weight                       Stat = "weight"
)

func NewStat(stat string) Stat {
	lookup := strings.ReplaceAll(strings.ToLower(stat), " ", "_")
	lookup = strings.ReplaceAll(lookup, "+", "_")
	switch lookup {
	case "points":
		return Points
	case "assists":
		return Assists
	case "rebounds":
		return Rebounds
	case "steals":
		return Steals
	case "blocks":
		return Blocks
	case "three_pointers_made":
		fallthrough
	case "3-pt_made":
		return ThreePointersMade
	case "three_pointers_attempted":
		return ThreePointersAttempted
	case "three_point_percentage":
		return ThreePointPercentage
	case "free_throws_made":
		return FreeThrowsMade
	case "free_throws_attempted":
		return FreeThrowsAttempted
	case "free_throws_percentage":
		return FreeThrowsPercentage
	case "field_goals_made":
		return FieldGoalsMade
	case "field_goals_attempted":
		return FieldGoalsAttempted
	case "field_goal_percentage":
		return FieldGoalPercentage
	case "effective_field_goal_percentage":
		return EffectiveFieldGoalPercentage
	case "true_shooting_percentage":
		return TrueShootingPercentage
	case "minutes":
		return Minutes
	case "offensive_rebounds":
		return OffensiveRebounds
	case "defensive_rebounds":
		return DefensiveRebounds
	case "assist_percentage":
		return AssistPercentage
	case "offensive_rebound_percentage":
		return OffensiveReboundPercentage
	case "defensive_rebound_percentage":
		return DefensiveReboundPercentage
	case "usage":
		return Usage
	case "turnovers":
		return Turnovers
	case "games_played":
		return GamesPlayed
	case "personal_fouls":
		return PersonalFouls
	case "personal_fouls_drawn":
		return PersonalFoulsDrawn
	case "points_rebounds_assists":
		fallthrough
	case "pts_rebs_asts":
		return PointsReboundsAssists
	case "points_rebounds":
		fallthrough
	case "pts_rebs":
		return PointsRebounds
	case "points_assists":
		fallthrough
	case "pts_asts":
		return PointsAssists
	case "rebounds_assists":
		fallthrough
	case "rebs_asts":
		return ReboundsAssists
	case "blocks_steals":
		fallthrough
	case "blks_stls":
		return BlocksSteals
	case "fantasy_score":
		return FantasyScore
	case "height":
		return Height
	case "weight":
		return Weight
	default:
		return ""
	}
}

func (g *PlayerGame) Score(stat Stat) float32 {
	switch stat {
	case Points:
		return float32(g.Points)
	case Assists:
		return float32(g.Assists)
	case Rebounds:
		return float32(g.Rebounds)
	case Steals:
		return float32(g.Steals)
	case Blocks:
		return float32(g.Blocks)
	case ThreePointersMade:
		return float32(g.ThreePointersMade)
	case ThreePointersAttempted:
		return float32(g.ThreePointersAttempted)
	case ThreePointPercentage:
		return g.ThreePointPercentage
	case FreeThrowsMade:
		return float32(g.FreeThrowsMade)
	case FreeThrowsAttempted:
		return float32(g.FreeThrowsAttempted)
	case FreeThrowsPercentage:
		return g.FreeThrowsPercentage
	case FieldGoalsMade:
		return float32(g.FieldGoalsMade)
	case FieldGoalsAttempted:
		return float32(g.FieldGoalsAttempted)
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
		return float32(min)
	case OffensiveRebounds:
		return float32(g.OffensiveRebounds)
	case DefensiveRebounds:
		return float32(g.DefensiveRebounds)
	case AssistPercentage:
		return g.AssistPercentage
	case OffensiveReboundPercentage:
		return g.OffensiveReboundPercentage
	case DefensiveReboundPercentage:
		return g.DefensiveReboundPercentage
	case Usage:
		return g.Usage
	case Turnovers:
		return float32(g.Turnovers)
	case PersonalFouls:
		return float32(g.PersonalFouls)
	case PersonalFoulsDrawn:
		return float32(g.PersonalFoulsDrawn)
	case PointsReboundsAssists:
		return float32(g.Points + g.Rebounds + g.Assists)
	case PointsRebounds:
		return float32(g.Points + g.Rebounds)
	case PointsAssists:
		return float32(g.Points + g.Assists)
	case ReboundsAssists:
		return float32(g.Rebounds + g.Assists)
	case BlocksSteals:
		return float32(g.Blocks + g.Steals)
	case FantasyScore:
		return float32(g.Points) + float32(g.Rebounds)*1.2 + float32(g.Assists)*1.5 + float32(g.Steals)*3 + float32(g.Blocks)*3 - float32(g.Turnovers)
	default:
		return 0
	}
}

func PlayerAverageStats() []Stat {
	return []Stat{
		Assists,
		Blocks,
		DefensiveRebounds,
		FieldGoalsAttempted,
		FieldGoalsMade,
		FreeThrowsAttempted,
		FreeThrowsMade,
		GamesPlayed,
		Height,
		Minutes,
		OffensiveRebounds,
		PersonalFoulsDrawn,
		PersonalFouls,
		Points,
		Rebounds,
		Steals,
		ThreePointersAttempted,
		ThreePointersMade,
		Turnovers,
		Weight,
	}
}

func AllStats() []Stat {
	return []Stat{
		Points,
		Assists,
		Rebounds,
		Steals,
		Blocks,
		ThreePointersMade,
		ThreePointersAttempted,
		ThreePointPercentage,
		FreeThrowsMade,
		FreeThrowsAttempted,
		FreeThrowsPercentage,
		FieldGoalsMade,
		FieldGoalsAttempted,
		FieldGoalPercentage,
		EffectiveFieldGoalPercentage,
		TrueShootingPercentage,
		Minutes,
		OffensiveRebounds,
		DefensiveRebounds,
		AssistPercentage,
		OffensiveReboundPercentage,
		DefensiveReboundPercentage,
		Usage,
		Turnovers,
		GamesPlayed,
		PersonalFouls,
		PersonalFoulsDrawn,
		PointsReboundsAssists,
		PointsRebounds,
		PointsAssists,
		ReboundsAssists,
		BlocksSteals,
		FantasyScore,
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
