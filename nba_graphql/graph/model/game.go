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
	Assists                      sql.NullInt16   `db:"assists" json:"assists" bson:"assists"`
	PotentialAssists             sql.NullInt16   `db:"potentialAssists" json:"potential_assists" bson:"potential_assists"`
	AssistConversionRate         sql.NullFloat64 `db:"assistConversionRate" json:"assist_conversion_rate" bson:"assist_conversion_rate"`
	Passes                       sql.NullInt16   `db:"passes" json:"passes" bson:"passes"`
	Date                         *time.Time      `db:"date" json:"date" bson:"date"`
	DefensiveReboundPercentage   sql.NullFloat64 `db:"defensiveReboundPct" json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	DefensiveRebounds            sql.NullInt16   `db:"defensiveRebounds" json:"defensive_rebounds" bson:"defensive_rebounds"`
	EffectiveFieldGoalPercentage sql.NullFloat64 `db:"effectiveFieldGoalPct" json:"effective_field_goal_percentage" bson:"effective_field_goal_percentage"`
	FieldGoalPercentage          sql.NullFloat64 `db:"fieldGoalPct" json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted          sql.NullInt16   `db:"fieldGoalsAttempted" json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade               sql.NullInt16   `db:"fieldGoalsMade" json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted          sql.NullInt16   `db:"freeThrowsAttempted" json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade               sql.NullInt16   `db:"freeThrowsMade" json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage         sql.NullFloat64 `db:"freeThrowPct" json:"free_throws_percentage" bson:"free_throws_percentage"` //change to freeThrowPercentage
	GameID                       string          `db:"gameID" json:"gameID" bson:"gameID"`
	HomeOrAway                   HomeOrAway      `db:"homeAway" json:"home_or_away" bson:"home_or_away"` //add "HOME" or "AWAY" constraint
	Margin                       sql.NullInt16   `db:"margin" json:"margin" bson:"margin"`
	Minutes                      float64         `db:"minutes" json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   sql.NullFloat64 `db:"offensiveReboundPct" json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            sql.NullInt16   `db:"offensiveRebounds" json:"offensive_rebounds" bson:"offensive_rebounds"`
	TeamID                       int             `db:"teamID" json:"team" bson:"team"`
	OpponentID                   int             `db:"opponentID" json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           sql.NullInt16   `db:"personalFoulsDrawn" json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                sql.NullInt16   `db:"personalFouls" json:"personal_fouls" bson:"personal_fouls"`
	Points                       sql.NullInt16   `db:"points" json:"points" bson:"points"`
	PlayerID                     int             `db:"playerID" json:"playerID" bson:"playerID"`
	PlayerRef                    *Player         `json:"playerRef" bson:"playerRef"`
	Playoffs                     bool            `db:"playoffs" json:"playoffs" bson:"playoffs"`
	Season                       SeasonOption    `db:"season" json:"season" bson:"season"`
	ThreePointPercentage         sql.NullFloat64 `db:"threePointPct" json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       sql.NullInt16   `db:"threePointersAttempted" json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            sql.NullInt16   `db:"threePointersMade" json:"three_pointers_made" bson:"three_pointers_made"`
	Rebounds                     sql.NullInt16   `db:"rebounds" json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       sql.NullFloat64 `db:"trueShootingPct" json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    sql.NullInt16   `db:"turnovers" json:"turnovers" bson:"turnovers"`
	Blocks                       sql.NullInt16   `db:"blocks" json:"blocks" bson:"blocks"`
	Steals                       sql.NullInt16   `db:"steals" json:"steals" bson:"steals"`
	Usage                        sql.NullFloat64 `db:"usage" json:"usage" bson:"usage"`
	Outcome                      string          `db:"outcome" json:"win_or_loss" bson:"win_or_loss"`

	CreatedAt *time.Time `db:"CreatedAt" json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt *time.Time `db:"UpdatedAt" json:"UpdatedAt" bson:"UpdatedAt"`
}

func (g *PlayerGame) Score(stat Stat) float64 {
	switch stat {
	case Points:
		return float64(g.Points.Int16)
	case Assists:
		return float64(g.Assists.Int16)
	case Rebounds:
		return float64(g.Rebounds.Int16)
	case Steals:
		return float64(g.Steals.Int16)
	case Blocks:
		return float64(g.Blocks.Int16)
	case ThreePointersMade:
		return float64(g.ThreePointersMade.Int16)
	case ThreePointersAttempted:
		return float64(g.ThreePointersAttempted.Int16)
	case ThreePointPercentage:
		return g.ThreePointPercentage.Float64
	case FreeThrowsMade:
		return float64(g.FreeThrowsMade.Int16)
	case FreeThrowsAttempted:
		return float64(g.FreeThrowsAttempted.Int16)
	case FreeThrowsPercentage:
		return g.FreeThrowsPercentage.Float64
	case FieldGoalsMade:
		return float64(g.FieldGoalsMade.Int16)
	case FieldGoalsAttempted:
		return float64(g.FieldGoalsAttempted.Int16)
	case FieldGoalPercentage:
		return g.FieldGoalPercentage.Float64
	case EffectiveFieldGoalPercentage:
		return g.EffectiveFieldGoalPercentage.Float64
	case TrueShootingPercentage:
		return g.TrueShootingPercentage.Float64
	case Minutes:
		return g.Minutes
	case OffensiveRebounds:
		return float64(g.OffensiveRebounds.Int16)
	case DefensiveRebounds:
		return float64(g.DefensiveRebounds.Int16)
	case AssistPercentage:
		return g.AssistPercentage.Float64
	case OffensiveReboundPercentage:
		return g.OffensiveReboundPercentage.Float64
	case DefensiveReboundPercentage:
		return g.DefensiveReboundPercentage.Float64
	case Usage:
		return g.Usage.Float64
	case Turnovers:
		return float64(g.Turnovers.Int16)
	case PersonalFouls:
		return float64(g.PersonalFouls.Int16)
	case PersonalFoulsDrawn:
		return float64(g.PersonalFoulsDrawn.Int16)
	case PointsReboundsAssists:
		return float64(g.Points.Int16 + g.Rebounds.Int16 + g.Assists.Int16)
	case PointsRebounds:
		return float64(g.Points.Int16 + g.Rebounds.Int16)
	case PointsAssists:
		return float64(g.Points.Int16 + g.Assists.Int16)
	case ReboundsAssists:
		return float64(g.Rebounds.Int16 + g.Assists.Int16)
	case BlocksSteals:
		return float64(g.Blocks.Int16 + g.Steals.Int16)
	case FantasyScore:
		return float64(g.Points.Int16) + float64(g.Rebounds.Int16)*1.2 + float64(g.Assists.Int16)*1.5 + float64(g.Steals.Int16)*3 + float64(g.Blocks.Int16)*3 - float64(g.Turnovers.Int16)
	case GamesPlayed:
		if g.Score(Minutes) > 0 {
			return 1
		}
		return 0
	case DoubleDouble:
		countDouble := 0
		if g.Points.Int16 >= 10 {
			countDouble++
		}
		if g.Rebounds.Int16 >= 10 {
			countDouble++
		}
		if g.Assists.Int16 >= 10 {
			countDouble++
		}
		if g.Steals.Int16 >= 10 {
			countDouble++
		}
		if g.Blocks.Int16 >= 10 {
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
