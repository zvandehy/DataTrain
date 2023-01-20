package cloverdata

import (
	"fmt"
	"strings"
)

func GetComparableStats() []string {
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
		"heightInches",
	}
}

func SQLUploadPlayerComparisons(date string) string {
	avgStdDevs := []string{}
	for _, stat := range GetComparableStats() {
		avgStdDevs = append(avgStdDevs, fmt.Sprintf("AVG_%[1]s, STDDEV_%[1]s", stat))
	}
	return fmt.Sprintf("REPLACE INTO standardized(`date`, `duration`, %s) %s", strings.Join(avgStdDevs, ", "), SQLPlayerComparisonsOnDate(date))
}

func SQLPlayerComparisonsOnDate(date string) string {
	comparableStats := GetComparableStats()
	duration := "2022-23"

	avgStdDevs := []string{}
	for _, stat := range comparableStats {
		avgStdDevs = append(avgStdDevs, fmt.Sprintf("AVG(%[1]s) AS AVG_%[1]s, stddev(%[1]s) AS STDDEV_%[1]s", stat))
	}
	return fmt.Sprintf(`SELECT Cast(%[1]s AS Date) AS date, %[2]s AS duration, %[3]s FROM playergames
		JOIN players USING
			(playerID)
		WHERE date < Cast(%[1]s AS Date)
		AND season = %[2]s;`, fmt.Sprintf("\"%s\"", date), fmt.Sprintf("\"%s\"", duration), strings.Join(avgStdDevs, ", "))
}
