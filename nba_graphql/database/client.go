package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var instance *mongo.Client

type NBADatabaseClient struct {
	Name string
	*mongo.Database
	conn    string
	Queries int
	Client  *mongo.Client
}

func ConnectDB(ctx context.Context, db string) (*NBADatabaseClient, error) {
	var connErr error
	config, err := util.LoadConfig(".")
	if err != nil {
		logrus.Fatal("cannot load configuration")
	}
	nbaClient := &NBADatabaseClient{conn: config.DBSource}
	once.Do(func() {
		client, connErr := mongo.NewClient(options.Client().ApplyURI(nbaClient.conn))
		if connErr != nil {
			logrus.Fatal("Can't create new mongo client: ", connErr)
		}
		connErr = client.Connect(ctx)
		if connErr != nil {
			logrus.Fatal("Can't connect to mongo client: ", connErr)
			return
		}
		instance = client
	})
	if connErr != nil {
		return nil, connErr
	}
	nbaClient.Name = db
	nbaClient.Client = instance
	nbaClient.Database = nbaClient.Client.Database(nbaClient.Name)
	logrus.Infof("Connected to DB: '%v/%v'", config.DBSource, nbaClient.Name)
	return nbaClient, nil
}

func (c *NBADatabaseClient) GetTeamsByAbr(ctx context.Context, abbreviations []string) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamsDB := c.Collection("teams")
	filter := bson.M{
		"abbreviation": bson.M{"$in": abbreviations},
	}
	cur, err := teamsDB.Find(ctx, filter)
	if len(abbreviations) < 5 {
		logrus.Printf("[%v] Query Teams By Abbreviations: %v\tTook %v", time.Now().Format(util.TIMENOW), abbreviations, time.Since(start))
	} else {
		logrus.Printf("[%v] Query %d Teams By Abbreviation\tTook %v", time.Now().Format(util.TIMENOW), len(abbreviations), time.Since(start))
	}
	return cur, err
}

func (c *NBADatabaseClient) GetTeamsById(ctx context.Context, teamIDs []int) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamsDB := c.Collection("teams")
	filter := bson.M{
		"teamID": bson.M{"$in": teamIDs},
	}
	cur, err := teamsDB.Find(ctx, filter)
	if len(teamIDs) < 5 {
		logrus.Printf("[%v] Query Teams By IDs: %v\tTook %v", time.Now().Format(util.TIMENOW), teamIDs, time.Since(start))
	} else {
		logrus.Printf("[%v] Query %d Teams By IDs\tTook %v", time.Now().Format(util.TIMENOW), len(teamIDs), time.Since(start))
	}
	return cur, err
}

func (c *NBADatabaseClient) GetTeamGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamGamesDB := c.Collection("teamgames")
	var teamIDs []int
	var seasons []string
	for _, in := range inputs {
		if in.TeamID != nil {
			teamIDs = append(teamIDs, *in.TeamID)
		}
		if in.Season != nil {
			seasons = append(seasons, *in.Season)
		}
	}
	//TODO: I think this isn't quite right
	filter := bson.M{
		"teamID": bson.M{"$in": teamIDs},
		"season": bson.M{"$in": seasons},
	}
	cur, err := teamGamesDB.Find(ctx, filter)

	logrus.Printf("[%v] Query %d TeamGames\tTook %v", time.Now().Format(util.TIMENOW), len(inputs), time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetPlayerGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	playerGamesDB := c.Collection("games")
	filter := createGameFilter(inputs)
	cur, err := playerGamesDB.Find(ctx, filter)
	if len(inputs) < 5 {
		logrus.Printf("[%v] Query PlayerGames From: %v\tTook %v", time.Now().Format(util.TIMENOW), inputs, time.Since(start))
	} else {
		logrus.Printf("[%v] Query %d PlayerGames\tTook %v", time.Now().Format(util.TIMENOW), len(inputs), time.Since(start))
	}
	return cur, err
}

