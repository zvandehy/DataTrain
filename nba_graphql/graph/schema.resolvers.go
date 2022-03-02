package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *playerResolver) Name(ctx context.Context, obj *model.Player) (string, error) {
	return obj.FirstName + " " + obj.LastName, nil
}

func (r *playerResolver) CurrentTeam(ctx context.Context, obj *model.Player) (*model.Team, error) {
	//logrus.Printf("Get TEAM for player %v", obj)
	t, err := dataloader.For(ctx).TeamByAbr.Load(obj.CurrentTeam)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return &model.Team{}, nil
	}
	return t, err
}

func (r *playerResolver) Games(ctx context.Context, obj *model.Player, input model.GameFilter) ([]*model.PlayerGame, error) {
	//logrus.Printf("Get Games filtered by %v for Player %v", input, obj)
	if obj.PlayerID == 0 {
		return []*model.PlayerGame{}, nil
	}
	input.PlayerID = &obj.PlayerID
	return dataloader.For(ctx).PlayerGameByFilter.Load(input)
}

func (r *playerGameResolver) Opponent(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	//logrus.Printf("Get Opponent from PlayerGame %v", obj)
	return dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
}

func (r *playerGameResolver) Player(ctx context.Context, obj *model.PlayerGame) (*model.Player, error) {
	//logrus.Printf("Get Player from PlayerGame %v", obj)
	return dataloader.For(ctx).PlayerByID.Load(obj.PlayerID)
}

func (r *playerGameResolver) PlayersInGame(ctx context.Context, obj *model.PlayerGame) (*model.PlayersInGame, error) {
	//logrus.Printf("Get PlayersInGame from PlayerGame %v", obj)
	gameCur, err := r.Db.GetPlayerGames(ctx, []model.GameFilter{{GameID: &obj.GameID}})
	if err != nil {
		return nil, err
	}
	defer gameCur.Close(ctx)
	var playerGames map[int]model.PlayerGame = make(map[int]model.PlayerGame)
	var playerFilters []model.PlayerFilter
	for gameCur.Next(ctx) {
		playerGame := model.PlayerGame{}
		err = gameCur.Decode(&playerGame)
		if err != nil {
			return nil, err
		}
		playerGames[playerGame.PlayerID] = playerGame
		playerFilters = append(playerFilters, model.PlayerFilter{PlayerID: &playerGame.PlayerID})
	}
	playerCur, err := r.Db.GetPlayers(ctx, playerFilters)
	if err != nil {
		return nil, err
	}
	defer playerCur.Close(ctx)
	var teamPlayers []*model.Player
	var oppPlayers []*model.Player
	for playerCur.Next(ctx) {
		player := &model.Player{}
		err := playerCur.Decode(&player)
		if err != nil {
			return nil, err
		}
		//team's opponent == player's opponent --> player is teammate
		if playerGames[player.PlayerID].OpponentID == obj.OpponentID {
			teamPlayers = append(teamPlayers, player)
		} else {
			oppPlayers = append(oppPlayers, player)
		}
	}
	if err := playerCur.Err(); err != nil {
		return nil, err
	}
	return &model.PlayersInGame{TeamPlayers: teamPlayers, OpponentPlayers: oppPlayers}, nil
}

func (r *playersInGameResolver) Team(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
	return obj.TeamPlayers, nil
}

func (r *playersInGameResolver) Opponent(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
	return obj.OpponentPlayers, nil
}

func (r *projectionResolver) Player(ctx context.Context, obj *model.Projection) (*model.Player, error) {
	//logrus.Printf("Get Player from Projection %v", obj)
	playerFilter := model.PlayerFilter{Name: &obj.PlayerName}
	p, err := dataloader.For(ctx).PlayerByFilter.Load(playerFilter)
	if p == nil {
		logrus.Warnf("Player %v is nil. Probably needs to be uploaded to the database.", *playerFilter.Name)
		name := strings.SplitN(*playerFilter.Name, " ", 2)
		return &model.Player{FirstName: name[0], LastName: name[1]}, nil
	}
	return p, err
}

func (r *projectionResolver) Opponent(ctx context.Context, obj *model.Projection) (*model.Team, error) {
	//logrus.Printf("Get TEAM for projection %v", obj)
	return dataloader.For(ctx).TeamByAbr.Load(obj.OpponentAbr)
}

func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	//logrus.Println("Get Players")
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
	//logrus.Printf("Get Players with filter  %v", input)
	cur, err := r.Db.GetPlayers(ctx, []model.PlayerFilter{input})
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
	//logrus.Printf("Get Player with filter  %v", input)
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
	//logrus.Println("Get Teams")
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
	//logrus.Printf("Get Teams with filter %v\n", input)
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
	//logrus.Printf("Get Team with filter %#v\n", input)
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
	//logrus.Printf("Get TeamGames with teamID %#v\n", input)
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

