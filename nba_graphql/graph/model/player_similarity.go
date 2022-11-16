package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	similarity "github.com/zvandehy/DataTrain/nba_graphql/math"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

// PlayerSnapshots in time holding the similarity matrix for each date range.
// The key is "<startDate>-<endDate>" and value is a similarity Matrix.
// The value is a matrix that holds the average stats of players and comparisons between them.
type PlayerSnapshots map[string]PlayerSimilarityMatrix

func NewPlayerSnapshots() *PlayerSnapshots {
	s := make(PlayerSnapshots)
	return &s
}

func (s *PlayerSnapshots) Key(startDate, endDate string, playerPoolFilter PlayerFilter) string {
	return fmt.Sprintf("%s-%s-%s", startDate, endDate, playerPoolFilter.Key())
}

func (s *PlayerSnapshots) AddSnapshot(startDate, endDate time.Time, playerFilter *PlayerFilter, players []*Player) {
	averages := []PlayerAverage{}
	for _, player := range players {
		var games []*PlayerGame
		for _, game := range player.GamesCache {
			// date, err := time.Parse(util.DATE_FORMAT, game.Date)
			// if err != nil {
			// 	logrus.Errorf("error parsing game date: %v", err)
			// }
			if game.Date.After(startDate) && game.Date.Before(endDate) {
				games = append(games, game)
			}
		}
		if len(games) > 0 {
			averages = append(averages, NewPlayerAverage(games, player))
		}
	}
	(*s)[fmt.Sprintf("%s-%s-%s", startDate.Format(util.DATE_FORMAT), endDate.Format(util.DATE_FORMAT), playerFilter.Key())] = *NewPlayerSimilarityMatrix(averages)
}

func (s *PlayerSnapshots) GetSimilarPlayers(playerID, limit int, startDate, endDate string, playerFilter *PlayerFilter, statsOfInterest []Stat) []*Player {
	key := s.Key(startDate, endDate, *playerFilter)
	if snapshot, ok := (*s)[key]; ok {
		similarPlayers, err := snapshot.GetNearestPlayers(playerID, limit, statsOfInterest)
		if err != nil {
			logrus.Errorf("error getting similar players from matrix '%v': %v", key, err)
			return []*Player{}
		}
		return similarPlayers
	}
	return []*Player{}
}

type PlayerSimilarityMatrix struct {
	Matrix map[int]SimilarityVector
}

type SimilarityVector struct {
	//Comparisons show the difference between the average of this player and the average of other players (e.g. thisPlayer.points - otherPlayer.points)
	Comparisons map[int]PlayerDiff
	//Average is the average stats of this player, so that they can be compared to other (new) players
	Average PlayerAverage
}

func (v *SimilarityVector) GetNearest(limit int, statsOfInterest []Stat) []PlayerDiff {
	// get limit number of nearest players using EuclideanDistance
	nearest := make([]PlayerDiff, 0, len(v.Comparisons))
	for _, diff := range v.Comparisons {
		nearest = append(nearest, diff)
	}
	sort.Slice(nearest, func(i, j int) bool {
		return EuclideanDistance(nearest[i], statsOfInterest) < EuclideanDistance(nearest[j], statsOfInterest)
	})
	if len(v.Comparisons) <= limit {
		return nearest
	}
	return nearest[:limit]
}

func (m *PlayerSimilarityMatrix) AddPlayers(players []PlayerAverage) error {
	for _, player := range players {
		if err := m.AddPlayer(player.Player.PlayerID, player); err != nil {
			return err
		}
	}
	return nil
}

type StatOfInterest struct {
	Scores []float64
	Stat   Stat
	Mean   float64
	StdDev float64
}

func NewStatOfInterest(stat Stat) *StatOfInterest {
	return &StatOfInterest{
		Scores: []float64{},
		Stat:   stat,
	}
}

func (s *StatOfInterest) Add(playerAverage PlayerAverage) {
	score := playerAverage.Score(s.Stat)
	s.Scores = append(s.Scores, score)
	s.StdDev, s.Mean = similarity.StdDevMean(s.Scores)
}

func (s *StatOfInterest) ZScore(playerAverage PlayerAverage) float64 {
	return similarity.ZScore(playerAverage.Score(s.Stat), s.Mean, s.StdDev)
}

func (m *PlayerSimilarityMatrix) AddNormalizedPlayers(players []PlayerAverage) {
	statsOfInterest := PlayerAverageStats()
	stats := make([]StatOfInterest, len(statsOfInterest))
	for i, input := range statsOfInterest {
		stat, err := NewStat(string(input))

		if err != nil {
			logrus.Warnf("Stat of interest not found: %v", stat)
			continue
		}
		soi := NewStatOfInterest(stat)
		for _, player := range players {
			soi.Add(player)
		}
		stats[i] = *soi
	}

	for _, player := range players {
		normalized := player.Normalize(stats...)
		m.AddPlayer(player.Player.PlayerID, normalized)
	}
}

func NewPlayerSimilarityMatrix(players []PlayerAverage) *PlayerSimilarityMatrix {
	m := &PlayerSimilarityMatrix{
		Matrix: make(map[int]SimilarityVector),
	}
	m.AddNormalizedPlayers(players)
	return m
}

// AddPlayer adds a player to the matrix and calculates the comparisons between the new player and the existing players.
func (m *PlayerSimilarityMatrix) AddPlayer(playerID int, playerAverage PlayerAverage) error {
	if _, ok := m.Matrix[playerID]; !ok {
		m.Matrix[playerID] = SimilarityVector{
			Comparisons: m.CompareAverages(playerID, playerAverage),
			Average:     playerAverage,
		}
		return nil
	}
	// Didn't add player because already exists or something went wrong
	return fmt.Errorf("player %d already exists in matrix", playerID)
}

func (m *PlayerSimilarityMatrix) CompareAverages(in int, averageIn PlayerAverage) map[int]PlayerDiff {
	comparisons := make(map[int]PlayerDiff, len(m.Matrix))
	for playerID, vector := range m.Matrix {
		if playerID == in {
			continue
		}
		// add comparison between averageIn and other player's average to the new player's vector
		comparisons[playerID] = vector.Average.Difference(averageIn)
		// add comparison between other player's average and averageIn to the other player's vector
		m.Matrix[playerID].Comparisons[in] = averageIn.Difference(vector.Average)
	}
	return comparisons
}

func (m *PlayerSimilarityMatrix) GetNearestPlayers(toPlayer int, limit int, statsOfInterest []Stat) (similarPlayers []*Player, err error) {
	if vector, ok := m.Matrix[toPlayer]; ok {
		for _, player := range vector.GetNearest(limit, statsOfInterest) {
			similarPlayers = append(similarPlayers, &player.Player)
		}
		return similarPlayers, nil
	}
	return nil, fmt.Errorf("player %d not found in matrix", toPlayer)
}
