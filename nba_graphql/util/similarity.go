package util

import (
	"math"
	"sort"
	"strconv"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

//SimilarPlayers uses the euclidean distance formula to calculate the similarity between the given list of players and the target player
// and returns the 10 closest players to the target player "toPlayer"
func SimilarPlayers(players []model.PlayerAverage, toPlayer model.PlayerAverage) []*model.Player {
	var playerDistances map[float64]model.Player = make(map[float64]model.Player, len(players))
	for _, player := range players {
		if player.Player.PlayerID == toPlayer.Player.PlayerID {
			continue
		}
		if toPlayer.GamesPlayed >= 30 {
			if player.GamesPlayed < 30 || player.Minutes < toPlayer.Minutes-5 {
				continue
			}
		}
		distance := EuclideanDistance(player, toPlayer)
		_, found := playerDistances[distance]
		for found {
			distance += 0.00001
			_, found = playerDistances[distance]
		}
		playerDistances[distance] = player.Player
	}
	var distances []float64
	for k := range playerDistances {
		distances = append(distances, k)
	}
	sort.Float64s(distances)
	var closestPlayers []*model.Player = make([]*model.Player, 0, 10)
	for i, distance := range distances {
		player := playerDistances[distance]
		closestPlayers = append(closestPlayers, &player)
		if i == 9 {
			break
		}
	}
	return closestPlayers
}

func EuclideanDistance(player1 model.PlayerAverage, player2 model.PlayerAverage) float64 {
	sumDistance := distance(player1.Points, player2.Points) +
		distance(player1.Assists, player2.Assists) +
		distance(player1.Rebounds, player2.Rebounds) +
		distance(player1.Steals, player2.Steals) +
		distance(player1.Blocks, player2.Blocks) +
		distance(player1.Turnovers, player2.Turnovers) +
		distance(player1.Minutes, player2.Minutes) +
		distance(player1.FieldGoalsMade, player2.FieldGoalsMade) +
		distance(player1.FieldGoalsAttempted, player2.FieldGoalsAttempted) +
		distance(player1.ThreePointersMade, player2.ThreePointersMade) +
		distance(player1.ThreePointersAttempted, player2.ThreePointersAttempted) +
		distance(player1.FreeThrowsMade, player2.FreeThrowsMade) +
		distance(player1.FreeThrowsAttempted, player2.FreeThrowsAttempted) +
		distance(player1.Usage, player2.Usage) +
		distance(heightInches(player1.Player.Height), heightInches(player2.Player.Height)) +
		distance(float64(player1.Player.Weight), float64(player2.Player.Weight))
	return math.Sqrt(sumDistance)
}

func distance(p float64, q float64) float64 {
	return math.Pow((p - q), 2)
}

func heightInches(height string) float64 {
	feet, err := strconv.ParseFloat(string(rune(height[0])), 64)
	if err != nil {
		return 0
	}
	inches, err := strconv.ParseFloat(string(rune(height[2])), 64)
	if err != nil {
		return 0
	}
	return feet*12 + inches
}
