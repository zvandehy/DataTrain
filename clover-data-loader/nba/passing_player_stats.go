package nba

const LeagueDashboardStatsUrl = "leaguedashptstats?"

type PlayerPassingStats struct {
	PLAYER_ID           float64 `json:"PLAYER_ID"`
	PLAYER_NAME         string  `json:"PLAYER_NAME"`
	TEAM_ID             float64 `json:"TEAM_ID"`
	TEAM_ABBREVIATION   string  `json:"TEAM_ABBREVIATION"`
	GP                  float64 `json:"GP"`
	W                   float64 `json:"W"`
	L                   float64 `json:"L"`
	MIN                 float64 `json:"MIN"`
	PASSES_MADE         float64 `json:"PASSES_MADE"`
	PASSES_RECEIVED     float64 `json:"PASSES_RECEIVED"`
	AST                 float64 `json:"AST"`
	FT_AST              float64 `json:"FT_AST"`
	SECONDARY_AST       float64 `json:"SECONDARY_AST"`
	POTENTIAL_AST       float64 `json:"POTENTIAL_AST"`
	AST_POINTS_CREATED  float64 `json:"AST_POINTS_CREATED"`
	AST_ADJ             float64 `json:"AST_ADJ"`
	AST_TO_PASS_PCT     float64 `json:"AST_TO_PASS_PCT"`
	AST_TO_PASS_PCT_ADJ float64 `json:"AST_TO_PASS_PCT_ADJ"`
}

func PassingStats(params DashboardParams) (*NBAResponse[PlayerPassingStats], error) {
	if params.LastNGames == "" {
		params.LastNGames = "0"
	}
	if params.LeagueID == "" {
		params.LeagueID = "00"
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
	if params.PerMode == "" {
		params.PerMode = "PerGame"
	}
	if params.PlayerOrTeam == "" {
		params.PlayerOrTeam = "Player"
	}
	if params.PtMeasureType == "" {
		params.PtMeasureType = "Passing"
	}
	if params.Season == "" {
		params.Season = "2022-23"
	}
	if params.SeasonType == "" {
		params.SeasonType = "Regular Season"
	}
	if params.TeamID == "" {
		params.TeamID = "0"
	}
	resp, err := NBAQuery[PlayerPassingStats](LeagueDashboardStatsUrl, ParamMap(params))
	return resp, err
}

type DashboardParams struct {
	LeagueID      string `json:"LeagueID"`
	PtMeasureType string `json:"PtMeasureType"` // e.g. "Passing"
	PlayerOrTeam  string `json:"PlayerOrTeam"`  // "Player" or "Team"
	Season        string `json:"Season"`        //Season e.g. 2022-23
	PerMode       string `json:"PerMode"`       // "Totals" or "PerGame"

	DateFrom         string `json:"DateFrom"`
	DateTo           string `json:"DateTo"`
	LastNGames       string `json:"LastNGames"`
	Month            string `json:"Month"`
	OpponentTeamID   string `json:"OpponentTeamID"`
	SeasonType       string `json:"SeasonType"`
	College          string `json:"College"`
	Conference       string `json:"Conference"`
	Country          string `json:"Country"`
	Division         string `json:"Division"`
	DraftPick        string `json:"DraftPick"`
	DraftYear        string `json:"DraftYear"`
	GameScope        string `json:"GameScope"`
	Height           string `json:"Height"`
	Location         string `json:"Location"`
	Outcome          string `json:"Outcome"`
	PORound          string `json:"PORound"`
	PlayerExperience string `json:"PlayerExperience"`
	PlayerPosition   string `json:"PlayerPosition"`
	SeasonSegment    string `json:"SeasonSegment"`
	StarterBench     string `json:"StarterBench"`
	TeamID           string `json:"TeamID"`
	VsConference     string `json:"VsConference"`
	VsDivision       string `json:"VsDivision"`
	Weight           string `json:"Weight"`
}

// https://github.com/swar/nba_api/blob/a5d90f5a8637f76bc8bdb60af94428ba04036f12/src/nba_api/stats/library/parameters.py#L1
