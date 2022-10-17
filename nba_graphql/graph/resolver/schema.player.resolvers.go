package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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
	panic(fmt.Errorf("(r *playerGameResolver) Player not implemented"))
}

// Prediction is the resolver for the prediction field.
func (r *playerGameResolver) Prediction(ctx context.Context, obj *model.PlayerGame, input model.ModelInput) (*model.AverageStats, error) {
	prediction := model.AverageStats{}
	// all gamebreakdowns
	breakdownAverages := make([]*model.PlayerAverage, len(input.GameBreakdowns))
	distributeExtraWeight := 0.0
	for i, breakdown := range input.GameBreakdowns {
		games := []*model.PlayerGame{}
		for _, game := range obj.PlayerRef.GamesCache {
			// only get games before the current game
			breakdown.Filter.EndDate = &obj.Date
			// if opponentMatch is true, only get games previously matched up against opponent
			if breakdown.Filter.OpponentMatch != nil && *breakdown.Filter.OpponentMatch {
				breakdown.Filter.OpponentID = &obj.OpponentID
			}
			input.GameBreakdowns[i].Filter = breakdown.Filter
			if breakdown.Filter.MatchPlayerGame(game) {
				games = append(games, game)
				if breakdown.Filter.LastX != nil && len(games) >= *breakdown.Filter.LastX {
					break
				}
			}
		}
		pAvg := model.NewPlayerAverage(games, obj.PlayerRef)

		// Don't calculate averages if there are no games, or less than lastX games
		if len(games) == 0 || (breakdown.Filter.LastX != nil && len(games) < *breakdown.Filter.LastX) {
			distributeExtraWeight += breakdown.Weight
			input.GameBreakdowns[i].Weight = 0
		} else {
			breakdownAverages[i] = &pAvg
		}
	}

	similarPlayerAverages := []*model.PlayerAverage{}
	if input.SimilarPlayerInput != nil {
		similarPlayers, err := r.Db.GetSimilarPlayersFromMatrix(ctx, obj.PlayerID, *input.SimilarPlayerInput, obj.Date)
		if err != nil {
			logrus.Errorf("Error getting similar players: %v", err)
			return nil, err
		}
		for _, player := range similarPlayers {
			// Analyze player similarity
			x := []*model.PlayerGame{}
			f := model.GameFilter{
				EndDate: &obj.Date,
			}
			for _, game := range player.GamesCache {
				if f.MatchPlayerGame(game) {
					x = append(x, game)
				}
			}
			avg := model.NewPlayerAverage(x, &player)
			fmt.Printf("Similar player %v: Pts: %.2g, Rebs: %.2g, Asts: %.2g, 3PM: %.2g, FGA: %.2g, MIN: %.2g, Height: %.2g\n", player.Name, avg.Points, avg.Rebounds, avg.Assists, avg.ThreePointersMade, avg.FieldGoalsAttempted, avg.Minutes, avg.Height)

			filter := model.GameFilter{
				EndDate:    &obj.Date,
				OpponentID: &obj.OpponentID,
			}
			if len(*input.SimilarPlayerInput.PlayerPoolFilter.Seasons) > 0 {
				filter.Seasons = input.SimilarPlayerInput.PlayerPoolFilter.Seasons
			}
			games := []*model.PlayerGame{}
			for _, game := range player.GamesCache {
				if filter.MatchPlayerGame(game) {
					games = append(games, game)
				}
			}
			// Don't calculate averages if there are no games vs opponent
			if len(games) > 0 {
				pAvg := model.NewPlayerAverage(games, &player)
				similarPlayerAverages = append(similarPlayerAverages, &pAvg)
			}
		}
	}
	similarPlayerWeights := input.SimilarPlayerInput.Weight / float64(len(similarPlayerAverages))

	validBreakdowns := 0
	for _, breakdown := range input.GameBreakdowns {
		if breakdown.Weight > 0 {
			validBreakdowns++
		}
	}

	//distribute extra weight evenly
	if distributeExtraWeight > 0 {
		for i, breakdown := range input.GameBreakdowns {
			if breakdown.Weight > 0 {
				input.GameBreakdowns[i].Weight += distributeExtraWeight / float64(validBreakdowns)
			}
		}
	}

	for i, avg := range breakdownAverages {
		if avg != nil {
			prediction.Assists += avg.Assists * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Blocks += avg.Blocks * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.DefensiveRebounds += avg.DefensiveRebounds * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.FieldGoalsAttempted += avg.FieldGoalsAttempted * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.FieldGoalsMade += avg.FieldGoalsMade * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.FreeThrowsAttempted += avg.FreeThrowsAttempted * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.FreeThrowsMade += avg.FreeThrowsMade * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.OffensiveRebounds += avg.OffensiveRebounds * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.PersonalFouls += avg.PersonalFouls * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.PersonalFoulsDrawn += avg.PersonalFoulsDrawn * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Points += avg.Points * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Rebounds += avg.Rebounds * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Steals += avg.Steals * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.ThreePointersAttempted += avg.ThreePointersAttempted * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.ThreePointersMade += avg.ThreePointersMade * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Turnovers += avg.Turnovers * (input.GameBreakdowns[i].Weight / 100.0)
			prediction.Minutes += avg.Minutes * (input.GameBreakdowns[i].Weight / 100.0)
			//height
			//weight
			fmt.Printf("Player %s (%v) pts: %v, rebs: %v, asts: %v\n", avg.Player.Name, avg.GamesPlayed, avg.Points, avg.Rebounds, avg.Assists)
		}
	}
	for _, avg := range similarPlayerAverages {
		prediction.Assists += avg.Assists * (similarPlayerWeights / 100.0)
		prediction.Blocks += avg.Blocks * (similarPlayerWeights / 100.0)
		prediction.DefensiveRebounds += avg.DefensiveRebounds * (similarPlayerWeights / 100.0)
		prediction.FieldGoalsAttempted += avg.FieldGoalsAttempted * (similarPlayerWeights / 100.0)
		prediction.FieldGoalsMade += avg.FieldGoalsMade * (similarPlayerWeights / 100.0)
		prediction.FreeThrowsAttempted += avg.FreeThrowsAttempted * (similarPlayerWeights / 100.0)
		prediction.FreeThrowsMade += avg.FreeThrowsMade * (similarPlayerWeights / 100.0)
		prediction.OffensiveRebounds += avg.OffensiveRebounds * (similarPlayerWeights / 100.0)
		prediction.PersonalFouls += avg.PersonalFouls * (similarPlayerWeights / 100.0)
		prediction.PersonalFoulsDrawn += avg.PersonalFoulsDrawn * (similarPlayerWeights / 100.0)
		prediction.Points += avg.Points * (similarPlayerWeights / 100.0)
		prediction.Rebounds += avg.Rebounds * (similarPlayerWeights / 100.0)
		prediction.Steals += avg.Steals * (similarPlayerWeights / 100.0)
		prediction.ThreePointersAttempted += avg.ThreePointersAttempted * (similarPlayerWeights / 100.0)
		prediction.ThreePointersMade += avg.ThreePointersMade * (similarPlayerWeights / 100.0)
		prediction.Turnovers += avg.Turnovers * (similarPlayerWeights / 100.0)
		prediction.Minutes += avg.Minutes * (similarPlayerWeights / 100.0)
		fmt.Printf("Similar %s (%v) pts: %v, rebs: %v, asts: %v\n", avg.Player.Name, avg.GamesPlayed, avg.Points, avg.Rebounds, avg.Assists)
	}
	prediction.Assists = math.RoundFloat(prediction.Assists, 2)
	prediction.Blocks = math.RoundFloat(prediction.Blocks, 2)
	prediction.DefensiveRebounds = math.RoundFloat(prediction.DefensiveRebounds, 2)
	prediction.FieldGoalsAttempted = math.RoundFloat(prediction.FieldGoalsAttempted, 2)
	prediction.FieldGoalsMade = math.RoundFloat(prediction.FieldGoalsMade, 2)
	prediction.FreeThrowsAttempted = math.RoundFloat(prediction.FreeThrowsAttempted, 2)
	prediction.FreeThrowsMade = math.RoundFloat(prediction.FreeThrowsMade, 2)
	prediction.OffensiveRebounds = math.RoundFloat(prediction.OffensiveRebounds, 2)
	prediction.PersonalFouls = math.RoundFloat(prediction.PersonalFouls, 2)
	prediction.PersonalFoulsDrawn = math.RoundFloat(prediction.PersonalFoulsDrawn, 2)
	prediction.Points = math.RoundFloat(prediction.Points, 2)
	prediction.Rebounds = math.RoundFloat(prediction.Rebounds, 2)
	prediction.Steals = math.RoundFloat(prediction.Steals, 2)
	prediction.ThreePointersAttempted = math.RoundFloat(prediction.ThreePointersAttempted, 2)
	prediction.ThreePointersMade = math.RoundFloat(prediction.ThreePointersMade, 2)
	prediction.Turnovers = math.RoundFloat(prediction.Turnovers, 2)
	prediction.Minutes = math.RoundFloat(prediction.Minutes, 2)
	fmt.Println()
	return &prediction, nil
}

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// PlayerGame returns generated.PlayerGameResolver implementation.
func (r *Resolver) PlayerGame() generated.PlayerGameResolver { return &playerGameResolver{r} }

type playerResolver struct{ *Resolver }
type playerGameResolver struct{ *Resolver }
