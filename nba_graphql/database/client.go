package database

import (
	"context"
	"fmt"
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
	// logrus.Printf("Query Players From: %v\n", inputs)
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

	return cur, err
}

func (c *NBADatabaseClient) GetProjections(ctx context.Context, input model.ProjectionFilter) (*mongo.Cursor, error) {
	logrus.Printf("Query Projections with: %v\n", input)
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

	return cur, err
}
