package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// Team is the resolver for the team field.
func (r *playerResolver) Team(ctx context.Context, obj *model.Player) (*model.Team, error) {
	if obj.CurrentTeam == "" && obj.TeamID == 0 {
		logrus.Errorf("Player is missing current team: %v", obj)
		return &model.Team{}, fmt.Errorf("Player is missing current team: %v", obj)
	}
	team, err := dataloader.For(ctx).TeamByAbr.Load(obj.CurrentTeam)
	if err != nil {
		logrus.Errorf("Error loading team '%v' from: %v", obj.CurrentTeam, obj)
		return nil, err
	}
	return team, err
}

// Games is the resolver for the games field.
func (r *playerResolver) Games(ctx context.Context, obj *model.Player, input model.GameFilter) ([]*model.PlayerGame, error) {
	games := []*model.PlayerGame{}
	for _, game := range obj.GamesCache {
		if input.MatchPlayerGame(game) {
			games = append(games, game)
		}
	}
	return games, nil
}

// Image is the resolver for the image field.
func (r *playerResolver) Image(ctx context.Context, obj *model.Player) (string, error) {
	return fmt.Sprintf("https://ak-static.cms.nba.com/wp-content/uploads/headshots/%s/latest/260x190/%d.png", strings.ToLower(obj.League), obj.PlayerID), nil
}

// AssistPercentage is the resolver for the assist_percentage field.
func (r *playerGameResolver) AssistPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.AssistPercentage.Float64, nil
}

// Assists is the resolver for the assists field.
func (r *playerGameResolver) Assists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Assists.Valid {
		return int(obj.Assists.Int16), nil
	}
	return 0, nil
}

// Passes is the resolver for the passes field.
func (r *playerGameResolver) Passes(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Passes.Valid {
		return int(obj.Passes.Int16), nil
	}
	return 0, nil
}

// PotentialAssists is the resolver for the potential_assists field.
func (r *playerGameResolver) PotentialAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.PotentialAssists.Valid {
		return int(obj.PotentialAssists.Int16), nil
	}
	return 0, nil
}

// AssistConversionRate is the resolver for the assist_conversion_rate field.
func (r *playerGameResolver) AssistConversionRate(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return float64(int(obj.AssistConversionRate.Float64)) / 100, nil
}

// Date is the resolver for the date field.
func (r *playerGameResolver) Date(ctx context.Context, obj *model.PlayerGame) (string, error) {
	return obj.Date.Format(util.DATE_FORMAT), nil
}

// DefensiveReboundPercentage is the resolver for the defensive_rebound_percentage field.
func (r *playerGameResolver) DefensiveReboundPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.DefensiveReboundPercentage.Float64, nil
}

// DefensiveRebounds is the resolver for the defensive_rebounds field.
func (r *playerGameResolver) DefensiveRebounds(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.DefensiveRebounds.Valid {
		return int(obj.DefensiveRebounds.Int16), nil
	}
	return 0, nil
}

// EffectiveFieldGoalPercentage is the resolver for the effective_field_goal_percentage field.
func (r *playerGameResolver) EffectiveFieldGoalPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.EffectiveFieldGoalPercentage.Float64, nil
}

// FieldGoalPercentage is the resolver for the field_goal_percentage field.
func (r *playerGameResolver) FieldGoalPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.FieldGoalPercentage.Float64, nil
}

// FieldGoalsAttempted is the resolver for the field_goals_attempted field.
func (r *playerGameResolver) FieldGoalsAttempted(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.FieldGoalsAttempted.Valid {
		return int(obj.FieldGoalsAttempted.Int16), nil
	}
	return 0, nil
}

// FieldGoalsMade is the resolver for the field_goals_made field.
func (r *playerGameResolver) FieldGoalsMade(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.FieldGoalsMade.Valid {
		return int(obj.FieldGoalsMade.Int16), nil
	}
	return 0, nil
}

// FreeThrowsAttempted is the resolver for the free_throws_attempted field.
func (r *playerGameResolver) FreeThrowsAttempted(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.FreeThrowsAttempted.Valid {
		return int(obj.FreeThrowsAttempted.Int16), nil
	}
	return 0, nil
}

// FreeThrowsMade is the resolver for the free_throws_made field.
func (r *playerGameResolver) FreeThrowsMade(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.FreeThrowsMade.Valid {
		return int(obj.FreeThrowsMade.Int16), nil
	}
	return 0, nil
}

// FreeThrowsPercentage is the resolver for the free_throws_percentage field.
func (r *playerGameResolver) FreeThrowsPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.FreeThrowsPercentage.Float64, nil
}

// HomeOrAway is the resolver for the home_or_away field.
func (r *playerGameResolver) HomeOrAway(ctx context.Context, obj *model.PlayerGame) (string, error) {
	return string(obj.HomeOrAway), nil
}

// Outcome is the resolver for the outcome field.
func (r *playerGameResolver) Outcome(ctx context.Context, obj *model.PlayerGame) (model.GameOutcome, error) {
	outcome := strings.ToLower(obj.Outcome)
	switch outcome[0] {
	case 'w':
		return model.GameOutcomeWin, nil
	case 'l':
		return model.GameOutcomeLoss, nil
	case 'p':
		return model.GameOutcomePending, nil
	default:
		return model.GameOutcomePending, fmt.Errorf("unknown outcome: %v", outcome)
	}
}

// Margin is the resolver for the margin field.
func (r *playerGameResolver) Margin(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Margin.Valid {
		return int(obj.Margin.Int16), nil
	}
	return 0, nil
}

// Minutes is the resolver for the minutes field.
func (r *playerGameResolver) Minutes(ctx context.Context, obj *model.PlayerGame) (string, error) {
	// convert minutes decimal to string
	minutes := int(obj.Minutes)
	seconds := int((obj.Minutes - float64(minutes)) * 60)
	return fmt.Sprintf("%d:%02d", minutes, seconds), nil
}

// OffensiveReboundPercentage is the resolver for the offensive_rebound_percentage field.
func (r *playerGameResolver) OffensiveReboundPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.OffensiveReboundPercentage.Float64, nil
}

