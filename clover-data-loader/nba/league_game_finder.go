package nba

import "reflect"

const LeagueGameFinderURL = "leaguegamefinder?"

// LeagueGameFinderResults is a struct that contains all the fields for the LeagueGameFinderResults endpoint.
type LeagueGameFinderResults struct {
	SeasonID         string  `json:"SEASON_ID"`
	PlayerID         float64 `json:"PLAYER_ID"`
	PlayerName       string  `json:"PLAYER_NAME"`
	TeamID           float64 `json:"TEAM_ID"`
	TeamAbbreviation string  `json:"TEAM_ABBREVIATION"`
	TeamName         string  `json:"TEAM_NAME"`
	GameID           string  `json:"GAME_ID"`
	GameDate         string  `json:"GAME_DATE"`
	Matchup          string  `json:"MATCHUP"`
	WL               string  `json:"WL"`
	Min              float64 `json:"MIN"`
	Pts              float64 `json:"PTS"`
	Fgm              float64 `json:"FGM"`
	Fga              float64 `json:"FGA"`
	FgPct            float64 `json:"FG_PCT"`
	Fg3M             float64 `json:"FG3M"`
	Fg3A             float64 `json:"FG3A"`
	Fg3Pct           float64 `json:"FG3_PCT"`
	Ftm              float64 `json:"FTM"`
	Fta              float64 `json:"FTA"`
	FtPct            float64 `json:"FT_PCT"`
	Oreb             float64 `json:"OREB"`
	Dreb             float64 `json:"DREB"`
	Reb              float64 `json:"REB"`
	Ast              float64 `json:"AST"`
	Stl              float64 `json:"STL"`
	Blk              float64 `json:"BLK"`
	Tov              float64 `json:"TOV"`
	Pf               float64 `json:"PF"`
	PlusMinus        float64 `json:"PLUS_MINUS"`
}

func NBAPlayerGameFinder(params *LeagueGameFinderParams) (*NBAResponse[LeagueGameFinderResults], error) {
	if params == nil {
		params = &LeagueGameFinderParams{}
	}
	params.LeagueID = NBA_LEAGUE_ID
	params.PlayerOrTeam = "P"
	return NBAQuery[LeagueGameFinderResults](LeagueGameFinderURL, params.ParamMap())
}

