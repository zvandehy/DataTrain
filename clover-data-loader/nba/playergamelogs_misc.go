package nba

type PlayerMiscGameLog struct {
	SEASON_YEAR        string  `json:"SEASON_YEAR"`
	PLAYER_ID          float64 `json:"PLAYER_ID"`
	PLAYER_NAME        string  `json:"PLAYER_NAME"`
	NICKNAME           string  `json:"NICKNAME"`
	TEAM_ID            float64 `json:"TEAM_ID"`
	TEAM_ABBREVIATION  string  `json:"TEAM_ABBREVIATION"`
	TEAM_NAME          string  `json:"TEAM_NAME"`
	GAME_ID            string  `json:"GAME_ID"`
	GAME_DATE          string  `json:"GAME_DATE"`
	MATCHUP            string  `json:"MATCHUP"`
	WL                 string  `json:"WL"`
	MIN                float64 `json:"MIN"`
	PTS_OFF_TOV        float64 `json:"PTS_OFF_TOV"`
	PTS_2ND_CHANCE     float64 `json:"PTS_2ND_CHANCE"`
	PTS_FB             float64 `json:"PTS_FB"`
	PTS_PAINT          float64 `json:"PTS_PAINT"`
	OPP_PTS_OFF_TOV    float64 `json:"OPP_PTS_OFF_TOV"`
	OPP_PTS_2ND_CHANCE float64 `json:"OPP_PTS_2ND_CHANCE"`
	OPP_PTS_FB         float64 `json:"OPP_PTS_FB"`
	OPP_PTS_PAINT      float64 `json:"OPP_PTS_PAINT"`
	BLK                float64 `json:"BLK"`
	BLKA               float64 `json:"BLKA"`
	PF                 float64 `json:"PF"`
	PFD                float64 `json:"PFD"`
	NBA_FANTASY_PTS    float64 `json:"NBA_FANTASY_PTS"`
}

func NBAMiscPlayerGameLog(params GameLogParams) (*NBAResponse[PlayerMiscGameLog], error) {
	params.LeagueID = NBA_LEAGUE_ID
	params.MeasureType = "Misc"

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
	resp, err := NBAQuery[PlayerMiscGameLog](PlayerGameLogsURL, ParamMap(params))
	return resp, err
}
