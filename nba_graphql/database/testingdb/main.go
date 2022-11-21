package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

// need list of games for a specific player

type Row struct {
	FirstName    string `db:"firstName"`
	LastName     string `db:"lastName"`
	PlayerName   string `db:"playername"`
	PlayerID     int    `db:"playerID"`
	Position     string `db:"position"`
	PlayerTeamID int    `db:"playerTeamID"`
	Height       string `db:"height"`
	Weight       int    `db:"weight"`

	League string `db:"league"`

	TeamID           int    `db:"teamID"`
	TeamName         string `db:"teamName"`
	TeamCity         string `db:"teamCity"`
	TeamAbbreviation string `db:"teamAbbreviation"`

	OpponentID           int    `db:"opponentID"`
	OpponentName         string `db:"opponentName"`
	OpponentCity         string `db:"opponentCity"`
	OpponentAbbreviation string `db:"opponentAbbreviation"`

	AssistPercentage             sql.NullFloat64 `db:"assistPct"`
	Assists                      sql.NullInt64   `db:"assists"`
	Date                         *time.Time      `db:"date"`
	DefensiveReboundPercentage   sql.NullFloat64 `db:"defensiveReboundPct"`
	DefensiveRebounds            sql.NullInt64   `db:"defensiveRebounds"`
	EffectiveFieldGoalPercentage sql.NullFloat64 `db:"effectiveFieldGoalPct"`
	FieldGoalPercentage          sql.NullFloat64 `db:"fieldGoalPct"`
	FieldGoalsAttempted          sql.NullInt64   `db:"fieldGoalsAttempted"`
	FieldGoalsMade               sql.NullInt64   `db:"fieldGoalsMade"`
	FreeThrowsAttempted          sql.NullInt64   `db:"freeThrowsAttempted"`
	FreeThrowsMade               sql.NullInt64   `db:"freeThrowsMade"`
	FreeThrowPercentage          sql.NullFloat64 `db:"freeThrowPct"`
	GameID                       string          `db:"gameID"`
	HomeOrAway                   string          `db:"homeAway"`
	Margin                       sql.NullInt64   `db:"margin"`
	Minutes                      float64         `db:"minutes"`
	OffensiveReboundPercentage   sql.NullFloat64 `db:"offensiveReboundPct"`
	OffensiveRebounds            sql.NullInt64   `db:"offensiveRebounds"`
	PersonalFoulsDrawn           sql.NullInt64   `db:"personalFoulsDrawn"`
	PersonalFouls                sql.NullInt64   `db:"personalFouls"`
	Points                       sql.NullInt64   `db:"points"`
	Playoffs                     bool            `db:"playoffs"`
	Season                       string          `db:"season"`
	ThreePointPercentage         sql.NullFloat64 `db:"threePointPct"`
	ThreePointersAttempted       sql.NullInt64   `db:"threePointersAttempted"`
	ThreePointersMade            sql.NullInt64   `db:"threePointersMade"`
	Rebounds                     sql.NullInt64   `db:"rebounds"`
	TrueShootingPercentage       sql.NullFloat64 `db:"trueShootingPct"`
	Turnovers                    sql.NullInt64   `db:"turnovers"`
	Blocks                       sql.NullInt64   `db:"blocks"`
	Steals                       sql.NullInt64   `db:"steals"`
	Usage                        sql.NullFloat64 `db:"usage"`
	Outcome                      string          `db:"outcome"`

	CreatedAt *time.Time `db:"CreatedAt"`
	UpdatedAt *time.Time `db:"UpdatedAt"`
}

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

func ParsePrizePickProposition(db database.BasketballRepository, schedule model.Schedule, prop model.PrizePicksData, itemIDToNameMap map[string]string) (proposition *model.DBProposition, err error) {
	if prop.Attributes.Is_promo {
		logrus.Warn("skipping promo prizepick %+v", prop)
		return nil, fmt.Errorf("skipping promo prizepick")
	}
	playerName, err := getPlayerName(prop, itemIDToNameMap)
	if err != nil {
		return nil, err
	}
	statType, err := getStatType(prop, itemIDToNameMap)
	if err != nil {
		return nil, err
	}
	target, err := getTarget(prop)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(time.RFC3339, prop.Attributes.Start_time)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve prizepicks date")
	}
	playerID, err := getPlayerID(db, playerName)
	if err != nil {
		return nil, err
	}
	opponentID, err := getTeamID(db, prop.Attributes.Description)
	if err != nil {
		return nil, err
	}

	gameID, err := getGameIDFromSchedule(schedule, prop.Attributes.Start_time, prop.Attributes.Description)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	teamABR, err := getTeamABRFromSchedule(schedule, prop.Attributes.Start_time, prop.Attributes.Description)
	if err != nil {
		return nil, err
	}
	teamID, err := getTeamID(db, teamABR)
	if err != nil {
		return nil, err
	}
	// save upcoming game in db
	game := &model.PlayerGame{
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
	_, err = db.SaveUpcomingGames(context.Background(), []*model.PlayerGame{game})
	if err != nil {
		logrus.Errorf("failed to save upcoming game %v", game)
	}
	// logrus.Infof("saved upcoming game %v %v", game.GameID, game.PlayerID)
	return proposition, nil
}

