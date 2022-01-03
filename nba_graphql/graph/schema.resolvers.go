package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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
