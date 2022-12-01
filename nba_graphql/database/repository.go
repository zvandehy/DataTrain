package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

var altNames map[string]string = map[string]string{
	"Nicolas Claxton":              "Nic Claxton",
	"OG Anunoby":                   "O.G. Anunoby",
	"Marcus Morris":                "Marcus Morris Sr.",
	"Nah'Shon Hyland":              "Bones Hyland",
	"Kevin \"Slim Reaper\" Durant": "Kevin Durant",
	"The Greek Freak":              "Giannis Antetokounmpo",
	"Greek Freak":                  "Giannis Antetokounmpo",
	"Ja Morantula":                 "Ja Morant",
	"James \"The Beard\" Harden":   "James Harden",
	"Ice Trae Young":               "Trae Young",
	"Ice Trae":                     "Trae Young",
}

type BasketballRepository interface {
	GetLeague() string
	CountQueries() int
	SetQueries(int)
	AddQuery()
	GetPlayers(ctx context.Context, withGames bool, playerFilters ...*model.PlayerFilter) ([]*model.Player, error)
	GetPropositions(ctx context.Context, propositionFilter *model.PropositionFilter) ([]*model.Proposition, error)
	SavePropositions(ctx context.Context, propositions []*model.Proposition) (int, error)
	SaveDBPropositions(ctx context.Context, propositions []*model.DBProposition) (int, error)
	SaveUpcomingGames(ctx context.Context, games []*model.PlayerGame) (int, error)
	// GetTeamsByID(ctx context.Context, teamIDs []int) ([]*model.Team, error)
	// GetTeamsByAbr(ctx context.Context, teamAbrs []string) ([]*model.Team, error)
	GetTeams(ctx context.Context, withGames bool, teamFilters ...*model.TeamFilter) ([]*model.Team, error)
	GetSimilarPlayers(ctx context.Context, similarToPlayerID int, input *model.SimilarPlayerInput, endDate *time.Time) ([]*model.Player, error)
	GetSimilarTeams(ctx context.Context, similarToTeamID int, input *model.SimilarTeamInput, endDate string) ([]*model.Team, error)
	GetPropositionsByPlayerGame(ctx context.Context, game model.PlayerGame) ([]*model.Proposition, error)
	GetPlayerGames(ctx context.Context, inputs ...model.GameFilter) ([]*model.PlayerGame, error)
	GetStandardizedPlayerStats(ctx context.Context, similarPlayerQuery model.SimilarPlayerQuery, toPlayerIDs ...int) ([]model.StandardizedPlayerStats, error)
}

var cacheSchedule *model.Schedule

func getSchedule() *model.Schedule {
	if cacheSchedule != nil {
		return cacheSchedule
	}
	resp, err := http.Get("https://cdn.nba.com/static/json/staticData/scheduleLeagueV2.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var schedule model.Schedule
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		panic(err)
	}
	cacheSchedule = &schedule
	return &schedule
}

