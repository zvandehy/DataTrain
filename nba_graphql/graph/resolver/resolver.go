package resolver

import (
	"context"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	database.BasketballRepository
}

// NewResolver returns a new instance of the resolver.
func NewResolver(db database.BasketballRepository) *Resolver {
	return &Resolver{BasketballRepository: db}
}

func (r *Resolver) GetBase(ctx context.Context, inputs []*model.GameBreakdownInput, game *model.PlayerGame) []*model.PlayerGame {
	//TODO: Add ALLTIME support
	seasons := []model.SeasonOption{}
	seasons = append(seasons, game.Season)
	filter := model.GameFilter{
		PlayerID: &game.PlayerID,
		Seasons:  &seasons,
	}
	for _, input := range inputs {
		if input.Filter.PreviousSeasonMatch != nil && *input.Filter.PreviousSeasonMatch {
			if game.Season == model.SEASON_2022_23 {
				seasons = append(seasons, model.SEASON_2021_22)
			} else if game.Season == model.SEASON_2021_22 {
				seasons = append(seasons, model.SEASON_2020_21)
			}
			filter.Seasons = &seasons
			break
		}
	}
	games, err := dataloader.For(ctx).PlayerGameByFilter.Load(filter)
	gamesBeforeThisGame := []*model.PlayerGame{}
	for _, g := range games {
		if g.Date.Before(*game.Date) {
			gamesBeforeThisGame = append(gamesBeforeThisGame, g)
		}
	}
	if err != nil {
		logrus.Errorf("Error getting base player games: %v", err)
		return []*model.PlayerGame{}
	}
	return gamesBeforeThisGame
}

func (r *Resolver) GetGamelogBreakdowns(ctx context.Context, inputs []*model.GameBreakdownInput, game *model.PlayerGame, target *float64, stat model.Stat) []*model.PropBreakdown {
	base := r.GetBase(ctx, inputs, game)
	if len(base) == 0 {
		return []*model.PropBreakdown{}
	}

	baseTotal := 0.0
	for _, game := range base {
		baseTotal += game.Score(stat)
	}
	baseAvg := float64(baseTotal) / float64(len(base))
	if target == nil {
		target = &baseAvg
	}
	date := game.Date.Format(util.DATE_FORMAT)
	breakdowns := []*model.PropBreakdown{}
	for _, input := range inputs {
		seasons := []model.SeasonOption{}
		filter := model.GameFilter{
			PlayerID: &game.PlayerID,
			EndDate:  &date,
		}
		if input.Filter.SeasonMatch != nil && *input.Filter.SeasonMatch {
			// query for games the player has played in the same season
			seasons = append(seasons, game.Season)
		}
		if input.Filter.PreviousSeasonMatch != nil && *input.Filter.PreviousSeasonMatch {
			// query for games the player has played in the previous season
			if game.Season == model.SEASON_2022_23 {
				seasons = append(seasons, model.SEASON_2021_22)
			} else if game.Season == model.SEASON_2021_22 {
				seasons = append(seasons, model.SEASON_2020_21)
			}
		}
		filter.Seasons = &seasons
		if input.Filter.GameTypeMatch != nil && *input.Filter.GameTypeMatch {
			// query for games the player has played in the same game type
			gameType := model.GameTypeRegularSeason
			if game.Playoffs {
				gameType = model.GameTypePlayoffs
			}
			filter.GameType = &gameType
		}
		if input.Filter.HomeOrAwayMatch != nil && *input.Filter.HomeOrAwayMatch {
			// query for games the player has played in the same home or away
			filter.HomeOrAway = &game.HomeOrAway
		}
		if input.Filter.OpponentMatch != nil && *input.Filter.OpponentMatch {
			// query for games the player has played in the same opponent
			filter.OpponentID = &game.OpponentID
		}

		derivedGames := []*model.PlayerGame{}
		for _, game := range base {
			if filter.MatchPlayerGame(game) {
				derivedGames = append(derivedGames, game)
			}
		}
		if input.Filter.LastX != nil && *input.Filter.LastX > 0 {
			// todo: make sure games is sorted by date and we are taking the most recent x games
			if len(derivedGames) > *input.Filter.LastX {
				derivedGames = derivedGames[len(derivedGames)-*input.Filter.LastX:]
			}
		}

		breakdown, err := CalculateGamelogBreakdown(input, derivedGames, baseAvg, *target, stat)
		if err != nil {
			logrus.Errorf("Error calculating gamelog breakdown: %v", err)
			continue
		}
		breakdowns = append(breakdowns, breakdown)
	}
	distribute := 0.0
	skipped := 0.0
	for _, breakdown := range breakdowns {
		if breakdown.Over+breakdown.Under+breakdown.Push == 0 {
			distribute += breakdown.Weight
			skipped++
			breakdown.Weight = 0
		}
	}
	if skipped > 0 {
		for _, breakdown := range breakdowns {
			if breakdown.Weight == 0 {
				continue
			}
			breakdown.Weight += distribute / skipped
		}
	}
	// for _, breakdown := range breakdowns {
	// 	logrus.Warnf("%+v", *breakdown)
	// }
	return breakdowns
}

