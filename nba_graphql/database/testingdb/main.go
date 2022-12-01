package main

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/database"
)

func main() {
	// UploadComparisons()
	// db, err := database.NewSQLClient("NBA")
	// if err != nil {
	// 	panic(err)
	// }
	// comparisons := []PlayerComparison{}
	// players := []model.Player{}
	// err = db.Select(&players, "SELECT * FROM players WHERE name IN ('Deandre Ayton')")
	// if err != nil {
	// 	fmt.Println("Error getting players")
	// 	panic(err)
	// }
	// fmt.Println(len(players), "players")
	// fmt.Println(players[0].PlayerID)
	// propositions := []model.Proposition{}
	// err = db.Select(&propositions, "SELECT propositions.* FROM propositions JOIN playergames USING (playerID, gameID) WHERE propositions.playerID = ? AND playergames.date = Cast(\"2021-11-28\" AS Date)", players[0].PlayerID)
	// if err != nil {
	// 	fmt.Println("Error getting propositions")
	// 	panic(err)
	// }

	// err = db.Select(&comparisons, "SELECT toPlayerID,  from comparisons WHERE fromPlayerID = ? AND endDate = Cast(\"2022-11-28\" AS Date)", players[0].PlayerID)
	// if err != nil {
	// 	fmt.Println("Error getting comparisons")
	// 	fmt.Println(players[0].PlayerID)
	// 	panic(err)
	// }
	// fmt.Println(len(comparisons))

	db, err := database.NewSQLClient("NBA")
	if err != nil {
		panic(err)
	}
	dates := []struct {
		Date *time.Time `db:"date"`
	}{}
	err = db.DB.Select(&dates, "SELECT DISTINCT date FROM playergames WHERE season=\"2022-23\" AND date>Cast(\"2022-10-18\" AS Date) ORDER BY date;")
	if err != nil {
		panic(err)
	}
	for _, date := range dates {
		res := db.DB.MustExec(InsertStandardizedSQL(GetStats(), date.Date.Format("2006-01-02"), "2022-23", GetStandardizeddForDate(GetStats(), date.Date.Format("2006-01-02"), "2022-23")))
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("Uploaded %2.2d league standardized stats for %s", rowsAffected, date.Date.Format("2006-01-02"))
	}

}

func UploadComparisons() {
	db, err := database.NewSQLClient("NBA")
	if err != nil {
		panic(err)
	}

	datesQuery := "SELECT date FROM propositions JOIN playergames USING (playerID, gameID) WHERE season=\"2022-23\" GROUP BY date ORDER BY date;"
	dates := []struct {
		Date *time.Time `db:"date"`
	}{}
	err = db.DB.Select(&dates, datesQuery)
	if err != nil {
		panic(err)
	}

	stats := GetStats()
	for _, date := range dates {
		if date.Date.Before(time.Date(2022, 11, 29, 0, 0, 0, 0, time.UTC)) {
			continue
		}
		res := db.DB.MustExec(InsertStandardizedSQL(stats, date.Date.Format("2006-01-02"), "2022-23", GetComparisonsForPropositionsSQL(stats, date.Date.Format("2006-01-02"), "2022-23")))
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("Uploaded %5.5d comparisons for %s", rowsAffected, date.Date)
	}
}