func createGameFilter(inputs []model.GameFilter) bson.M {
	filters := bson.A{}
	for _, input := range inputs {
		inputAsFilter := bson.M{}
		if input.GameID != nil {
			inputAsFilter["gameID"] = *input.GameID
		}
		if input.TeamID != nil {
			inputAsFilter["teamID"] = *input.TeamID
		}
		if input.OpponentID != nil {
			inputAsFilter["opponent"] = *input.OpponentID
		}
		if input.PlayerID != nil {
			inputAsFilter["playerID"] = *input.PlayerID
		}
		if input.Season != nil {
			inputAsFilter["season"] = *input.Season
		}
		if input.GameID != nil {
			inputAsFilter["gameID"] = *input.GameID
		}
		if input.StartDate != nil && input.EndDate != nil {
			if *input.StartDate == *input.EndDate {
				inputAsFilter["date"] = *input.StartDate
			} else {
				inputAsFilter["date"] = bson.M{"$gte": *input.StartDate, "$lt": *input.EndDate}
			}
		} else if input.StartDate != nil {
			inputAsFilter["date"] = bson.M{"$gte": *input.StartDate}
		} else if input.EndDate != nil {
			inputAsFilter["date"] = bson.M{"$lt": *input.EndDate}
		}
		filters = append(filters, inputAsFilter)
	}
	return bson.M{"$or": filters}

	// applyFilters := make(map[string]bson.M, 4)
	// var gameIDs []string
	// var teamIDs []int
	// var opponentIDs []int
	// var playerIDs []int
	// var seasons []string
	// for _, in := range inputs {
	// 	if in.GameID != nil {
	// 		gameIDs = append(gameIDs, *in.GameID)
	// 	}
	// 	if in.PlayerID != nil {
	// 		playerIDs = append(playerIDs, *in.PlayerID)
	// 	}
	// 	if in.TeamID != nil {
	// 		teamIDs = append(teamIDs, *in.TeamID)
	// 	}
	// 	if in.OpponentID != nil {
	// 		opponentIDs = append(opponentIDs, *in.OpponentID)
	// 	}
	// 	if in.Season != nil {
	// 		seasons = append(seasons, *in.Season)
	// 	}

	// }
	// if len(gameIDs) == 0 {
	// 	applyFilters["gameID"] = bson.M{"$nin": []string{""}}
	// } else {
	// 	applyFilters["gameID"] = bson.M{"$in": gameIDs}
	// }
	// if len(playerIDs) == 0 {
	// 	applyFilters["playerID"] = bson.M{"$nin": []string{""}}
	// } else {
	// 	applyFilters["playerID"] = bson.M{"$in": playerIDs}
	// }
	// if len(teamIDs) == 0 {
	// 	applyFilters["teamID"] = bson.M{"$nin": []string{""}}
	// } else {
	// 	applyFilters["teamID"] = bson.M{"$in": teamIDs}
	// }
	// if len(opponentIDs) == 0 {
	// 	applyFilters["opponent"] = bson.M{"$nin": []string{""}}
	// } else {
	// 	applyFilters["opponent"] = bson.M{"$in": opponentIDs}
	// }
	// if len(seasons) == 0 {
	// 	applyFilters["season"] = bson.M{"$nin": []string{""}}
	// } else {
	// 	applyFilters["season"] = bson.M{"$in": seasons}
	// }
	// //get all the players from this game
	// filter := bson.M{
	// 	"gameID":   applyFilters["gameID"],
	// 	"player":   applyFilters["playerID"],
	// 	"teamID":   applyFilters["teamID"],
	// 	"season":   applyFilters["season"],
	// 	"opponent": applyFilters["opponent"],
	// }
	// return filter
}