func (r *Resolver) GetSimilarPlayerBreakdowns(ctx context.Context, input *model.SimilarPlayerInput, game *model.PlayerGame, stat model.Stat) []*model.PropBreakdown {

	// TODO: Allow getting games for this season, this and last season, or all time
	similarPlayers, err := r.GetSimilarPlayers(ctx, game.PlayerID, input, game.Date)
	if err != nil || len(similarPlayers) == 0 {
		logrus.Errorf("Error getting similar players: %v", err)
		return []*model.PropBreakdown{}
	}

	breakdowns := []*model.PropBreakdown{}
	skipped := 0
	for _, similarPlayer := range similarPlayers {
		// similar players query returns the player with their base games for this game
		similarBase := similarPlayer.GamesCache
		vsOpp := []*model.PlayerGame{}
		derivedTotal := 0.0
		baseTotal := 0.0
		for _, game := range similarBase {
			score := game.Score(stat)
			baseTotal += score
			if game.OpponentID == game.OpponentID {
				derivedTotal += score
				vsOpp = append(vsOpp, game)
			}
		}
		baseAvg := baseTotal / float64(len(similarBase))
		if len(vsOpp) == 0 {
			skipped++
			continue
		}
		derivedAvg := derivedTotal / float64(len(similarBase))
		diff := derivedAvg - baseAvg
		diffPct := diff / baseAvg
		breakdown := &model.PropBreakdown{
			Name:           fmt.Sprintf("%s vs %s", similarPlayer.Name, "Opponent"),
			Over:           0,
			Under:          0,
			Push:           0,
			OverPct:        0.0,
			UnderPct:       0.0,
			PushPct:        0.0,
			Base:           baseAvg,
			DerivedAverage: derivedAvg,
			PctChange:      diffPct,
			DerivedGames:   vsOpp,
			Weight:         input.Weight,
		}
		// over/under for similar players determnined by if the player scored higher or lower than their base average
		for _, game := range vsOpp {
			score := game.Score(stat)
			if score > baseAvg {
				breakdown.Over++
			}
			if score < baseAvg {
				breakdown.Under++
			}
			if score == baseAvg {
				breakdown.Push++
			}
		}
		breakdown.OverPct = float64(breakdown.Over) / float64(len(vsOpp))
		breakdown.UnderPct = float64(breakdown.Under) / float64(len(vsOpp))
		breakdown.PushPct = float64(breakdown.Push) / float64(len(vsOpp))
		breakdowns = append(breakdowns, breakdown)
	}
	if len(breakdowns) == 0 {
		logrus.Warnf("No similar players have played the opponent: %d vs %d", game.PlayerID, game.OpponentID)
		return []*model.PropBreakdown{}
	}
	for _, breakdown := range breakdowns {
		breakdown.Weight = (input.Weight / float64(len(breakdowns)))
	}
	return breakdowns
}

