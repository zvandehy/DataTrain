package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zvandehy/nba_graphql/database"
	"github.com/zvandehy/nba_graphql/graph/generated"
	"github.com/zvandehy/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

//TODO: Need to revisit GraphQL, data models, and what the best design of these endpoints would be
func (r *queryResolver) Games(ctx context.Context, playerIDs []int) ([]*model.PlayerGame, error) {
	client, err := database.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database("nba").Collection("games")

	filter := bson.M{"player": bson.M{"$in": playerIDs}}

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
	client, err := database.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database("nba").Collection("teams")
	filter := bson.M{}
	cur, err := db.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	teams := make([]*model.Team, 0, 30)
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

func (r *queryResolver) Prizepicks(ctx context.Context) ([]*model.PlayerProp, error) {
	//TODO: Optimize performance by not overfetching the same player's PlayerGames
	//TODO: Optimize performance with concurrency
	//TODO: Optimize performance with cacheing
	fmt.Println("----------------------")
	t := time.Now()
	url := "https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=7"
	res, err := http.Get(url)
	fmt.Printf("Prizepicks REQ: %v\n", time.Since(t))
	if err != nil {
		return nil, err
	}
	t = time.Now()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var prizepicks model.PrizePicks
	if err := json.Unmarshal(bytes, &prizepicks); err != nil {
		return nil, err
	}

	playerNames := make(map[string][]string)
	for _, prop := range prizepicks.Data {
		for _, p := range prizepicks.Included {
			if p.ID == prop.Relationships.Player.Data.ID {
				split := strings.SplitN(p.Attributes.Name, " ", 2)
				firstName := split[0]
				lastName := split[1]
				playerNames["first"] = append(playerNames["first"], firstName)
				playerNames["last"] = append(playerNames["last"], lastName)
			}
		}
	}
	fmt.Printf("Parse Prizepicks: %v\n", time.Since(t))
	fmt.Printf("\t Players: %v\n", len(playerNames["first"]))
	t = time.Now()
	client, err := database.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connect DB: %v\n", time.Since(t))
	t = time.Now()
	playersDB := client.Database("nba").Collection("players")
	fmt.Printf("Connect Players: %v\n", time.Since(t))
	t = time.Now()
	filter := bson.M{
		"first_name": bson.M{"$in": playerNames["first"]},
		"last_name":  bson.M{"$in": playerNames["last"]},
	}
	cur, err := playersDB.Find(ctx, filter)
	fmt.Printf("Get Players: %v\n", time.Since(t))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	t = time.Now()
	playerIDs := make([]int, 0, len(playerNames["first"]))
	players := make([]*model.Player, 0, len(playerNames["first"]))
	for cur.Next(ctx) {
		player := &model.Player{}
		err := cur.Decode(&player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
		playerIDs = append(playerIDs, player.PlayerID)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	sort.Sort(model.ByFirstName(players))
	fmt.Printf("Parse players: %v\n", time.Since(t))
	t = time.Now()
	gamesDB := client.Database("nba").Collection("games")
	fmt.Printf("Connect Games: %v\n", time.Since(t))
	t = time.Now()
	filter = bson.M{"player": bson.M{"$in": playerIDs}, "season": "2021-22"}

	cur, err = gamesDB.Find(ctx, filter)
	fmt.Printf("Get Games: %v\n", time.Since(t))
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	t = time.Now()
	games := make(map[int][]*model.PlayerGame, 500)
	for cur.Next(ctx) {
		game := &model.PlayerGame{}
		err := cur.Decode(&game)
		if err != nil {
			return nil, err
		}
		games[game.Player] = append(games[game.Player], game)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	fmt.Printf("Parse Games: %v\n", time.Since(t))
	x := 0
	for _, gs := range games {
		x += len(gs)
		// fmt.Printf("PLAYER: %v has %d\n", playerID, len(games))
	}
	fmt.Printf("TOTAL: %d\n", x)
	t = time.Now()
	playerProps := make([]*model.PlayerProp, 0, len(prizepicks.Data))
	ch := make(chan *model.PlayerProp, len(prizepicks.Data))
	wg := sync.WaitGroup{}
	fmt.Printf("\tData:%d\n", len(prizepicks.Data))
	fmt.Printf("\tIncluded:%d\n", len(prizepicks.Included))
	for _, prop := range prizepicks.Data {
		wg.Add(1)
		go createPlayerProp(prop, prizepicks.Included, players, games, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for prop := range ch {
		if prop.Player.PlayerID != 0 && len(prop.PlayerGames) != 0 {
			playerProps = append(playerProps, prop)
		}
	}
	fmt.Printf("Parse Remaining: %v\n", time.Since(t))
	return playerProps, nil
}

func createPlayerProp(prop model.PrizePicksData, included []model.PrizePicksIncluded, players []*model.Player, games map[int][]*model.PlayerGame, ch chan *model.PlayerProp, wg *sync.WaitGroup) {
	var playerName string
	var statType string
	for _, p := range included {
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
		// return nil, fmt.Errorf("error retrieving prizepick player name")
		fmt.Println("OOPS")
		wg.Done()
	}
	if statType == "" {
		// return nil, fmt.Errorf("error retrieving prizepick stat type")
		wg.Done()
	}
	target, err := strconv.ParseFloat(prop.Attributes.Line_score, 64)
	if err != nil {
		// return nil, err
		wg.Done()
	}
	propPlayer := findPlayerByName(playerName, players)
	ch <- &model.PlayerProp{Target: target, Type: statType, Opponent: prop.Attributes.Description, Player: propPlayer, PlayerGames: games[propPlayer.PlayerID]}
	wg.Done()
}

func findPlayerByName(name string, players []*model.Player) *model.Player {
	split := strings.SplitN(name, " ", 2)
	firstName := split[0]
	lastName := split[1]
	for _, player := range players {
		if firstName == player.FirstName && lastName == player.LastName {
			return player
		}
	}
	return &model.Player{}
	// i := sort.Search(len(players), func(i int) bool { return players[i].FirstName == firstName && players[i].LastName == lastName })
	// if i < len(players) && players[i].FirstName == firstName && players[i].LastName == lastName {
	// 	return players[i]
	// } else {
	// 	return &model.Player{}
	// }
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
