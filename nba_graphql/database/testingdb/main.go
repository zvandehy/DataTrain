package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
)

func main() {
	db, err := database.NewSQLClient("NBA")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PlanetScale!")

	rows, err := db.Client.Query("SELECT * FROM players LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

}
