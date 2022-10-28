package model

type Schedule struct {
	LeagueSchedule struct {
		GameDates []struct {
			Games []struct {
				GameDateEst string `json:"gameDateEst"`
				AwayTeam    struct {
					TeamTriCode string `json:"triCode"`
				} `json:"awayTeam"`
				HomeTeam struct {
					TeamTriCode string `json:"triCode"`
				} `json:"homeTeam"`
				GameID string `json:"gameId"`
			} `json:"games"`
		} `json:"gameDates"`
	} `json:"leagueSchedule"`
}
