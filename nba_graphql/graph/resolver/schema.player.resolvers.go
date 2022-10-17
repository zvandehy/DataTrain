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
	"github.com/zvandehy/DataTrain/nba_graphql/math"
)

// Team is the resolver for the team field.
func (r *playerResolver) Team(ctx context.Context, obj *model.Player) (*model.Team, error) {
	if obj.CurrentTeam == "" {
		logrus.Errorf("CurrentTeam is empty for player %v", obj)
		return &model.Team{}, fmt.Errorf("CurrentTeam is empty for player %v", obj)
	}
	t, err := dataloader.For(ctx).TeamByAbr.Load(obj.CurrentTeam)
	if err != nil {
		logrus.Errorf("Error loading team '%v' from: %v", obj.CurrentTeam, obj)
		return nil, err
	}
	if t == nil {
		return &model.Team{}, nil
	}
	return t, err
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

// Projections is the resolver for the projections field.
func (r *playerResolver) Projections(ctx context.Context, obj *model.Player, input model.ProjectionFilter) ([]*model.Projection, error) {
	panic(fmt.Errorf("(r *playerResolver) Projections not implemented"))
}

// Image is the resolver for the image field.
func (r *playerResolver) Image(ctx context.Context, obj *model.Player) (string, error) {
	return fmt.Sprintf("https://ak-static.cms.nba.com/wp-content/uploads/headshots/%s/latest/260x190/%d.png", strings.ToLower(obj.League), obj.PlayerID), nil
}

// Opponent is the resolver for the opponent field.
func (r *playerGameResolver) Opponent(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	team, err := dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
	if err != nil {
		logrus.Errorf("Error loading opponent '%v' from: %v", obj.OpponentID, obj)
		return nil, err
	}
	if team == nil {
		return &model.Team{}, nil
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
	if team == nil {
		return &model.Team{}, nil
	}
	return team, err
}

// TeamStats is the resolver for the teamStats field.
func (r *playerGameResolver) TeamStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
	panic(fmt.Errorf("(r *playerGameResolver) TeamStats not implemented"))
}

// Player is the resolver for the player field.
func (r *playerGameResolver) Player(ctx context.Context, obj *model.PlayerGame) (*model.Player, error) {
	return obj.PlayerRef, nil
}

// PointsRebounds is the resolver for the points_rebounds field.
func (r *playerGameResolver) PointsRebounds(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return obj.Points + obj.Rebounds, nil
}

// PointsAssists is the resolver for the points_assists field.
func (r *playerGameResolver) PointsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return obj.Points + obj.Assists, nil
}

// PointsReboundsAssists is the resolver for the points_rebounds_assists field.
func (r *playerGameResolver) PointsReboundsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return obj.Points + obj.Rebounds + obj.Assists, nil
}

// ReboundsAssists is the resolver for the rebounds_assists field.
func (r *playerGameResolver) ReboundsAssists(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return obj.Rebounds + obj.Assists, nil
}

// BlocksSteals is the resolver for the blocks_steals field.
func (r *playerGameResolver) BlocksSteals(ctx context.Context, obj *model.PlayerGame) (int, error) {
	return obj.Blocks + obj.Steals, nil
}

// FantasyScore is the resolver for the fantasy_score field.
func (r *playerGameResolver) FantasyScore(ctx context.Context, obj *model.PlayerGame) (float64, error) {
	return float64(obj.Points) + float64(obj.Rebounds)*1.2 + float64(obj.Assists)*1.5 + float64(obj.Steals)*3.0 + float64(obj.Blocks)*3.0 - float64(obj.Turnovers), nil
}