// func getPropStatus(foundProps []Proposition, sportsbook string, statType string, target float64, playerID int, gameID int) string {
// 	for _, prop := range foundProps {
// 		if prop.Sportsbook == sportsbook && prop.StatType == statType && prop.Target == target && prop.PlayerID == playerID && prop.GameID == gameID {
// 			return prop.Status
// 		}
// 	}
// 	return ""

// }

func getPlayerID(db database.BasketballRepository, playerName string) (int, error) {
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

func getTeamID(db database.BasketballRepository, opponentAbr string) (int, error) {
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

func main() {
	db, err := database.NewSQLClient("NBA")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PlanetScale!")

	mongodb, err := database.ConnectDB(context.Background(), "nba")
	if err != nil {
		panic(err)
	}

	cur, err := mongodb.Collection("projections").Find(context.Background(), bson.M{"propositions.propType": bson.M{"$exists": 1}})
	if err != nil {
		panic(err)
	}
	defer cur.Close(context.Background())

	// countskipstarttime := 0
	countskipproposiitons := 0
	propositions := make([]*model.Proposition, 0)
	gameInputs := []model.GameFilter{}
	playerInputs := []*model.PlayerFilter{}
	for cur.Next(context.Background()) {

		var projection struct {
			Date         string `bson:"date"`
			PlayerName   string `bson:"playername"`
			OpponentAbr  string `bson:"opponent"`
			Propositions []struct {
				Sportsbook     string     `bson:"sportsbook"`
				Target         float64    `bson:"target"`
				Type           string     `bson:"propType"`
				LastModifiedAt *time.Time `bson:"lastModified"`
			} `bson:"propositions"`
			StartTime *time.Time `bson:"startTime"`
		}
		err := cur.Decode(&projection)
		if err != nil {
			panic(err)
		}
		if projection.Propositions == nil || len(projection.Propositions) == 0 {
			countskipproposiitons++
			continue
		}

		for _, proposition := range projection.Propositions {
			// if projection.PlayerName contains "starters"
			if strings.Contains(strings.ToLower(projection.PlayerName), "starters") {
				fmt.Println("skipping starters prop")
				continue
			}
			projection.PlayerName = strings.TrimSpace(projection.PlayerName)
			name, ok := altNames[projection.PlayerName]
			if !ok {
				name = projection.PlayerName
			} else {
				projection.PlayerName = name
			}
			propositions = append(propositions, &model.Proposition{
				PlayerID:     0,
				GameID:       "",
				OpponentID:   0,
				PlayerName:   projection.PlayerName,
				TypeRaw:      proposition.Type,
				Target:       proposition.Target,
				Sportsbook:   model.SportsbookOption(proposition.Sportsbook),
				LastModified: proposition.LastModifiedAt,
			})

			gameInputs = append(gameInputs, model.GameFilter{
				StartDate: &projection.Date,
				EndDate:   &projection.Date,
			})
			playerInputs = append(playerInputs, &model.PlayerFilter{
				Name: &name,
			})
		}
	}

	fmt.Println(len(propositions), " received from mongo")

	games, err := db.GetPlayerGames(context.Background(), gameInputs...)
	if err != nil {
		panic(err)
	}
	players, err := db.GetPlayers(context.Background(), false, playerInputs...)
	if err != nil {
		panic(err)
	}

	parsedProps := make([]*model.Proposition, 0)
	for i, proposition := range propositions {
		player, err := findPlayer(players, proposition.PlayerName)
		if err != nil {
			fmt.Printf("%v skipping %s on %s\n", err, proposition.PlayerName, *gameInputs[i].EndDate)
			continue
		}
		proposition.PlayerID = player.PlayerID

		game, err := findGame(games, proposition.PlayerID, *gameInputs[i].EndDate)
		if err != nil {
			fmt.Printf("%v skipping %s on %s\n", err, proposition.PlayerName, *gameInputs[i].EndDate)
			continue
		} else {
			proposition.GameID = game.GameID
			proposition.OpponentID = game.OpponentID
		}
		parsedProps = append(parsedProps, proposition)
	}
	fmt.Println("FOUND PROPS ", len(parsedProps))
	saveProps := make([]*model.Proposition, 50)
	for i, proposition := range parsedProps {
		saveProps[i%50] = proposition
		if i%50 == 49 {
			x, err := db.SavePropositions(context.Background(), saveProps)
			if err != nil {
				panic(err)
			}
			fmt.Println("saved ", x)
			saveProps = make([]*model.Proposition, 50)
			i = -1
		}
	}
}

func findGame(games []*model.PlayerGame, playerID int, date string) (model.PlayerGame, error) {
	for _, game := range games {
		gameDate := game.Date.Format("2006-01-02")
		if game.PlayerID == playerID && gameDate == date {
			return *game, nil
		}
	}
	return model.PlayerGame{}, errors.New("game not found")
}

func findPlayer(players []*model.Player, name string) (model.Player, error) {
	for _, player := range players {
		if player.Name == name {
			return *player, nil
		}
	}

	return model.Player{}, fmt.Errorf("player not found: %v", name)
}

func getgames(db *database.SQLClient) {
	rows := []Row{}
	if err := db.Select(&rows, `SELECT
		p.name as playername,
		playerID,
		firstName,
		lastName,
		position,
		p.teamID as playerTeamID,
		height,
		weight,
		p.league,

		t.teamID as teamID,
		t.name as teamName,
		t.city as teamCity,
		t.abbreviation as teamAbbreviation,

		pg.opponentID,
		opp.name as opponentName,
		opp.city as opponentCity,
		opp.abbreviation as opponentAbbreviation,

		assistPct,
		assists,
		date,
		defensiveReboundPct,
		defensiveRebounds,
		effectiveFieldGoalPct,
		fieldGoalPct,
		fieldGoalsAttempted,
		fieldGoalsMade,
		freeThrowsAttempted,
		freeThrowsMade,
		freeThrowPct,
		gameID,
		homeAway,
		margin,
		minutes,
		offensiveReboundPct,
		offensiveRebounds,
		personalFoulsDrawn,
		personalFouls,
		points,
		playoffs,
		pg.season,
		threePointPct,
		threePointersAttempted,
		threePointersMade,
		rebounds,
		trueShootingPct,
		turnovers,
		blocks,
		steals,
		usage,
		pg.outcome

		FROM players p
		JOIN playergames pg USING (playerID)
		JOIN teams t ON pg.teamID = t.teamID
		JOIN teams opp ON pg.opponentID = opp.teamID;
		`); err != nil {
		panic(err)
	}
	fmt.Println("Successfully queried players!")

	playerGames := make(map[int][]Row)
	for _, row := range rows {
		playerGames[row.PlayerID] = append(playerGames[row.PlayerID], row)
	}
	// fmt.Printf("? has %d games\n", rows[0].PlayerName, len(playerGames[rows[0].PlayerID]))
	// for _, game := range playerGames[rows[0].PlayerID] {
	// 	fmt.Println(game)
	// }

	// for _, row := range rows {
	// 	//reflect print each field and value
	// 	fmt.Printf("%+v\n\n", row)
	// 	// fmt.Println(player.Name, player.PlayerID, player.CurrentTeam, player.TeamID, player.Position, player.Weight, player.Height)
	// }

	fmt.Println(len(rows), " TOTAL GAMES")
	fmt.Println(len(playerGames), " PLAYERS")

	// q := "SELECT players.name, playergames.points, playergames.gameID FROM players JOIN playergames ON players.playerID = playergames.playerID WHERE playergames.minutes > 15"
	// rows, err := db.Client.Query(q)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var name string
	// 	var points int
	// 	var gameID int
	// 	if err := rows.Scan(&name, &points, &gameID); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("? scored %d points in game %d\n", name, points, gameID)
	// }
	// if err := rows.Err(); err != nil {
	// 	panic(err)
	// }
}

func getSchedule() *model.Schedule {
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
	return &schedule
}
