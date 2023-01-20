package nba

import "github.com/sirupsen/logrus"

const LeagueBiosURL = "leaguedashplayerbiostats?"

// LeagueGame is a struct that contains all the fields for the LeagueGame endpoint.
type CommonPlayer struct {
	League               string  `json:"league"`
	PLAYER_ID            int     `json:"PLAYER_ID"`
	PLAYER_NAME          string  `json:"PLAYER_NAME"`
	TEAM_ID              int     `json:"TEAM_ID"`
	TEAM_ABBREVIATION    string  `json:"TEAM_ABBREVIATION"`
	AGE                  float64 `json:"AGE"`
	PLAYER_HEIGHT        string  `json:"PLAYER_HEIGHT"`
	PLAYER_HEIGHT_INCHES float64 `json:"PLAYER_HEIGHT_INCHES"`
	PLAYER_WEIGHT        string  `json:"PLAYER_WEIGHT"`
	COLLEGE              string  `json:"COLLEGE"`
	COUNTRY              string  `json:"COUNTRY"`
	DRAFT_YEAR           string  `json:"DRAFT_YEAR"`
	DRAFT_ROUND          string  `json:"DRAFT_ROUND"`
	DRAFT_NUMBER         string  `json:"DRAFT_NUMBER"`
	GP                   float64 `json:"GP"`
	PTS                  float64 `json:"PTS"`
	REB                  float64 `json:"REB"`
	AST                  float64 `json:"AST"`
	NET_RATING           float64 `json:"NET_RATING"`
	OREB_PCT             float64 `json:"OREB_PCT"`
	DREB_PCT             float64 `json:"DREB_PCT"`
	USG_PCT              float64 `json:"USG_PCT"`
	TS_PCT               float64 `json:"TS_PCT"`
	AST_PCT              float64 `json:"AST_PCT"`
}

func NBAPlayers() (*NBAResponse[CommonPlayer], error) {
	params := CommonPlayerInfoParams{
		LeagueID:       NBA_LEAGUE_ID,
		LastNGames:     "0",
		Month:          "0",
		OpponentTeamID: "0",
		PORound:        "0",
		PerMode:        "PerGame",
		Period:         "0",
		Season:         "2022-23",
		SeasonType:     "Regular Season",
	}
	resp, err := NBAQuery[CommonPlayer](LeagueBiosURL, ParamMap(params))
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(resp.ResultSets[0].RowSet); i++ {
		resp.ResultSets[0].RowSet[i].League = "NBA"
	}
	return resp, nil
}

type CommonPlayerInfoParams struct {
	College          string `json:"College"`
	Conference       string `json:"Conference"`
	Country          string `json:"Country"`
	DateFrom         string `json:"DateFrom"`
	DateTo           string `json:"DateTo"`
	Division         string `json:"Division"`
	DraftPick        string `json:"DraftPick"`
	DraftYear        string `json:"DraftYear"`
	GameScope        string `json:"GameScope"`
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
	PerGame          string `json:"PerGame"`
	Period           string `json:"Period"`
	PlayerExperience string `json:"PlayerExperience"`
	PlayerPosition   string `json:"PlayerPosition"`
	Season           string `json:"Season"`
	SeasonSegment    string `json:"SeasonSegment"`
	SeasonType       string `json:"SeasonType"`
	ShotClockRange   string `json:"ShotClockRange"`
	StarterBench     string `json:"StarterBench"`
	TeamID0          string `json:"TeamID0"`
	VsConference     string `json:"VsConference"`
	VsDivision       string `json:"VsDivision"`
	Weight           string `json:"Weight"`
}

func FindLeagueBioByID(leagueBios []CommonPlayer, playerID int) CommonPlayer {
	for _, leagueBio := range leagueBios {
		if leagueBio.PLAYER_ID == playerID {
			return leagueBio
		}
	}
	logrus.Warnf("Could not find player bio with ID %d", playerID)
	return CommonPlayer{}
}