// OffensiveRebounds is the resolver for the offensive_rebounds field.
func (r *playerGameResolver) OffensiveRebounds(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.OffensiveRebounds.Valid {
		return int(obj.OffensiveRebounds.Int16), nil
	}
	return 0, nil
}

// Opponent is the resolver for the opponent field.
func (r *playerGameResolver) Opponent(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	team, err := dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
	if err != nil {
		logrus.Errorf("Error loading opponent '%v' from: %v", obj.OpponentID, obj)
		return nil, err
	}
	return team, err
}

// OpponentStats is the resolver for the opponentStats field.
func (r *playerGameResolver) OpponentStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
	panic(fmt.Errorf("(r *playerGameResolver) OpponentStats not implemented"))
}

// Team is the resolver for the team field.
func (r *playerGameResolver) Team(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	team, err := dataloader.For(ctx).TeamByID.Load(obj.TeamID)
	if err != nil {
		logrus.Errorf("Error loading team '%v' from: %v", obj.TeamID, obj)
		return nil, err
	}
	return team, err
}

// TeamStats is the resolver for the teamStats field.
func (r *playerGameResolver) TeamStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
	panic(fmt.Errorf("(r *playerGameResolver) TeamStats not implemented"))
}

// PersonalFoulsDrawn is the resolver for the personal_fouls_drawn field.
func (r *playerGameResolver) PersonalFoulsDrawn(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.PersonalFoulsDrawn.Valid {
		return int(obj.PersonalFoulsDrawn.Int16), nil
	}
	return 0, nil
}

// PersonalFouls is the resolver for the personal_fouls field.
func (r *playerGameResolver) PersonalFouls(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.PersonalFouls.Valid {
		return int(obj.PersonalFouls.Int16), nil
	}
	return 0, nil
}

// Points is the resolver for the points field.
func (r *playerGameResolver) Points(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Points.Valid {
		return int(obj.Points.Int16), nil
	}
	return 0, nil
}

// Player is the resolver for the player field.
func (r *playerGameResolver) Player(ctx context.Context, obj *model.PlayerGame) (*model.Player, error) {
	player, err := dataloader.For(ctx).PlayerByID.Load(obj.PlayerID)
	if err != nil {
		logrus.Errorf("Error loading player '%v' from: %v | %v", obj.PlayerID, obj, err)
		return nil, err
	}
	if player == nil {
		logrus.Errorf("Unable to find player '%v' from: %v", obj.PlayerID, obj)
		return nil, fmt.Errorf("unable to find player '%v' from game: %v %v", obj.PlayerID, obj.GameID, obj.Date)
	}
	return player, err
}

// ThreePointPercentage is the resolver for the three_point_percentage field.
func (r *playerGameResolver) ThreePointPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.ThreePointPercentage.Float64, nil
}

// ThreePointersAttempted is the resolver for the three_pointers_attempted field.
func (r *playerGameResolver) ThreePointersAttempted(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.ThreePointersAttempted.Valid {
		return int(obj.ThreePointersAttempted.Int16), nil
	}
	return 0, nil
}

// ThreePointersMade is the resolver for the three_pointers_made field.
func (r *playerGameResolver) ThreePointersMade(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.ThreePointersMade.Valid {
		return int(obj.ThreePointersMade.Int16), nil
	}
	return 0, nil
}

// Rebounds is the resolver for the rebounds field.
func (r *playerGameResolver) Rebounds(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Rebounds.Valid {
		return int(obj.Rebounds.Int16), nil
	}
	return 0, nil
}

// TrueShootingPercentage is the resolver for the true_shooting_percentage field.
func (r *playerGameResolver) TrueShootingPercentage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.TrueShootingPercentage.Float64, nil
}

// Turnovers is the resolver for the turnovers field.
func (r *playerGameResolver) Turnovers(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Turnovers.Valid {
		return int(obj.Turnovers.Int16), nil
	}
	return 0, nil
}

// Usage is the resolver for the usage field.
func (r *playerGameResolver) Usage(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return obj.Usage.Float64, nil
}

// Blocks is the resolver for the blocks field.
func (r *playerGameResolver) Blocks(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Blocks.Valid {
		return int(obj.Blocks.Int16), nil
	}
	return 0, nil
}

// Steals is the resolver for the steals field.
func (r *playerGameResolver) Steals(ctx context.Context, obj *model.PlayerGame) (int, error) {
	if obj.Steals.Valid {
		return int(obj.Steals.Int16), nil
	}
	return 0, nil
}

// PointsRebounds is the resolver for the points_rebounds field.
func (r *playerGameResolver) PointsRebounds(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return int(obj.Points.Int16 + obj.Rebounds.Int16), nil
}

// PointsAssists is the resolver for the points_assists field.
func (r *playerGameResolver) PointsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return int(obj.Points.Int16 + obj.Assists.Int16), nil
}

// PointsReboundsAssists is the resolver for the points_rebounds_assists field.
func (r *playerGameResolver) PointsReboundsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return int(obj.Points.Int16 + obj.Rebounds.Int16 + obj.Assists.Int16), nil
}

// ReboundsAssists is the resolver for the rebounds_assists field.
func (r *playerGameResolver) ReboundsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return int(obj.Rebounds.Int16 + obj.Assists.Int16), nil
}

// BlocksSteals is the resolver for the blocks_steals field.
func (r *playerGameResolver) BlocksSteals(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return int(obj.Blocks.Int16 + obj.Steals.Int16), nil
}

// FantasyScore is the resolver for the fantasy_score field.
func (r *playerGameResolver) FantasyScore(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return float64(obj.Points.Int16) + float64(obj.Rebounds.Int16)*1.2 + float64(obj.Assists.Int16)*1.5 + float64(obj.Steals.Int16)*3.0 + float64(obj.Blocks.Int16)*3.0 - float64(obj.Turnovers.Int16), nil
}

