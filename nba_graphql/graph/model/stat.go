package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type Stat string

func (s Stat) String() string {
	return string(s)
}

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
	Passes                       Stat = "passes"
	PotentialAssists             Stat = "potential_assists"
	AssistConversionRate         Stat = "assist_conversion_rate"

	//TeamStats
	GamesWon    Stat = "games_won"
	GamesLost   Stat = "games_lost"
	OppPoints   Stat = "opponent_points"
	OppAssists  Stat = "opponent_assists"
	OppRebounds Stat = "opponent_rebounds"
)

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func NewStat(stat string) (Stat, error) {
	if stat == "" {
		return "", fmt.Errorf("empty stat")
	}
	// lookup := strings.ReplaceAll(strings.ToLower(stat), " ", "_")
	// lookup = strings.ReplaceAll(lookup, "+", "_")
	// lookup = strings.ReplaceAll(lookup, "-", "_")
	// lookup = strings.ReplaceAll(lookup, "_", "")
	// lookup = strings.ReplaceAll(lookup, "-", "")
	// lookup = strings.TrimSpace(lookup)
	lookup := strip(strings.ToLower(stat))
	lookup = strings.ReplaceAll(lookup, " ", "")
	switch lookup {
	case "points":
		return Points, nil
	case "assists":
		return Assists, nil
	case "rebounds":
		return Rebounds, nil
	case "steals":
		return Steals, nil
	// TODO: This might actually be blocks against...
	case "blocked_shots":
		fallthrough
	case "blockedshots":
		fallthrough
	case "blocks":
		return Blocks, nil
	case "three_pointers_made":
		fallthrough
	case "threepointersmade":
		fallthrough
	case "3-pt_made":
		fallthrough
	case "3_pt_made":
		fallthrough
	case "3ptmade":
		return ThreePointersMade, nil
	case "three_pointers_attempted":
		fallthrough
	case "threepointersattempted":
		return ThreePointersAttempted, nil
	case "three_point_percentage":
		fallthrough
	case "threepointpercentage":
		return ThreePointPercentage, nil
	case "free_throws_made":
		fallthrough
	case "freethrowsmade":
		return FreeThrowsMade, nil
	case "free_throws_attempted":
		fallthrough
	case "freethrowsattempted":
		return FreeThrowsAttempted, nil
	case "free_throws_percentage":
		fallthrough
	case "freethrowspercentage":
		return FreeThrowsPercentage, nil
	case "field_goals_made":
		fallthrough
	case "fieldgoalsmade":
		return FieldGoalsMade, nil
	case "field_goals_attempted":
		fallthrough
	case "fieldgoalsattempted":
		return FieldGoalsAttempted, nil
	case "field_goal_percentage":
		fallthrough
	case "fieldgoalpercentage":
		return FieldGoalPercentage, nil
	case "effective_field_goal_percentage":
		fallthrough
	case "effectivefieldgoalpercentage":
		return EffectiveFieldGoalPercentage, nil
	case "true_shooting_percentage":
		fallthrough
	case "trueshootingpercentage":
		return TrueShootingPercentage, nil
	case "minutes":
		return Minutes, nil
	case "offensive_rebounds":
		fallthrough
	case "offensiverebounds":
		return OffensiveRebounds, nil
	case "defensive_rebounds":
		fallthrough
	case "defensiverebounds":
		return DefensiveRebounds, nil
	case "assist_percentage":
		fallthrough
	case "assistpercentage":
		return AssistPercentage, nil
	case "offensive_rebound_percentage":
		fallthrough
	case "offensivereboundpercentage":
		return OffensiveReboundPercentage, nil
	case "defensive_rebound_percentage":
		fallthrough
	case "defensivereboundpercentage":
		return DefensiveReboundPercentage, nil
	case "usage":
		return Usage, nil
	case "turnovers":
		return Turnovers, nil
	case "gamesplayed":
		fallthrough
	case "games_played":
		return GamesPlayed, nil
	case "personal_fouls":
		fallthrough
	case "personalfouls":
		return PersonalFouls, nil
	case "personal_fouls_drawn":
		fallthrough
	case "personalfoulsdrawn":
		return PersonalFoulsDrawn, nil
	case "points_rebounds_assists":
		fallthrough
	case "pointsreboundsassists":
		fallthrough
	case "pts_rebs_asts":
		fallthrough
	case "ptsrebsasts":
		return PointsReboundsAssists, nil
	case "pointsrebounds":
		fallthrough
	case "points_rebounds":
		fallthrough
	case "pts_rebs":
		fallthrough
	case "ptsrebs":
		return PointsRebounds, nil
	case "points_assists":
		fallthrough
	case "pointsassists":
		fallthrough
	case "pts_asts":
		fallthrough
	case "ptsasts":
		return PointsAssists, nil
	case "rebounds_assists":
		fallthrough
	case "reboundsassists":
		fallthrough
	case "rebs_asts":
		fallthrough
	case "rebsasts":
		return ReboundsAssists, nil
	case "blocks_steals":
		fallthrough
	case "blks_stls":
		fallthrough
	case "blksstls":
		fallthrough
	case "blockssteals":
		return BlocksSteals, nil
	case "fantasy_score":
		fallthrough
	case "fantasyscore":
		return FantasyScore, nil
	case "height":
		return Height, nil
	case "weight":
		return Weight, nil
	case "double-double":
		fallthrough
	case "double_double":
		fallthrough
	case "doubledouble":
		return DoubleDouble, nil
	// Team Stats
	case "games_won":
		fallthrough
	case "gameswon":
		return GamesWon, nil
	case "games_lost":
		fallthrough
	case "gameslost":
		return GamesLost, nil
	case "opponent_points":
		fallthrough
	case "opponentpoints":
		fallthrough
	case "opppoints":
		fallthrough
	case "opp_points":
		return OppPoints, nil
	case "opponent_assists":
		fallthrough
	case "opponentassists":
		fallthrough
	case "oppassists":
		fallthrough
	case "opp_assists":
		return OppAssists, nil
	case "opponent_rebounds":
		fallthrough
	case "opponentrebounds":
		fallthrough
	case "opprebounds":
		fallthrough
	case "opp_rebounds":
		return OppRebounds, nil
	default:
		logrus.Errorf("Unknown new stat: '%s'", stat)
		return "", fmt.Errorf("Unknown new stat: '%s'", stat)
	}
}

func (s *Stat) UnmarshalJSON(data []byte) error {
	v, err := NewStat(string(data[:]))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

func (s *Stat) UnmarshalBSON(data []byte) error {
	v, err := NewStat(string(data[:]))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

func (s *Stat) UnmarshalGQL(v interface{}) error {
	stat, ok := v.(string)
	if !ok {
		return fmt.Errorf("Stat must be a string")
	}

	x, err := NewStat(stat)
	if err != nil {
		return err
	}
	*s = x
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
	// Team Stats
	case "games_won":
		fallthrough
	case "gameswon":
		io.WriteString(w, `"games_won"`)
	case "games_lost":
		fallthrough
	case "gameslost":
		io.WriteString(w, `"games_lost"`)
	case "opponent_points":
		fallthrough
	case "opponentpoints":
		io.WriteString(w, `"opponent_points"`)
	case "opponent_assists":
		fallthrough
	case "opponentassists":
		io.WriteString(w, `"opponent_assists"`)
	case "opponent_rebounds":
		fallthrough
	case "opponentrebounds":
		io.WriteString(w, `"opponent_rebounds"`)
	default:
		logrus.Errorf("Stat.MarshalGQL: unknown Stat: '%s'", s)
		io.WriteString(w, string(s))
	}
}
