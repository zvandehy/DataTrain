package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
	"github.com/zvandehy/DataTrain/nba_graphql/dataloader"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/resolver"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

const defaultPort = "8080"

func main() {
	logrus.Info("STARTING SERVER")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://clover-analytics.fly.dev", "https://www.clover-analytics.com/", "https://www.clover-analytics.com"},
		AllowCredentials: true,
		Debug:            true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		AllowOriginFunc: func(origin string) bool {
			return true
		}, //overrides allowed origins
	}).Handler)
	// nbaClient, err := database.ConnectDB(context.Background(), "nba")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// wnbaClient, err := database.ConnectDB(context.Background(), "wnba")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	nbaClient, err := database.NewSQLClient("NBA")
	if err != nil {
		log.Fatal(err)
	}
	// wnbaClient, err := database.NewSQLClient("WNBA")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// wait for the database to be ready
	time.Sleep(time.Second * 5)

	// call this function every Tick
	go func() {
		database.Getprizepicks(nbaClient)
		tick := time.Tick(10 * time.Minute)
		for range tick {
			database.Getprizepicks(nbaClient)
		}
	}()

	nbaServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver.NewResolver(nbaClient)}))
	// wnbaServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver.NewResolver(wnbaClient)}))

	router.Handle("/nba", playground.Handler("GraphQL playground", "/nba/query"))
	router.Handle("/nba/query", timer(dataloader.Middleware(nbaClient, nbaServer), nbaClient))
	// router.Handle("/wnba", playground.Handler("GraphQL playground", "/wnba/query"))
	// router.Handle("/wnba/query", timer(dataloader.Middleware(wnbaClient, wnbaServer), wnbaClient))

	log.Printf("connect to http://localhost:%s/nba for NBA Playground\n", port)
	// log.Printf("connect to http://localhost:%s/wnba for WNBA Playground\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func timer(h http.Handler, db database.BasketballRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//skip heartbeat logs
		if r.ContentLength == 1468 {
			h.ServeHTTP(w, r)
			return
		}
		db.SetQueries(0)
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		logrus.Printf("[%v] Request took %v to resolve %v queries.", time.Now().Format(util.TIMENOW), duration, db.CountQueries())
	})
}

// func getprizepicks(nbaClient *database.MongoClient) {
// 	leagueID := 7
// 	if nbaClient.Name == "wnba" {
// 		leagueID = 3
// 	}
// 	start := time.Now()
// 	var projections []*model.Projection
// 	url := fmt.Sprintf("https://partner-api.prizepicks.com/projections?single_stat=True&per_page=1000&league_id=%d", leagueID)
// 	res, err := http.Get(url)
// 	if err != nil {
// 		logrus.Warnf("couldn't retrieve prizepicks projections for today: %v", err)
// 		return
// 	}
// 	bytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		logrus.Warnf("couldn't read prizepicks projections for today: %v", err)
// 		return
// 	}
// 	var prizepicks model.PrizePicks
// 	if err := json.Unmarshal(bytes, &prizepicks); err != nil {
// 		logrus.Warnf("couldn't decode prizepicks projections for today: %v", err)
// 		return
// 	}
// 	for _, prop := range prizepicks.Data {
// 		if prop.Attributes.Is_promo {
// 			logrus.Warn("skipping promo")
// 			continue
// 		}
// 		projections, err = model.ParsePrizePick(prop, prizepicks.Included, projections)
// 		if err != nil {
// 			logrus.Warnf("couldn't parse prizepicks projections for today: %v", err)
// 			return
// 		}
// 	}
// 	projectionsDB := nbaClient.Collection("projections")
// 	insertCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	for _, projection := range projections {
// 		if val, ok := model.PlayerNames[projection.PlayerName]; ok {
// 			projection.PlayerName = val
// 		}
// 		// if projection with same playername and date exists, add the propositions to the list of propositions
// 		filter := bson.M{
// 			"playerName": projection.PlayerName,
// 			"date":       projection.Date,
// 		}
// 		var projectionFound *model.Projection
// 		err := projectionsDB.FindOne(insertCtx, filter).Decode(&projectionFound)
// 		if err != nil && err != mongo.ErrNoDocuments {
// 			logrus.Warnf("error finding projection for %v: %v", projection.PlayerName, err)
// 		}
// 		if projectionFound != nil {
// 			for _, prop := range projection.Props {
// 				//if projectionFound.Propositions contains the proposition, don't add it
// 				exists := false
// 				for _, propFound := range projectionFound.Props {
// 					if prop.Type == propFound.Type && prop.Target == propFound.Target && prop.Sportsbook == propFound.Sportsbook {
// 						exists = true
// 						break
// 					}
// 				}
// 				if !exists {
// 					projectionFound.Props = append(projectionFound.Props, prop)
// 				}
// 			}
// 			projection.Props = projectionFound.Props
// 		}
// 		//upsert
// 		res, err := projectionsDB.UpdateOne(insertCtx, filter, bson.M{"$set": projection}, options.Update().SetUpsert(true))
// 		if err != nil {
// 			if err != nil {
// 				logrus.Warn(err)
// 			}
// 		}

// 		if res.UpsertedCount > 0 {
// 			logrus.Printf("INSERTED '%v'", projection.PlayerName)
// 		}
// 		// if res.ModifiedCount > 0 {
// 		// 	logrus.Printf("UPDATED '%v'", projection.PlayerName)
// 		// }
// 	}
// 	cancel()
// 	logrus.Printf(util.TimeLog(fmt.Sprintf("Retrieved %d %s projections from PrizePicks", len(projections), nbaClient.Name), start))
// }

// func cachePlayers(nbaClient *database.MongoClient) {
// 	nbaClient.PlayerCache = make(map[string][]*model.Player)
// 	nbaClient.TeamCache = make(map[string][]*model.Team)
// 	logrus.Info("CACHEING PLAYERS")
// 	// TODO: After automating game data collection, the cache should be updated
// 	caches := [][]model.SeasonOption{
// 		{model.SEASON_2020_21},
// 		{model.SEASON_2021_22},
// 		{model.SEASON_2022_23},
// 		{model.SEASON_2020_21, model.SEASON_2021_22},
// 		{model.SEASON_2021_22, model.SEASON_2022_23},
// 		{model.SEASON_2020_21, model.SEASON_2021_22, model.SEASON_2022_23},
// 	}
// 	for _, cache := range caches {
// 		players, err := nbaClient.GetPlayers(context.Background(), true, &model.PlayerFilter{Seasons: &cache})
// 		if err != nil {
// 			logrus.Errorf("Error getting players for cache: %v", err)
// 		}
// 		nbaClient.PlayerCache[fmt.Sprintf("%v", cache)] = players
// 		logrus.Info("Cached players for: ", cache)
// 	}

// 	nbaClient.PlayerSimilarity = *model.NewPlayerSnapshots()
// 	nbaClient.TeamSimilarity = *model.NewTeamSnapshots()
// }
