package util

import (
	"fmt"
	"regexp"
	"time"
)

const (
	TIMENOW = "15:04:05.000"
)

var min, _ = time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")

func TIME_MINIMUM() time.Time {
	return min
}

var max, _ = time.Parse(time.RFC3339, "3000-01-01T00:00:00Z")

func TIME_MAXIMUM() time.Time {
	return max
}

const (
	SEASON_START_2022_23 = "2022-10-18"
	SEASON_END_2022_23   = "2023-04-10"
	SEASON_START_2021_22 = "2021-10-19"
	SEASON_END_2021_22   = "2022-04-11"
	SEASON_START_2020_21 = "2020-12-22"
	SEASON_END_2020_21   = "2021-04-17"
)

func TimeLog(msg string, start time.Time) string {
	return fmt.Sprintf("[%v] %s\n\tElapsed: %v", time.Now().Format(TIMENOW), msg, time.Since(start))
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
