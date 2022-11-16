package graph

// // This file will be automatically regenerated based on the schema, any resolver implementations
// // will be copied through when generating and any unknown code will be moved to the end.

// import (
// 	"context"
// 	"fmt"

// 	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
// 	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
// )

// // Team is the resolver for the team field.
// func (r *playerResolver) Team(ctx context.Context, obj *model.Player) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Projections is the resolver for the projections field.
// func (r *playerResolver) Projections(ctx context.Context, obj *model.Player, input model.ProjectionFilter) ([]*model.Projection, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Opponent is the resolver for the opponent field.
// func (r *playerGameResolver) Opponent(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // OpponentStats is the resolver for the opponentStats field.
// func (r *playerGameResolver) OpponentStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Team is the resolver for the team field.
// func (r *playerGameResolver) Team(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // TeamStats is the resolver for the teamStats field.
// func (r *playerGameResolver) TeamStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Player is the resolver for the player field.
// func (r *playerGameResolver) Player(ctx context.Context, obj *model.PlayerGame) (*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Team is the resolver for the team field.
// func (r *playersInGameResolver) Team(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Opponent is the resolver for the opponent field.
// func (r *playersInGameResolver) Opponent(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Confidence is the resolver for the confidence field.
// func (r *predictionResolver) Confidence(ctx context.Context, obj *model.Prediction) (*float64, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Estimation is the resolver for the estimation field.
// func (r *predictionResolver) Estimation(ctx context.Context, obj *model.Prediction) (*float64, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Player is the resolver for the player field.
// func (r *projectionResolver) Player(ctx context.Context, obj *model.Projection) (*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Opponent is the resolver for the opponent field.
// func (r *projectionResolver) Opponent(ctx context.Context, obj *model.Projection) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Result is the resolver for the result field.
// func (r *projectionResolver) Result(ctx context.Context, obj *model.Projection) (*model.PlayerGame, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Prediction is the resolver for the prediction field.
// func (r *propositionResolver) Prediction(ctx context.Context, obj *model.Proposition) (*model.Prediction, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // LastModified is the resolver for the lastModified field.
// func (r *propositionResolver) LastModified(ctx context.Context, obj *model.Proposition) (string, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Players is the resolver for the players field.
// func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // FilterPlayers is the resolver for the filterPlayers field.
// func (r *queryResolver) FilterPlayers(ctx context.Context, input model.PlayerFilter) ([]*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Player is the resolver for the player field.
// func (r *queryResolver) Player(ctx context.Context, input model.PlayerFilter) (*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Teams is the resolver for the teams field.
// func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // FilterTeams is the resolver for the filterTeams field.
// func (r *queryResolver) FilterTeams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Team is the resolver for the team field.
// func (r *queryResolver) Team(ctx context.Context, input model.TeamFilter) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // TeamGames is the resolver for the teamGames field.
// func (r *queryResolver) TeamGames(ctx context.Context, input model.GameFilter) ([]*model.TeamGame, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // PlayerGames is the resolver for the playerGames field.
// func (r *queryResolver) PlayerGames(ctx context.Context, input model.GameFilter) ([]*model.PlayerGame, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Projections is the resolver for the projections field.
// func (r *queryResolver) Projections(ctx context.Context, input model.ProjectionFilter) ([]*model.Projection, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Players is the resolver for the players field.
// func (r *teamResolver) Players(ctx context.Context, obj *model.Team) ([]*model.Player, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Opponent is the resolver for the opponent field.
// func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
// 	panic(fmt.Errorf("not implemented"))
// }

// // Player returns generated.PlayerResolver implementation.
// func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// // PlayerGame returns generated.PlayerGameResolver implementation.
// func (r *Resolver) PlayerGame() generated.PlayerGameResolver { return &playerGameResolver{r} }

// // PlayersInGame returns generated.PlayersInGameResolver implementation.
// func (r *Resolver) PlayersInGame() generated.PlayersInGameResolver { return &playersInGameResolver{r} }

// // Prediction returns generated.PredictionResolver implementation.
// func (r *Resolver) Prediction() generated.PredictionResolver { return &predictionResolver{r} }

// // Projection returns generated.ProjectionResolver implementation.
// func (r *Resolver) Projection() generated.ProjectionResolver { return &projectionResolver{r} }

// // Proposition returns generated.PropositionResolver implementation.
// func (r *Resolver) Proposition() generated.PropositionResolver { return &propositionResolver{r} }

// // Query returns generated.QueryResolver implementation.
// func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// // Team returns generated.TeamResolver implementation.
// func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// // TeamGame returns generated.TeamGameResolver implementation.
// func (r *Resolver) TeamGame() generated.TeamGameResolver { return &teamGameResolver{r} }

// type playerResolver struct{ *Resolver }
// type playerGameResolver struct{ *Resolver }
// type playersInGameResolver struct{ *Resolver }
// type predictionResolver struct{ *Resolver }
// type projectionResolver struct{ *Resolver }
// type propositionResolver struct{ *Resolver }
// type queryResolver struct{ *Resolver }
// type teamResolver struct{ *Resolver }
// type teamGameResolver struct{ *Resolver }