func (c *NBADatabaseClient) GetPlayers(ctx context.Context, inputs []model.PlayerFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	playersDB := c.Collection("players")
	var playerIDs []int
	var names []string
	for _, in := range inputs {
		if in.PlayerID != nil {
			playerIDs = append(playerIDs, *in.PlayerID)
		}
		if in.Name != nil {
			names = append(names, *in.Name)
		}
	}
	var filter bson.M
	if len(playerIDs) == 0 {
		if len(names) == 0 {
			filter = bson.M{}
		} else {
			filter = bson.M{"name": bson.M{"$in": names}}
		}
	} else {
		if len(names) == 0 {
			filter = bson.M{"playerID": bson.M{"$in": playerIDs}}
		} else {
			filter = bson.M{"$or": bson.A{
				bson.M{"playerID": bson.M{"$in": playerIDs}},
				bson.M{"$and": bson.M{"name": bson.M{"$in": names}}},
			},
			}
		}
	}
	cur, err := playersDB.Find(ctx, filter)
	logrus.Printf("[%v] Query %d Players\tTook %v", time.Now().Format(util.TIMENOW), len(inputs), time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetProjections(ctx context.Context, input model.ProjectionFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	projectionDB := c.Collection("projections")
	filter := bson.M{}

	if input.PlayerName != nil && *input.PlayerName != "" {
		filter["playername"] = *input.PlayerName
	}

	//filter date between input.StartDate and input.EndDate if they are set
	if input.StartDate != nil && input.EndDate != nil {
		if *input.StartDate == *input.EndDate {
			filter["date"] = *input.StartDate
		} else {
			filter["date"] = bson.M{"$gte": *input.StartDate, "$lt": *input.EndDate}
		}
	} else if input.StartDate != nil {
		filter["date"] = bson.M{"$gte": *input.StartDate}
	} else if input.EndDate != nil {
		filter["date"] = bson.M{"$lt": *input.EndDate}
	}
	cur, err := projectionDB.Find(ctx, filter, options.Find().SetSort(bson.M{"date": 1}))
	logrus.Printf("[%v] Query Projections From: %v \tTook %v", time.Now().Format(util.TIMENOW), input, time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetAverages(ctx context.Context, inputs []model.GameFilter) (*[]model.PlayerAverage, error) {
	start := time.Now()
	c.Queries++
	filter := createGameFilter(inputs)
	matchGames := bson.M{"$match": filter}
	averageStats := bson.M{"$group": bson.M{"_id": "$playerID",
		"games_played":             bson.M{"$count": bson.M{}},
		"all_minutes":              bson.M{"$push": "$minutes"},
		"assists":                  bson.M{"$avg": "$assists"},
		"blocks":                   bson.M{"$avg": "$blocks"},
		"field_goals_attempted":    bson.M{"$avg": "$field_goals_attempted"},
		"field_goals_made":         bson.M{"$avg": "$field_goals_made"},
		"free_throws_attempted":    bson.M{"$avg": "$free_throws_attempted"},
		"free_throws_made":         bson.M{"$avg": "$free_throws_made"},
		"offensive_rebounds":       bson.M{"$avg": "$offensive_rebounds"},
		"defensive_rebounds":       bson.M{"$avg": "$defensive_rebounds"},
		"personal_fouls":           bson.M{"$avg": "$personal_fouls"},
		"personal_fouls_drawn":     bson.M{"$avg": "$personal_fouls_drawn"},
		"points":                   bson.M{"$avg": "$points"},
		"rebounds":                 bson.M{"$avg": "$total_rebounds"},
		"steals":                   bson.M{"$avg": "$steals"},
		"three_pointers_attempted": bson.M{"$avg": "$three_pointers_attempted"},
		"three_pointers_made":      bson.M{"$avg": "$three_pointers_made"},
		"turnovers":                bson.M{"$avg": "$turnovers"},
	}}
	lookupPlayer := bson.M{"$lookup": bson.M{"from": "players", "localField": "_id", "foreignField": "playerID", "as": "player"}}
	unwindPlayer := bson.M{"$unwind": "$player"}
	// must have played more than 1 game and average more than 5 minutes per game
	// and must have a height and weight
	matchPlayersWhoPlay := bson.M{"$match": bson.M{"games_played": bson.M{"$gt": 0}, "player.height": bson.M{"$ne": ""}, "player.weight": bson.M{"$gt": 0}}}

	cur, err := c.Collection("games").Aggregate(ctx, bson.A{matchGames, averageStats, lookupPlayer, unwindPlayer, matchPlayersWhoPlay})
	if err != nil {
		return nil, err
	}
	var averages []model.PlayerAverage
	if err = cur.All(ctx, &averages); err != nil {
		return nil, err
	}
	for i := 0; i < len(averages); i++ {
		avg := averages[i]
		minutes, err := avg.AverageMinutes()
		if err != nil {
			return nil, err
		}
		averages[i].Minutes = minutes
		if averages[i].Minutes < 5 {
			// remove
			averages = append(averages[:i], averages[i+1:]...)
			i--
		}
	}
	logrus.Printf("[%v] Query %d Player Averages\tTook %v\n", time.Now().Format(util.TIMENOW), len(inputs), time.Since(start))
	return &averages, nil

}

func (c *NBADatabaseClient) GetTeamAverages(ctx context.Context, inputs []model.GameFilter) (*[]model.TeamAverage, error) {
	start := time.Now()
	c.Queries++
	filter := createGameFilter(inputs)

	matchGames := bson.M{"$match": filter}
	averageStats := bson.M{"$group": bson.M{"_id": "$teamID",
		"wins_and_losses":      bson.M{"$push": "$win_or_loss"},
		"points":               bson.M{"$avg": "$points"},
		"opponent_points":      bson.M{"$avg": "$opponent_points"},
		"assists":              bson.M{"$avg": "$assists"},
		"opponent_assists":     bson.M{"$avg": "$opponent_assists"},
		"rebounds":             bson.M{"$avg": "$rebounds"},
		"opponent_rebounds":    bson.M{"$avg": "$opponent_rebounds"},
		"steals":               bson.M{"$avg": "$steals"},
		"blocks":               bson.M{"$avg": "$blocks"},
		"turnovers":            bson.M{"$avg": "$turnovers"},
		"three_pointers_made":  bson.M{"$avg": "$three_pointers_made"},
		"personal_fouls_drawn": bson.M{"$avg": "$personal_fouls_drawn"},
		"personal_fouls":       bson.M{"$avg": "$personal_fouls"},
	}}
	lookupTeam := bson.M{"$lookup": bson.M{"from": "teams", "localField": "_id", "foreignField": "teamID", "as": "team"}}
	unwindTeam := bson.M{"$unwind": "$team"}

	cur, err := c.Collection("teamgames").Aggregate(ctx, bson.A{matchGames, averageStats, lookupTeam, unwindTeam})
	if err != nil {
		return nil, err
	}
	var averages []model.TeamAverage
	if err = cur.All(ctx, &averages); err != nil {
		return nil, err
	}
	//count wins and losses
	for i := 0; i < len(averages); i++ {
		avg := averages[i]
		wins := 0.0
		losses := 0.0
		for _, a := range avg.WinsAndLosses {
			if a == "win" {
				wins++
			} else if a == "loss" {
				losses++
			}
		}
		averages[i].GamesWon = wins
		averages[i].GamesLost = losses
	}

	logrus.Printf("[%v] Query %d Team Averages\tTook %v\n", time.Now().Format(util.TIMENOW), len(inputs), time.Since(start))
	return &averages, nil

}

func (c *NBADatabaseClient) GetPlayerInjuries(ctx context.Context, playerIDs []int) ([]*model.Injury, error) {
	start := time.Now()
	c.Queries++
	cur, err := c.Collection("injuries").Find(ctx, bson.M{"playerID": bson.M{"$in": playerIDs}})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var injuries []*model.Injury
	cur.All(ctx, &injuries)
	logrus.Infof("Received Player Injuries for %d players\tTook %v", len(playerIDs), time.Since(start))
	return injuries, nil
}

func (c *NBADatabaseClient) GetTeamInjuries(ctx context.Context, teamIDs []int) ([]*model.Injury, error) {
	start := time.Now()
	c.Queries++

	lookupPlayers := bson.M{"$lookup": bson.M{
		"from":         "players",
		"localField":   "playerID",
		"foreignField": "playerID",
		"as":           "player",
	}}
	unwindPlayer := bson.M{"$unwind": bson.M{
		"path":                       "$player",
		"preserveNullAndEmptyArrays": true,
	}}
	lookupTeam := bson.M{"$lookup": bson.M{
		"from":         "teams",
		"localField":   "player.teamABR",
		"foreignField": "abbreviation",
		"as":           "team",
	}}
	unwindTeam := bson.M{"$unwind": bson.M{
		"path":                       "$team",
		"preserveNullAndEmptyArrays": true,
	}}
	matchTeam := bson.M{"$match": bson.M{"team.teamID": bson.M{"$in": teamIDs}}}

	cur, err := c.Collection("injuries").Aggregate(ctx, bson.A{lookupPlayers, unwindPlayer, lookupTeam, unwindTeam, matchTeam})
	if err != nil {
		logrus.Warnf("error getting injury data for teams: %v", err)
		return nil, err
	}
	var injuries []*model.Injury
	err = cur.All(ctx, &injuries)
	if err != nil {
		logrus.Warnf("error getting injury data for teams: %v", err)
		return nil, err
	}
	logrus.Infof("Received Injuries for %d teams\tTook %v", len(teamIDs), time.Since(start))
	return injuries, nil

}