type LeagueStandardized struct {
	Date                          *time.Time `db:"date"`
	Duration                      string     `db:"duration"`
	AVG_assists                   float64    `db:"AVG_assists"`
	STDDEV_assists                float64    `db:"STDDEV_assists"`
	AVG_defensiveRebounds         float64    `db:"AVG_defensiveRebounds"`
	STDDEV_defensiveRebounds      float64    `db:"STDDEV_defensiveRebounds"`
	AVG_offensiveRebounds         float64    `db:"AVG_offensiveRebounds"`
	STDDEV_offensiveRebounds      float64    `db:"STDDEV_offensiveRebounds"`
	AVG_fieldGoalsAttempted       float64    `db:"AVG_fieldGoalsAttempted"`
	STDDEV_fieldGoalsAttempted    float64    `db:"STDDEV_fieldGoalsAttempted"`
	AVG_fieldGoalsMade            float64    `db:"AVG_fieldGoalsMade"`
	STDDEV_fieldGoalsMade         float64    `db:"STDDEV_fieldGoalsMade"`
	AVG_freeThrowsAttempted       float64    `db:"AVG_freeThrowsAttempted"`
	STDDEV_freeThrowsAttempted    float64    `db:"STDDEV_freeThrowsAttempted"`
	AVG_freeThrowsMade            float64    `db:"AVG_freeThrowsMade"`
	STDDEV_freeThrowsMade         float64    `db:"STDDEV_freeThrowsMade"`
	AVG_personalFoulsDrawn        float64    `db:"AVG_personalFoulsDrawn"`
	STDDEV_personalFoulsDrawn     float64    `db:"STDDEV_personalFoulsDrawn"`
	AVG_personalFouls             float64    `db:"AVG_personalFouls"`
	STDDEV_personalFouls          float64    `db:"STDDEV_personalFouls"`
	AVG_points                    float64    `db:"AVG_points"`
	STDDEV_points                 float64    `db:"STDDEV_points"`
	AVG_threePointersAttempted    float64    `db:"AVG_threePointersAttempted"`
	STDDEV_threePointersAttempted float64    `db:"STDDEV_threePointersAttempted"`
	AVG_threePointersMade         float64    `db:"AVG_threePointersMade"`
	STDDEV_threePointersMade      float64    `db:"STDDEV_threePointersMade"`
	AVG_rebounds                  float64    `db:"AVG_rebounds"`
	STDDEV_rebounds               float64    `db:"STDDEV_rebounds"`
	AVG_turnovers                 float64    `db:"AVG_turnovers"`
	STDDEV_turnovers              float64    `db:"STDDEV_turnovers"`
	AVG_blocks                    float64    `db:"AVG_blocks"`
	STDDEV_blocks                 float64    `db:"STDDEV_blocks"`
	AVG_steals                    float64    `db:"AVG_steals"`
	STDDEV_steals                 float64    `db:"STDDEV_steals"`
	AVG_potentialAssists          float64    `db:"AVG_potentialAssists"`
	STDDEV_potentialAssists       float64    `db:"STDDEV_potentialAssists"`
	AVG_passes                    float64    `db:"AVG_passes"`
	STDDEV_passes                 float64    `db:"STDDEV_passes"`
	AVG_minutes                   float64    `db:"AVG_minutes"`
	STDDEV_minutes                float64    `db:"STDDEV_minutes"`
	AVG_heightInches              float64    `db:"AVG_heightInches"`
	STDDEV_heightInches           float64    `db:"STDDEV_heightInches"`
	AVG_weight                    float64    `db:"AVG_weight"`
	STDDEV_weight                 float64    `db:"STDDEV_weight"`
}

func InsertStandardizedSQL(stats []string, date string, duration string, getSql string) string {
	avgStdDevs := []string{}
	for _, stat := range stats {
		avgStdDevs = append(avgStdDevs, fmt.Sprintf("AVG_%[1]s, STDDEV_%[1]s", stat))
	}
	return fmt.Sprintf("REPLACE INTO standardized(`date`, `duration`, %s) %s", strings.Join(avgStdDevs, ", "), getSql)
}

func GetStandardizeddForDate(stats []string, date string, duration string) string {
	avgStdDevs := []string{}
	for _, stat := range stats {
		avgStdDevs = append(avgStdDevs, fmt.Sprintf("AVG(%[1]s) AS AVG_%[1]s, stddev(%[1]s) AS STDDEV_%[1]s", stat))
	}
	return fmt.Sprintf(`SELECT  Cast(%[1]s AS Date)     AS date
       ,%[2]s                      AS duration
       ,%[3]s
FROM playergames
JOIN players USING
(playerID
)
WHERE date < Cast(%[1]s AS Date)
AND season = %[2]s
AND minutes > 10;`, fmt.Sprintf("\"%s\"", date), fmt.Sprintf("\"%s\"", duration), strings.Join(avgStdDevs, ", "))
}

