package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/zvandehy/nba_graphql/database"
	"github.com/zvandehy/nba_graphql/graph/generated"
	"github.com/zvandehy/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Games(ctx context.Context, playerID int) ([]*model.PlayerGame, error) {
	client, err := database.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database("nba").Collection("games")
	filter := bson.M{"player": playerID}
	cur, err := db.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	games := make([]*model.PlayerGame, 0, 100)
	for cur.Next(ctx) {
		game := &model.PlayerGame{}
		err := cur.Decode(&game)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *queryResolver) LastNGames(ctx context.Context, playerID int, n int) ([]*model.PlayerGame, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TeamGames(ctx context.Context, playerID int) ([]*model.TeamGame, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Prizepicks(ctx context.Context) ([]*model.PlayerProp, error) {
	url := "https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=7"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var prizepicks model.PrizePicks
	if err := json.Unmarshal(bytes, &prizepicks); err != nil {
		return nil, err
	}
	playerProps := make([]*model.PlayerProp, len(prizepicks.Data))
	for i, prop := range prizepicks.Data {
		var playerName string
		var statType string
		for _, p := range prizepicks.Included {
			if p.ID == prop.Relationships.Player.Data.ID {
				playerName = p.Attributes.Name
			}
			if p.ID == prop.Relationships.StatType.Data.ID {
				statType = p.Attributes.Name
			}
			if statType != "" && playerName != "" {
				break
			}
		}
		if playerName == "" {
			return nil, fmt.Errorf("error retrieving prizepick player name")
		}
		if statType == "" {
			return nil, fmt.Errorf("error retrieving prizepick stat type")
		}
		target, err := strconv.ParseFloat(prop.Attributes.Line_score, 64)
		if err != nil {
			return nil, err
		}
		playerProps[i] = &model.PlayerProp{Target: target, Type: statType, Opponent: prop.Attributes.Description, PlayerName: playerName}
	}
	return playerProps, nil
}

func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	client, err := database.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	playersDB := client.Database("nba").Collection("players")
	filter := bson.M{}
	cur, err := playersDB.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	players := make([]*model.Player, 0, 10)
	for cur.Next(ctx) {
		player := &model.Player{}
		err := cur.Decode(&player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
