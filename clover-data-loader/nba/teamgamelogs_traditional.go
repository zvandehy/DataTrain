package nba

const TeamGameLogsURL = "teamgamelogs?"

type TeamTraditionalGameLog struct {
	SEASON_YEAR       string  `json:"SEASON_YEAR"`
	TEAM_ID           float64 `json:"TEAM_ID"`
	TEAM_ABBREVIATION string  `json:"TEAM_ABBREVIATION"`
	TEAM_NAME         string  `json:"TEAM_NAME"`
	GAME_ID           string  `json:"GAME_ID"`
	GAME_DATE         string  `json:"GAME_DATE"`
	MATCHUP           string  `json:"MATCHUP"`
	WL                string  `json:"WL"`
	MIN               float64 `json:"MIN"`
	FGM               float64 `json:"FGM"`
	FGA               float64 `json:"FGA"`
	FG_PCT            float64 `json:"FG_PCT"`
	FG3M              float64 `json:"FG3M"`
	FG3A              float64 `json:"FG3A"`
	FG3_PCT           float64 `json:"FG3_PCT"`
	FTM               float64 `json:"FTM"`
	FTA               float64 `json:"FTA"`
	FT_PCT            float64 `json:"FT_PCT"`
	OREB              float64 `json:"OREB"`
	DREB              float64 `json:"DREB"`
	REB               float64 `json:"REB"`
	AST               float64 `json:"AST"`
	TOV               float64 `json:"TOV"`
	STL               float64 `json:"STL"`
	BLK               float64 `json:"BLK"`
	BLKA              float64 `json:"BLKA"`
	PF                float64 `json:"PF"`
	PFD               float64 `json:"PFD"`
	PTS               float64 `json:"PTS"`
	PLUS_MINUS        float64 `json:"PLUS_MINUS"`
}

func NBATeamTraditionalGameLog(params GameLogParams) (*NBAResponse[TeamTraditionalGameLog], error) {
	params.LeagueID = NBA_LEAGUE_ID
	params.MeasureType = "Base"

	if params.LastNGames == "" {
		params.LastNGames = "0"
	}
	if params.Month == "" {
		params.Month = "0"
	}
	if params.OpponentTeamID == "" {
		params.OpponentTeamID = "0"
	}
	if params.PORound == "" {
		params.PORound = "0"
	}
	if params.PaceAdjust == "" {
		params.PaceAdjust = "N"
	}
	if params.PlusMinus == "" {
		params.PlusMinus = "N"
	}
	if params.Rank == "" {
		params.Rank = "N"
	}
	if params.PerMode == "" {
		params.PerMode = "Totals"
	}
	if params.Period == "" {
		params.Period = "0"
	}
	if params.Season == "" {
		params.Season = "2022-23"
	}
	if params.SeasonType == "" {
		params.SeasonType = "Regular Season"
	}
	resp, err := NBAQuery[TeamTraditionalGameLog](TeamGameLogsURL, ParamMap(params))
	return resp, err
}
