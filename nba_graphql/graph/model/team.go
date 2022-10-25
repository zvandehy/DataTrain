package model

import (
	"math"
	"strings"

	"github.com/sirupsen/logrus"
	similarity "github.com/zvandehy/DataTrain/nba_graphql/math"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

type Team struct {
	Name         string      `json:"name" bson:"name"`
	TeamID       int         `json:"teamID" bson:"teamID"`
	Abbreviation string      `json:"abbreviation" bson:"abbreviation"`
	Location     string      `json:"location" bson:"city"`
	NumWins      int         `json:"numWins" bson:"numWins"`
	NumLoss      int         `json:"numLoss" bson:"numLoss"`
	League       string      `json:"league" bson:"league"`
	GamesCache   []*TeamGame `json:"gamesCache" bson:"gamesCache"`
}

func (t Team) String() string {
	return util.Print(t)
}

type TeamGame struct {
	Assists                              int     `json:"assists" bson:"assists"`
	Blocks                               int     `json:"blocks" bson:"blocks"`
	Date                                 string  `json:"date" bson:"date"`
	DefensiveRating                      float64 `json:"defensive_rating" bson:"defensive_rating"`
	DefensiveRebounds                    int     `json:"defensive_rebounds" bson:"defensive_rebounds"`
	DefensiveReboundPercentage           float64 `json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	FieldGoalPercentage                  float64 `json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted                  int     `json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade                       int     `json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted                  int     `json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade                       int     `json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage                 float64 `json:"free_throws_percentage" bson:"free_throws_percentage"`
	GameID                               string  `json:"gameID" bson:"gameID"`
	HomeOrAway                           string  `json:"home_or_away" bson:"home_or_away"`
	Margin                               int     `json:"margin" bson:"margin"`
	OffensiveRebounds                    int     `json:"offensive_rebounds" bson:"offensive_rebounds"`
	OffensiveReboundPercentage           float64 `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	TeamID                               int     `json:"teamID" bson:"teamID"`
	OpponentID                           int     `json:"opponent" bson:"opponent"`
	OpponentAssists                      int     `json:"opponent_assists" bson:"opponent_assists"`
	OpponentEffectiveFieldGoalPercentage float64 `json:"opponent_effective_field_goal_percentage" bson:"opponent_effective_field_goal_percentage"`
	OpponentFieldGoalsAttempted          int     `json:"opponent_field_goals_attempted" bson:"opponent_field_goals_attempted"`
	OpponentFreeThrowsAttempted          int     `json:"opponent_free_throws_attempted" bson:"opponent_free_throws_attempted"`
	OpponentPoints                       int     `json:"opponent_points" bson:"opponent_points"`
	OpponentRebounds                     int     `json:"opponent_rebounds" bson:"opponent_rebounds"`
	OpponentThreePointersAttempted       int     `json:"opponent_three_pointers_attempted" bson:"opponent_three_pointers_attempted"`
	OpponentThreePointersMade            int     `json:"opponent_three_pointers_made" bson:"opponent_three_pointers_made"`
	PlusMinusPerHundred                  float64 `json:"plus_minus_per_hundred" bson:"plus_minus_per_hundred"`
	Points                               int     `json:"points" bson:"points"`
	Playoffs                             bool    `json:"playoffs" bson:"playoffs"`
	Possessions                          int     `json:"possessions" bson:"possessions"`
	PersonalFouls                        int     `json:"personal_fouls" bson:"personal_fouls"`
	PersonalFoulsDrawn                   int     `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	Rebounds                             int     `json:"rebounds" bson:"rebounds"`
	Season                               string  `json:"season" bson:"season"`
	Steals                               int     `json:"steals" bson:"steals"`
	ThreePointersAttempted               int     `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade                    int     `json:"three_pointers_made" bson:"three_pointers_made"`
	Turnovers                            int     `json:"turnovers" bson:"turnovers"`
	WinOrLoss                            string  `json:"win_or_loss" bson:"win_or_loss"`
}

func (t TeamGame) String() string {
	return util.Print(t)
}

type TeamAverage struct {
	WinsAndLosses      []string `json:"wins_and_losses" bson:"wins_and_losses"`
	GamesPlayed        float64  `json:"games_played" bson:"games_played"`
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

func NewTeamAverage(games []*TeamGame, team *Team) TeamAverage {
	average := TeamAverage{}
	average.Team = *team

	for _, game := range games {
		average.Assists += float64(game.Assists)
		average.Blocks += float64(game.Blocks)
		average.PersonalFoulsDrawn += float64(game.PersonalFoulsDrawn)
		average.PersonalFouls += float64(game.PersonalFouls)
		average.Points += float64(game.Points)
		average.Rebounds += float64(game.Rebounds)
		average.Steals += float64(game.Steals)
		average.ThreePointersMade += float64(game.ThreePointersMade)
		average.Turnovers += float64(game.Turnovers)
		average.OppAssists += float64(game.OpponentAssists)
		average.OppPoints += float64(game.OpponentPoints)
		average.OppRebounds += float64(game.OpponentRebounds)
		if strings.ToLower(game.WinOrLoss)[0] == 'w' {
			average.GamesWon++
		}
		if strings.ToLower(game.WinOrLoss)[0] == 'l' {
			average.GamesLost++
		}
		average.GamesPlayed++
	}

	average.Assists /= float64(len(games))
	average.Blocks /= float64(len(games))
	average.PersonalFoulsDrawn /= float64(len(games))
	average.PersonalFouls /= float64(len(games))
	average.Points /= float64(len(games))
	average.Rebounds /= float64(len(games))
	average.Steals /= float64(len(games))
	average.ThreePointersMade /= float64(len(games))
	average.Turnovers /= float64(len(games))
	average.OppAssists /= float64(len(games))
	average.OppPoints /= float64(len(games))
	average.OppRebounds /= float64(len(games))
	return average
}

func (t *TeamAverage) Score(stat Stat) float64 {
	switch NewStat(string(stat)) {
	case Points:
		return float64(t.Points)
	case Assists:
		return float64(t.Assists)
	case Rebounds:
		return float64(t.Rebounds)
	case Steals:
		return float64(t.Steals)
	case Blocks:
		return float64(t.Blocks)
	case ThreePointersMade:
		return float64(t.ThreePointersMade)
	case Turnovers:
		return float64(t.Turnovers)
	case PersonalFouls:
		return float64(t.PersonalFouls)
	case PersonalFoulsDrawn:
		return float64(t.PersonalFoulsDrawn)
	case OppPoints:
		return float64(t.OppPoints)
	case OppAssists:
		return float64(t.OppAssists)
	case OppRebounds:
		return float64(t.OppRebounds)
	case GamesWon:
		return float64(t.GamesWon)
	case GamesLost:
		return float64(t.GamesLost)
	case GamesPlayed:
		return float64(t.GamesPlayed)
	default:
		logrus.Warnf("Unknown team stat: '%s'", stat)
		return 0
	}
}

func (t *TeamAverage) Normalize(stats ...StatOfInterest) TeamAverage {
	normalized := *t
	for _, stat := range stats {
		switch stat.Stat {
		case Points:
			normalized.Points = stat.ZScoreTeam(*t)
		case Rebounds:
			normalized.Rebounds = stat.ZScoreTeam(*t)
		case Assists:
			normalized.Assists = stat.ZScoreTeam(*t)
		case Steals:
			normalized.Steals = stat.ZScoreTeam(*t)
		case Blocks:
			normalized.Blocks = stat.ZScoreTeam(*t)
		case ThreePointersMade:
			normalized.ThreePointersMade = stat.ZScoreTeam(*t)
		case PersonalFouls:
			normalized.PersonalFouls = stat.ZScoreTeam(*t)
		case PersonalFoulsDrawn:
			normalized.PersonalFoulsDrawn = stat.ZScoreTeam(*t)
		case Turnovers:
			normalized.Turnovers = stat.ZScoreTeam(*t)
		case GamesPlayed:
			normalized.GamesPlayed = stat.ZScoreTeam(*t)
		case GamesWon:
			normalized.GamesWon = stat.ZScoreTeam(*t)
		case GamesLost:
			normalized.GamesLost = stat.ZScoreTeam(*t)
		case OppPoints:
			normalized.OppPoints = stat.ZScoreTeam(*t)
		case OppAssists:
			normalized.OppAssists = stat.ZScoreTeam(*t)
		case OppRebounds:
			normalized.OppRebounds = stat.ZScoreTeam(*t)
		}
	}
	return normalized
}

type TeamDiff TeamAverage

func (t *TeamAverage) Difference(fromTeam TeamAverage) TeamDiff {
	d := TeamDiff{
		Assists:            similarity.RoundFloat(fromTeam.Assists-t.Assists, 2),
		Blocks:             similarity.RoundFloat(fromTeam.Blocks-t.Blocks, 2),
		GamesPlayed:        similarity.RoundFloat(fromTeam.GamesPlayed-t.GamesPlayed, 2),
		Points:             similarity.RoundFloat(fromTeam.Points-t.Points, 2),
		Rebounds:           similarity.RoundFloat(fromTeam.Rebounds-t.Rebounds, 2),
		Steals:             similarity.RoundFloat(fromTeam.Steals-t.Steals, 2),
		ThreePointersMade:  similarity.RoundFloat(fromTeam.ThreePointersMade-t.ThreePointersMade, 2),
		Turnovers:          similarity.RoundFloat(fromTeam.Turnovers-t.Turnovers, 2),
		PersonalFoulsDrawn: similarity.RoundFloat(fromTeam.PersonalFoulsDrawn-t.PersonalFoulsDrawn, 2),
		PersonalFouls:      similarity.RoundFloat(fromTeam.PersonalFouls-t.PersonalFouls, 2),
		OppAssists:         similarity.RoundFloat(fromTeam.OppAssists-t.OppAssists, 2),
		OppPoints:          similarity.RoundFloat(fromTeam.OppPoints-t.OppPoints, 2),
		OppRebounds:        similarity.RoundFloat(fromTeam.OppRebounds-t.OppRebounds, 2),
		GamesLost:          similarity.RoundFloat(fromTeam.GamesLost-t.GamesLost, 2),
		GamesWon:           similarity.RoundFloat(fromTeam.GamesWon-t.GamesWon, 2),
		Team:               t.Team}
	return d
}

func TeamAverageStats() []Stat {
	return []Stat{
		GamesWon,
		GamesLost,
		Points,
		OppPoints,
		Assists,
		OppAssists,
		Rebounds,
		OppRebounds,
		Steals,
		Blocks,
		Turnovers,
		ThreePointersMade,
		PersonalFouls,
		PersonalFoulsDrawn,
	}
}

func TeamEuclideanDistance(diff TeamDiff, statsOfInterest []Stat) float64 {
	// TODO: could add user-inputed weights for different stats
	sum := 0.0
	for _, stat := range statsOfInterest {
		switch stat {
		case Assists:
			sum += math.Pow(diff.Assists, 2)
		case Blocks:
			sum += math.Pow(diff.Blocks, 2)
		case GamesPlayed:
			sum += math.Pow(diff.GamesPlayed, 2)
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
		case ThreePointersMade:
			sum += math.Pow(diff.ThreePointersMade, 2)
		case Turnovers:
			sum += math.Pow(diff.Turnovers, 2)
		case OppAssists:
			sum += math.Pow(diff.OppAssists, 2)
		case OppPoints:
			sum += math.Pow(diff.OppPoints, 2)
		case OppRebounds:
			sum += math.Pow(diff.OppRebounds, 2)
		case GamesLost:
			sum += math.Pow(diff.GamesLost, 2)
		case GamesWon:
			sum += math.Pow(diff.GamesWon, 2)
		default:
			sum += 0.0
		}
	}
	return math.Sqrt(sum)
}

type PlayersInGame struct {
	TeamPlayers     []*Player `json:"team" bson:"team"`
	OpponentPlayers []*Player `json:"opponent" bson:"opponent"`
}
