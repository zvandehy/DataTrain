package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var instance *NBADatabaseClient

type NBADatabaseClient struct {
	*mongo.Client
	conn string
}

func ConnectDB(ctx context.Context) (*NBADatabaseClient, error) {
	var connErr error
	once.Do(func() {
		//TODO: Create config file for mongoDB access
		instance = &NBADatabaseClient{conn: "mongodb+srv://datatrain:nbawinners@datatrain.i5rgk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"}
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
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"abbreviation": bson.M{"$in": abbreviations},
	}
	cur, err := teamsDB.Find(ctx, filter)
	return cur, err
}

func (c *NBADatabaseClient) GetTeamsById(ctx context.Context, teamIDs []int) (*mongo.Cursor, error) {
	logrus.Printf("Query Teams By IDs: %v\n", teamIDs)
	teamsDB := c.Database("nba").Collection("teams")
	filter := bson.M{
		"teamID": bson.M{"$in": teamIDs},
	}
	cur, err := teamsDB.Find(ctx, filter)
	return cur, err
}

func (c *NBADatabaseClient) GetTeamGames(ctx context.Context, inputs []model.GameFilter) (*mongo.Cursor, error) {
	logrus.Printf("Query Team Games: %v\n", inputs)
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
	fmt.Println(teamIDs, seasons)
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
	playersDB := c.Database("nba").Collection("players")
	var playerIDs []int
	for _, in := range inputs {
		if in.PlayerID != nil {
			playerIDs = append(playerIDs, *in.PlayerID)
		}
	}
	filter := bson.M{
		"playerID": bson.M{"$in": playerIDs},
	}
	cur, err := playersDB.Find(ctx, filter)
	return cur, err
}