func GetComparisonsForPlayersSQL(stats []string, date string, duration string, playerIDs ...int) string {
	if len(playerIDs) == 0 {
		panic("No playerIDs provided")
	}
	diffs := []string{}
	toPlayerStats := []string{}
	leagueStats := []string{}
	for _, stat := range stats {
		diffs = append(diffs, statZScoreDiffSQL(stat))
		toPlayerStats = append(toPlayerStats, statToPlayerSQL(stat))
		leagueStats = append(leagueStats, statLeagueAvgStdDevSQL(stat))
	}
	strPlayerIDs := []string{}
	for _, playerID := range playerIDs {
		strPlayerIDs = append(strPlayerIDs, fmt.Sprintf("%d", playerID))
	}
	return fmt.Sprintf(`
		SELECT playerID AS fromPlayerID
       ,toPlayerID
       ,%[1]s AS duration
       ,Cast(%[2]s AS Date) AS endDate
	   ,%[3]s
	   ,((heightInches-LEAGUE_AVG_HEIGHT)/LEAGUE_STD_HEIGHT)-((toPlayer_HEIGHT-LEAGUE_AVG_HEIGHT)/LEAGUE_STD_HEIGHT) AS height
       ,((weight-LEAGUE_AVG_WEIGHT)/LEAGUE_STD_WEIGHT)-((toPlayer_WEIGHT-LEAGUE_AVG_WEIGHT)/LEAGUE_STD_WEIGHT) AS weight
	   FROM players
		JOIN playergames USING
		(playerID
		)
		JOIN
		(
			SELECT  playerID AS toPlayerID
				,%[4]s
				,heightInches AS toPlayer_HEIGHT
				,weight AS toPlayer_WEIGHT
			FROM players
			JOIN playergames USING
			(playerID
			)
			WHERE date < Cast(%[2]s AS Date)
			AND season = %[1]s
			GROUP BY  playerID
			ORDER BY AVG(minutes) DESC
			LIMIT 320
		) AS toPlayers
		JOIN
		(
			SELECT  %[5]s
				,AVG(heightInches)              AS LEAGUE_AVG_HEIGHT
				,stddev(heightInches)           AS LEAGUE_STD_HEIGHT
				,AVG(weight)                    AS LEAGUE_AVG_WEIGHT
				,stddev(weight)                 AS LEAGUE_STD_WEIGHT
			FROM playergames
			JOIN players USING
			(playerID
			)
			WHERE date < Cast(%[2]s AS Date)
			AND season = %[1]s 
		) AS LEAGUE
		WHERE date < Cast(%[2]s AS Date)
		AND season = %[1]s
		AND playerID IN ( %[6]s)
		AND toPlayerID <> playerID
		GROUP BY playerID
				,toPlayerID
				,toPlayer_POINTS
				,LEAGUE_AVG_POINTS`, fmt.Sprintf("\"%s\"", duration), fmt.Sprintf("\"%s\"", date), strings.Join(diffs, ", "), strings.Join(toPlayerStats, ", "), strings.Join(leagueStats, ", "), strings.Join(strPlayerIDs, ", "))
}

