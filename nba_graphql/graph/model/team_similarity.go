package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	similarity "github.com/zvandehy/DataTrain/nba_graphql/math"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

type TeamSnapshots map[string]TeamSimilarityMatrix

func NewTeamSnapshots() *TeamSnapshots {
	s := make(TeamSnapshots)
	return &s
}

func (m *TeamSnapshots) Key(startDate, endDate string) string {
	return fmt.Sprintf("%s-%s", startDate, endDate)
}

func (s *TeamSnapshots) AddSnapshot(startDate, endDate time.Time, teams []*Team) {
	averages := []TeamAverage{}
	for _, team := range teams {
		var games []*TeamGame
		for _, game := range team.GamesCache {
			date, err := time.Parse(util.DATE_TIME_FORMAT, game.Date)
			if err != nil {
				logrus.Errorf("error parsing team game date: %v", err)
			}
			if date.After(startDate) && date.Before(endDate) {
				games = append(games, game)
			}
		}
		if len(games) > 0 {
			averages = append(averages, NewTeamAverage(games, team))
		} else {
			logrus.Errorf("no games for team %s", team.Name)
		}
	}
	key := s.Key(startDate.Format(util.DATE_FORMAT), endDate.Format(util.DATE_FORMAT))
	(*s)[key] = *NewTeamSimilarityMatrix(averages)
}

func (s *TeamSnapshots) GetSimilarTeams(teamID, limit int, startDate, endDate string, statsOfInterest []Stat) []*Team {
	key := s.Key(startDate, endDate)
	if snapshot, ok := (*s)[key]; ok {
		similarTeams, err := snapshot.GetNearestTeams(teamID, limit, statsOfInterest)
		if err != nil {
			logrus.Errorf("error getting similar teams from matrix '%v': %v", key, err)
			logrus.Errorf("%+v", *s)
			return []*Team{}
		}
		return similarTeams
	}
	logrus.Errorf("snapshot '%v' not found in:\n\t%+v", key, *s)
	return []*Team{}
}

type TeamSimilarityMatrix struct {
	Matrix map[int]TeamSimilarityVector
}

type TeamSimilarityVector struct {
	//Comparisons show the difference between the average of this player and the average of other players (e.g. thisPlayer.points - otherPlayer.points)
	Comparisons map[int]TeamDiff
	//Average is the average stats of this player, so that they can be compared to other (new) players
	Average TeamAverage
}

func (v *TeamSimilarityVector) GetNearest(limit int, statsOfInterest []Stat) []TeamDiff {
	// get limit number of nearest teams using EuclideanDistance
	nearest := make([]TeamDiff, 0, len(v.Comparisons))
	for _, diff := range v.Comparisons {
		nearest = append(nearest, diff)
	}
	sort.Slice(nearest, func(i, j int) bool {
		return TeamEuclideanDistance(nearest[i], statsOfInterest) < TeamEuclideanDistance(nearest[j], statsOfInterest)
	})
	if len(v.Comparisons) <= limit {
		return nearest
	}
	return nearest[:limit]
}

func (m *TeamSimilarityMatrix) AddTeams(teams []TeamAverage) error {
	for _, team := range teams {
		if err := m.AddTeam(team.Team.TeamID, team); err != nil {
			return err
		}
	}
	return nil
}

func (s *StatOfInterest) AddTeam(teamAverage TeamAverage) {
	score := teamAverage.Score(s.Stat)
	s.Scores = append(s.Scores, score)
	s.StdDev, s.Mean = similarity.StdDevMean(s.Scores)
}

func (s *StatOfInterest) ZScoreTeam(teamAverage TeamAverage) float64 {
	return similarity.ZScore(teamAverage.Score(s.Stat), s.Mean, s.StdDev)
}

func (m *TeamSimilarityMatrix) AddNormalizedTeams(teams []TeamAverage) {
	statsOfInterest := TeamAverageStats()
	stats := make([]StatOfInterest, len(statsOfInterest))
	for i, input := range statsOfInterest {
		stat, err := NewStat(string(input))
		if err != nil {
			logrus.Warning("Stat of interest not found: ", stat)
			continue
		}
		soi := NewStatOfInterest(stat)
		for _, team := range teams {
			soi.AddTeam(team)
		}
		stats[i] = *soi
	}

	for _, team := range teams {
		// fmt.Printf("Adding: %12.12v\t [OppPts:%5.5v, OppRebs:%5.5v, OppAsts:%5.5v]\n", team.Team.Name, team.OppPoints, team.OppRebounds, team.OppAssists)
		normalized := team.Normalize(stats...)
		m.AddTeam(team.Team.TeamID, normalized)
	}
}

func NewTeamSimilarityMatrix(teams []TeamAverage) *TeamSimilarityMatrix {
	m := &TeamSimilarityMatrix{
		Matrix: make(map[int]TeamSimilarityVector),
	}
	m.AddNormalizedTeams(teams)
	return m
}

// AddTeam adds a team to the matrix and calculates the comparisons between the new team and the existing teams.
func (m *TeamSimilarityMatrix) AddTeam(teamID int, teamAverage TeamAverage) error {
	if _, ok := m.Matrix[teamID]; !ok {
		m.Matrix[teamID] = TeamSimilarityVector{
			Comparisons: m.CompareAverages(teamID, teamAverage),
			Average:     teamAverage,
		}
		return nil
	}
	// Didn't add player because already exists or something went wrong
	return fmt.Errorf("team %d already exists in matrix", teamID)
}

func (m *TeamSimilarityMatrix) CompareAverages(in int, averageIn TeamAverage) map[int]TeamDiff {
	comparisons := make(map[int]TeamDiff, len(m.Matrix))
	for teamID, vector := range m.Matrix {
		if teamID == in {
			continue
		}
		// add comparison between averageIn and other teams's average to the new team's vector
		comparisons[teamID] = vector.Average.Difference(averageIn)
		// add comparison between other team's average and averageIn to the other team's vector
		m.Matrix[teamID].Comparisons[in] = averageIn.Difference(vector.Average)
	}
	return comparisons
}

func (m *TeamSimilarityMatrix) GetNearestTeams(toTeam int, limit int, statsOfInterest []Stat) (similarTeams []*Team, err error) {
	if vector, ok := m.Matrix[toTeam]; ok {
		for _, team := range vector.GetNearest(limit, statsOfInterest) {
			similarTeams = append(similarTeams, &team.Team)
		}
		return similarTeams, nil
	}
	return nil, fmt.Errorf("team %d not found in matrix", toTeam)
}