// Prediction is the resolver for the prediction field.
func (r *playerGameResolver) Prediction(ctx context.Context, obj *model.PlayerGame, input model.ModelInput) (*model.PredictionBreakdown, error) {
	propositions, err := r.Db.GetPropositionsByGame(ctx, obj)
	if err != nil {
		logrus.Errorf("Error loading propositions for game %v", obj)
		return nil, err
	}

	// TODO: Get earliest season for gamebreakdowns... separately? (currently crashes if no similar player input provided)
	startDate, err := input.SimilarPlayerInput.PlayerPoolFilter.GetEarliestSeasonStartDate()
	if err != nil {
		logrus.Errorf("Error getting earliest season start date: %v", err)
		return nil, err
	}
	startDateStr := startDate.Format("2006-01-02")

	totalPrediction := model.AverageStats{}
	gameBreakdownFragments := []*model.PredictionFragment{}
	similarPlayerFragments := []*model.PredictionFragment{}

	// Player GameLog Breakdowns
	// Base (Filter, Games, Avg) for all of this player's games from the start of the range to the game date
	baseFilter := model.GameFilter{
		StartDate: &startDateStr,
		EndDate:   &obj.Date,
	}
	playerBaseGames := []*model.PlayerGame{}
	for _, game := range obj.PlayerRef.GamesCache {
		if baseFilter.MatchPlayerGame(game) {
			playerBaseGames = append(playerBaseGames, game)
		}
	}
	playerBase := &model.AverageStats{}
	if len(playerBaseGames) > 0 {
		b := model.NewPlayerAverage(playerBaseGames, obj.PlayerRef)
		playerBase = b.AverageStats()
	}

	distributeExtraWeight := 0.0
	for i := range input.GameBreakdowns {
		games := []*model.PlayerGame{}
		for _, game := range obj.PlayerRef.GamesCache {
			// only get games before the current game
			input.GameBreakdowns[i].Filter.EndDate = &obj.Date
			// if opponentMatch is true, only get games previously matched up against opponent
			if input.GameBreakdowns[i].Filter.OpponentMatch != nil && *input.GameBreakdowns[i].Filter.OpponentMatch {
				input.GameBreakdowns[i].Filter.OpponentID = &obj.OpponentID
			}
			// TODO: add more filters, like home/away, playeoff/regular season, etc.

			// if the game matches the filter, add it to the list of games for this breakdown
			if input.GameBreakdowns[i].Filter.MatchPlayerGame(game) {
				games = append(games, game)
				if input.GameBreakdowns[i].Filter.LastX != nil && len(games) >= *input.GameBreakdowns[i].Filter.LastX {
					break
				}
			}
		}
		// player average for the games in this breakdown
		pAvg := model.NewPlayerAverage(games, obj.PlayerRef)

		// Don't calculate averages if there are no games, or less than lastX games
		if len(games) == 0 || (input.GameBreakdowns[i].Filter.LastX != nil && len(games) < *input.GameBreakdowns[i].Filter.LastX) {
			distributeExtraWeight += input.GameBreakdowns[i].Weight
			input.GameBreakdowns[i].Weight = 0
		} else {
			derived := pAvg.AverageStats()
			frag := &model.PredictionFragment{
				Name:      input.GameBreakdowns[i].Name,
				Derived:   derived,
				Base:      playerBase,
				PctChange: playerBase.PercentChange(derived),
				Weight:    input.GameBreakdowns[i].Weight,
			}
			for _, proposition := range propositions {
				analysis := model.PropositionSummary{}
				for _, game := range games {
					propScore := game.Score(proposition.Type)
					if propScore > proposition.Target {
						analysis.NumOver++
					}
					if propScore < proposition.Target {
						analysis.NumUnder++
					}
					if propScore == proposition.Target {
						analysis.NumPush++
					}
				}
				analysis.PctOver = float64(analysis.NumOver) / float64(len(games))
				analysis.PctUnder = float64(analysis.NumUnder) / float64(len(games))
				analysis.PctPush = float64(analysis.NumPush) / float64(len(games))
				prop := *proposition
				prop.Analysis = &analysis
				frag.Propositions = append(frag.Propositions, &prop)
			}
			gameBreakdownFragments = append(gameBreakdownFragments, frag)
		}
	}

	countSimilarPlayersWithGamesVsOpp := 0
	if input.SimilarPlayerInput != nil {
		// gets X similar players to the current player, where X is defined by the input limit
		similarPlayers, err := r.Db.GetSimilarPlayersFromMatrix(ctx, obj.PlayerID, input.SimilarPlayerInput, obj.Date)
		if err != nil {
			logrus.Errorf("Error getting similar players: %v", err)
			return nil, err
		}
		for _, player := range similarPlayers {
			// get the similar player's games from the start of the range to the game date
			baseGames := []*model.PlayerGame{}
			for _, game := range player.GamesCache {
				if baseFilter.MatchPlayerGame(game) {
					baseGames = append(baseGames, game)
				}
			}
			b := model.NewPlayerAverage(baseGames, &player)
			baseAvg := b.AverageStats()

			// get the similar player's games vs the matchup opponent
			opponentFilter := model.GameFilter{
				EndDate:    &obj.Date,
				OpponentID: &obj.OpponentID,
			}
			if len(*input.SimilarPlayerInput.PlayerPoolFilter.Seasons) > 0 {
				opponentFilter.Seasons = input.SimilarPlayerInput.PlayerPoolFilter.Seasons
			}
			matchupGames := []*model.PlayerGame{}
			for _, game := range player.GamesCache {
				if opponentFilter.MatchPlayerGame(game) {
					matchupGames = append(matchupGames, game)
				}
			}
			// Don't calculate averages if there are no games vs opponent
			if len(matchupGames) > 0 {
				countSimilarPlayersWithGamesVsOpp++
				pAvg := model.NewPlayerAverage(matchupGames, &player)
				derived := pAvg.AverageStats()
				similarPlayerFragments = append(similarPlayerFragments, &model.PredictionFragment{
					Name:      fmt.Sprintf("%v vs Opp", player.Name),
					Derived:   derived,
					Base:      baseAvg,
					PctChange: baseAvg.PercentChange(derived),
					Weight:    input.SimilarPlayerInput.Weight,
				})
			} else {
				similarPlayerFragments = append(similarPlayerFragments, &model.PredictionFragment{
					Name:      fmt.Sprintf("%v vs Opp (None)", player.Name),
					Derived:   &model.AverageStats{},
					Base:      baseAvg,
					PctChange: &model.AverageStats{},
					Weight:    0,
				})
			}
		}
	}

	// if there are no similar players with games vs the opponent, don't use the similar player input and distribute that weight across the other inputs
	if countSimilarPlayersWithGamesVsOpp == 0 {
		distributeExtraWeight += input.SimilarPlayerInput.Weight
	} else {
		// distribute similar player weights evenly between all similar players with games vs Opp
		similarPlayerWeights := input.SimilarPlayerInput.Weight / float64(countSimilarPlayersWithGamesVsOpp)
		for i := range similarPlayerFragments {
			if similarPlayerFragments[i].Weight > 0 {
				similarPlayerFragments[i].Weight = math.RoundFloat(similarPlayerWeights, 2)
			}
		}
	}

	//distribute extra weight evenly among game breakdowns
	// TODO: what about similar teams?
	if distributeExtraWeight > 0 {
		for i := range gameBreakdownFragments {
			gameBreakdownFragments[i].Weight += math.RoundFloat(distributeExtraWeight/float64(len(gameBreakdownFragments)), 2)
		}
	}

	// adjust the totalPrediction according to each game breakdown's average and weight
	for _, fragment := range gameBreakdownFragments {
		totalPrediction.Assists += fragment.Derived.Assists * (fragment.Weight / 100.0)
		totalPrediction.Blocks += fragment.Derived.Blocks * (fragment.Weight / 100.0)
		totalPrediction.DefensiveRebounds += fragment.Derived.DefensiveRebounds * (fragment.Weight / 100.0)
		totalPrediction.FieldGoalsAttempted += fragment.Derived.FieldGoalsAttempted * (fragment.Weight / 100.0)
		totalPrediction.FieldGoalsMade += fragment.Derived.FieldGoalsMade * (fragment.Weight / 100.0)
		totalPrediction.FreeThrowsAttempted += fragment.Derived.FreeThrowsAttempted * (fragment.Weight / 100.0)
		totalPrediction.FreeThrowsMade += fragment.Derived.FreeThrowsMade * (fragment.Weight / 100.0)
		totalPrediction.OffensiveRebounds += fragment.Derived.OffensiveRebounds * (fragment.Weight / 100.0)
		totalPrediction.PersonalFouls += fragment.Derived.PersonalFouls * (fragment.Weight / 100.0)
		totalPrediction.PersonalFoulsDrawn += fragment.Derived.PersonalFoulsDrawn * (fragment.Weight / 100.0)
		totalPrediction.Points += fragment.Derived.Points * (fragment.Weight / 100.0)
		totalPrediction.Rebounds += fragment.Derived.Rebounds * (fragment.Weight / 100.0)
		totalPrediction.Steals += fragment.Derived.Steals * (fragment.Weight / 100.0)
		totalPrediction.ThreePointersAttempted += fragment.Derived.ThreePointersAttempted * (fragment.Weight / 100.0)
		totalPrediction.ThreePointersMade += fragment.Derived.ThreePointersMade * (fragment.Weight / 100.0)
		totalPrediction.Turnovers += fragment.Derived.Turnovers * (fragment.Weight / 100.0)
		totalPrediction.Minutes += fragment.Derived.Minutes * (fragment.Weight / 100.0)
	}

	// adjust the totalPrediction according to each similar player's percent change and weight
	// using what the estimation would be without similar players

	wouldBeEstimate := &model.AverageStats{}
	wouldBeWeightAdded := input.SimilarPlayerInput.Weight / float64(len(gameBreakdownFragments))
	//TODO: what about similar teams?
	for _, fragment := range gameBreakdownFragments {
		wouldBeEstimate.Assists += fragment.Base.Assists * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Blocks += fragment.Base.Blocks * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.DefensiveRebounds += fragment.Base.DefensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.FieldGoalsAttempted += fragment.Base.FieldGoalsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.FieldGoalsMade += fragment.Base.FieldGoalsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.FreeThrowsAttempted += fragment.Base.FreeThrowsAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.FreeThrowsMade += fragment.Base.FreeThrowsMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.OffensiveRebounds += fragment.Base.OffensiveRebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.PersonalFouls += fragment.Base.PersonalFouls * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.PersonalFoulsDrawn += fragment.Base.PersonalFoulsDrawn * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Points += fragment.Base.Points * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Rebounds += fragment.Base.Rebounds * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Steals += fragment.Base.Steals * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.ThreePointersAttempted += fragment.Base.ThreePointersAttempted * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.ThreePointersMade += fragment.Base.ThreePointersMade * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Turnovers += fragment.Base.Turnovers * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
		wouldBeEstimate.Minutes += fragment.Base.Minutes * ((fragment.Weight + wouldBeWeightAdded) / 100.0)
	}

	for _, fragment := range similarPlayerFragments {
		totalPrediction.Assists += (wouldBeEstimate.Assists + ((fragment.PctChange.Assists / 100) * wouldBeEstimate.Assists)) * (fragment.Weight / 100.0)
		totalPrediction.Blocks += (wouldBeEstimate.Blocks + ((fragment.PctChange.Blocks / 100) * wouldBeEstimate.Blocks)) * (fragment.Weight / 100.0)
		totalPrediction.DefensiveRebounds += (wouldBeEstimate.DefensiveRebounds + ((fragment.PctChange.DefensiveRebounds / 100) * wouldBeEstimate.DefensiveRebounds)) * (fragment.Weight / 100.0)
		totalPrediction.FieldGoalsAttempted += (wouldBeEstimate.FieldGoalsAttempted + ((fragment.PctChange.FieldGoalsAttempted / 100) * wouldBeEstimate.FieldGoalsAttempted)) * (fragment.Weight / 100.0)
		totalPrediction.FieldGoalsMade += (wouldBeEstimate.FieldGoalsMade + ((fragment.PctChange.FieldGoalsMade / 100) * wouldBeEstimate.FieldGoalsMade)) * (fragment.Weight / 100.0)
		totalPrediction.FreeThrowsAttempted += (wouldBeEstimate.FreeThrowsAttempted + ((fragment.PctChange.FreeThrowsAttempted / 100) * wouldBeEstimate.FreeThrowsAttempted)) * (fragment.Weight / 100.0)
		totalPrediction.FreeThrowsMade += (wouldBeEstimate.FreeThrowsMade + ((fragment.PctChange.FreeThrowsMade / 100) * wouldBeEstimate.FreeThrowsMade)) * (fragment.Weight / 100.0)
		totalPrediction.OffensiveRebounds += (wouldBeEstimate.OffensiveRebounds + ((fragment.PctChange.OffensiveRebounds / 100) * wouldBeEstimate.OffensiveRebounds)) * (fragment.Weight / 100.0)
		totalPrediction.PersonalFouls += (wouldBeEstimate.PersonalFouls + ((fragment.PctChange.PersonalFouls / 100) * wouldBeEstimate.PersonalFouls)) * (fragment.Weight / 100.0)
		totalPrediction.PersonalFoulsDrawn += (wouldBeEstimate.PersonalFoulsDrawn + ((fragment.PctChange.PersonalFoulsDrawn / 100) * wouldBeEstimate.PersonalFoulsDrawn)) * (fragment.Weight / 100.0)
		totalPrediction.Points += (wouldBeEstimate.Points + ((fragment.PctChange.Points / 100) * wouldBeEstimate.Points)) * (fragment.Weight / 100.0)
		totalPrediction.Rebounds += (wouldBeEstimate.Rebounds + ((fragment.PctChange.Rebounds / 100) * wouldBeEstimate.Rebounds)) * (fragment.Weight / 100.0)
		totalPrediction.Steals += (wouldBeEstimate.Steals + ((fragment.PctChange.Steals / 100) * wouldBeEstimate.Steals)) * (fragment.Weight / 100.0)
		totalPrediction.ThreePointersAttempted += (wouldBeEstimate.ThreePointersAttempted + ((fragment.PctChange.ThreePointersAttempted / 100) * wouldBeEstimate.ThreePointersAttempted)) * (fragment.Weight / 100.0)
		totalPrediction.ThreePointersMade += (wouldBeEstimate.ThreePointersMade + ((fragment.PctChange.ThreePointersMade / 100) * wouldBeEstimate.ThreePointersMade)) * (fragment.Weight / 100.0)
		totalPrediction.Turnovers += (wouldBeEstimate.Turnovers + ((fragment.PctChange.Turnovers / 100) * wouldBeEstimate.Turnovers)) * (fragment.Weight / 100.0)
		totalPrediction.Minutes += (wouldBeEstimate.Minutes + ((fragment.PctChange.Minutes / 100) * wouldBeEstimate.Minutes)) * (fragment.Weight / 100.0)
	}
	totalPrediction.Assists = math.RoundFloat(totalPrediction.Assists, 2)
	totalPrediction.Blocks = math.RoundFloat(totalPrediction.Blocks, 2)
	totalPrediction.DefensiveRebounds = math.RoundFloat(totalPrediction.DefensiveRebounds, 2)
	totalPrediction.FieldGoalsAttempted = math.RoundFloat(totalPrediction.FieldGoalsAttempted, 2)
	totalPrediction.FieldGoalsMade = math.RoundFloat(totalPrediction.FieldGoalsMade, 2)
	totalPrediction.FreeThrowsAttempted = math.RoundFloat(totalPrediction.FreeThrowsAttempted, 2)
	totalPrediction.FreeThrowsMade = math.RoundFloat(totalPrediction.FreeThrowsMade, 2)
	totalPrediction.OffensiveRebounds = math.RoundFloat(totalPrediction.OffensiveRebounds, 2)
	totalPrediction.PersonalFouls = math.RoundFloat(totalPrediction.PersonalFouls, 2)
	totalPrediction.PersonalFoulsDrawn = math.RoundFloat(totalPrediction.PersonalFoulsDrawn, 2)
	totalPrediction.Points = math.RoundFloat(totalPrediction.Points, 2)
	totalPrediction.Rebounds = math.RoundFloat(totalPrediction.Rebounds, 2)
	totalPrediction.Steals = math.RoundFloat(totalPrediction.Steals, 2)
	totalPrediction.ThreePointersAttempted = math.RoundFloat(totalPrediction.ThreePointersAttempted, 2)
	totalPrediction.ThreePointersMade = math.RoundFloat(totalPrediction.ThreePointersMade, 2)
	totalPrediction.Turnovers = math.RoundFloat(totalPrediction.Turnovers, 2)
	totalPrediction.Minutes = math.RoundFloat(totalPrediction.Minutes, 2)

	fragments := []*model.PredictionFragment{}
	fragments = append(fragments, gameBreakdownFragments...)
	fragments = append(fragments, similarPlayerFragments...)
	breakdown := &model.PredictionBreakdown{
		WeightedTotal: &totalPrediction,
		Fragments:     fragments,
	}
	return breakdown, nil
}

// LastModified is the resolver for the lastModified field.
func (r *propositionResolver) LastModified(ctx context.Context, obj *model.Proposition) (string, error) {
	if obj.LastModified == nil {
		return "", nil
	}
	return obj.LastModified.Format("2006-01-02"), nil
}

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// PlayerGame returns generated.PlayerGameResolver implementation.
func (r *Resolver) PlayerGame() generated.PlayerGameResolver { return &playerGameResolver{r} }

// Proposition returns generated.PropositionResolver implementation.
func (r *Resolver) Proposition() generated.PropositionResolver { return &propositionResolver{r} }

type playerResolver struct{ *Resolver }
type playerGameResolver struct{ *Resolver }
type propositionResolver struct{ *Resolver }
