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
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Player is the resolver for the player field.
func (r *injuryResolver) Player(ctx context.Context, obj *model.Injury) (*model.Player, error) {
	return dataloader.For(ctx).PlayerByID.Load(obj.PlayerID)
}

// CurrentTeam is the resolver for the currentTeam field.
func (r *playerResolver) CurrentTeam(ctx context.Context, obj *model.Player) (*model.Team, error) {
	//logrus.Printf("Get TEAM for player %v", obj)
	if obj.CurrentTeam == "" {
		logrus.Errorf("CurrentTeam is empty for player %v", obj)
		return nil, fmt.Errorf("CurrentTeam is empty for player %v", obj)
	}
	t, err := dataloader.For(ctx).TeamByAbr.Load(obj.CurrentTeam)
	if err != nil {
		logrus.Errorf("Error loading team '%v' from: %v", obj.CurrentTeam, obj)
		return nil, err
	}
	if t == nil {
		return &model.Team{}, nil
	}
	return t, err
}

// Games is the resolver for the games field.
func (r *playerResolver) Games(ctx context.Context, obj *model.Player, input model.GameFilter) ([]*model.PlayerGame, error) {
	// logrus.Printf("Get Games filtered by %v for Player %v", input, obj)
	if obj.PlayerID == 0 {
		return []*model.PlayerGame{}, nil
	}
	input.PlayerID = &obj.PlayerID
	games, err := dataloader.For(ctx).PlayerGameByFilter.Load(input)
	if err != nil {
		return games, err
	}
	sort.SliceStable(games, func(i, j int) bool {
		a, err := time.Parse("2006-01-02", games[i].Date)
		if err != nil {
			return false
		}
		b, err := time.Parse("2006-01-02", games[j].Date)
		if err != nil {
			return false
		}
		return a.After(b)
	})
	if input.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *input.StartDate)
		if err != nil {
			return nil, err
		}
		var gamesToKeep []*model.PlayerGame
		for _, game := range games {
			gameDate, err := time.Parse("2006-01-02", game.Date)
			if err != nil {
				return nil, err
			}
			if gameDate.After(startDate) {
				gamesToKeep = append(gamesToKeep, game)
			}
		}
		games = gamesToKeep
	}
	if input.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			return nil, err
		}
		var gamesToKeep []*model.PlayerGame
		for _, game := range games {
			gameDate, err := time.Parse("2006-01-02", game.Date)
			if err != nil {
				return nil, err
			}
			if gameDate.Before(endDate) {
				gamesToKeep = append(gamesToKeep, game)
			}
		}
		games = gamesToKeep
	}
	return games, err
}

// Injuries is the resolver for the injuries field.
func (r *playerResolver) Injuries(ctx context.Context, obj *model.Player) ([]*model.Injury, error) {
	return dataloader.For(ctx).PlayerInjuryLoader.Load(obj.PlayerID)
}

// Projections is the resolver for the projections field.
func (r *playerResolver) Projections(ctx context.Context, obj *model.Player, input model.ProjectionFilter) ([]*model.Projection, error) {
	start := time.Now()
	var allProjections []*model.Projection
	// if strings.ToLower(*input.Sportsbook) != "prizepicks" {
	// 	return nil, fmt.Errorf("unsupported Sportsbook: %s. Current support only exists for: %v", *input.Sportsbook, []string{"PrizePicks"})
	// }

	if input.StartDate == nil && input.EndDate == nil {
		// today := time.Now().Format("2006-01-02")
		// input.StartDate = &today
	}
	input.PlayerName = &obj.Name

	cur, err := r.Db.GetProjections(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get player projections: %v", err)
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &allProjections)
	if err != nil {
		return nil, fmt.Errorf("failed to get player projections: %v", err)
	}
	if time.Since(start) > (time.Second * 5) {
		logrus.Warnf("Received %d player projections after %v", len(allProjections), time.Since(start))
	}
	duplicates := make(map[string][]*model.Projection, len(allProjections)/2)
	for _, projection := range allProjections {
		//TODO: Potential bug with "Tacos" / discounted / promo projections
		if projection.OpponentAbr == "" {
			continue
		}
		key := fmt.Sprintf("%s+%s", projection.PlayerName, projection.Date)
		if _, ok := duplicates[key]; !ok {
			duplicates[key] = []*model.Projection{projection}
		} else {
			duplicates[key] = append(duplicates[key], projection)
		}
	}
	//TODO: need to fix having same code in 2 places
	filteredProjections := []*model.Projection{}
	if input.Sportsbook != nil {
		logrus.Warn("filtering projections: ", *input.Sportsbook)
		for _, projection := range allProjections {
			filteredProps := []*model.Proposition{}
			for _, prop := range projection.Propositions {
				if strings.ToLower(prop.Sportsbook) == strings.ToLower(*input.Sportsbook) {
					filteredProps = append(filteredProps, prop)
				}
			}

			// log difference
			logrus.Warn("filtered projections: ", len(filteredProps), " of ", len(projection.Propositions))
			projection.Propositions = filteredProps
			if len(projection.Propositions) > 0 {
				filteredProjections = append(filteredProjections, projection)
			}
		}
	} else {
		return allProjections, nil
	}
	return filteredProjections, nil
}

