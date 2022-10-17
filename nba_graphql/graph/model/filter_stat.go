package model

import (
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

func (f StatFilter) String() string {
	return util.Print(f)
}

func (f *StatFilter) MatchPlayer(player *Player) bool {
	if player.GamesCache == nil || len(player.GamesCache) == 0 {
		return false
	}
	var games []*PlayerGame
	for _, game := range player.GamesCache {
		if f.Period == nil || f.Period.MatchGame(game) {
			games = append(games, game)
		}
	}
	avgPerGame := NewPlayerAverage(games, player)
	switch f.Mode {
	case StatModePerGame:
		score := avgPerGame.Score(f.Stat)
		return f.Operator.Evaluate(score, f.Value)
	case StatModePerMinute:
		score := avgPerGame.Score(f.Stat) / avgPerGame.Score(Minutes)
		return f.Operator.Evaluate(score, f.Value)
	case StatModePer36:
		score := avgPerGame.Score(f.Stat) / avgPerGame.Score(Minutes) * 36
		return f.Operator.Evaluate(score, f.Value)
	case StatModeTotal:
		score := 0.0
		for _, game := range player.GamesCache {
			score += game.Score(f.Stat)
		}
		return f.Operator.Evaluate(score, f.Value)
	default:
		logrus.Warn("StatFilter.MatchPlayer: unknown mode")
		return true
	}
}

func (f *StatFilter) MatchGame(game *PlayerGame) bool {
	if f.Period != nil && !f.Period.MatchGame(game) {
		return false
	}
	switch f.Mode {
	case StatModeTotal:
		fallthrough
	case StatModePerGame:
		score := game.Score(f.Stat)
		return f.Operator.Evaluate(score, f.Value)
	case StatModePerMinute:
		score := game.Score(f.Stat) / game.Score(Minutes)
		return f.Operator.Evaluate(score, f.Value)
	case StatModePer36:
		score := game.Score(f.Stat) / game.Score(Minutes) * 36
		return f.Operator.Evaluate(score, f.Value)
	default:
		logrus.Warn("StatFilter.MatchPlayer: unknown mode")
		return true
	}
}
