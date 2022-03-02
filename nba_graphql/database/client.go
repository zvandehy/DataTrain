package database

import (
	"context"
	"strings"
	"sync"

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
	logrus.Printf("Query Teams By Abbreviations: %v\n", abbreviations)
	c.Queries++
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"abbreviation": bson.M{"$in": abbreviations},
	}
	cur, err := teamsDB.Find(ctx, filter)

	return cur, err
}

func (c *NBADatabaseClient) GetTeamsById(ctx context.Context, teamIDs []int) (*mongo.Cursor, error) {
	logrus.Printf("Query Teams By IDs: %v\n", teamIDs)
	c.Queries++
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"teamID": bson.M{"$in": teamIDs},
	}
	cur, err := teamsDB.Find(ctx, filter)

	return cur, err
}

func (c *NBADatabaseClient) GetTeamGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	logrus.Printf("Query Team Games: %v\n", inputs)
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

	return cur, err
}

func (c *NBADatabaseClient) GetPlayerGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	logrus.Printf("Query PlayerGames From: %v\n", inputs)
	c.Queries++
	playerGamesDB := c.Database("nba").Collection("games")
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
	cur, err := playerGamesDB.Find(ctx, filter)

	return cur, err
}

func (c *NBADatabaseClient) GetPlayers(ctx context.Context, inputs []model.PlayerFilter) (*mongo.Cursor, error) {
	logrus.Printf("Query Players From: %v\n", inputs)
	c.Queries++
	playersDB := c.Database("nba").Collection("players")
	applyFilters := make(map[string]bson.M, 4)
	var playerIDs []int
	var first_names []string
	var last_names []string
	for _, in := range inputs {
		if in.PlayerID != nil {
			playerIDs = append(playerIDs, *in.PlayerID)
		}
		if in.Name != nil {
			first_names = append(first_names, strings.SplitN(*in.Name, " ", 2)[0])
			last_names = append(last_names, strings.SplitN(*in.Name, " ", 2)[1])
		}
	}
	if len(playerIDs) == 0 {
		applyFilters["playerID"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["playerID"] = bson.M{"$in": playerIDs}
	}
	if len(first_names) == 0 {
		applyFilters["first_name"] = bson.M{"$nin": []string{""}}
		applyFilters["last_name"] = bson.M{"$nin": []string{""}}
	} else {
		applyFilters["first_name"] = bson.M{"$in": first_names}
		applyFilters["last_name"] = bson.M{"$in": last_names}
	}
	//get all the players from this game
	filter := bson.M{"$or": bson.A{
		bson.M{"player": applyFilters["playerID"]},
		bson.M{"$and": bson.A{bson.M{"first_name": applyFilters["first_name"]},
			bson.M{"last_name": applyFilters["last_name"]},
		}},
	}}
	cur, err := playersDB.Find(ctx, filter)

	return cur, err
}