func CalculateGamelogBreakdown(input *model.GameBreakdownInput, derivedGames []*model.PlayerGame, baseAvg float64, target float64, stat model.Stat) (*model.PropBreakdown, error) {
	breakdown := model.PropBreakdown{
		Name:           input.Name,
		Over:           0,
		Under:          0,
		Push:           0,
		OverPct:        0.0,
		UnderPct:       0.0,
		PushPct:        0.0,
		DerivedAverage: 0.0,
		Weight:         input.Weight, // TODO: Handle weight of gamelog breakdown with 0 games
		DerivedGames:   derivedGames,
		PctChange:      0.0,
		Base:           baseAvg,
		StdDev:         0.0,
	}
	if len(derivedGames) == 0 {
		return &breakdown, nil
	}
	sum := 0.0
	count := 0
	for _, game := range derivedGames {
		if game == nil || game.Outcome == model.GameOutcomePending.String() {
			continue
		}
		score := game.Score(stat)
		if score > target {
			breakdown.Over++
		} else if score < target {
			breakdown.Under++
		} else {
			breakdown.Push++
		}
		count++
		sum += score
	}
	breakdown.OverPct = float64(breakdown.Over) / float64(count)
	breakdown.UnderPct = float64(breakdown.Under) / float64(count)
	breakdown.PushPct = float64(breakdown.Push) / float64(count)
	breakdown.DerivedAverage = sum / float64(count)
	breakdown.StdDev = StandardDeviation(derivedGames, breakdown.DerivedAverage, stat)
	breakdown.PctChange = (breakdown.DerivedAverage - baseAvg) / baseAvg
	return &breakdown, nil
}

func StandardDeviation(games []*model.PlayerGame, mean float64, stat model.Stat) float64 {
	sum := 0.0
	for _, game := range games {
		sum += math.Pow(game.Score(stat)-mean, 2)
	}
	return math.Sqrt(sum / float64(len(games)))
}

// TODO: Weigh the variance by the assigned weight
func PoolVariance(datasets [][]float64) float64 {
	sum := 0.0
	means := []float64{}
	weights := []int{}
	n := 0
	for _, dataset := range datasets {
		sum += Variance(dataset) * float64(len(dataset))
		means = append(means, Mean(dataset))
		n += len(dataset)
		weights = append(weights, len(dataset))
	}
	meanOfVariances := sum / float64(n)
	varianceOfMeans := 0.0
	pooledMean := Mean(datasets...)
	for i := 0; i < len(means); i++ {
		varianceOfMeans += math.Pow(means[i]-pooledMean, 2) * float64(weights[i])
	}
	varianceOfMeans /= float64(n)
	return math.Round((meanOfVariances+varianceOfMeans)*100) / 100.0

}

func SumOfSquares(dataset []float64, weights ...int) float64 {
	if len(weights) > 0 && len(dataset) != len(weights) {
		logrus.Warnf("SumOfSquares: dataset and weights must be the same length")
		return 0.0
	}
	mean := Mean(dataset)
	sum := 0.0
	if len(weights) > 0 {
		for i, value := range dataset {
			sum += math.Pow(value-mean, 2) * float64(weights[i])
		}
	} else {
		for _, value := range dataset {
			sum += math.Pow(value-mean, 2)
		}
	}
	return sum
}

func Variance(dataset []float64) float64 {
	sum := SumOfSquares(dataset)
	return sum / float64(len(dataset))
}

func Mean(datasets ...[]float64) float64 {
	sum := 0.0
	total := 0
	for _, dataset := range datasets {
		total += len(dataset)
		for _, x := range dataset {
			sum += x
		}
	}
	return sum / float64(total)
}
