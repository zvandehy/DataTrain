package nba

import "github.com/sirupsen/logrus"

const LeagueDashboardDefenseURL = "leaguedashptdefend?"

type PlayerDefenseStats struct {
	CLOSE_DEF_PERSON_ID           int     `json:"CLOSE_DEF_PERSON_ID"`
	PLAYER_NAME                   string  `json:"PLAYER_NAME"`
	PLAYER_LAST_TEAM_ID           int     `json:"PLAYER_LAST_TEAM_ID"`
	PLAYER_LAST_TEAM_ABBREVIATION string  `json:"PLAYER_LAST_TEAM_ABBREVIATION"`
	PLAYER_POSITION               string  `json:"PLAYER_POSITION"`
	AGE                           float64 `json:"AGE"`
	GP                            float64 `json:"GP"`
	G                             float64 `json:"G"`
	FREQ                          float64 `json:"FREQ"`
	D_FGM                         float64 `json:"D_FGM"`
	D_FGA                         float64 `json:"D_FGA"`
	D_FG_PCT                      float64 `json:"D_FG_PCT"`
	NORMAL_FG_PCT                 float64 `json:"NORMAL_FG_PCT"`
	PCT_PLUSMINUS                 float64 `json:"PCT_PLUSMINUS"`
}

func NBAPlayerDefense(params DefenseDashboardParams) (*NBAResponse[PlayerDefenseStats], error) {
	if params.DefenseCategory == "" {
		params.DefenseCategory = "Overall"
	}
	if params.LastNGames == "" {
		params.LastNGames = "0"
	}
	if params.LeagueID == "" {
		params.LeagueID = NBA_LEAGUE_ID
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
	if params.Period == "" {
		params.Period = "0"
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
	resp, err := NBAQuery[PlayerDefenseStats](LeagueDashboardDefenseURL, ParamMap(params))
	return resp, err
}

type DefenseDashboardParams struct {
	College          string `json:"College"`
	Conference       string `json:"Conference"`
	Country          string `json:"Country"`
	DateFrom         string `json:"DateFrom"`
	DateTo           string `json:"DateTo"`
	DefenseCategory  string `json:"DefenseCategory"`
	Division         string `json:"Division"`
	DraftPick        string `json:"DraftPick"`
	DraftYear        string `json:"DraftYear"`
	GameSegment      string `json:"GameSegment"`
	Height           string `json:"Height"`
	LastNGames       string `json:"LastNGames"`
	LeagueID         string `json:"LeagueID"`
	Location         string `json:"Location"`
	Month            string `json:"Month"`
	OpponentTeamID   string `json:"OpponentTeamID"`
	Outcome          string `json:"Outcome"`
	PORound          string `json:"PORound"`
	PerMode          string `json:"PerMode"`
	Period           string `json:"Period"`
	PlayerExperience string `json:"PlayerExperience"`
	PlayerPosition   string `json:"PlayerPosition"`
	Season           string `json:"Season"`
	SeasonSegment    string `json:"SeasonSegment"`
	SeasonType       string `json:"SeasonType"`
	StarterBench     string `json:"StarterBench"`
	TeamID           string `json:"TeamID"`
	VsConference     string `json:"VsConference"`
	VsDivision       string `json:"VsDivision"`
	Weight           string `json:"Weight"`
}

func FindDefenseByID(defense []PlayerDefenseStats, playerID int) PlayerDefenseStats {
	for _, d := range defense {
		if d.CLOSE_DEF_PERSON_ID == playerID {
			return d
		}
	}
	logrus.Warnf("Could not find defense stats for player %d", playerID)
	return PlayerDefenseStats{}
}

// https://github.com/swar/nba_api/blob/a5d90f5a8637f76bc8bdb60af94428ba04036f12/src/nba_api/stats/library/parameters.py#L1
