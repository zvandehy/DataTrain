package model

type Schedule struct {
	LeagueSchedule struct {
		GameDates []struct {
			Games []struct {
				GameDateEst string `json:"gameDateEst"`
				AwayTeam    struct {
					TeamTriCode string `json:"teamTriCode"`
				} `json:"awayTeam"`
				HomeTeam struct {
					TeamTriCode string `json:"teamTriCode"`
				} `json:"homeTeam"`
				GameID string `json:"gameId"`
			} `json:"games"`
		} `json:"gameDates"`
	} `json:"leagueSchedule"`
}
