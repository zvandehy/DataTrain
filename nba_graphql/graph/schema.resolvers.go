package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *playerResolver) CurrentTeam(ctx context.Context, obj *model.Player) (*model.Team, error) {
	logrus.Printf("Get TEAM for player %v", obj)
	return dataloader.For(ctx).TeamByAbr.Load(obj.CurrentTeam)
}

func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	logrus.Println("Get Players")
	playersDB := r.Db.Database("nba").Collection("players")
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

func (r *queryResolver) FilterPlayers(ctx context.Context, input model.PlayerFilter) ([]*model.Player, error) {
	logrus.Printf("Get Players with filter  %#v", input)
	playersDB := r.Db.Database("nba").Collection("players")
	filter := bson.M{
		"last_name": input.Name,
	}
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

func (r *queryResolver) Player(ctx context.Context, input model.PlayerFilter) (*model.Player, error) {
	logrus.Printf("Get Player with filter  %#v", input)
	playersDB := r.Db.Database("nba").Collection("players")
	filter := bson.M{
		"playerID": input.PlayerID,
	}
	opts := options.FindOne().SetSort(bson.D{{"playerID", 1}})
	var player *model.Player
	err := playersDB.FindOne(ctx, filter, opts).Decode(&player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
	logrus.Println("Get Teams")
	teamsDB := r.Db.Database("nba").Collection("teams")
	filter := bson.M{}
	cur, err := teamsDB.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	teams := make([]*model.Team, 0, 10)
	for cur.Next(ctx) {
		team := &model.Team{}
		err := cur.Decode(&team)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *queryResolver) FilterTeams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
	logrus.Printf("Get Teams with filter %v\n", input)
	teamsDB := r.Db.Database("nba").Collection("teams")
	filter := bson.M{
		"teamID": input.TeamID,
		// "$or": []interface{}{
		// 	bson.M{"name": input.Name},
		// 	bson.M{"abbreviation": input.Abbreviation},
		// },
	}
	cur, err := teamsDB.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	teams := make([]*model.Team, 0, 10)
	for cur.Next(ctx) {
		team := &model.Team{}
		err := cur.Decode(&team)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *queryResolver) Team(ctx context.Context, input model.TeamFilter) (*model.Team, error) {
	logrus.Printf("Get Team with filter %#v\n", input)
	teamsDB := r.Db.Database("nba").Collection("teams")
	filter := bson.M{
		"abbreviation": input.Abbreviation,
		// "$or": []interface{}{
		// 	bson.M{"name": input.Name},
		// 	bson.M{"teamID": input.TeamID},
		// 	bson.M{"abbreviation": input.Abbreviation},
		// },
	}
	var team *model.Team
	opts := options.FindOne().SetSort(bson.D{{"abbreviation", 1}})
	err := teamsDB.FindOne(ctx, filter, opts).Decode(&team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *queryResolver) TeamGames(ctx context.Context, input model.GameFilter) ([]*model.TeamGame, error) {
	logrus.Printf("Get TeamGames with teamID %#v\n", input)
	cur, err := r.Db.GetTeamGames(ctx, []model.GameFilter{input})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	games := make([]*model.TeamGame, 0, 10)
	for cur.Next(ctx) {
		teamGame := &model.TeamGame{}
		err := cur.Decode(&teamGame)
		if err != nil {
			return nil, err
		}
		games = append(games, teamGame)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
	logrus.Printf("Get Team Games for opponent %v", obj)
	return dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
}

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// TeamGame returns generated.TeamGameResolver implementation.
func (r *Resolver) TeamGame() generated.TeamGameResolver { return &teamGameResolver{r} }

type playerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamGameResolver struct{ *Resolver }