func Getprizepicks(nbaClient BasketballRepository) {
	// TODO: THIS IS SLOWWWW
	leagueID := 7
	if strings.ToLower(nbaClient.GetLeague()) == "wnba" {
		leagueID = 3
	}
	start := time.Now()
	var propositions []*model.DBProposition
	var games []*model.PlayerGame
	url := fmt.Sprintf("https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=%d", leagueID)
	res, err := http.Get(url)
	if err != nil {
		logrus.Warnf("couldn't retrieve prizepicks projections for today: %v", err)
		return
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Warnf("couldn't read prizepicks projections for today: %v", err)
		return
	}
	// TODO: not the most elegant parsing, but it works
	var prizepicks model.PrizePicks
	if err := json.Unmarshal(bytes, &prizepicks); err != nil {
		logrus.Warnf("couldn't decode prizepicks projections for today: %v", err)
		return
	}
	var idToName map[string]string = make(map[string]string)
	for _, inc := range prizepicks.Included {
		idToName[inc.ID] = inc.Attributes.Name
	}
	schedule := getSchedule()
	countgames := 0
	countprops := 0
	for _, prop := range prizepicks.Data {
		if prop.Attributes.Is_promo {
			logrus.Warn("skipping promo")
			continue
		}
		p, game, err := ParsePrizePickProposition(nbaClient, *schedule, prop, idToName)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "prop": prop}).Warn("couldn't parse prizepicks proposition")
			continue
		}
		propositions = append(propositions, p)
		countprops++
		contains := false
		for _, g := range games {
			if g.PlayerID == game.PlayerID && g.GameID == game.GameID {
				contains = true
				break
			}
		}
		if !contains {
			games = append(games, game)
			countgames++
		}
		fmt.Println(len(propositions), p.PlayerName, p.Target, p.StatType)
		if len(propositions) >= 50 {
			//upsert
			x, err := nbaClient.SaveDBPropositions(context.Background(), propositions)
			if err != nil {
				if err != nil {
					logrus.Warn(err)
				}
			}
			fmt.Println("Uploaded ", x, " propositions")
			propositions = []*model.DBProposition{}
		}
	}
	if len(propositions) == 0 {
		logrus.Info("No prizepicks propositions to upload")
		return
	}
	fmt.Println("REMAINING TO IMPORT: ", len(propositions))
	//upsert
	_, err = nbaClient.SaveDBPropositions(context.Background(), propositions)
	if err != nil {
		if err != nil {
			logrus.Warn(err)
		}
	}
	fmt.Printf("Saving %d upcoming games\n", len(games))
	date := games[0].Date.Format(util.DATE_FORMAT)
	foundGames, err := nbaClient.GetPlayerGames(context.Background(), model.GameFilter{StartDate: &date, EndDate: &date})
	if err != nil {
		logrus.Warn(err)
	}
	fmt.Println("Found player games: ", len(foundGames), " Player Prop games: ", len(games))
	//remove foundGames from games
	for _, foundGame := range foundGames {
		for i, game := range games {
			if foundGame.PlayerID == game.PlayerID && foundGame.GameID == game.GameID {
				games = append(games[:i], games[i+1:]...)
				break
			}
		}
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].GameID < games[j].GameID
	})
	for _, game := range games {
		fmt.Println(game.PlayerID, game.GameID, game.Date, game.OpponentID)
	}

	if len(games) > 0 {
		_, err = nbaClient.SaveUpcomingGames(context.Background(), games)
		if err != nil {
			logrus.Errorf("failed to save upcoming games %+v", games)
		}
	}
	fmt.Println("Imported ", countprops, " propositions")
	logrus.Printf(util.TimeLog(fmt.Sprintf("Retrieved %s propositions from PrizePicks and inserted %d games", nbaClient.GetLeague(), len(games)), start))
}

func ParsePrizePickProposition(db BasketballRepository, schedule model.Schedule, prop model.PrizePicksData, itemIDToNameMap map[string]string) (proposition *model.DBProposition, game *model.PlayerGame, err error) {
	if prop.Attributes.Is_promo {
		logrus.Warn("skipping promo prizepick %+v", prop)
		return nil, nil, fmt.Errorf("skipping promo prizepick")
	}
	playerName, err := getPlayerName(prop, itemIDToNameMap)
	if err != nil {
		return nil, nil, err
	}
	statType, err := getStatType(prop, itemIDToNameMap)
	if err != nil {
		return nil, nil, err
	}
	target, err := getTarget(prop)
	if err != nil {
		return nil, nil, err
	}
	playerID, err := getPlayerID(db, playerName)
	if err != nil {
		return nil, nil, err
	}
	opponentID, err := getTeamID(db, prop.Attributes.Description)
	if err != nil {
		return nil, nil, err
	}
	gameID, err := getGameIDFromSchedule(schedule, prop.Attributes.Start_time, prop.Attributes.Description)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()

	proposition = &model.DBProposition{
		PlayerID:     playerID,
		GameID:       gameID,
		OpponentID:   opponentID,
		Sportsbook:   "PrizePicks",
		Target:       target,
		StatType:     statType,
		LastModified: &now,
		PlayerName:   playerName,
	}

	homeAway, err := getHomeTeam(schedule, prop.Attributes.Start_time, prop.Attributes.Description)
	if err != nil {
		return nil, nil, err
	}
	teamABR, err := getTeamABRFromSchedule(schedule, prop.Attributes.Start_time, prop.Attributes.Description)
	if err != nil {
		return nil, nil, err
	}
	teamID, err := getTeamID(db, teamABR)
	if err != nil {
		return nil, nil, err
	}
	dateSplit := strings.Split(prop.Attributes.Start_time, "T")
	date, err := time.Parse(util.DATE_FORMAT, dateSplit[0])
	if err != nil {
		return nil, nil, err
	}
	// save upcoming game in db
	game = &model.PlayerGame{
		PlayerID:   playerID,
		GameID:     gameID,
		OpponentID: opponentID,
		Date:       &date,
		HomeOrAway: homeAway,
		TeamID:     teamID,
		// Season:     schedule.Season,
		Season: "2022-23",
		// Playoffs:   schedule.Playoffs,
		Playoffs: false,
		Outcome:  "PENDING",
	}
	// logrus.Infof("saved upcoming game %v %v", game.GameID, game.PlayerID)
	return proposition, game, nil
}

