package model

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type PlayerGame struct {
	AssistPercentage             sql.NullFloat64 `db:"assistPct" json:"assist_percentage" bson:"assist_percentage"`
	Assists                      int             `db:"assists" json:"assists" bson:"assists"`
	PotentialAssists             int             `db:"potentialAssists" json:"potential_assists" bson:"potential_assists"`
	AssistConversionRate         sql.NullFloat64 `db:"assistConversionRate" json:"assist_conversion_rate" bson:"assist_conversion_rate"`
	Passes                       int             `db:"passes" json:"passes" bson:"passes"`
	Date                         *time.Time      `db:"date" json:"date" bson:"date"`
	DefensiveReboundPercentage   sql.NullFloat64 `db:"defensiveReboundPct" json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	DefensiveRebounds            int             `db:"defensiveRebounds" json:"defensive_rebounds" bson:"defensive_rebounds"`
	EffectiveFieldGoalPercentage sql.NullFloat64 `db:"effectiveFieldGoalPct" json:"effective_field_goal_percentage" bson:"effective_field_goal_percentage"`
	FieldGoalPercentage          sql.NullFloat64 `db:"fieldGoalPct" json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted          int             `db:"fieldGoalsAttempted" json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade               int             `db:"fieldGoalsMade" json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted          int             `db:"freeThrowsAttempted" json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade               int             `db:"freeThrowsMade" json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage         sql.NullFloat64 `db:"freeThrowPct" json:"free_throws_percentage" bson:"free_throws_percentage"` //change to freeThrowPercentage
	GameID                       string          `db:"gameID" json:"gameID" bson:"gameID"`
	HomeOrAway                   string          `db:"homeAway" json:"home_or_away" bson:"home_or_away"` //add "HOME" or "AWAY" constraint
	Margin                       int             `db:"margin" json:"margin" bson:"margin"`
	Minutes                      float64         `db:"minutes" json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   sql.NullFloat64 `db:"offensiveReboundPct" json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            int             `db:"offensiveRebounds" json:"offensive_rebounds" bson:"offensive_rebounds"`
	TeamID                       int             `db:"teamID" json:"team" bson:"team"`
	OpponentID                   int             `db:"opponentID" json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           int             `db:"personalFoulsDrawn" json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                int             `db:"personalFouls" json:"personal_fouls" bson:"personal_fouls"`
	Points                       int             `db:"points" json:"points" bson:"points"`
	PlayerID                     int             `db:"playerID" json:"playerID" bson:"playerID"`
	PlayerRef                    *Player         `json:"playerRef" bson:"playerRef"`
	Playoffs                     bool            `db:"playoffs" json:"playoffs" bson:"playoffs"`
	Season                       SeasonOption    `db:"season" json:"season" bson:"season"`
	ThreePointPercentage         sql.NullFloat64 `db:"threePointPct" json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       int             `db:"threePointersAttempted" json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            int             `db:"threePointersMade" json:"three_pointers_made" bson:"three_pointers_made"`
	Rebounds                     int             `db:"rebounds" json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       sql.NullFloat64 `db:"trueShootingPct" json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    int             `db:"turnovers" json:"turnovers" bson:"turnovers"`
	Blocks                       int             `db:"blocks" json:"blocks" bson:"blocks"`
	Steals                       int             `db:"steals" json:"steals" bson:"steals"`
	Usage                        sql.NullFloat64 `db:"usage" json:"usage" bson:"usage"`
	Outcome                      string          `db:"outcome" json:"win_or_loss" bson:"win_or_loss"`

	CreatedAt *time.Time `db:"CreatedAt" json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt *time.Time `db:"UpdatedAt" json:"UpdatedAt" bson:"UpdatedAt"`
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
		return g.ThreePointPercentage.Float64
	case FreeThrowsMade:
		return float64(g.FreeThrowsMade)
	case FreeThrowsAttempted:
		return float64(g.FreeThrowsAttempted)
	case FreeThrowsPercentage:
		return g.FreeThrowsPercentage.Float64
	case FieldGoalsMade:
		return float64(g.FieldGoalsMade)
	case FieldGoalsAttempted:
		return float64(g.FieldGoalsAttempted)
	case FieldGoalPercentage:
		return g.FieldGoalPercentage.Float64
	case EffectiveFieldGoalPercentage:
		return g.EffectiveFieldGoalPercentage.Float64
	case TrueShootingPercentage:
		return g.TrueShootingPercentage.Float64
	case Minutes:
		return g.Minutes
	case OffensiveRebounds:
		return float64(g.OffensiveRebounds)
	case DefensiveRebounds:
		return float64(g.DefensiveRebounds)
	case AssistPercentage:
		return g.AssistPercentage.Float64
	case OffensiveReboundPercentage:
		return g.OffensiveReboundPercentage.Float64
	case DefensiveReboundPercentage:
		return g.DefensiveReboundPercentage.Float64
	case Usage:
		return g.Usage.Float64
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
	case DoubleDouble:
		countDouble := 0
		if g.Points >= 10 {
			countDouble++
		}
		if g.Rebounds >= 10 {
			countDouble++
		}
		if g.Assists >= 10 {
			countDouble++
		}
		if g.Steals >= 10 {
			countDouble++
		}
		if g.Blocks >= 10 {
			countDouble++
		}
		if countDouble >= 2 {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ConvertMinutesToFloat(minutes string) (float64, error) {
	//test regex for "00:00"
	re := regexp.MustCompile(`^([0-9]{2}):([0-9]{2})$`)
	if !re.MatchString(minutes) {
		return 0, fmt.Errorf("invalid minutes format: %s", minutes)
	}
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