type LeagueGameFinderParams struct {
	LeagueID     string `json:"LeagueID"`
	PlayerOrTeam string `json:"PlayerOrTeam"`

	Conference      string `json:"Conference"`
	DateFrom        string `json:"DateFrom"` //MM/DD/YYYY
	DateTo          string `json:"DateTo"`   //MM/DD/YYYY
	Division        string `json:"Division"`
	DraftNumber     string `json:"DraftNumber"`
	DraftRound      string `json:"DraftRound"`
	DraftTeamID     string `json:"DraftTeamID"`
	DraftYear       string `json:"DraftYear"`
	EqAST           string `json:"EqAST"`
	EqBLK           string `json:"EqBLK"`
	EqDD            string `json:"EqDD"`
	EqDREB          string `json:"EqDREB"`
	EqFG3A          string `json:"EqFG3A"`
	EqFG3M          string `json:"EqFG3M"`
	EqFG3_PCT       string `json:"EqFG3_PCT"`
	EqFGA           string `json:"EqFGA"`
	EqFGM           string `json:"EqFGM"`
	EqFG_PCT        string `json:"EqFG_PCT"`
	EqFTA           string `json:"EqFTA"`
	EqFTM           string `json:"EqFTM"`
	EqFT_PCT        string `json:"EqFT_PCT"`
	EqMINUTES       string `json:"EqMINUTES"`
	EqOREB          string `json:"EqOREB"`
	EqPF            string `json:"EqPF"`
	EqPTS           string `json:"EqPTS"`
	EqREB           string `json:"EqREB"`
	EqSTL           string `json:"EqSTL"`
	EqTD            string `json:"EqTD"`
	EqTOV           string `json:"EqTOV"`
	GameID          string `json:"GameID"`
	GtAST           string `json:"GtAST"`
	GtBLK           string `json:"GtBLK"`
	GtDD            string `json:"GtDD"`
	GtDREB          string `json:"GtDREB"`
	GtFG3A          string `json:"GtFG3A"`
	GtFG3M          string `json:"GtFG3M"`
	GtFG3_PCT       string `json:"GtFG3_PCT"`
	GtFGA           string `json:"GtFGA"`
	GtFGM           string `json:"GtFGM"`
	GtFG_PCT        string `json:"GtFG_PCT"`
	GtFTA           string `json:"GtFTA"`
	GtFTM           string `json:"GtFTM"`
	GtFT_PCT        string `json:"GtFT_PCT"`
	GtMINUTES       string `json:"GtMINUTES"`
	GtOREB          string `json:"GtOREB"`
	GtPF            string `json:"GtPF"`
	GtPTS           string `json:"GtPTS"`
	GtREB           string `json:"GtREB"`
	GtSTL           string `json:"GtSTL"`
	GtTD            string `json:"GtTD"`
	GtTOV           string `json:"GtTOV"`
	Location        string `json:"Location"`
	LtAST           string `json:"LtAST"`
	LtBLK           string `json:"LtBLK"`
	LtDD            string `json:"LtDD"`
	LtDREB          string `json:"LtDREB"`
	LtFG3A          string `json:"LtFG3A"`
	LtFG3M          string `json:"LtFG3M"`
	LtFG3_PCT       string `json:"LtFG3_PCT"`
	LtFGA           string `json:"LtFGA"`
	LtFGM           string `json:"LtFGM"`
	LtFG_PCT        string `json:"LtFG_PCT"`
	LtFTA           string `json:"LtFTA"`
	LtFTM           string `json:"LtFTM"`
	LtFT_PCT        string `json:"LtFT_PCT"`
	LtMINUTES       string `json:"LtMINUTES"`
	LtOREB          string `json:"LtOREB"`
	LtPF            string `json:"LtPF"`
	LtPTS           string `json:"LtPTS"`
	LtREB           string `json:"LtREB"`
	LtSTL           string `json:"LtSTL"`
	LtTD            string `json:"LtTD"`
	LtTOV           string `json:"LtTOV"`
	Outcome         string `json:"Outcome"`
	PORound         string `json:"PORound"`
	PlayerID        string `json:"PlayerID"`
	RookieYear      string `json:"RookieYear"`
	Season          string `json:"Season"`
	SeasonSegment   string `json:"SeasonSegment"`
	SeasonType      string `json:"SeasonType"`
	StarterBench    string `json:"StarterBench"`
	TeamID          string `json:"TeamID"`
	VsConference    string `json:"VsConference"`
	VsDivision      string `json:"VsDivision"`
	VsTeamID        string `json:"VsTeamID"`
	YearsExperience string `json:"YearsExperience"`
}

func (p *LeagueGameFinderParams) ParamMap() map[string]string {
	params := map[string]string{}

	// if p.LeagueID != "" {
	// 	params["LeagueID"] = p.LeagueID
	// } else {
	// 	params["LeagueID"] = NBA_LEAGUE_ID
	// }

	// if p.PlayerOrTeam != "" {
	// 	params["PlayerOrTeam"] = p.PlayerOrTeam
	// } else {
	// 	params["PlayerOrTeam"] = "P"
	// }

	// if p.DateFrom != "" {
	// 	params["DateFrom"] = p.DateFrom
	// }

	// if p.DateTo != "" {
	// 	params["DateTo"] = p.DateTo
	// }

	//for each param in the struct, add it to the map
	val := reflect.ValueOf(*p)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			if val.Field(i).String() != "" {
				params[tag] = val.Field(i).String()
			}
		}
	}

	return params
}

// https://github.com/swar/nba_api/blob/a5d90f5a8637f76bc8bdb60af94428ba04036f12/src/nba_api/stats/library/parameters.py#L1
