package nba

const PlayerGameLogsURL = "playergamelogs?"

type PlayerAdvancedGameLog struct {
	SEASON_YEAR        string  `json:"SEASON_YEAR"`
	PLAYER_ID          int     `json:"PLAYER_ID"`
	PLAYER_NAME        string  `json:"PLAYER_NAME"`
	NICKNAME           string  `json:"NICKNAME"`
	TEAM_ID            int     `json:"TEAM_ID"`
	TEAM_ABBREVIATION  string  `json:"TEAM_ABBREVIATION"`
	TEAM_NAME          string  `json:"TEAM_NAME"`
	GAME_ID            string  `json:"GAME_ID"`
	GAME_DATE          string  `json:"GAME_DATE"`
	MATCHUP            string  `json:"MATCHUP"`
	WL                 string  `json:"WL"`
	MIN                float64 `json:"MIN"`
	E_OFF_RATING       float64 `json:"E_OFF_RATING"`
	OFF_RATING         float64 `json:"OFF_RATING"`
	SP_WORK_OFF_RATING float64 `json:"sp_work_OFF_RATING"`
	E_DEF_RATING       float64 `json:"E_DEF_RATING"`
	DEF_RATING         float64 `json:"DEF_RATING"`
	SP_WORK_DEF_RATING float64 `json:"sp_work_DEF_RATING"`
	E_NET_RATING       float64 `json:"E_NET_RATING"`
	NET_RATING         float64 `json:"NET_RATING"`
	SP_WORK_NET_RATING float64 `json:"sp_work_NET_RATING"`
	AST_PCT            float64 `json:"AST_PCT"`
	AST_TO             float64 `json:"AST_TO"`
	AST_RATIO          float64 `json:"AST_RATIO"`
	OREB_PCT           float64 `json:"OREB_PCT"`
	DREB_PCT           float64 `json:"DREB_PCT"`
	REB_PCT            float64 `json:"REB_PCT"`
	TM_TOV_PCT         float64 `json:"TM_TOV_PCT"`
	E_TOV_PCT          float64 `json:"E_TOV_PCT"`
	EFG_PCT            float64 `json:"EFG_PCT"`
	TS_PCT             float64 `json:"TS_PCT"`
	USG_PCT            float64 `json:"USG_PCT"`
	E_USG_PCT          float64 `json:"E_USG_PCT"`
	E_PACE             float64 `json:"E_PACE"`
	PACE               float64 `json:"PACE"`
	PACE_PER40         float64 `json:"PACE_PER40"`
	SP_WORK_PACE       float64 `json:"sp_work_PACE"`
	PIE                float64 `json:"PIE"`
	POSS               float64 `json:"POSS"`
	FGM                float64 `json:"FGM"`
	FGA                float64 `json:"FGA"`
	FGM_PG             float64 `json:"FGM_PG"`
	FGA_PG             float64 `json:"FGA_PG"`
	FG_PCT             float64 `json:"FG_PCT"`
}

// TODO:  Change this to not be pointer
func NBAAdvancedPlayerGameLog(params GameLogParams) (*NBAResponse[PlayerAdvancedGameLog], error) {
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
	resp, err := NBAQuery[PlayerAdvancedGameLog](PlayerGameLogsURL, ParamMap(params))
	return resp, err
}

type GameLogParams struct {
	DateFrom       string `json:"DateFrom"`
	DateTo         string `json:"DateTo"`
	GameSegment    string `json:"GameSegment"`
	LastNGames     string `json:"LastNGames"`
	LeagueID       string `json:"LeagueID"`
	Location       string `json:"Location"`
	MeasureType    string `json:"MeasureType"`
	Month          string `json:"Month"`
	OpponentTeamID string `json:"OpponentTeamID"`
	Outcome        string `json:"Outcome"`
	PORound        string `json:"PORound"`
	PaceAdjust     string `json:"PaceAdjust"`
	PerMode        string `json:"PerMode"`
	Period         string `json:"Period"`
	PlusMinus      string `json:"PlusMinus"`
	Rank           string `json:"Rank"`
	Season         string `json:"Season"`
	SeasonSegment  string `json:"SeasonSegment"`
	SeasonType     string `json:"SeasonType"`
	ShotClockRange string `json:"ShotClockRange"`
	VsConference   string `json:"VsConference"`
	VsDivision     string `json:"VsDivision"`
}