// SimilarPlayers is the resolver for the similarPlayers field.
func (r *playerResolver) SimilarPlayers(ctx context.Context, obj *model.Player, input model.GameFilter) ([]*model.Player, error) {
	input.PlayerID = &obj.PlayerID
	return dataloader.For(ctx).SimilarPlayerLoader.Load(input)
}

// Opponent is the resolver for the opponent field.
func (r *playerGameResolver) Opponent(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	//logrus.Printf("Get Opponent from PlayerGame %v", obj)
	start := time.Now()
	if obj.OpponentID == 0 {
		return nil, fmt.Errorf("OpponentID doesn't exist on player game")
	}
	team, err := dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
	if time.Since(start) > (time.Second * 5) {
		logrus.Warnf("Received Team (from Playergame) after %v", time.Since(start))
	}
	return team, err
}

// OpponentStats is the resolver for the opponentStats field.
func (r *playerGameResolver) OpponentStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
	return dataloader.For(ctx).OpponentGameByPlayerGame.Load(*obj)
}

// Team is the resolver for the team field.
func (r *playerGameResolver) Team(ctx context.Context, obj *model.PlayerGame) (*model.Team, error) {
	start := time.Now()
	team, err := dataloader.For(ctx).TeamByID.Load(obj.TeamID)
	if time.Since(start) > (time.Second * 5) {
		logrus.Warnf("Received Opponent (from PlayerGame) after %v", time.Since(start))
	}
	return team, err
}

// TeamStats is the resolver for the teamStats field.
func (r *playerGameResolver) TeamStats(ctx context.Context, obj *model.PlayerGame) (*model.TeamGame, error) {
	return dataloader.For(ctx).TeamGameByPlayerGame.Load(*obj)
}

// Player is the resolver for the player field.
func (r *playerGameResolver) Player(ctx context.Context, obj *model.PlayerGame) (*model.Player, error) {
	//logrus.Printf("Get Player from PlayerGame %v", obj)
	return dataloader.For(ctx).PlayerByID.Load(obj.PlayerID)
}

// PlayersInGame is the resolver for the playersInGame field.
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