func GetComparisonsForPropositionsSQL(stats []string, date string, duration string) string {
	diffs := []string{}
	toPlayerStats := []string{}
	leagueStats := []string{}
	for _, stat := range stats {
		diffs = append(diffs, statZScoreDiffSQL(stat))
		toPlayerStats = append(toPlayerStats, statToPlayerSQL(stat))
		leagueStats = append(leagueStats, statLeagueAvgStdDevSQL(stat))
	}
	return fmt.Sprintf(`
		SELECT playerID AS fromPlayerID
       ,toPlayerID
       ,%[1]s AS duration
       ,Cast(%[2]s AS Date) AS endDate
	   ,%[3]s
	   ,((heightInches-LEAGUE_AVG_HEIGHT)/LEAGUE_STD_HEIGHT)-((toPlayer_HEIGHT-LEAGUE_AVG_HEIGHT)/LEAGUE_STD_HEIGHT) AS height
       ,((weight-LEAGUE_AVG_WEIGHT)/LEAGUE_STD_WEIGHT)-((toPlayer_WEIGHT-LEAGUE_AVG_WEIGHT)/LEAGUE_STD_WEIGHT) AS weight
	   FROM players
		JOIN playergames USING
		(playerID
		)
		JOIN
		(
			SELECT  playerID AS toPlayerID
				,%[4]s
				,heightInches AS toPlayer_HEIGHT
				,weight AS toPlayer_WEIGHT
			FROM players
			JOIN playergames USING
			(playerID
			)
			WHERE date < Cast(%[2]s AS Date)
			AND season = %[1]s
			GROUP BY  playerID
			ORDER BY AVG(minutes) DESC
			LIMIT 320
		) AS toPlayers
		JOIN
		(
			SELECT  %[5]s
				,AVG(heightInches)              AS LEAGUE_AVG_HEIGHT
				,stddev(heightInches)           AS LEAGUE_STD_HEIGHT
				,AVG(weight)                    AS LEAGUE_AVG_WEIGHT
				,stddev(weight)                 AS LEAGUE_STD_WEIGHT
			FROM playergames
			JOIN players USING
			(playerID
			)
			WHERE date < Cast(%[2]s AS Date)
			AND season = %[1]s 
		) AS LEAGUE
		WHERE date < Cast(%[2]s AS Date)
		AND season = %[1]s
		AND playerID IN ( SELECT playerID FROM propositions JOIN playergames USING (playerID, gameID ) WHERE season = %[1]s AND date = Cast(%[2]s AS Date) GROUP BY playerID)
		AND toPlayerID <> playerID
		GROUP BY playerID
				,toPlayerID
				,toPlayer_POINTS
				,LEAGUE_AVG_POINTS`, fmt.Sprintf("\"%s\"", duration), fmt.Sprintf("\"%s\"", date), strings.Join(diffs, ", "), strings.Join(toPlayerStats, ", "), strings.Join(leagueStats, ", "))
}

func statZScoreDiffSQL(stat string) string {
	return fmt.Sprintf("((AVG(%[1]s)-LEAGUE_AVG_%[2]s)/LEAGUE_STD_%[2]s)-((toPlayer_POINTS-LEAGUE_AVG_%[2]s)/LEAGUE_STD_%[2]s) AS %[1]s", stat, strings.ToUpper(stat))
}

func statToPlayerSQL(stat string) string {
	return fmt.Sprintf("AVG(%[1]s) AS toPlayer_%[2]s", stat, strings.ToUpper(stat))
}

func statLeagueAvgStdDevSQL(stat string) string {
	return fmt.Sprintf("AVG(%[1]s) AS LEAGUE_AVG_%[2]s, stddev(%[1]s) AS LEAGUE_STD_%[2]s", stat, strings.ToUpper(stat))
}

func GetStats() []string {
	return []string{
		"assists",
		"defensiveRebounds",
		"offensiveRebounds",
		"fieldGoalsAttempted",
		"fieldGoalsMade",
		"freeThrowsAttempted",
		"freeThrowsMade",
		"personalFoulsDrawn",
		"personalFouls",
		"points",
		"threePointersAttempted",
		"threePointersMade",
		"rebounds",
		"turnovers",
		"blocks",
		"steals",
		"potentialAssists",
		"passes",
		"minutes",
		"weight",
		"heightInches"}
}
