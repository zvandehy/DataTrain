package main

import (
	"context"
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
	"github.com/zvandehy/DataTrain/nba_graphql/graph"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()

	//add comment
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://datatrain-mp34k.ondigitalocean.app", "https://clover-analytics.fly.dev", "https://www.clover-analytics.com/"},
		AllowCredentials: true,
		Debug:            true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// AllowOriginFunc: func(origin string) bool {
		// 	return true
		// }, //overrides allowed origins
	}).Handler)
	nbaClient, err := database.ConnectDB(context.Background(), "nba")
	if err != nil {
		log.Fatal(err)
	}
	wnbaClient, err := database.ConnectDB(context.Background(), "wnba")
	if err != nil {
		log.Fatal(err)
	}
	// wait for the database to be ready
	time.Sleep(time.Second * 5)

	nbaServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Db: nbaClient}}))
	wnbaServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Db: wnbaClient}}))

	router.Handle("/nba", playground.Handler("GraphQL playground", "/nba/query"))
	router.Handle("/nba/query", timer(dataloader.Middleware(nbaClient, nbaServer), nbaClient))
	router.Handle("/wnba", playground.Handler("GraphQL playground", "/wnba/query"))
	router.Handle("/wnba/query", timer(dataloader.Middleware(wnbaClient, wnbaServer), wnbaClient))

	log.Printf("connect to http://localhost:%s/nba for NBA UI\n", port)
	log.Printf("connect to http://localhost:%s/wnba for WNBA UI\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func timer(h http.Handler, db *database.NBADatabaseClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//skip heartbeat logs
		if r.ContentLength == 1468 {
			h.ServeHTTP(w, r)
			return
		}
		db.Queries = 0
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		logrus.Printf("Request took %v to resolve %v queries.", duration, db.Queries)
	})
}