// Prediction is the resolver for the prediction field.
func (r *playerGameResolver) Prediction(ctx context.Context, obj *model.PlayerGame, input model.ModelInput) (*model.GamePrediction, error) {
	// TODO: ADD DATALOADER
	// propositions, err := r.GetPropositionsByPlayerGame(ctx, *obj)
	// if err != nil {
	// 	logrus.Errorf("Error loading propositions for game %v", obj)
	// 	return nil, err
	// }

	// startDate := &[]time.Time{util.SEASON_DATE(util.SEASON_START_2022_23)}[0]
	// for _, breakdown := range input.GameBreakdowns {
	// 	if breakdown.Filter != nil && breakdown.Filter.StartDate != nil {
	// 		start, err := time.Parse(util.DATE_FORMAT, *breakdown.Filter.StartDate)
	// 		if err != nil {
	// 			logrus.Errorf("Error parsing start date %v", *breakdown.Filter.StartDate)
	// 			return nil, err
	// 		}
	// 		if start.Before(*startDate) {
	// 			startDate = &start
	// 		}
	// 	}
	// 	if breakdown.Filter.Seasons != nil && len(*breakdown.Filter.Seasons) > 0 {
	// 		for _, season := range *breakdown.Filter.Seasons {
	// 			seasonStart := util.SEASON_DATE(string(season))
	// 			if seasonStart.Before(*startDate) {
	// 				startDate = &seasonStart
	// 			}
	// 		}
	// 	}
	// }
	// if input.SimilarPlayerInput != nil {
	// 	var err error
	// 	startDate, err = input.SimilarPlayerInput.PlayerPoolFilter.GetEarliestSeasonStartDate()
	// 	if err != nil {
	// 		logrus.Errorf("Error getting earliest season start date: %v", err)
	// 		return nil, err
	// 	}
	// }
	// startDateStr := startDate.Format(util.DATE_FORMAT)

	// totalPrediction := model.AverageStats{}
	// gameBreakdownFragments := []*model.PredictionFragment{}
	// similarPlayerFragments := []*model.PredictionFragment{}
	// similarTeamFragments := []*model.PredictionFragment{}

	// // Player GameLog Breakdowns
	// // Base (Filter, Games, Avg) for all of this player's games from the start of the range to the game date
	// date := obj.Date.Format(util.DATE_FORMAT)
	// baseFilter := model.GameFilter{
	// 	StartDate: &startDateStr,
	// 	EndDate:   &date,
	// }
	// playerBaseGames := []*model.PlayerGame{}
	// for _, game := range obj.PlayerRef.GamesCache {
	// 	if baseFilter.MatchPlayerGame(game) {
	// 		playerBaseGames = append(playerBaseGames, game)
	// 	}
	// }
	// // sort playerBaseGames by game.date
	// sort.Slice(playerBaseGames, func(i, j int) bool {
	// 	// a, err := time.Parse(util.DATE_FORMAT, playerBaseGames[i].Date)
	// 	// if err != nil {
	// 	// 	logrus.Errorf("Error parsing date %v", playerBaseGames[i].Date)
	// 	// 	return false
	// 	// }
	// 	// b, err := time.Parse(util.DATE_FORMAT, playerBaseGames[j].Date)
	// 	// if err != nil {
	// 	// 	logrus.Errorf("Error parsing date %v", playerBaseGames[j].Date)
	// 	// 	return false
	// 	// }
	// 	return playerBaseGames[i].Date.Before(*playerBaseGames[j].Date)
	// })
	// playerBase := &model.AverageStats{}
	// if len(playerBaseGames) > 0 {
	// 	b := model.NewPlayerAverage(playerBaseGames, obj.PlayerRef)
	// 	playerBase = b.AverageStats()
	// }

	// distributeExtraWeight := 0.0
	// for i := range input.GameBreakdowns {
	// 	games := []*model.PlayerGame{}
	// 	for _, game := range playerBaseGames {
	// 		// only get games before the current game
	// 		d := game.Date.Format(util.DATE_FORMAT)
	// 		input.GameBreakdowns[i].Filter.EndDate = &d
	// 		// if opponentMatch is true, only get games previously matched up against opponent
	// 		if input.GameBreakdowns[i].Filter.OpponentMatch != nil && *input.GameBreakdowns[i].Filter.OpponentMatch {
	// 			input.GameBreakdowns[i].Filter.OpponentID = &obj.OpponentID
	// 		}
	// 		// TODO: add more filters, like home/away, playeoff/regular season, etc.

	// 		// if the game matches the filter, add it to the list of games for this breakdown
	// 		if input.GameBreakdowns[i].Filter.MatchPlayerGame(game) {
	// 			games = append(games, game)
	// 			if input.GameBreakdowns[i].Filter.LastX != nil && len(games) >= *input.GameBreakdowns[i].Filter.LastX {
	// 				break
	// 			}
	// 		}
	// 	}
	// 	// player average for the games in this breakdown
	// 	pAvg := model.NewPlayerAverage(games, obj.PlayerRef)

	// 	// Don't calculate averages if there are no games, or less than lastX games
	// 	if len(games) == 0 || (input.GameBreakdowns[i].Filter.LastX != nil && len(games) < *input.GameBreakdowns[i].Filter.LastX) {
	// 		distributeExtraWeight += input.GameBreakdowns[i].Weight
	// 		input.GameBreakdowns[i].Weight = 0
	// 	} else {
	// 		derived := pAvg.AverageStats()
	// 		frag := &model.PredictionFragment{
	// 			Name:         input.GameBreakdowns[i].Name,
	// 			Derived:      derived,
	// 			DerivedGames: games,
	// 			Base:         playerBase,
	// 			PctChange:    playerBase.PercentChange(derived),
	// 			Weight:       input.GameBreakdowns[i].Weight,
	// 		}
	// 		for _, proposition := range propositions {
	// 			analysis := model.PropositionSummary{}
	// 			for _, game := range games {
	// 				propScore := game.Score(proposition.Type)
	// 				if propScore > proposition.Target {
	// 					analysis.NumOver++
	// 				}
	// 				if propScore < proposition.Target {
	// 					analysis.NumUnder++
	// 				}
	// 				if propScore == proposition.Target {
	// 					analysis.NumPush++
	// 				}
	// 			}
	// 			analysis.PctOver = float64(analysis.NumOver) / float64(len(games))
	// 			analysis.PctUnder = float64(analysis.NumUnder) / float64(len(games))
	// 			analysis.PctPush = float64(analysis.NumPush) / float64(len(games))
	// 			prop := *proposition
	// 			prop.Analysis = &analysis
	// 			frag.Propositions = append(frag.Propositions, &prop)
	// 		}
	// 		gameBreakdownFragments = append(gameBreakdownFragments, frag)
	// 	}
	// }

	// countSimilarTeamsWithGames := 0
	// // Similar Team Breakdown
	// if input.SimilarTeamInput != nil {
	// 	// only get games before the current game (or the inputted end date if that is earlier)
	// 	d := obj.Date.Format(util.DATE_FORMAT)
	// 	if input.SimilarTeamInput.Period.EndDate == nil {
	// 		input.SimilarTeamInput.Period.EndDate = &d
	// 	} else {
	// 		end, err := time.Parse(util.DATE_FORMAT, *input.SimilarTeamInput.Period.EndDate)
	// 		if err != nil {
	// 			logrus.Errorf("Error parsing end date: %v", err)
	// 			input.SimilarTeamInput.Period.EndDate = &d
	// 		}
	// 		// date, err := time.Parse(util.DATE_FORMAT, d)
	// 		// if err != nil {
	// 		// 	logrus.Errorf("Error parsing game date: %v", err)
	// 		// 	return nil, err
	// 		// }
	// 		if end.After(*obj.Date) {
	// 			input.SimilarTeamInput.Period.EndDate = &d
	// 		}
	// 	}
	// 	similarTeams, err := r.GetSimilarTeams(ctx, obj.OpponentID, input.SimilarTeamInput, d)
	// 	if err != nil {
	// 		logrus.Errorf("Error getting similar teams: %v", err)
	// 		return nil, err
	// 	}
	// 	for _, team := range similarTeams {
	// 		games := []*model.PlayerGame{}
	// 		for _, game := range obj.PlayerRef.GamesCache {
	// 			if input.SimilarTeamInput.Period.MatchGame(game) && game.OpponentID == team.TeamID {
	// 				games = append(games, game)
	// 			}
	// 		}
	// 		// Don't calculate averages if there are no games vs team
	// 		if len(games) > 0 {
	// 			countSimilarTeamsWithGames++
	// 			// player average for the games vs similar team
	// 			pAvg := model.NewPlayerAverage(games, obj.PlayerRef)
	// 			derived := pAvg.AverageStats()
	// 			frag := &model.PredictionFragment{
	// 				Name:         fmt.Sprintf("vs %v", team.Abbreviation),
	// 				Derived:      derived,
	// 				DerivedGames: games,
	// 				Base:         playerBase,
	// 				PctChange:    playerBase.PercentChange(derived),
	// 				Weight:       input.SimilarTeamInput.Weight,
	// 			}
	// 			for _, proposition := range propositions {
	// 				analysis := model.PropositionSummary{}
	// 				for _, game := range games {
	// 					propScore := game.Score(proposition.Type)
	// 					if propScore > proposition.Target {
	// 						analysis.NumOver++
	// 					}
	// 					if propScore < proposition.Target {
	// 						analysis.NumUnder++
	// 					}
	// 					if propScore == proposition.Target {
	// 						analysis.NumPush++
	// 					}
	// 				}
	// 				analysis.PctOver = float64(analysis.NumOver) / float64(len(games))
	// 				analysis.PctUnder = float64(analysis.NumUnder) / float64(len(games))
	// 				analysis.PctPush = float64(analysis.NumPush) / float64(len(games))
	// 				prop := *proposition
	// 				prop.Analysis = &analysis
	// 				frag.Propositions = append(frag.Propositions, &prop)
	// 			}
	// 			similarTeamFragments = append(similarTeamFragments, frag)
	// 		} else {
	// 			similarTeamFragments = append(similarTeamFragments, &model.PredictionFragment{
	// 				Name:         fmt.Sprintf("vs %v (None)", team.Abbreviation),
	// 				Derived:      &model.AverageStats{},
	// 				DerivedGames: []*model.PlayerGame{},
	// 				Base:         playerBase,
	// 				PctChange:    &model.AverageStats{},
	// 				Weight:       0,
	// 			})
	// 		}
	// 	}
	// 	// if there are no games vs similar teams, don't use the similar team input and distribute that weight across the other inputs
	// 	if countSimilarTeamsWithGames == 0 {
	// 		distributeExtraWeight += input.SimilarTeamInput.Weight
	// 	} else {
	// 		// distribute similar team weights evenly between all similar team breakdowns
	// 		similarTeamWeights := input.SimilarTeamInput.Weight / float64(countSimilarTeamsWithGames)
	// 		for i := range similarTeamFragments {
	// 			if similarTeamFragments[i].Weight > 0 {
	// 				similarTeamFragments[i].Weight = similarity.RoundFloat(similarTeamWeights, 2)
	// 			}
	// 		}
	// 	}
	// }

	// countSimilarPlayersWithGamesVsOpp := 0
	// if input.SimilarPlayerInput != nil {
	// 	d := obj.Date.Format(util.DATE_FORMAT)
	// 	// gets X similar players to the current player, where X is defined by the input limit
	// 	// TODO: Refactor function to use time.Time instead of string
	// 	similarPlayers, err := r.GetSimilarPlayers(ctx, obj.PlayerID, input.SimilarPlayerInput, d)
	// 	if err != nil {
	// 		logrus.Errorf("Error getting similar players: %v", err)
	// 		return nil, err
	// 	}
	// 	for _, player := range similarPlayers {
	// 		// get the similar player's games from the start of the range to the game date
	// 		baseGames := []*model.PlayerGame{}
	// 		for _, game := range player.GamesCache {
	// 			if baseFilter.MatchPlayerGame(game) {
	// 				baseGames = append(baseGames, game)
	// 			}
	// 		}
	// 		b := model.NewPlayerAverage(baseGames, player)
	// 		baseAvg := b.AverageStats()

	// 		// get the similar player's games vs the matchup opponent
	// 		opponentFilter := model.GameFilter{
	// 			EndDate:    &d,
	// 			OpponentID: &obj.OpponentID,
	// 		}
	// 		if len(*input.SimilarPlayerInput.PlayerPoolFilter.Seasons) > 0 {
	// 			opponentFilter.Seasons = input.SimilarPlayerInput.PlayerPoolFilter.Seasons
	// 		}
	// 		matchupGames := []*model.PlayerGame{}
	// 		for _, game := range player.GamesCache {
	// 			if opponentFilter.MatchPlayerGame(game) {
	// 				matchupGames = append(matchupGames, game)
	// 			}
	// 		}
	// 		// Don't calculate averages if there are no games vs opponent
	// 		if len(matchupGames) > 0 {
	// 			countSimilarPlayersWithGamesVsOpp++
	// 			pAvg := model.NewPlayerAverage(matchupGames, player)
	// 			derived := pAvg.AverageStats()
	// 			similarPlayerFragments = append(similarPlayerFragments, &model.PredictionFragment{
	// 				Name:         fmt.Sprintf("%v (%s) vs Opp", player.Name, player.Position),
	// 				Derived:      derived,
	// 				DerivedGames: matchupGames,
	// 				Base:         baseAvg,
	// 				PctChange:    baseAvg.PercentChange(derived),
	// 				Weight:       input.SimilarPlayerInput.Weight,
	// 			})
	// 		} else {
	// 			similarPlayerFragments = append(similarPlayerFragments, &model.PredictionFragment{
	// 				Name:         fmt.Sprintf("%v (%s) vs Opp (None)", player.Name, player.Position),
	// 				Derived:      &model.AverageStats{},
	// 				DerivedGames: matchupGames,
	// 				Base:         baseAvg,
	// 				PctChange:    &model.AverageStats{},
	// 				Weight:       0,
	// 			})
	// 		}
	// 	}
	// 	// if there are no similar players with games vs the opponent, don't use the similar player input and distribute that weight across the other inputs
	// 	if countSimilarPlayersWithGamesVsOpp == 0 {
	// 		distributeExtraWeight += input.SimilarPlayerInput.Weight
	// 	} else {
	// 		// distribute similar player weights evenly between all similar players with games vs Opp
	// 		similarPlayerWeights := input.SimilarPlayerInput.Weight / float64(countSimilarPlayersWithGamesVsOpp)
	// 		for i := range similarPlayerFragments {
	// 			if similarPlayerFragments[i].Weight > 0 {
	// 				similarPlayerFragments[i].Weight = similarity.RoundFloat(similarPlayerWeights, 2)
	// 			}
	// 		}
	// 	}
	// }

	// //distribute extra weight evenly among all valid breakdowns
	// if distributeExtraWeight > 0 {
	// 	distributeBetween := len(gameBreakdownFragments)
	// 	if countSimilarTeamsWithGames > 0 {
	// 		distributeBetween += countSimilarTeamsWithGames
	// 	}
	// 	if len(similarPlayerFragments) > 0 {
	// 		distributeBetween += countSimilarPlayersWithGamesVsOpp
	// 	}
	// 	extraWeight := similarity.RoundFloat(distributeExtraWeight/float64(distributeBetween), 2)
	// 	for i := range gameBreakdownFragments {
	// 		gameBreakdownFragments[i].Weight += extraWeight
	// 	}
	// 	for i := range similarTeamFragments {
	// 		if similarTeamFragments[i].Weight > 0 {
	// 			similarTeamFragments[i].Weight += extraWeight
	// 		}
	// 	}
	// 	for i := range similarPlayerFragments {
	// 		if similarPlayerFragments[i].Weight > 0 {
	// 			similarPlayerFragments[i].Weight += extraWeight
	// 		}
	// 	}
	// }

	// // adjust the totalPrediction according to each game breakdown's average and weight
	// for _, fragment := range gameBreakdownFragments {
	// 	totalPrediction.Assists += fragment.Derived.Assists * (fragment.Weight / 100.0)
	// 	totalPrediction.Blocks += fragment.Derived.Blocks * (fragment.Weight / 100.0)
	// 	totalPrediction.DefensiveRebounds += fragment.Derived.DefensiveRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.FieldGoalsAttempted += fragment.Derived.FieldGoalsAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.FieldGoalsMade += fragment.Derived.FieldGoalsMade * (fragment.Weight / 100.0)
	// 	totalPrediction.FreeThrowsAttempted += fragment.Derived.FreeThrowsAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.FreeThrowsMade += fragment.Derived.FreeThrowsMade * (fragment.Weight / 100.0)
	// 	totalPrediction.OffensiveRebounds += fragment.Derived.OffensiveRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.PersonalFouls += fragment.Derived.PersonalFouls * (fragment.Weight / 100.0)
	// 	totalPrediction.PersonalFoulsDrawn += fragment.Derived.PersonalFoulsDrawn * (fragment.Weight / 100.0)
	// 	totalPrediction.Points += fragment.Derived.Points * (fragment.Weight / 100.0)
	// 	totalPrediction.Rebounds += fragment.Derived.Rebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.Steals += fragment.Derived.Steals * (fragment.Weight / 100.0)
	// 	totalPrediction.ThreePointersAttempted += fragment.Derived.ThreePointersAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.ThreePointersMade += fragment.Derived.ThreePointersMade * (fragment.Weight / 100.0)
	// 	totalPrediction.Turnovers += fragment.Derived.Turnovers * (fragment.Weight / 100.0)
	// 	totalPrediction.Minutes += fragment.Derived.Minutes * (fragment.Weight / 100.0)
	// 	totalPrediction.FantasyScore += fragment.Derived.FantasyScore * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsAssists += fragment.Derived.PointsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsRebounds += fragment.Derived.PointsRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsReboundsAssists += fragment.Derived.PointsReboundsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.ReboundsAssists += fragment.Derived.ReboundsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.BlocksSteals += fragment.Derived.BlocksSteals * (fragment.Weight / 100.0)
	// 	totalPrediction.DoubleDouble += fragment.Derived.DoubleDouble * (fragment.Weight / 100.0)
	// }

	// for _, fragment := range similarTeamFragments {
	// 	totalPrediction.Assists += fragment.Derived.Assists * (fragment.Weight / 100.0)
	// 	totalPrediction.Blocks += fragment.Derived.Blocks * (fragment.Weight / 100.0)
	// 	totalPrediction.DefensiveRebounds += fragment.Derived.DefensiveRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.FieldGoalsAttempted += fragment.Derived.FieldGoalsAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.FieldGoalsMade += fragment.Derived.FieldGoalsMade * (fragment.Weight / 100.0)
	// 	totalPrediction.FreeThrowsAttempted += fragment.Derived.FreeThrowsAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.FreeThrowsMade += fragment.Derived.FreeThrowsMade * (fragment.Weight / 100.0)
	// 	totalPrediction.OffensiveRebounds += fragment.Derived.OffensiveRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.PersonalFouls += fragment.Derived.PersonalFouls * (fragment.Weight / 100.0)
	// 	totalPrediction.PersonalFoulsDrawn += fragment.Derived.PersonalFoulsDrawn * (fragment.Weight / 100.0)
	// 	totalPrediction.Points += fragment.Derived.Points * (fragment.Weight / 100.0)
	// 	totalPrediction.Rebounds += fragment.Derived.Rebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.Steals += fragment.Derived.Steals * (fragment.Weight / 100.0)
	// 	totalPrediction.ThreePointersAttempted += fragment.Derived.ThreePointersAttempted * (fragment.Weight / 100.0)
	// 	totalPrediction.ThreePointersMade += fragment.Derived.ThreePointersMade * (fragment.Weight / 100.0)
	// 	totalPrediction.Turnovers += fragment.Derived.Turnovers * (fragment.Weight / 100.0)
	// 	totalPrediction.Minutes += fragment.Derived.Minutes * (fragment.Weight / 100.0)
	// 	totalPrediction.FantasyScore += fragment.Derived.FantasyScore * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsAssists += fragment.Derived.PointsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsRebounds += fragment.Derived.PointsRebounds * (fragment.Weight / 100.0)
	// 	totalPrediction.PointsReboundsAssists += fragment.Derived.PointsReboundsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.ReboundsAssists += fragment.Derived.ReboundsAssists * (fragment.Weight / 100.0)
	// 	totalPrediction.BlocksSteals += fragment.Derived.BlocksSteals * (fragment.Weight / 100.0)
	// 	totalPrediction.DoubleDouble += fragment.Derived.DoubleDouble * (fragment.Weight / 100.0)
	// }

	// if input.SimilarPlayerInput != nil {
	// 	wouldBeEstimate := &model.AverageStats{}
	// 	wouldBeWeightAdded := input.SimilarPlayerInput.Weight / float64(len(gameBreakdownFragments)+countSimilarTeamsWithGames)
	// 	//TODO: what about similar teams?
	// 	for _, fragment := range gameBreakdownFragments {
	// 		wouldBeEstimate.Assists += fragment.Base.Assists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Blocks += fragment.Base.Blocks * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.DefensiveRebounds += fragment.Base.DefensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FieldGoalsAttempted += fragment.Base.FieldGoalsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FieldGoalsMade += fragment.Base.FieldGoalsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FreeThrowsAttempted += fragment.Base.FreeThrowsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FreeThrowsMade += fragment.Base.FreeThrowsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.OffensiveRebounds += fragment.Base.OffensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PersonalFouls += fragment.Base.PersonalFouls * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PersonalFoulsDrawn += fragment.Base.PersonalFoulsDrawn * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Points += fragment.Base.Points * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Rebounds += fragment.Base.Rebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Steals += fragment.Base.Steals * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ThreePointersAttempted += fragment.Base.ThreePointersAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ThreePointersMade += fragment.Base.ThreePointersMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Turnovers += fragment.Base.Turnovers * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Minutes += fragment.Base.Minutes * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FantasyScore += fragment.Base.FantasyScore * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsAssists += fragment.Base.PointsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsRebounds += fragment.Base.PointsRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsReboundsAssists += fragment.Base.PointsReboundsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ReboundsAssists += fragment.Base.ReboundsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.BlocksSteals += fragment.Base.BlocksSteals * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.DoubleDouble += fragment.Base.DoubleDouble * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 	}
	// 	for _, fragment := range similarTeamFragments {
	// 		wouldBeEstimate.Assists += fragment.Base.Assists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Blocks += fragment.Base.Blocks * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.DefensiveRebounds += fragment.Base.DefensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FieldGoalsAttempted += fragment.Base.FieldGoalsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FieldGoalsMade += fragment.Base.FieldGoalsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FreeThrowsAttempted += fragment.Base.FreeThrowsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FreeThrowsMade += fragment.Base.FreeThrowsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.OffensiveRebounds += fragment.Base.OffensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PersonalFouls += fragment.Base.PersonalFouls * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PersonalFoulsDrawn += fragment.Base.PersonalFoulsDrawn * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Points += fragment.Base.Points * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Rebounds += fragment.Base.Rebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Steals += fragment.Base.Steals * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ThreePointersAttempted += fragment.Base.ThreePointersAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ThreePointersMade += fragment.Base.ThreePointersMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Turnovers += fragment.Base.Turnovers * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.Minutes += fragment.Base.Minutes * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.FantasyScore += fragment.Base.FantasyScore * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsAssists += fragment.Base.PointsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsRebounds += fragment.Base.PointsRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.PointsReboundsAssists += fragment.Base.PointsReboundsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.ReboundsAssists += fragment.Base.ReboundsAssists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.BlocksSteals += fragment.Base.BlocksSteals * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 		wouldBeEstimate.DoubleDouble += fragment.Base.DoubleDouble * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	// 	}

	// 	for _, fragment := range similarPlayerFragments {
	// 		totalPrediction.Assists += (wouldBeEstimate.Assists + ((fragment.PctChange.Assists / 100) * wouldBeEstimate.Assists)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Blocks += (wouldBeEstimate.Blocks + ((fragment.PctChange.Blocks / 100) * wouldBeEstimate.Blocks)) * (fragment.Weight / 100.0)
	// 		totalPrediction.DefensiveRebounds += (wouldBeEstimate.DefensiveRebounds + ((fragment.PctChange.DefensiveRebounds / 100) * wouldBeEstimate.DefensiveRebounds)) * (fragment.Weight / 100.0)
	// 		totalPrediction.FieldGoalsAttempted += (wouldBeEstimate.FieldGoalsAttempted + ((fragment.PctChange.FieldGoalsAttempted / 100) * wouldBeEstimate.FieldGoalsAttempted)) * (fragment.Weight / 100.0)
	// 		totalPrediction.FieldGoalsMade += (wouldBeEstimate.FieldGoalsMade + ((fragment.PctChange.FieldGoalsMade / 100) * wouldBeEstimate.FieldGoalsMade)) * (fragment.Weight / 100.0)
	// 		totalPrediction.FreeThrowsAttempted += (wouldBeEstimate.FreeThrowsAttempted + ((fragment.PctChange.FreeThrowsAttempted / 100) * wouldBeEstimate.FreeThrowsAttempted)) * (fragment.Weight / 100.0)
	// 		totalPrediction.FreeThrowsMade += (wouldBeEstimate.FreeThrowsMade + ((fragment.PctChange.FreeThrowsMade / 100) * wouldBeEstimate.FreeThrowsMade)) * (fragment.Weight / 100.0)
	// 		totalPrediction.OffensiveRebounds += (wouldBeEstimate.OffensiveRebounds + ((fragment.PctChange.OffensiveRebounds / 100) * wouldBeEstimate.OffensiveRebounds)) * (fragment.Weight / 100.0)
	// 		totalPrediction.PersonalFouls += (wouldBeEstimate.PersonalFouls + ((fragment.PctChange.PersonalFouls / 100) * wouldBeEstimate.PersonalFouls)) * (fragment.Weight / 100.0)
	// 		totalPrediction.PersonalFoulsDrawn += (wouldBeEstimate.PersonalFoulsDrawn + ((fragment.PctChange.PersonalFoulsDrawn / 100) * wouldBeEstimate.PersonalFoulsDrawn)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Points += (wouldBeEstimate.Points + ((fragment.PctChange.Points / 100) * wouldBeEstimate.Points)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Rebounds += (wouldBeEstimate.Rebounds + ((fragment.PctChange.Rebounds / 100) * wouldBeEstimate.Rebounds)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Steals += (wouldBeEstimate.Steals + ((fragment.PctChange.Steals / 100) * wouldBeEstimate.Steals)) * (fragment.Weight / 100.0)
	// 		totalPrediction.ThreePointersAttempted += (wouldBeEstimate.ThreePointersAttempted + ((fragment.PctChange.ThreePointersAttempted / 100) * wouldBeEstimate.ThreePointersAttempted)) * (fragment.Weight / 100.0)
	// 		totalPrediction.ThreePointersMade += (wouldBeEstimate.ThreePointersMade + ((fragment.PctChange.ThreePointersMade / 100) * wouldBeEstimate.ThreePointersMade)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Turnovers += (wouldBeEstimate.Turnovers + ((fragment.PctChange.Turnovers / 100) * wouldBeEstimate.Turnovers)) * (fragment.Weight / 100.0)
	// 		totalPrediction.Minutes += (wouldBeEstimate.Minutes + ((fragment.PctChange.Minutes / 100) * wouldBeEstimate.Minutes)) * (fragment.Weight / 100.0)
	// 		totalPrediction.FantasyScore += (wouldBeEstimate.FantasyScore + ((fragment.PctChange.FantasyScore / 100) * wouldBeEstimate.FantasyScore)) * (fragment.Weight / 100.0)
	// 		totalPrediction.PointsAssists += (wouldBeEstimate.PointsAssists + ((fragment.PctChange.PointsAssists / 100) * wouldBeEstimate.PointsAssists)) * (fragment.Weight / 100.0)
	// 		totalPrediction.PointsRebounds += (wouldBeEstimate.PointsRebounds + ((fragment.PctChange.PointsRebounds / 100) * wouldBeEstimate.PointsRebounds)) * (fragment.Weight / 100.0)
	// 		totalPrediction.PointsReboundsAssists += (wouldBeEstimate.PointsReboundsAssists + ((fragment.PctChange.PointsReboundsAssists / 100) * wouldBeEstimate.PointsReboundsAssists)) * (fragment.Weight / 100.0)
	// 		totalPrediction.ReboundsAssists += (wouldBeEstimate.ReboundsAssists + ((fragment.PctChange.ReboundsAssists / 100) * wouldBeEstimate.ReboundsAssists)) * (fragment.Weight / 100.0)
	// 		totalPrediction.BlocksSteals += (wouldBeEstimate.BlocksSteals + ((fragment.PctChange.BlocksSteals / 100) * wouldBeEstimate.BlocksSteals)) * (fragment.Weight / 100.0)
	// 		totalPrediction.DoubleDouble += (wouldBeEstimate.DoubleDouble + ((fragment.PctChange.DoubleDouble / 100) * wouldBeEstimate.DoubleDouble)) * (fragment.Weight / 100.0)
	// 	}
	// }
	// totalPrediction.Assists = similarity.RoundFloat(totalPrediction.Assists, 2)
	// totalPrediction.Blocks = similarity.RoundFloat(totalPrediction.Blocks, 2)
	// totalPrediction.DefensiveRebounds = similarity.RoundFloat(totalPrediction.DefensiveRebounds, 2)
	// totalPrediction.FieldGoalsAttempted = similarity.RoundFloat(totalPrediction.FieldGoalsAttempted, 2)
	// totalPrediction.FieldGoalsMade = similarity.RoundFloat(totalPrediction.FieldGoalsMade, 2)
	// totalPrediction.FreeThrowsAttempted = similarity.RoundFloat(totalPrediction.FreeThrowsAttempted, 2)
	// totalPrediction.FreeThrowsMade = similarity.RoundFloat(totalPrediction.FreeThrowsMade, 2)
	// totalPrediction.OffensiveRebounds = similarity.RoundFloat(totalPrediction.OffensiveRebounds, 2)
	// totalPrediction.PersonalFouls = similarity.RoundFloat(totalPrediction.PersonalFouls, 2)
	// totalPrediction.PersonalFoulsDrawn = similarity.RoundFloat(totalPrediction.PersonalFoulsDrawn, 2)
	// totalPrediction.Points = similarity.RoundFloat(totalPrediction.Points, 2)
	// totalPrediction.Rebounds = similarity.RoundFloat(totalPrediction.Rebounds, 2)
	// totalPrediction.Steals = similarity.RoundFloat(totalPrediction.Steals, 2)
	// totalPrediction.ThreePointersAttempted = similarity.RoundFloat(totalPrediction.ThreePointersAttempted, 2)
	// totalPrediction.ThreePointersMade = similarity.RoundFloat(totalPrediction.ThreePointersMade, 2)
	// totalPrediction.Turnovers = similarity.RoundFloat(totalPrediction.Turnovers, 2)
	// totalPrediction.Minutes = similarity.RoundFloat(totalPrediction.Minutes, 2)
	// totalPrediction.FantasyScore = similarity.RoundFloat(totalPrediction.FantasyScore, 2)
	// totalPrediction.PointsAssists = similarity.RoundFloat(totalPrediction.PointsAssists, 2)
	// totalPrediction.PointsRebounds = similarity.RoundFloat(totalPrediction.PointsRebounds, 2)
	// totalPrediction.PointsReboundsAssists = similarity.RoundFloat(totalPrediction.PointsReboundsAssists, 2)
	// totalPrediction.ReboundsAssists = similarity.RoundFloat(totalPrediction.ReboundsAssists, 2)
	// totalPrediction.BlocksSteals = similarity.RoundFloat(totalPrediction.BlocksSteals, 2)
	// totalPrediction.DoubleDouble = similarity.RoundFloat(totalPrediction.DoubleDouble, 2)

	// predictionAccuracy := &model.AverageStats{}
	// if obj.Outcome != "" && strings.ToLower(obj.Outcome)[0] != 'p' {
	// 	predictionAccuracy.Assists = similarity.RoundFloat(totalPrediction.Assists-float64(obj.Assists), 2)
	// 	predictionAccuracy.Blocks = similarity.RoundFloat(totalPrediction.Blocks-float64(obj.Blocks), 2)
	// 	predictionAccuracy.DefensiveRebounds = similarity.RoundFloat(totalPrediction.DefensiveRebounds-float64(obj.DefensiveRebounds), 2)
	// 	predictionAccuracy.FieldGoalsAttempted = similarity.RoundFloat(totalPrediction.FieldGoalsAttempted-float64(obj.FieldGoalsAttempted), 2)
	// 	predictionAccuracy.FieldGoalsMade = similarity.RoundFloat(totalPrediction.FieldGoalsMade-float64(obj.FieldGoalsMade), 2)
	// 	predictionAccuracy.FreeThrowsAttempted = similarity.RoundFloat(totalPrediction.FreeThrowsAttempted-float64(obj.FreeThrowsAttempted), 2)
	// 	predictionAccuracy.FreeThrowsMade = similarity.RoundFloat(totalPrediction.FreeThrowsMade-float64(obj.FreeThrowsMade), 2)
	// 	predictionAccuracy.OffensiveRebounds = similarity.RoundFloat(totalPrediction.OffensiveRebounds-float64(obj.OffensiveRebounds), 2)
	// 	predictionAccuracy.PersonalFouls = similarity.RoundFloat(totalPrediction.PersonalFouls-float64(obj.PersonalFouls), 2)
	// 	predictionAccuracy.PersonalFoulsDrawn = similarity.RoundFloat(totalPrediction.PersonalFoulsDrawn-float64(obj.PersonalFoulsDrawn), 2)
	// 	predictionAccuracy.Points = similarity.RoundFloat(totalPrediction.Points-float64(obj.Points), 2)
	// 	predictionAccuracy.Rebounds = similarity.RoundFloat(totalPrediction.Rebounds-float64(obj.Rebounds), 2)
	// 	predictionAccuracy.Steals = similarity.RoundFloat(totalPrediction.Steals-float64(obj.Steals), 2)
	// 	predictionAccuracy.ThreePointersAttempted = similarity.RoundFloat(totalPrediction.ThreePointersAttempted-float64(obj.ThreePointersAttempted), 2)
	// 	predictionAccuracy.ThreePointersMade = similarity.RoundFloat(totalPrediction.ThreePointersMade-float64(obj.ThreePointersMade), 2)
	// 	predictionAccuracy.Turnovers = similarity.RoundFloat(totalPrediction.Turnovers-float64(obj.Turnovers), 2)
	// 	predictionAccuracy.Minutes = similarity.RoundFloat(totalPrediction.Minutes-obj.Score(model.Minutes), 2)
	// 	predictionAccuracy.FantasyScore = similarity.RoundFloat(totalPrediction.FantasyScore-obj.Score(model.FantasyScore), 2)
	// 	predictionAccuracy.PointsAssists = similarity.RoundFloat(totalPrediction.PointsAssists-obj.Score(model.PointsAssists), 2)
	// 	predictionAccuracy.PointsRebounds = similarity.RoundFloat(totalPrediction.PointsRebounds-obj.Score(model.PointsRebounds), 2)
	// 	predictionAccuracy.PointsReboundsAssists = similarity.RoundFloat(totalPrediction.PointsReboundsAssists-obj.Score(model.PointsReboundsAssists), 2)
	// 	predictionAccuracy.ReboundsAssists = similarity.RoundFloat(totalPrediction.ReboundsAssists-obj.Score(model.ReboundsAssists), 2)
	// 	predictionAccuracy.BlocksSteals = similarity.RoundFloat(totalPrediction.BlocksSteals-obj.Score(model.BlocksSteals), 2)
	// 	predictionAccuracy.DoubleDouble = similarity.RoundFloat(totalPrediction.DoubleDouble-obj.Score(model.DoubleDouble), 2)
	// }

	// fragments := []*model.PredictionFragment{}
	// fragments = append(fragments, gameBreakdownFragments...)
	// fragments = append(fragments, similarTeamFragments...)
	// fragments = append(fragments, similarPlayerFragments...)
	// breakdown := &model.PredictionBreakdown{
	// 	WeightedTotal:      &totalPrediction,
	// 	PredictionAccuracy: predictionAccuracy,
	// 	Fragments:          fragments,
	// }
	// return breakdown, nil
	panic("propositions not implemented")
}

// Propositions is the resolver for the propositions field.
func (r *playerGameResolver) Propositions(ctx context.Context, obj *model.PlayerGame, input *model.PropositionFilter) ([]*model.Proposition, error) {
	panic(fmt.Errorf("not implemented"))
}

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// PlayerGame returns generated.PlayerGameResolver implementation.
func (r *Resolver) PlayerGame() generated.PlayerGameResolver { return &playerGameResolver{r} }

type playerResolver struct{ *Resolver }
type playerGameResolver struct{ *Resolver }
