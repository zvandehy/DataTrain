package database

import (
	"context"
	"fmt"
	"strings"
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
var instance *NBADatabaseClient

type NBADatabaseClient struct {
	*mongo.Client
	conn    string
	Queries int
}

func ConnectDB(ctx context.Context) (*NBADatabaseClient, error) {
	var connErr error
	once.Do(func() {
		config, err := util.LoadConfig(".")
		if err != nil {
			logrus.Fatal("cannot load configuration")
		}
		instance = &NBADatabaseClient{conn: config.DBSource}
		client, connErr := mongo.NewClient(options.Client().ApplyURI(instance.conn))
		if connErr != nil {
			return
		}
		connErr = client.Connect(ctx)
		if connErr != nil {
			return
		}
		//TODO: Should the initialized client be Database("nba")?
		instance.Client = client
	})
	if connErr != nil {
		return nil, connErr
	}
	logrus.Println("Connected to DB")
	return instance, nil
}

func (c *NBADatabaseClient) GetTeamsByAbr(ctx context.Context, abbreviations []string) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"abbreviation": bson.M{"$in": abbreviations},
	}
	cur, err := teamsDB.Find(ctx, filter)
	if len(abbreviations) < 5 {
		logrus.Printf("Query Teams By Abbreviations: %v\tTook %v", abbreviations, time.Since(start))
	} else {
		logrus.Printf("Query %d Teams By Abbreviation\tTook %v", len(abbreviations), time.Since(start))
	}
	return cur, err
}

func (c *NBADatabaseClient) GetTeamsById(ctx context.Context, teamIDs []int) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"teamID": bson.M{"$in": teamIDs},
	}
	cur, err := teamsDB.Find(ctx, filter)
	if len(teamIDs) < 5 {
		logrus.Printf("Query Teams By IDs: %v\tTook %v", teamIDs, time.Since(start))
	} else {
		logrus.Printf("Query %d Teams By IDs\tTook %v", len(teamIDs), time.Since(start))
	}
	return cur, err
}

func (c *NBADatabaseClient) GetTeamGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	teamGamesDB := c.Database("nba").Collection("teamgames")
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

	logrus.Printf("Query %d TeamGames\tTook %v", len(inputs), time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetPlayerGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	playerGamesDB := c.Database("nba").Collection("games")
	filter := createGameFilter(inputs)
	cur, err := playerGamesDB.Find(ctx, filter)
	if len(inputs) < 5 {
		logrus.Printf("Query PlayerGames From: %v\tTook %v", inputs, time.Since(start))
	} else {
		logrus.Printf("Query %d PlayerGames\tTook %v", len(inputs), time.Since(start))
	}
	return cur, err
}

