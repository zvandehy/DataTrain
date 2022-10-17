package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// Propositions is the resolver for the propositions field.
func (r *projectionResolver) Propositions(ctx context.Context, obj *model.Projection, input *model.PropositionFilter) ([]*model.Proposition, error) {
	propositions := []*model.Proposition{}
	for _, p := range obj.Props {
		if p.Match(*input) {
			propositions = append(propositions, p)
		}
	}
	return propositions, nil
}

// Result is the resolver for the result field.
func (r *projectionResolver) Result(ctx context.Context, obj *model.Projection) (*model.PlayerGame, error) {
	panic(fmt.Errorf("(r *projectionResolver) Result not implemented"))
}

// Prediction is the resolver for the prediction field.
func (r *propositionResolver) Prediction(ctx context.Context, obj *model.Proposition, input model.ModelInput) (*model.Prediction, error) {
	// start := util.SEASON_START_2021_22 //TODO: make this dynamic
	// // TODO: Need to somehow ensure that the similarity matrix has been loaded
	// r.Db.PlayerSimilarity.GetSimilarPlayers(obj.ProjectionRef.Player.PlayerID, input.SimilarPlayersLimit, start, obj.ProjectionRef.Date)

	// panic(fmt.Errorf("not implemented"))
	fmt.Println(util.Print(input))
	return &model.Prediction{
		Model:               *input.Model,
		OverUnderPrediction: model.OutputOver,
		Confidence:          0.5,
		Estimation:          10,
	}, nil
}

// LastModified is the resolver for the lastModified field.
func (r *propositionResolver) LastModified(ctx context.Context, obj *model.Proposition) (string, error) {
	if obj.LastModified == nil {
		return "", nil
	}
	return obj.LastModified.Format("2006-01-02"), nil
}

// Projection returns generated.ProjectionResolver implementation.
func (r *Resolver) Projection() generated.ProjectionResolver { return &projectionResolver{r} }

// Proposition returns generated.PropositionResolver implementation.
func (r *Resolver) Proposition() generated.PropositionResolver { return &propositionResolver{r} }

type projectionResolver struct{ *Resolver }
type propositionResolver struct{ *Resolver }
