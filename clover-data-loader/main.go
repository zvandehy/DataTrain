package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/nba"
)

func main() {
	data, err := nba.NBAPlayerGameFinder(&nba.LeagueGameFinderParams{
		DateFrom: "12/01/2019",
		DateTo:   "12/31/2019",
	},
	)
	if err != nil {
		logrus.Fatalf("couldn't get LeagueGameFinder data: %v", err)
	}
	fmt.Println(data.Parameters)
	fmt.Println(data.Resource)
	fmt.Println(len(data.ResultSets))
	fmt.Println(data.ResultSets[0].Headers)
	fmt.Println(data.ResultSets[0].Name)
	fmt.Println(len(data.ResultSets[0].RowSet))
	fmt.Println(data.ResultSets[0].RowSet[0])
}
