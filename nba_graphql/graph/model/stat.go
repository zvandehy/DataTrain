package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

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
	DoubleDouble                 Stat = "double_double"
)

func NewStat(stat string) Stat {
	lookup := strings.ReplaceAll(strings.ToLower(stat), " ", "_")
	lookup = strings.ReplaceAll(lookup, "+", "_")
	lookup = strings.ReplaceAll(lookup, "-", "_")
	switch lookup {
	case "points":
		return Points
	case "assists":
		return Assists
	case "rebounds":
		return Rebounds
	case "steals":
		return Steals
	// TODO: This might actually be blocks against...
	case "blocked_shots":
		fallthrough
	case "blocks":
		return Blocks
	case "three_pointers_made":
		fallthrough
	case "threepointersmade":
		fallthrough
	case "3-pt_made":
		fallthrough
	case "3_pt_made":
		return ThreePointersMade
	case "three_pointers_attempted":
		fallthrough
	case "threepointersattempted":
		return ThreePointersAttempted
	case "three_point_percentage":
		fallthrough
	case "threepointpercentage":
		return ThreePointPercentage
	case "free_throws_made":
		fallthrough
	case "freethrowsmade":
		return FreeThrowsMade
	case "free_throws_attempted":
		fallthrough
	case "freethrowsattempted":
		return FreeThrowsAttempted
	case "free_throws_percentage":
		fallthrough
	case "freethrowspercentage":
		return FreeThrowsPercentage
	case "field_goals_made":
		fallthrough
	case "fieldgoalsmade":
		return FieldGoalsMade
	case "field_goals_attempted":
		fallthrough
	case "fieldgoalsattempted":
		return FieldGoalsAttempted
	case "field_goal_percentage":
		fallthrough
	case "fieldgoalpercentage":
		return FieldGoalPercentage
	case "effective_field_goal_percentage":
		fallthrough
	case "effectivefieldgoalpercentage":
		return EffectiveFieldGoalPercentage
	case "true_shooting_percentage":
		fallthrough
	case "trueshootingpercentage":
		return TrueShootingPercentage
	case "minutes":
		return Minutes
	case "offensive_rebounds":
		fallthrough
	case "offensiverebounds":
		return OffensiveRebounds
	case "defensive_rebounds":
		fallthrough
	case "defensiverebounds":
		return DefensiveRebounds
	case "assist_percentage":
		fallthrough
	case "assistpercentage":
		return AssistPercentage
	case "offensive_rebound_percentage":
		fallthrough
	case "offensivereboundpercentage":
		return OffensiveReboundPercentage
	case "defensive_rebound_percentage":
		fallthrough
	case "defensivereboundpercentage":
		return DefensiveReboundPercentage
	case "usage":
		return Usage
	case "turnovers":
		return Turnovers
	case "gamesplayed":
		fallthrough
	case "games_played":
		return GamesPlayed
	case "personal_fouls":
		fallthrough
	case "personalfouls":
		return PersonalFouls
	case "personal_fouls_drawn":
		fallthrough
	case "personalfoulsdrawn":
		return PersonalFoulsDrawn
	case "points_rebounds_assists":
		fallthrough
	case "pointsreboundsassists":
		fallthrough
	case "pts_rebs_asts":
		return PointsReboundsAssists
	case "pointsrebounds":
		fallthrough
	case "points_rebounds":
		fallthrough
	case "pts_rebs":
		return PointsRebounds
	case "points_assists":
		fallthrough
	case "pointsassists":
		fallthrough
	case "pts_asts":
		return PointsAssists
	case "rebounds_assists":
		fallthrough
	case "reboundsassists":
		fallthrough
	case "rebs_asts":
		return ReboundsAssists
	case "blocks_steals":
		fallthrough
	case "blks_stls":
		fallthrough
	case "blockssteals":
		return BlocksSteals
	case "fantasy_score":
		fallthrough
	case "fantasyscore":
		return FantasyScore
	case "height":
		return Height
	case "weight":
		return Weight
	case "double-double":
		fallthrough
	case "double_double":
		fallthrough
	case "doubledouble":
		return DoubleDouble
	default:
		logrus.Errorf("Unknown stat: '%s'", stat)
		return ""
	}
}

func (s *Stat) UnmarshalJSON(data []byte) error {
	*s = NewStat(string(data[:]))
	return nil
}

func (s *Stat) UnmarshalBSON(data []byte) error {
	*s = NewStat(string(data[:]))
	return nil
}

func (s *Stat) UnmarshalGQL(v interface{}) error {
	stat, ok := v.(string)
	if !ok {
		return fmt.Errorf("Stat must be a string")
	}

	*s = NewStat(stat)
	return nil
}

