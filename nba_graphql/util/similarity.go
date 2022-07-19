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
	players = normalizePlayer(players)
	for i := range players {
		if players[i].Player.PlayerID == toPlayer.Player.PlayerID {
			toPlayer = players[i]
			break
		}
	}
	var playerDistances map[float64]model.Player = make(map[float64]model.Player, len(players))
	for _, player := range players {
		if player.Player.PlayerID == toPlayer.Player.PlayerID {
			continue
		}
		if player.Player.Height == "" {
			continue
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

func SimilarTeams(teams []model.TeamAverage, toTeam model.TeamAverage) []*model.Team {
	teams = normalizeTeam(teams)
	for i := range teams {
		if teams[i].Team.TeamID == toTeam.Team.TeamID {
			toTeam = teams[i]
			break
		}
	}
	var teamDistances map[float64]model.Team = make(map[float64]model.Team, len(teams))
	for _, team := range teams {
		if team.Team.TeamID == toTeam.Team.TeamID {
			continue
		}
		distance := TeamEuclideanDistance(team, toTeam)
		_, found := teamDistances[distance]
		for found {
			distance += 0.00001
			_, found = teamDistances[distance]
		}
		teamDistances[distance] = team.Team
	}
	var distances []float64
	for k := range teamDistances {
		distances = append(distances, k)
	}
	sort.Float64s(distances)
	//TODO: add a limit to the filter
	var closestTeams []*model.Team = make([]*model.Team, 0, 4)
	for i, distance := range distances {
		team := teamDistances[distance]
		closestTeams = append(closestTeams, &team)
		if i == 3 {
			break
		}
	}
	return closestTeams
}

func EuclideanDistance(player1 model.PlayerAverage, player2 model.PlayerAverage) float64 {
	sumDistance := distance("Points", player1.Points, player2.Points) +
		distance("Assists", player1.Assists, player2.Assists) +
		distance("Rebounds", player1.Rebounds, player2.Rebounds) +
		distance("Steals", player1.Steals, player2.Steals) +
		distance("Blocks", player1.Blocks, player2.Blocks) +
		distance("Turnovers", player1.Turnovers, player2.Turnovers) +
		distance("Minutes", player1.Minutes, player2.Minutes) +
		distance("FieldGoalsMade", player1.FieldGoalsMade, player2.FieldGoalsMade) +
		distance("FieldGoalsAttempted", player1.FieldGoalsAttempted, player2.FieldGoalsAttempted) +
		distance("ThreePointersMade", player1.ThreePointersMade, player2.ThreePointersMade) +
		distance("ThreePointersAttempted", player1.ThreePointersAttempted, player2.ThreePointersAttempted) +
		distance("FreeThrowsMade", player1.FreeThrowsMade, player2.FreeThrowsMade) +
		distance("FreeThrowsAttempted", player1.FreeThrowsAttempted, player2.FreeThrowsAttempted) +
		distance("Usage", player1.Usage, player2.Usage) +
		distance("Height", player1.Height, player2.Height) +
		distance("Weight", player1.Weight, player2.Weight)
	// fmt.Println(sumDistance)
	return math.Sqrt(sumDistance)
}

func TeamEuclideanDistance(team1 model.TeamAverage, team2 model.TeamAverage) float64 {
	sumDistance := distance("Points", team1.Points, team2.Points)*1.5 +
		distance("OppPoints", team1.OppPoints, team2.OppPoints)*1.5 +
		distance("Assists", team1.Assists, team2.Assists) +
		distance("OppAssists", team1.OppAssists, team2.OppAssists) +
		distance("OppRebounds", team1.OppRebounds, team2.OppRebounds) +
		distance("Rebounds", team1.Rebounds, team2.Rebounds) +
		distance("Steals", team1.Steals, team2.Steals) +
		distance("Blocks", team1.Blocks, team2.Blocks) +
		distance("Turnovers", team1.Turnovers, team2.Turnovers) +
		distance("ThreePointersMade", team1.ThreePointersMade, team2.ThreePointersMade) +
		distance("PersonalFouls", team1.PersonalFouls, team2.PersonalFouls) +
		distance("PersonalFoulsDrawn", team1.PersonalFoulsDrawn, team2.PersonalFoulsDrawn) +
		distance("GamesWon", team1.GamesWon, team2.GamesWon)*1.5 +
		distance("GamesLost", team1.GamesLost, team2.GamesLost)*1.5
	// fmt.Println(sumDistance)
	return math.Sqrt(sumDistance)
}

func distance(s string, p float64, q float64) float64 {
	r := math.Pow((p - q), 2)
	return r
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

func mean(values []float64) float64 {
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func std(values []float64) float64 {
	mean := mean(values)
	var sum float64
	for _, value := range values {
		sum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(sum / float64(len(values)))
}

func normalize(values []float64) []float64 {
	mean := mean(values)
	std := std(values)
	var normalized []float64
	for _, value := range values {
		normalized = append(normalized, (value-mean)/std)
	}
	return normalized
}

func normalizePlayer(players []model.PlayerAverage) []model.PlayerAverage {
	points := make([]float64, 0, len(players))
	gamesPlayed := make([]float64, 0, len(players))
	assists := make([]float64, 0, len(players))
	rebounds := make([]float64, 0, len(players))
	steals := make([]float64, 0, len(players))
	blocks := make([]float64, 0, len(players))
	turnovers := make([]float64, 0, len(players))
	minutes := make([]float64, 0, len(players))
	fieldGoalsMakes := make([]float64, 0, len(players))
	fieldGoalsAttempts := make([]float64, 0, len(players))
	threePointersMakes := make([]float64, 0, len(players))
	threePointersAttempts := make([]float64, 0, len(players))
	freeThrowsMakes := make([]float64, 0, len(players))
	freeThrowsAttempts := make([]float64, 0, len(players))
	usages := make([]float64, 0, len(players))
	heights := make([]float64, 0, len(players))
	weights := make([]float64, 0, len(players))
	for _, player := range players {
		points = append(points, player.Points)
		gamesPlayed = append(gamesPlayed, player.GamesPlayed)
		assists = append(assists, player.Assists)
		rebounds = append(rebounds, player.Rebounds)
		steals = append(steals, player.Steals)
		blocks = append(blocks, player.Blocks)
		turnovers = append(turnovers, player.Turnovers)
		minutes = append(minutes, player.Minutes)
		fieldGoalsMakes = append(fieldGoalsMakes, player.FieldGoalsMade)
		fieldGoalsAttempts = append(fieldGoalsAttempts, player.FieldGoalsAttempted)
		threePointersMakes = append(threePointersMakes, player.ThreePointersMade)
		threePointersAttempts = append(threePointersAttempts, player.ThreePointersAttempted)
		freeThrowsMakes = append(freeThrowsMakes, player.FreeThrowsMade)
		freeThrowsAttempts = append(freeThrowsAttempts, player.FreeThrowsAttempted)
		usages = append(usages, player.Usage)
		heights = append(heights, heightInches(player.Player.Height))
		weights = append(weights, float64(player.Player.Weight))
	}
	points = normalize(points)
	gamesPlayed = normalize(gamesPlayed)
	assists = normalize(assists)
	rebounds = normalize(rebounds)
	steals = normalize(steals)
	blocks = normalize(blocks)
	turnovers = normalize(turnovers)
	minutes = normalize(minutes)
	fieldGoalsMakes = normalize(fieldGoalsMakes)
	fieldGoalsAttempts = normalize(fieldGoalsAttempts)
	threePointersMakes = normalize(threePointersMakes)
	threePointersAttempts = normalize(threePointersAttempts)
	freeThrowsMakes = normalize(freeThrowsMakes)
	freeThrowsAttempts = normalize(freeThrowsAttempts)
	usages = normalize(usages)
	heights = normalize(heights)
	weights = normalize(weights)
	for i := range players {
		players[i].Points = points[i]
		players[i].GamesPlayed = gamesPlayed[i]
		players[i].Assists = assists[i]
		players[i].Rebounds = rebounds[i]
		players[i].Steals = steals[i]
		players[i].Blocks = blocks[i]
		players[i].Turnovers = turnovers[i]
		players[i].Minutes = minutes[i]
		players[i].FieldGoalsMade = fieldGoalsMakes[i]
		players[i].FieldGoalsAttempted = fieldGoalsAttempts[i]
		players[i].ThreePointersMade = threePointersMakes[i]
		players[i].ThreePointersAttempted = threePointersAttempts[i]
		players[i].FreeThrowsMade = freeThrowsMakes[i]
		players[i].FreeThrowsAttempted = freeThrowsAttempts[i]
		players[i].Usage = usages[i]
		players[i].Height = heights[i]
		players[i].Weight = weights[i]
	}
	return players
}

func normalizeTeam(teams []model.TeamAverage) []model.TeamAverage {
	gamesWon := make([]float64, 0, len(teams))
	gamesLost := make([]float64, 0, len(teams))
	points := make([]float64, 0, len(teams))
	oppPoints := make([]float64, 0, len(teams))
	assists := make([]float64, 0, len(teams))
	oppAssists := make([]float64, 0, len(teams))
	rebounds := make([]float64, 0, len(teams))
	oppRebounds := make([]float64, 0, len(teams))
	steals := make([]float64, 0, len(teams))
	blocks := make([]float64, 0, len(teams))
	turnovers := make([]float64, 0, len(teams))
	threePointersMakes := make([]float64, 0, len(teams))
	personalFouls := make([]float64, 0, len(teams))
	personalFoulsDrawn := make([]float64, 0, len(teams))

	for _, team := range teams {
		gamesWon = append(gamesWon, team.GamesWon)
		gamesLost = append(gamesLost, team.GamesLost)
		points = append(points, team.Points)
		assists = append(assists, team.Assists)
		rebounds = append(rebounds, team.Rebounds)
		oppPoints = append(oppPoints, team.OppPoints)
		oppAssists = append(oppAssists, team.OppAssists)
		oppRebounds = append(oppRebounds, team.OppRebounds)
		steals = append(steals, team.Steals)
		blocks = append(blocks, team.Blocks)
		turnovers = append(turnovers, team.Turnovers)
		threePointersMakes = append(threePointersMakes, team.ThreePointersMade)
		personalFouls = append(personalFouls, team.PersonalFouls)
		personalFoulsDrawn = append(personalFoulsDrawn, team.PersonalFoulsDrawn)
	}
	gamesWon = normalize(gamesWon)
	gamesLost = normalize(gamesLost)
	points = normalize(points)
	oppPoints = normalize(oppPoints)
	assists = normalize(assists)
	oppAssists = normalize(oppAssists)
	rebounds = normalize(rebounds)
	oppRebounds = normalize(oppRebounds)
	steals = normalize(steals)
	blocks = normalize(blocks)
	turnovers = normalize(turnovers)
	threePointersMakes = normalize(threePointersMakes)
	personalFouls = normalize(personalFouls)
	personalFoulsDrawn = normalize(personalFoulsDrawn)
	for i := range teams {
		teams[i].GamesWon = gamesWon[i]
		teams[i].GamesLost = gamesLost[i]
		teams[i].Points = points[i]
		teams[i].Assists = assists[i]
		teams[i].Rebounds = rebounds[i]
		teams[i].OppPoints = oppPoints[i]
		teams[i].OppAssists = oppAssists[i]
		teams[i].OppRebounds = oppRebounds[i]
		teams[i].Steals = steals[i]
		teams[i].Blocks = blocks[i]
		teams[i].Turnovers = turnovers[i]
		teams[i].ThreePointersMade = threePointersMakes[i]
		teams[i].PersonalFouls = personalFouls[i]
		teams[i].PersonalFoulsDrawn = personalFoulsDrawn[i]
	}
	return teams
}
