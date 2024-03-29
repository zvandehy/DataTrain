package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// LastModified is the resolver for the lastModified field.
func (r *propositionResolver) LastModified(ctx context.Context, obj *model.Proposition) (string, error) {
	if obj.LastModified == nil {
		return "", nil
	}
	return obj.LastModified.Format(time.RFC3339), nil
}

// Prediction is the resolver for the prediction field.
func (r *propositionResolver) Prediction(ctx context.Context, obj *model.Proposition, input model.ModelInput) (*model.PropPrediction, error) {
	propPrediction := model.PropPrediction{
		Estimation:         0.0,
		EstimationAccuracy: nil,
		Significance:       0.0,
		CumulativeOver:     0,
		CumulativeUnder:    0,
		CumulativePush:     0,
		CumulativeOverPct:  0.0,
		CumulativeUnderPct: 0.0,
		CumulativePushPct:  0.0,
		Wager:              model.WagerOver,
		WagerOutcome:       model.WagerOutcomePending,
		Breakdowns:         []*model.PropBreakdown{},
		StdDev:             0.0,
	}
	gamelogBreakdowns := r.GetGamelogBreakdowns(ctx, input.GameBreakdowns, obj.Game, &obj.Target, obj.Type)
	similarPlayerBreakdowns := []model.PropBreakdown{}
	if input.SimilarPlayerInput != nil {
		similarPlayerBreakdowns = r.GetSimilarPlayerBreakdowns(ctx, input.SimilarPlayerInput, obj.Game, obj.Type)
	}
	// similarTeamBreakdowns := r.GetSimilarTeamBreakdown(ctx, input.SimilarPlayerInput, &obj.Game, &obj.Target, obj.Type)
	if input.SimilarTeamInput != nil {
		return nil, fmt.Errorf("SimilarTeamInput not implemented")
	}

	varianceDatasets := [][]float64{}

	gamelogCumulativeWeight := 0.0
	for _, breakdown := range gamelogBreakdowns {
		gamelogCumulativeWeight += breakdown.Weight
	}
	similarPlayerCumulativeWeight := 0.0
	for _, breakdown := range similarPlayerBreakdowns {
		similarPlayerCumulativeWeight += breakdown.Weight
	}
	similarTeamCumulativeWeight := 0.0
	// for _, breakdown := range similarTeamBreakdowns {
	// 	similarTeamCumulativeWeight += breakdown.Weight
	// }

	validGamelogBreakdowns := 0.0
	for _, breakdown := range gamelogBreakdowns {
		if breakdown.Weight > 0 && len(breakdown.DerivedGames) > 0 {
			validGamelogBreakdowns++
		}
	}

	distribute := 100.0 - gamelogCumulativeWeight - similarPlayerCumulativeWeight - similarTeamCumulativeWeight
	if distribute > 0 {
		countValidBreakdowns := validGamelogBreakdowns
		for _, breakdown := range similarPlayerBreakdowns {
			if breakdown.Weight > 0 && len(breakdown.DerivedGames) > 0 {
				countValidBreakdowns++
			}
		}
		// for _, breakdown := range similarTeamBreakdowns {
		// 	if breakdown.Weight > 0 && len(breakdown.DerivedGames) > 0{
		// 		countValidBreakdowns++
		// 	}
		// }
		if countValidBreakdowns == 0 {
			return nil, fmt.Errorf("no valid breakdowns")
		}
		distributePerBreakdown := distribute / float64(countValidBreakdowns)
		for i, breakdown := range gamelogBreakdowns {
			if breakdown.Weight > 0 {
				gamelogBreakdowns[i].Weight += distributePerBreakdown
			}
		}
		for i, breakdown := range similarPlayerBreakdowns {
			if breakdown.Weight > 0 {
				similarPlayerBreakdowns[i].Weight += distributePerBreakdown
			}
		}
		// todo: similar team distribution
	}

	estimationWithoutSimilarPlayers := 0.0
	distributeWeightWithoutSimilarPlayers := 0.0
	if input.SimilarPlayerInput != nil {
		distributeWeightWithoutSimilarPlayers = input.SimilarPlayerInput.Weight / validGamelogBreakdowns
	}
	for _, breakdown := range gamelogBreakdowns {
		propPrediction.CumulativeOver += breakdown.Over
		propPrediction.CumulativeUnder += breakdown.Under
		propPrediction.CumulativePush += breakdown.Push
		contribution := (breakdown.Weight / 100.0) * breakdown.DerivedAverage
		breakdown.Contribution = contribution
		propPrediction.Estimation += contribution
		estimationWithoutSimilarPlayers += ((breakdown.Weight + distributeWeightWithoutSimilarPlayers) / 100.0) * breakdown.DerivedAverage
		propPrediction.Breakdowns = append(propPrediction.Breakdowns, breakdown)
		if len(breakdown.DerivedGames) > 0 {
			dataset := []float64{}
			for _, game := range breakdown.DerivedGames {
				dataset = append(dataset, game.Score(obj.Type))
			}
			varianceDatasets = append(varianceDatasets, dataset)
		}
	}
	//TODO: similar teams
	if len(similarPlayerBreakdowns) > 0 {
		similarPlayerDerivedSum := 0.0
		similarPlayerBaseSum := 0.0
		similarPlayerBreakdown := &model.PropBreakdown{
			Weight:       input.SimilarPlayerInput.Weight,
			Name:         "Similar Players vs Opponent",
			DerivedGames: []*model.PlayerGame{},
		}
		for _, breakdown := range similarPlayerBreakdowns {
			similarPlayerDerivedSum += breakdown.DerivedAverage
			similarPlayerBaseSum += breakdown.Base
			propPrediction.CumulativeOver += breakdown.Over
			propPrediction.CumulativeUnder += breakdown.Under
			propPrediction.CumulativePush += breakdown.Push
			similarPlayerBreakdown.Over += breakdown.Over
			similarPlayerBreakdown.Under += breakdown.Under
			similarPlayerBreakdown.Push += breakdown.Push
			similarPlayerBreakdown.DerivedAverage += breakdown.DerivedAverage
			similarPlayerBreakdown.DerivedGamesCount += breakdown.DerivedGamesCount
			similarPlayerBreakdown.DerivedGames = append(similarPlayerBreakdown.DerivedGames, breakdown.DerivedGames...)
			// propPrediction.Breakdowns = append(propPrediction.Breakdowns, breakdown)
			// dataset := []float64{}
			// for _, game := range breakdown.DerivedGames {
			// 	dataset = append(dataset, game.Score(obj.Type))
			// }
			// varianceDatasets = append(varianceDatasets, dataset)
		}
		similarPlayerBreakdown.DerivedAverage = similarPlayerDerivedSum / float64(similarPlayerBreakdown.DerivedGamesCount)
		similarPlayerBreakdown.Base = similarPlayerBaseSum / float64(similarPlayerBreakdown.DerivedGamesCount)
		similarPlayerPctChange := (similarPlayerDerivedSum - similarPlayerBaseSum) / similarPlayerBaseSum
		expectedDifference := estimationWithoutSimilarPlayers*similarPlayerPctChange + estimationWithoutSimilarPlayers
		contribution := expectedDifference * (input.SimilarPlayerInput.Weight / 100.0)
		similarPlayerBreakdown.Contribution = contribution
		similarPlayerBreakdown.PctChange = similarPlayerPctChange * 100.0
		similarPlayerBreakdown.OverPct = float64(similarPlayerBreakdown.Over) / float64(similarPlayerBreakdown.DerivedGamesCount)
		similarPlayerBreakdown.UnderPct = float64(similarPlayerBreakdown.Under) / float64(similarPlayerBreakdown.DerivedGamesCount)
		similarPlayerBreakdown.PushPct = float64(similarPlayerBreakdown.Push) / float64(similarPlayerBreakdown.DerivedGamesCount)
		similarPlayerStdDev := 0.0
		for _, game := range similarPlayerBreakdown.DerivedGames {
			similarPlayerStdDev += math.Pow(game.Score(obj.Type)-similarPlayerBreakdown.DerivedAverage, 2)
		}
		similarPlayerStdDev = math.Sqrt(similarPlayerStdDev / float64(len(similarPlayerBreakdown.DerivedGames)))
		similarPlayerBreakdown.StdDev = similarPlayerStdDev
		similarPlayerBreakdown.DerivedAverage = similarPlayerDerivedSum / float64(len(similarPlayerBreakdowns))
		propPrediction.Breakdowns = append(propPrediction.Breakdowns, similarPlayerBreakdown)
		propPrediction.Estimation += contribution
	}
	propPrediction.CumulativeOverPct = float64(propPrediction.CumulativeOver) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)
	propPrediction.CumulativeUnderPct = float64(propPrediction.CumulativeUnder) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)
	propPrediction.CumulativePushPct = float64(propPrediction.CumulativePush) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)

	// if propPrediction.Estimation > obj.Target {
	// 	propPrediction.Wager = model.WagerOver
	// } else {
	// 	propPrediction.Wager = model.WagerUnder
	// }

	propPrediction.Estimation = math.Round(propPrediction.Estimation*100) / 100
	if propPrediction.EstimationAccuracy != nil {
		*propPrediction.EstimationAccuracy = math.Round(*propPrediction.EstimationAccuracy*100) / 100
	}
	if len(varianceDatasets) == 0 {
		logrus.Error("No variance datasets")
	} else {
		if math.IsNaN(PoolVariance(varianceDatasets)) {
			for _, dataset := range varianceDatasets {
				logrus.Warn("Dataset length: ", len(dataset))
				logrus.Warn("Variance is: ", Variance(dataset))
				logrus.Warn("Mean is:", Mean(dataset))
			}
			log.Fatal("NaN variance", len(varianceDatasets))
		}
	}
	pooledStdDev := math.Round(math.Sqrt(PoolVariance(varianceDatasets))*100) / 100.0
	propPrediction.StdDev = pooledStdDev
	significance, wager := model.PValue(propPrediction.Estimation, pooledStdDev, obj.Target)
	propPrediction.Significance = math.Round((significance*100)*100) / 100
	propPrediction.Wager = wager

	if obj.ActualResult != nil {
		if *obj.ActualResult > obj.Target {
			if propPrediction.Wager == model.WagerOver {
				propPrediction.WagerOutcome = model.WagerOutcomeHit
			} else {
				propPrediction.WagerOutcome = model.WagerOutcomeMiss
			}
		}
		if *obj.ActualResult < obj.Target {
			if propPrediction.Wager == model.WagerUnder {
				propPrediction.WagerOutcome = model.WagerOutcomeHit
			} else {
				propPrediction.WagerOutcome = model.WagerOutcomeMiss
			}
		}
		if *obj.ActualResult == obj.Target {
			propPrediction.WagerOutcome = model.WagerOutcomePush
		}
		// TODO: use actual statistics to calculate accuracy instead of pct diff
		accuracy := (propPrediction.Estimation - *obj.ActualResult) / *obj.ActualResult
		if *obj.ActualResult == 0 {
			accuracy = 100
		}
		propPrediction.EstimationAccuracy = &accuracy
	} else {
		propPrediction.WagerOutcome = model.WagerOutcomePending
	}
	if math.IsNaN(propPrediction.Estimation) || math.IsNaN(propPrediction.Significance) {
		logrus.Errorf("Prop Prediction IsNaN: %v %v %v %v", obj.PlayerName, obj.Type, propPrediction.Estimation, propPrediction.Significance)
		propPrediction.Estimation = 0
		propPrediction.Significance = 0
	}
	return &propPrediction, nil
}

// Proposition returns generated.PropositionResolver implementation.
func (r *Resolver) Proposition() generated.PropositionResolver { return &propositionResolver{r} }

type propositionResolver struct{ *Resolver }
