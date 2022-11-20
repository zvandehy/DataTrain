package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
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
		CumulativeOver:     0,
		CumulativeUnder:    0,
		CumulativePush:     0,
		CumulativeOverPct:  0.0,
		CumulativeUnderPct: 0.0,
		CumulativePushPct:  0.0,
		Wager:              model.WagerOver,
		WagerOutcome:       model.WagerOutcomePending,
		Breakdowns:         []*model.PropBreakdown{},
	}
	gamelogBreakdowns := r.GetGamelogBreakdowns(ctx, input.GameBreakdowns, obj.Game, &obj.Target, obj.Type)
	similarplayerbreakdowns := []*model.PropBreakdown{}
	if input.SimilarPlayerInput != nil {
		similarplayerbreakdowns = r.GetSimilarPlayerBreakdowns(ctx, input.SimilarPlayerInput, obj.Game, obj.Type)
	}

	gamelogCumulativeWeight := 0.0
	for _, breakdown := range gamelogBreakdowns {
		gamelogCumulativeWeight += breakdown.Weight
	}
	similarplayerCumulativeWeight := 0.0
	for _, breakdown := range similarplayerbreakdowns {
		similarplayerCumulativeWeight += breakdown.Weight
	}
	distribute := 0.0
	if gamelogCumulativeWeight == 0 {
		for _, gameInput := range input.GameBreakdowns {
			distribute += gameInput.Weight
		}
	}
	if similarplayerCumulativeWeight == 0 && input.SimilarPlayerInput != nil {
		distribute += input.SimilarPlayerInput.Weight
	}

	// similarTeambreakdown := r.GetSimilarTeamBreakdown(ctx, input.SimilarPlayerInput, &obj.Game, &obj.Target, obj.Type)
	if input.SimilarTeamInput != nil {
		return nil, fmt.Errorf("SimilarTeamInput not implemented")
	}

	if distribute > 0 {
		logrus.Errorf("Distributing extra weight not implemented: %d", distribute)
	}

	estimationWithoutSimilarPlayers := 0.0
	distributeWeightWithoutSimilarPlayers := 0.0
	if input.SimilarPlayerInput != nil {
		distributeWeightWithoutSimilarPlayers = (input.SimilarPlayerInput.Weight / float64(len(gamelogBreakdowns))) / 100.0
	}
	for _, breakdown := range gamelogBreakdowns {
		propPrediction.CumulativeOver += breakdown.Over * int(breakdown.Weight)
		propPrediction.CumulativeUnder += breakdown.Under * int(breakdown.Weight)
		propPrediction.CumulativePush += breakdown.Push * int(breakdown.Weight)
		propPrediction.Estimation += (breakdown.Weight / 100.0) * breakdown.DerivedAverage
		estimationWithoutSimilarPlayers += ((breakdown.Weight / 100.0) + distributeWeightWithoutSimilarPlayers) * breakdown.DerivedAverage

		propPrediction.Breakdowns = append(propPrediction.Breakdowns, breakdown)
	}
	//TODO: similar teams
	for _, breakdown := range similarplayerbreakdowns {
		propPrediction.CumulativeOver += breakdown.Over * int(breakdown.Weight)
		propPrediction.CumulativeUnder += breakdown.Under * int(breakdown.Weight)
		propPrediction.CumulativePush += breakdown.Push * int(breakdown.Weight)
		propPrediction.Estimation = estimationWithoutSimilarPlayers + ((breakdown.PctChange/100.0)*estimationWithoutSimilarPlayers)*(breakdown.Weight/100.0)
	}

	propPrediction.CumulativeOverPct = float64(propPrediction.CumulativeOver) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)
	propPrediction.CumulativeUnderPct = float64(propPrediction.CumulativeUnder) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)
	propPrediction.CumulativePushPct = float64(propPrediction.CumulativePush) / float64(propPrediction.CumulativeOver+propPrediction.CumulativeUnder+propPrediction.CumulativePush)

	if propPrediction.Estimation > obj.Target {
		propPrediction.Wager = model.WagerOver
	} else {
		propPrediction.Wager = model.WagerUnder
	}

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

	propPrediction.Estimation = math.Round(propPrediction.Estimation*100) / 100
	roundedAccuracy := math.Round(*propPrediction.EstimationAccuracy*100) / 100.0
	propPrediction.EstimationAccuracy = &roundedAccuracy
	return &propPrediction, nil
}

// Proposition returns generated.PropositionResolver implementation.
func (r *Resolver) Proposition() generated.PropositionResolver { return &propositionResolver{r} }

type propositionResolver struct{ *Resolver }