func createGameFilter(inputs []model.GameFilter) bson.M {
	applyFilters := make(map[string]bson.M, 4)
	var gameIDs []string
	var teamIDs []int
	var playerIDs []int
	var seasons []string
	for _, in := range inputs {
		if in.GameID != nil {
			gameIDs = append(gameIDs, *in.GameID)
		}
		if in.PlayerID != nil {
			playerIDs = append(playerIDs, *in.PlayerID)
		}
		if in.TeamID != nil {
			teamIDs = append(teamIDs, *in.TeamID)
		}
		if in.Season != nil {
			seasons = append(seasons, *in.Season)
		}
	}
	if len(gameIDs) == 0 {
		applyFilters["gameID"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["gameID"] = bson.M{"$in": gameIDs}
	}
	if len(playerIDs) == 0 {
		applyFilters["playerID"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["playerID"] = bson.M{"$in": playerIDs}
	}
	if len(teamIDs) == 0 {
		applyFilters["teamID"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["teamID"] = bson.M{"$in": teamIDs}
	}
	if len(seasons) == 0 {
		applyFilters["season"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["season"] = bson.M{"$in": seasons}
	}
	//get all the players from this game
	filter := bson.M{
		"gameID": applyFilters["gameID"],
		"player": applyFilters["playerID"],
		"teamID": applyFilters["teamID"],
		"season": applyFilters["season"],
	}
	return filter
}

func (c *NBADatabaseClient) GetPlayers(ctx context.Context, inputs []model.PlayerFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	playersDB := c.Database("nba").Collection("players")
	var playerIDs []int
	var first_names []string
	var last_names []string
	for _, in := range inputs {
		if in.PlayerID != nil {
			playerIDs = append(playerIDs, *in.PlayerID)
		}
		if in.Name != nil {
			if len(strings.SplitN(*in.Name, " ", 2)) < 2 {
				return nil, fmt.Errorf("invalid Player Name: %s", *in.Name)
			}
			first_names = append(first_names, strings.SplitN(*in.Name, " ", 2)[0])
			last_names = append(last_names, strings.SplitN(*in.Name, " ", 2)[1])
		}
	}
	var filter bson.M
	if len(playerIDs) == 0 {
		if len(first_names) == 0 {
			filter = bson.M{}
		} else {
			filter = bson.M{"$and": bson.A{
				bson.M{"first_name": bson.M{"$in": first_names}},
				bson.M{"last_name": bson.M{"$in": last_names}},
			}}
		}
	} else {
		if len(first_names) == 0 {
			filter = bson.M{"playerID": bson.M{"$in": playerIDs}}
		} else {
			filter = bson.M{"$or": bson.A{
				bson.M{"playerID": bson.M{"$in": playerIDs}},
				bson.M{"$and": bson.A{
					bson.M{"first_name": bson.M{"$in": first_names}},
					bson.M{"last_name": bson.M{"$in": last_names}},
				},
				},
			}}
		}
	}

	cur, err := playersDB.Find(ctx, filter)
	logrus.Printf("Query %d Players\tTook %v", len(inputs), time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetProjections(ctx context.Context, input model.ProjectionFilter) (*mongo.Cursor, error) {
	start := time.Now()
	c.Queries++
	projectionDB := c.Database("nba").Collection("projections")
	filter := bson.M{"sportsbook": "PrizePicks"}

	if input.PlayerName != nil && *input.PlayerName != "" {
		filter["playername"] = *input.PlayerName
	}

	//filter date between input.StartDate and input.EndDate if they are set
	if input.StartDate != nil && input.EndDate != nil {
		filter["date"] = bson.M{"$gte": *input.StartDate, "$lte": *input.EndDate}
	} else if input.StartDate != nil {
		filter["date"] = bson.M{"$gte": *input.StartDate}
	} else if input.EndDate != nil {
		filter["date"] = bson.M{"$lte": *input.EndDate}
	}
	cur, err := projectionDB.Find(ctx, filter, options.Find().SetSort(bson.M{"date": 1}))
	logrus.Printf("Query Projections From: %v \tTook %v", input, time.Since(start))
	return cur, err
}

func (c *NBADatabaseClient) GetAverages(ctx context.Context, inputs []model.GameFilter) (*[]model.PlayerAverage, error) {
	start := time.Now()
	c.Queries++
	filter := createGameFilter(inputs)

	matchGames := bson.M{"$match": filter}
	averageStats := bson.M{"$group": bson.M{"_id": "$player",
		"games_played":             bson.M{"$count": bson.M{}},
		"points":                   bson.M{"$avg": "$points"},
		"assists":                  bson.M{"$avg": "$assists"},
		"rebounds":                 bson.M{"$avg": "$total_rebounds"},
		"steals":                   bson.M{"$avg": "$steals"},
		"blocks":                   bson.M{"$avg": "$blocks"},
		"turnovers":                bson.M{"$avg": "$turnovers"},
		"minutes":                  bson.M{"$avg": "$minutes"},
		"field_goals_made":         bson.M{"$avg": "$field_goals_made"},
		"field_goals_attempted":    bson.M{"$avg": "$field_goals_attempted"},
		"three_pointers_made":      bson.M{"$avg": "$three_pointers_made"},
		"three_pointers_attempted": bson.M{"$avg": "$three_pointers_attempted"},
		"free_throws_made":         bson.M{"$avg": "$free_throws_made"},
		"free_throws_attempted":    bson.M{"$avg": "$free_throws_attempted"},
	}}
	lookupPlayer := bson.M{"$lookup": bson.M{"from": "players", "localField": "_id", "foreignField": "playerID", "as": "player"}}
	unwindPlayer := bson.M{"$unwind": "$player"}

	cur, err := c.Database("nba").Collection("games").Aggregate(ctx, bson.A{matchGames, averageStats, lookupPlayer, unwindPlayer})
	if err != nil {
		return nil, err
	}
	var averages []model.PlayerAverage
	if err = cur.All(ctx, &averages); err != nil {
		return nil, err
	}
	logrus.Printf("Query %d Player Averages\tTook %v", len(inputs), time.Since(start))
	return &averages, nil

}

func (c *NBADatabaseClient) GetPlayerInjuries(ctx context.Context, playerIDs []int) ([]*model.Injury, error) {
	start := time.Now()
	c.Queries++
	cur, err := c.Database("nba").Collection("injuries").Find(ctx, bson.M{"playerID": bson.M{"$in": playerIDs}})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var injuries []*model.Injury
	cur.All(ctx, &injuries)
	logrus.Infof("Received Player Injuries for %d players\tTook %v", len(playerIDs), time.Since(start))
	return injuries, nil
}