func (s Stat) MarshalGQL(w io.Writer) {
	switch s {
	case "points":
		io.WriteString(w, `"points"`)
	case "assists":
		io.WriteString(w, `"assists"`)
	case "rebounds":
		io.WriteString(w, `"rebounds"`)
	case "steals":
		io.WriteString(w, `"steals"`)
	case "blocks":
		io.WriteString(w, `"blocks"`)
	case "three_pointers_made":
		fallthrough
	case "threepointersmade":
		fallthrough
	case "3-pt_made":
		io.WriteString(w, `"three_pointers_made"`)
	case "three_pointers_attempted":
		fallthrough
	case "threepointersattempted":
		io.WriteString(w, `"three_pointers_attempted"`)
	case "three_point_percentage":
		fallthrough
	case "threepointpercentage":
		io.WriteString(w, `"three_point_percentage"`)
	case "free_throws_made":
		fallthrough
	case "freethrowsmade":
		io.WriteString(w, `"free_throws_made"`)
	case "free_throws_attempted":
		fallthrough
	case "freethrowsattempted":
		io.WriteString(w, `"free_throws_attempted"`)
	case "free_throws_percentage":
		fallthrough
	case "freethrowspercentage":
		io.WriteString(w, `"free_throws_percentage"`)
	case "field_goals_made":
		fallthrough
	case "fieldgoalsmade":
		io.WriteString(w, `"field_goals_made"`)
	case "field_goals_attempted":
		fallthrough
	case "fieldgoalsattempted":
		io.WriteString(w, `"field_goals_attempted"`)
	case "field_goal_percentage":
		fallthrough
	case "fieldgoalpercentage":
		io.WriteString(w, `"field_goal_percentage"`)
	case "effective_field_goal_percentage":
		fallthrough
	case "effectivefieldgoalpercentage":
		io.WriteString(w, `"effective_field_goal_percentage"`)
	case "true_shooting_percentage":
		fallthrough
	case "trueshootingpercentage":
		io.WriteString(w, `"true_shooting_percentage"`)
	case "minutes":
		io.WriteString(w, `"minutes"`)
	case "offensive_rebounds":
		fallthrough
	case "offensiverebounds":
		io.WriteString(w, `"offensive_rebounds"`)
	case "defensive_rebounds":
		fallthrough
	case "defensiverebounds":
		io.WriteString(w, `"defensive_rebounds"`)
	case "assist_percentage":
		fallthrough
	case "assistpercentage":
		io.WriteString(w, `"assist_percentage"`)
	case "offensive_rebound_percentage":
		fallthrough
	case "offensivereboundpercentage":
		io.WriteString(w, `"offensive_rebound_percentage"`)
	case "defensive_rebound_percentage":
		fallthrough
	case "defensivereboundpercentage":
		io.WriteString(w, `"defensive_rebound_percentage"`)
	case "usage":
		io.WriteString(w, `"usage"`)
	case "turnovers":
		io.WriteString(w, `"turnovers"`)
	case "games_played":
		fallthrough
	case "gamesplayed":
		io.WriteString(w, `"games_played"`)
	case "personal_fouls":
		fallthrough
	case "personalfouls":
		io.WriteString(w, `"personal_fouls"`)
	case "personal_fouls_drawn":
		fallthrough
	case "personalfoulsdrawn":
		io.WriteString(w, `"personal_fouls_drawn"`)
	case "points_rebounds_assists":
		fallthrough
	case "pointsreboundsassists":
		fallthrough
	case "pts_rebs_asts":
		io.WriteString(w, `"points_rebounds_assists"`)
	case "points_rebounds":
		fallthrough
	case "pointsrebounds":
		fallthrough
	case "pts_rebs":
		io.WriteString(w, `"points_rebounds"`)
	case "points_assists":
		fallthrough
	case "pointsassists":
		fallthrough
	case "pts_asts":
		io.WriteString(w, `"points_assists"`)
	case "rebounds_assists":
		fallthrough
	case "reboundsassists":
		fallthrough
	case "rebs_asts":
		io.WriteString(w, `"rebounds_assists"`)
	case "blocks_steals":
		fallthrough
	case "blockssteals":
		fallthrough
	case "blks_stls":
		io.WriteString(w, `"blocks_steals"`)
	case "fantasy_score":
		fallthrough
	case "fantasyscore":
		io.WriteString(w, `"fantasy_score"`)
	case "height":
		io.WriteString(w, `"height"`)
	case "weight":
		io.WriteString(w, `"weight"`)
	case "double_double":
		fallthrough
	case "doubledouble":
		io.WriteString(w, `"double_double"`)
	default:
		logrus.Errorf("Stat.MarshalGQL: unknown Stat: '%s'", s)
		io.WriteString(w, string(s))
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
		FantasyScore,
		PointsRebounds,
		PointsAssists,
		PointsReboundsAssists,
		ReboundsAssists,
		BlocksSteals,
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
