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
	// router.Use(cors.New(cors.Options{
	// 	// AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return true
	// 	// }, //overrides allowed origins
	// }).Handler)

	mongoClient, err := database.ConnectDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Db: mongoClient}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", timer(dataloader.Middleware(mongoClient, srv), mongoClient))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func timer(h http.Handler, db *database.NBADatabaseClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db.Queries = 0
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		logrus.Printf("Request took %v to resolve %v queries.", duration, db.Queries)
	})
}