func getPlayerName(prop model.PrizePicksData, itemIDToNameMap map[string]string) (string, error) {
	// get playername and statType from id to name mapping
	if val, ok := itemIDToNameMap[prop.Relationships.Player.Data.ID]; ok {
		val = strings.TrimSpace(val)
		if val == "" {
			return "", fmt.Errorf("empty player name")
		}
		if v, ok := altNames[val]; ok {
			return v, nil
		}
		return val, nil
	} else {
		return "", errors.New("could not find player name")
	}
}

func getStatType(prop model.PrizePicksData, itemIDToNameMap map[string]string) (string, error) {
	// get playername and statType from id to name mapping
	if val, ok := itemIDToNameMap[prop.Relationships.StatType.Data.ID]; ok {
		var err error
		statType, err := model.NewStat(val)
		if err != nil {
			return "", err
		}
		if statType == "" {
			err := fmt.Errorf("error retrieving prizepick stat type for %+v", prop)
			logrus.Error(err)
			return "", err
		}
		return statType.String(), nil
	} else {
		return "", errors.New("could not find stat type")
	}
}

func getTarget(prop model.PrizePicksData) (float64, error) {
	target, err := strconv.ParseFloat(prop.Attributes.Line_score, 64)
	if err != nil {
		return 0.0, errors.Wrap(err, "failed to retrieve prizepicks target")
	}
	return target, nil
}

func getPlayerID(db BasketballRepository, playerName string) (int, error) {
	if val, ok := altNames[playerName]; ok {
		playerName = val
	}
	players, err := db.GetPlayers(context.Background(), false, &model.PlayerFilter{Name: &playerName})
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	if len(players) == 0 {
		logrus.Errorf("could not find player %s", playerName)
		return 0, fmt.Errorf("could not find player %s", playerName)
	}
	if len(players) > 1 {
		logrus.Warnf("found multiple players with name %s", playerName)
	}
	return players[0].PlayerID, nil
}

func getTeamID(db BasketballRepository, opponentAbr string) (int, error) {
	teamfilter := &model.TeamFilter{Abbreviation: &opponentAbr}
	teams, err := db.GetTeams(context.Background(), false, teamfilter)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	if len(teams) == 0 {
		logrus.Errorf("could not find opponent %s", opponentAbr)
		return 0, fmt.Errorf("could not find opponent %s", opponentAbr)
	}
	if len(teams) > 1 {
		logrus.Warnf("found multiple teams with name %s", opponentAbr)
	}
	return teams[0].TeamID, nil
}

func getGameIDFromSchedule(schedule model.Schedule, startTime string, opponentAbr string) (string, error) {
	dateSlice := strings.SplitN(startTime, "T", 2)
	startdate := dateSlice[0]
	for _, gamedate := range schedule.LeagueSchedule.GameDates {
		split := strings.SplitN(gamedate.Games[0].GameDateEst, "T", 2)
		if split[0] == startdate {
			for _, game := range gamedate.Games {
				if game.GameID == "" {
					continue
				}
				if game.HomeTeam.TeamTriCode == opponentAbr || game.AwayTeam.TeamTriCode == opponentAbr {
					return game.GameID, nil
				}
			}
		}
	}
	return "", fmt.Errorf("could not find game vs %s on %v", opponentAbr, startTime)
}

func getTeamABRFromSchedule(schedule model.Schedule, startTime string, opponentAbr string) (string, error) {
	dateSlice := strings.SplitN(startTime, "T", 2)
	startdate := dateSlice[0]
	for _, gamedate := range schedule.LeagueSchedule.GameDates {
		split := strings.SplitN(gamedate.Games[0].GameDateEst, "T", 2)
		if split[0] == startdate {
			for _, game := range gamedate.Games {
				if game.GameID == "" {
					continue
				}
				if game.HomeTeam.TeamTriCode == opponentAbr {
					return game.AwayTeam.TeamTriCode, nil
				}
				if game.AwayTeam.TeamTriCode == opponentAbr {
					return game.HomeTeam.TeamTriCode, nil
				}
			}
		}
	}
	return "", fmt.Errorf("could not find game vs %s on %v", opponentAbr, startTime)
}

func getHomeTeam(schedule model.Schedule, startTime string, opponentAbr string) (model.HomeOrAway, error) {
	dateSlice := strings.SplitN(startTime, "T", 2)
	startdate := dateSlice[0]
	for _, gamedate := range schedule.LeagueSchedule.GameDates {
		split := strings.SplitN(gamedate.Games[0].GameDateEst, "T", 2)
		if split[0] == startdate {
			for _, game := range gamedate.Games {
				if game.GameID == "" {
					continue
				}
				if game.HomeTeam.TeamTriCode == opponentAbr {
					return model.HomeOrAwayAway, nil
				}
				if game.AwayTeam.TeamTriCode == opponentAbr {
					return model.HomeOrAwayHome, nil
				}
			}
		}
	}
	return "", fmt.Errorf("could not find game vs %s on %v", opponentAbr, startTime)
}