// Projections is the resolver for the projections field.
func (r *playerGameResolver) Projections(ctx context.Context, obj *model.PlayerGame) ([]*model.Projection, error) {
	cur, err := r.Db.GetPlayers(ctx, []model.PlayerFilter{{PlayerID: &obj.PlayerID}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var player model.Player
	cur.Next(ctx)
	if err != nil {
		return nil, err
	}
	err = cur.Decode(&player)
	if err != nil {
		return nil, err
	}
	playername := player.Name
	fmt.Println(playername)
	cur, err = r.Db.Collection("projections").Find(ctx, bson.M{"playername": playername, "date": obj.Date})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	fmt.Println(model.ProjectionFilter{PlayerID: &obj.PlayerID, StartDate: &obj.Date, EndDate: &obj.Date})
	defer cur.Close(ctx)
	var projections []*model.Projection
	cur.All(ctx, &projections)
	return projections, nil
}

// Team is the resolver for the team field.
func (r *playersInGameResolver) Team(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
	return obj.TeamPlayers, nil
}

// Opponent is the resolver for the opponent field.
func (r *playersInGameResolver) Opponent(ctx context.Context, obj *model.PlayersInGame) ([]*model.Player, error) {
	return obj.OpponentPlayers, nil
}

// Player is the resolver for the player field.
func (r *projectionResolver) Player(ctx context.Context, obj *model.Projection) (*model.Player, error) {
	//logrus.Printf("Get Player from Projection %v", obj)
	if obj.PlayerName == "" {
		logrus.Fatalf("PlayerName is empty: %#v", obj)
		return nil, fmt.Errorf("cannot get player from projection without player name")
	}
	playerFilter := model.PlayerFilter{Name: &obj.PlayerName}

	if val, ok := model.PlayerNames[obj.PlayerName]; ok {
		playerFilter.Name = &val
	}
	p, err := dataloader.For(ctx).PlayerByFilter.Load(playerFilter)
	if err != nil {
		logrus.Warnf("err when loading player for projection: %v", err)
		return &model.Player{Name: *playerFilter.Name}, nil
	}
	if p == nil {
		logrus.Warnf("Player %v is nil. Probably needs to be uploaded to the database.", *playerFilter.Name)
		return &model.Player{Name: *playerFilter.Name}, nil
	}
	return p, err
}

// Opponent is the resolver for the opponent field.
func (r *projectionResolver) Opponent(ctx context.Context, obj *model.Projection) (*model.Team, error) {
	if obj.OpponentAbr == "" {
		logrus.Fatalf("OpponentAbr is empty: %#v", obj)
		return nil, fmt.Errorf("cannot get opponent from projection without opponent name")
	}
	t, err := dataloader.For(ctx).TeamByAbr.Load(obj.OpponentAbr)
	if err != nil {
		logrus.Errorf("Error loading team '%v' from: %v", obj.OpponentAbr, obj)
		return nil, err
	}
	if t == nil {
		return &model.Team{}, nil
	}
	return t, err
}

// Result is the resolver for the result field.
func (r *projectionResolver) Result(ctx context.Context, obj *model.Projection) (*model.PlayerGame, error) {
	if len(strings.SplitN(obj.PlayerName, " ", 2)) != 2 {
		return nil, fmt.Errorf("cannot get result from projection without player name")
	}
	var player *model.Player
	err := r.Db.Collection("players").FindOne(ctx, bson.M{"first_name": strings.SplitN(obj.PlayerName, " ", 2)[0], "last_name": strings.SplitN(obj.PlayerName, " ", 2)[1]}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if player == nil {
		return nil, fmt.Errorf("didn't find a player with name %v", obj.PlayerName)
	}
	var game model.PlayerGame
	err = r.Db.Collection("games").FindOne(ctx, bson.M{"player": player.PlayerID, "date": obj.Date}).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Warnf("no game found for player %v on date %v", player.Name, obj.Date)
			return nil, nil
		}
		logrus.Error(err)
		return nil, err
	}
	return &game, nil
}

// Players is the resolver for the players field.
func (r *queryResolver) Players(ctx context.Context) ([]*model.Player, error) {
	fmt.Println(r.Db.Name)
	cur, err := r.Db.GetPlayers(ctx, []model.PlayerFilter{})
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

// FilterPlayers is the resolver for the filterPlayers field.
func (r *queryResolver) FilterPlayers(ctx context.Context, input model.PlayerFilter) ([]*model.Player, error) {
	// logrus.Printf("Get Players with filter  %v", input)
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

// Player is the resolver for the player field.
func (r *queryResolver) Player(ctx context.Context, input model.PlayerFilter) (*model.Player, error) {
	//logrus.Printf("Get Player with filter  %v", input)

	cur, err := r.Db.GetPlayers(ctx, []model.PlayerFilter{input})
	if err != nil {
		return nil, fmt.Errorf("error getting player: %v", err)
	}
	defer cur.Close(ctx)
	var player *model.Player
	cur.Next(ctx)
	//get first result
	err = cur.Decode(&player)
	if err != nil {
		return nil, fmt.Errorf("error decoding player: %v", err)
	}

	return player, nil
}

// Teams is the resolver for the teams field.
func (r *queryResolver) Teams(ctx context.Context) ([]*model.Team, error) {
	//logrus.Println("Get Teams")
	teamsDB := r.Db.Collection("teams")
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

// FilterTeams is the resolver for the filterTeams field.
func (r *queryResolver) FilterTeams(ctx context.Context, input model.TeamFilter) ([]*model.Team, error) {
	//logrus.Printf("Get Teams with filter %v\n", input)
	teamsDB := r.Db.Collection("teams")
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

// Team is the resolver for the team field.
func (r *queryResolver) Team(ctx context.Context, input model.TeamFilter) (*model.Team, error) {
	//logrus.Printf("Get Team with filter %#v\n", input)
	teamsDB := r.Db.Collection("teams")
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

// TeamGames is the resolver for the teamGames field.
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

// PlayerGames is the resolver for the playerGames field.
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

// Projections is the resolver for the projections field.
func (r *queryResolver) Projections(ctx context.Context, input model.ProjectionFilter) ([]*model.Projection, error) {
	// TODO: Fix timeouts so that can query for all projections -- maybe paging?
	//TODO: Refactor projection queries into own functions
	go func() {
		start := time.Now()
		var projections []*model.Projection
		url := "https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=7"
		res, err := http.Get(url)
		if err != nil {
			logrus.Warnf("couldn't retrieve prizepicks projections for today: %v", err)
			return
		}
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logrus.Warnf("couldn't read prizepicks projections for today: %v", err)
			return
		}
		var prizepicks model.PrizePicks
		if err := json.Unmarshal(bytes, &prizepicks); err != nil {
			logrus.Warnf("couldn't decode prizepicks projections for today: %v", err)
			return
		}
		for _, prop := range prizepicks.Data {
			if prop.Attributes.Is_promo {
				logrus.Warn("skipping promo")
				continue
			}
			projections, err = model.ParsePrizePick(prop, prizepicks.Included, projections)
			if err != nil {
				logrus.Warnf("couldn't parse prizepicks projections for today: %v", err)
				return
			}
		}
		projectionsDB := r.Db.Client.Database("nba").Collection("projections")
		insertCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		for _, projection := range projections {
			if val, ok := model.PlayerNames[projection.PlayerName]; ok {
				projection.PlayerName = val
			}
			//if projection with same playername and date exists, add the propositions to the list of propositions
			filter := bson.M{
				"playerName": projection.PlayerName,
				"date":       projection.Date,
			}
			var projectionFound *model.Projection
			err := projectionsDB.FindOne(insertCtx, filter).Decode(&projectionFound)
			if err != nil && err != mongo.ErrNoDocuments {
				logrus.Warnf("error finding projection for %v: %v", projection.PlayerName, err)
			}
			if projectionFound != nil {
				for _, prop := range projection.Propositions {
					//if projectionFound.Propositions contains the proposition, don't add it
					exists := false
					for _, propFound := range projectionFound.Propositions {
						if prop.Type == propFound.Type && prop.Target == propFound.Target && prop.Sportsbook == propFound.Sportsbook {
							exists = true
							break
						}
					}
					if !exists {
						projectionFound.Propositions = append(projectionFound.Propositions, prop)
					}
				}
				projection.Propositions = projectionFound.Propositions
			}
			//upsert
			res, err := projectionsDB.UpdateOne(insertCtx, filter, bson.M{"$set": projection}, options.Update().SetUpsert(true))
			if err != nil {
				if err != nil {
					logrus.Warn(err)
				}
				if res.UpsertedCount > 0 {
					logrus.Printf("INSERTED %v", projection.PlayerName)
				}
				if res.ModifiedCount > 0 {
					logrus.Printf("UPDATED %v", projection.PlayerName)
				}

			}
		}
		cancel()
		logrus.Printf("DONE retrieving NBA prizepicks projections\tTook %v", time.Since(start))
	}()
	go func() {
		start := time.Now()
		var projections []*model.Projection
		url := "https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=3"
		res, err := http.Get(url)
		if err != nil {
			logrus.Warnf("couldn't retrieve prizepicks projections for today: %v", err)
			return
		}
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logrus.Warnf("couldn't read prizepicks projections for today: %v", err)
			return
		}
		var prizepicks model.PrizePicks
		if err := json.Unmarshal(bytes, &prizepicks); err != nil {
			logrus.Warnf("couldn't decode prizepicks projections for today: %v", err)
			return
		}
		for _, prop := range prizepicks.Data {
			if prop.Attributes.Is_promo {
				logrus.Warn("skipping promo")
				continue
			}
			projections, err = model.ParsePrizePick(prop, prizepicks.Included, projections)
			if err != nil {
				logrus.Warnf("couldn't parse prizepicks projections for today: %v", err)
				return
			}
		}
		projectionsDB := r.Db.Client.Database("wnba").Collection("projections")
		insertCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		for _, projection := range projections {
			if val, ok := model.PlayerNames[projection.PlayerName]; ok {
				projection.PlayerName = val
			}
			//if projection with same playername and date exists, add the propositions to the list of propositions
			filter := bson.M{
				"playerName": projection.PlayerName,
				"date":       projection.Date,
			}
			var projectionFound *model.Projection
			err := projectionsDB.FindOne(insertCtx, filter).Decode(&projectionFound)
			if err != nil && err != mongo.ErrNoDocuments {
				logrus.Warnf("error finding projection for %v: %v", projection.PlayerName, err)
			}
			if projectionFound != nil {
				for _, prop := range projection.Propositions {
					//if projectionFound.Propositions contains the proposition, don't add it
					exists := false
					for _, propFound := range projectionFound.Propositions {
						if prop.Type == propFound.Type && prop.Target == propFound.Target && prop.Sportsbook == propFound.Sportsbook {
							exists = true
							break
						}
					}
					if !exists {
						projectionFound.Propositions = append(projectionFound.Propositions, prop)
					}
				}
				projection.Propositions = projectionFound.Propositions
			}
			//upsert
			res, err := projectionsDB.UpdateOne(insertCtx, filter, bson.M{"$set": projection}, options.Update().SetUpsert(true))
			if err != nil {
				if err != nil {
					logrus.Warn(err)
				}
				if res.UpsertedCount > 0 {
					logrus.Printf("INSERTED %v", projection.PlayerName)
				}
				if res.ModifiedCount > 0 {
					logrus.Printf("UPDATED %v", projection.PlayerName)
				}

			}
		}
		cancel()
		logrus.Printf("DONE retrieving WNBA prizepicks projections\tTook %v", time.Since(start))
	}()
	go func() {
		start := time.Now()
		var projections []*model.Projection
		url := "https://api.underdogfantasy.com/beta/v2/over_under_lines"
		res, err := http.Get(url)
		if err != nil {
			logrus.Warnf("couldn't retrieve underdog fantasy projections for today: %v", err)
			return
		}
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logrus.Warnf("couldn't read underdog fantasy projections for today: %v", err)
			return
		}
		var underdog model.UnderdogFantasy
		if err := json.Unmarshal(bytes, &underdog); err != nil {
			logrus.Warnf("couldn't decode underdog fantasy projections for today: %v", err)
			return
		}
		projections, err = model.ParseUnderdogProjection(underdog, "wnba")
		if err != nil {
			logrus.Warnf("couldn't parse underdog fantasy projections for today: %v", err)
			return
		}
		projectionsDB := r.Db.Client.Database("wnba").Collection("projections")
		insertCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		for _, projection := range projections {
			if val, ok := model.PlayerNames[projection.PlayerName]; ok {
				projection.PlayerName = val
			}
			//if projection with same playername and date exists, add the propositions to the list of propositions
			filter := bson.M{
				"playerName": projection.PlayerName,
				"date":       projection.Date,
			}
			var projectionFound *model.Projection
			err := projectionsDB.FindOne(insertCtx, filter).Decode(&projectionFound)
			if err != nil && err != mongo.ErrNoDocuments {
				logrus.Warnf("error finding projection for %v: %v", projection.PlayerName, err)
			}
			if projectionFound != nil {
				for _, prop := range projection.Propositions {
					//if projectionFound.Propositions contains the proposition, don't add it
					exists := false
					for _, propFound := range projectionFound.Propositions {
						if prop.Type == propFound.Type && prop.Target == propFound.Target && prop.Sportsbook == propFound.Sportsbook {
							exists = true
							break
						}
					}
					if !exists {
						projectionFound.Propositions = append(projectionFound.Propositions, prop)
					}
				}
				projection.Propositions = projectionFound.Propositions
			}
			//upsert
			res, err := projectionsDB.UpdateOne(insertCtx, filter, bson.M{"$set": projection}, options.Update().SetUpsert(true))
			if err != nil {
				if err != nil {
					logrus.Warn(err)
				}
				if res.UpsertedCount > 0 {
					logrus.Printf("INSERTED %v", projection.PlayerName)
				}
				if res.ModifiedCount > 0 {
					logrus.Printf("UPDATED %v", projection.PlayerName)
				}

			}
		}
		cancel()
		logrus.Printf("DONE retrieving WNBA prizepicks projections\tTook %v", time.Since(start))
	}()

	start := time.Now()
	var allProjections []*model.Projection
	// if input.Sportsbook != nil && strings.ToLower(*input.Sportsbook) != "prizepicks" {
	// 	return nil, fmt.Errorf("unsupported Sportsbook: %s. Current support only exists for: %v", *input.Sportsbook, []string{"PrizePicks"})
	// }

	if input.StartDate == nil && input.EndDate == nil {
		today := time.Now().Format("2006-01-02")
		input.StartDate = &today
	}

	cur, err := r.Db.GetProjections(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get projections: %v", err)
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &allProjections)
	if err != nil {
		return nil, fmt.Errorf("failed to get projections: %v", err)
	}
	if time.Since(start) > (time.Second * 5) {
		logrus.Warnf("Received %d projections after %v", len(allProjections), time.Since(start))
	}
	// remove projections that have player containing "Starters
	for i := 0; i < len(allProjections); i++ {
		if strings.Contains(allProjections[i].PlayerName, "Starters") {
			allProjections = append(allProjections[:i], allProjections[i+1:]...)
			i--
		}
	}
	filteredProjections := []*model.Projection{}
	if input.Sportsbook != nil {
		logrus.Warn("filtering projections: ", *input.Sportsbook)
		for _, projection := range allProjections {
			filteredProps := []*model.Proposition{}
			for _, prop := range projection.Propositions {
				if strings.EqualFold(prop.Sportsbook, *input.Sportsbook) {
					filteredProps = append(filteredProps, prop)
				}
			}

			// log difference
			logrus.Warn("filtered projections: ", len(filteredProps), " of ", len(projection.Propositions))
			projection.Propositions = filteredProps
			if len(projection.Propositions) > 0 {
				filteredProjections = append(filteredProjections, projection)
			}
		}
	} else {
		return allProjections, nil
	}
	return filteredProjections, nil
	// duplicates := make(map[string][]*model.Projection, len(allProjections)/2)
	// for _, projection := range allProjections {
	// 	//TODO: Potential bug with "Tacos" / discounted projections
	// 	if projection.OpponentAbr == "" {
	// 		continue
	// 	}
	// 	key := fmt.Sprintf("%s+%s", projection.PlayerName, projection.Date)
	// 	if _, ok := duplicates[key]; !ok {
	// 		duplicates[key] = []*model.Projection{projection}
	// 	} else {
	// 		duplicates[key] = append(duplicates[key], projection)
	// 	}
	// }
	// var uniqueProjections []*model.Projection
	// for _, projections := range duplicates {
	// 	best := model.GetBestProjection(projections)
	// 	if best.PlayerName == "" {
	// 		logrus.Fatalf("projection has BLANK playername: %#v", best)
	// 	}
	// 	uniqueProjections = append(uniqueProjections, best)
	// }

	// return allProjections, nil
}

// Games is the resolver for the games field.
func (r *teamResolver) Games(ctx context.Context, obj *model.Team, input model.GameFilter) ([]*model.TeamGame, error) {
	//logrus.Printf("Get Games from team %v filtered by %v", obj, input)
	//TODO: Add dataloader for situation where player.games.team.games is called
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

// Players is the resolver for the players field.
func (r *teamResolver) Players(ctx context.Context, obj *model.Team) ([]*model.Player, error) {
	//logrus.Printf("Get Players from Team %v", obj)
	input := model.PlayerFilter{TeamID: &obj.TeamID}
	return r.Query().FilterPlayers(ctx, input)
}

// Injuries is the resolver for the injuries field.
func (r *teamResolver) Injuries(ctx context.Context, obj *model.Team) ([]*model.Injury, error) {
	return dataloader.For(ctx).TeamInjuryLoader.Load(obj.TeamID)
}

// SimilarTeams is the resolver for the similarTeams field.
func (r *teamResolver) SimilarTeams(ctx context.Context, obj *model.Team, input model.GameFilter) ([]*model.Team, error) {
	input.TeamID = &obj.TeamID
	return dataloader.For(ctx).SimilarTeamLoader.Load(input)
}

// Opponent is the resolver for the opponent field.
func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
	//logrus.Printf("Get Opponent from TeamGame %v", obj)
	start := time.Now()
	team, err := dataloader.For(ctx).TeamByID.Load(obj.OpponentID)
	if time.Since(start) > (time.Second * 5) {
		logrus.Warnf("Received Opponent (from teamGame) after %v", time.Since(start))
	}
	return team, err
}

// PlayersInGame is the resolver for the playersInGame field.
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

// Injury returns generated.InjuryResolver implementation.
func (r *Resolver) Injury() generated.InjuryResolver { return &injuryResolver{r} }

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

type injuryResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type playerGameResolver struct{ *Resolver }
type playersInGameResolver struct{ *Resolver }
type projectionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
type teamGameResolver struct{ *Resolver }
