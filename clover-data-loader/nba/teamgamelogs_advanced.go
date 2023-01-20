package nba

type TeamAdvancedGameLog struct {
	SEASON_YEAR       string  `json:"SEASON_YEAR"`
	TEAM_ID           float64 `json:"TEAM_ID"`
	TEAM_ABBREVIATION string  `json:"TEAM_ABBREVIATION"`
	TEAM_NAME         string  `json:"TEAM_NAME"`
	GAME_ID           string  `json:"GAME_ID"`
	GAME_DATE         string  `json:"GAME_DATE"`
	MATCHUP           string  `json:"MATCHUP"`
	WL                string  `json:"WL"`
	MIN               float64 `json:"MIN"`
	E_OFF_RATING      float64 `json:"E_OFF_RATING"`
	OFF_RATING        float64 `json:"OFF_RATING"`
	E_DEF_RATING      float64 `json:"E_DEF_RATING"`
	DEF_RATING        float64 `json:"DEF_RATING"`
	E_NET_RATING      float64 `json:"E_NET_RATING"`
	NET_RATING        float64 `json:"NET_RATING"`
	AST_PCT           float64 `json:"AST_PCT"`
	AST_TO            float64 `json:"AST_TO"`
	AST_RATIO         float64 `json:"AST_RATIO"`
	OREB_PCT          float64 `json:"OREB_PCT"`
	DREB_PCT          float64 `json:"DREB_PCT"`
	REB_PCT           float64 `json:"REB_PCT"`
	TM_TOV_PCT        float64 `json:"TM_TOV_PCT"`
	EFG_PCT           float64 `json:"EFG_PCT"`
	TS_PCT            float64 `json:"TS_PCT"`
	E_PACE            float64 `json:"E_PACE"`
	PACE              float64 `json:"PACE"`
	PACE_PER40        float64 `json:"PACE_PER40"`
	POSS              float64 `json:"POSS"`
	PIE               float64 `json:"PIE"`
}

func NBATeamAdvancedGameLog(params GameLogParams) (*NBAResponse[TeamAdvancedGameLog], error) {
	params.LeagueID = NBA_LEAGUE_ID
	params.MeasureType = "Advanced"

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
	resp, err := NBAQuery[TeamAdvancedGameLog](TeamGameLogsURL, ParamMap(params))
	return resp, err
}