func (r *queryResolver) PlayerGames(ctx context.Context, input model.GameFilter) ([]*model.PlayerGame, error) {
	//logrus.Printf("Get Games filtered by %v", input)
	cur, err := r.Db.GetPlayerGames(ctx, []model.GameFilter{input})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var playerGames []*model.PlayerGame
	for cur.Next(ctx) {
		playerGame := model.PlayerGame{}
		err = cur.Decode(&playerGame)
		if err != nil {
			return nil, err
		}
		playerGames = append(playerGames, &playerGame)
	}
	return playerGames, nil
}

func (r *queryResolver) Projections(ctx context.Context, sportsbook string) ([]*model.Projection, error) {
	var projections []*model.Projection
	if strings.ToLower(sportsbook) != "prizepicks" {
		return nil, fmt.Errorf("unsupported Sportsbook: %s. Current support only exists for: %v", sportsbook, []string{"PrizePicks"})
	}
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
	for _, prop := range prizepicks.Data {
		projections, err = model.ParsePrizePick(prop, prizepicks.Included, projections)
		if err != nil {
			return nil, fmt.Errorf("failed to parse prizepick projection: %v", err)
		}
	}
	return projections, nil
}

func (r *teamResolver) Games(ctx context.Context, obj *model.Team, input model.GameFilter) ([]*model.TeamGame, error) {
	//logrus.Printf("Get Games from team %v filtered by %v", obj, input)
	input.TeamID = &obj.TeamID
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

func (r *teamResolver) Players(ctx context.Context, obj *model.Team) ([]*model.Player, error) {
	//logrus.Printf("Get Players from Team %v", obj)
	input := model.PlayerFilter{TeamID: &obj.TeamID}
	return r.Query().FilterPlayers(ctx, input)
}

func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
	//logrus.Printf("Get Opponent from TeamGame %v", obj)
	return dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
}

func (r *teamGameResolver) PlayersInGame(ctx context.Context, obj *model.TeamGame) (*model.PlayersInGame, error) {
	//logrus.Printf("Get Players from Game %v", obj)
	//TODO: Abstract this with PlayerGamesResolver PlayersInGame()
	gameCur, err := r.Db.GetPlayerGames(ctx, []model.GameFilter{{GameID: &obj.GameID}})
	if err != nil {
		return nil, err
	}
	defer gameCur.Close(ctx)
	var playerGames map[int]model.PlayerGame = make(map[int]model.PlayerGame)
	var playerFilters []model.PlayerFilter
	for gameCur.Next(ctx) {
		playerGame := model.PlayerGame{}
		err = gameCur.Decode(&playerGame)
		if err != nil {
			return nil, err
		}
		playerGames[playerGame.PlayerID] = playerGame
		playerFilters = append(playerFilters, model.PlayerFilter{PlayerID: &playerGame.PlayerID})
	}
	playerCur, err := r.Db.GetPlayers(ctx, playerFilters)
	if err != nil {
		return nil, err
	}
	defer playerCur.Close(ctx)
	var teamPlayers []*model.Player
	var oppPlayers []*model.Player
	for playerCur.Next(ctx) {
		player := &model.Player{}
		err := playerCur.Decode(&player)
		if err != nil {
			return nil, err
		}
		//team's opponent == player's opponent --> player is teammate
		if playerGames[player.PlayerID].OpponentID == obj.OpponentID {
			teamPlayers = append(teamPlayers, player)
		} else {
			oppPlayers = append(oppPlayers, player)
		}
	}
	if err := playerCur.Err(); err != nil {
		return nil, err
	}
	return &model.PlayersInGame{TeamPlayers: teamPlayers, OpponentPlayers: oppPlayers}, nil
}

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// PlayerGame returns generated.PlayerGameResolver implementation.
func (r *Resolver) PlayerGame() generated.PlayerGameResolver { return &playerGameResolver{r} }

// PlayersInGame returns generated.PlayersInGameResolver implementation.
func (r *Resolver) PlayersInGame() generated.PlayersInGameResolver { return &playersInGameResolver{r} }

// Projection returns generated.ProjectionResolver implementation.
func (r *Resolver) Projection() generated.ProjectionResolver { return &projectionResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// TeamGame returns generated.TeamGameResolver implementation.
func (r *Resolver) TeamGame() generated.TeamGameResolver { return &teamGameResolver{r} }

type playerResolver struct{ *Resolver }
type playerGameResolver struct{ *Resolver }
type playersInGameResolver struct{ *Resolver }
type projectionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
type teamGameResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *projectionResolver) PropType(ctx context.Context, obj *model.Projection) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
