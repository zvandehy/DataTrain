package model

type Team struct {
	Name         string `json:"name" bson:"name"`
	TeamID       int    `json:"teamID" bson:"teamID"`
	Abbreviation string `json:"abbreviation" bson:"abbreviation"`
	Location     string `json:"location" bson:"city"`
	NumWins      int    `json:"numWins" bson:"numWins"`
	NumLoss      int    `json:"numLoss" bson:"numLoss"`
}

func (t Team) String() string {
	return Print(t)
}

type TeamGame struct {
	Date                                 string  `json:"date" bson:"date"`
	DefensiveRating                      float64 `json:"defensive_rating" bson:"defensive_rating"`
	DefensiveReboundPercentage           float64 `json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	FieldGoalPercentage                  float64 `json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted                  int     `json:"field_goals_attempted" bson:"field_goals_attempted"`
	GameID                               string  `json:"gameID" bson:"gameID"`
	HomeOrAway                           string  `json:"home_or_away" bson:"home_or_away"`
	OffensiveReboundPercentage           float64 `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OpponentID                           int     `json:"opponent" bson:"opponent"`
	OpponentAssists                      int     `json:"opponent_assists" bson:"opponent_assists"`
	OpponentEffectiveFieldGoalPercentage float64 `json:"opponent_effective_field_goal_percentage" bson:"opponent_effective_field_goal_percentage"`
	OpponentFieldGoalsAttempted          int     `json:"opponent_field_goals_attempted" bson:"opponent_field_goals_attempted"`
	OpponentFreeThrowsAttempted          int     `json:"opponent_free_throws_attempted" bson:"opponent_free_throws_attempted"`
	OpponentPoints                       int     `json:"opponent_points" bson:"opponent_points"`
	OpponentRebounds                     int     `json:"opponent_rebounds" bson:"opponent_rebounds"`
	OpponentThreePointersAttempted       int     `json:"opponent_three_pointers_attempted" bson:"opponent_three_pointers_attempted"`
	OpponentThreePointersMade            int     `json:"opponent_three_pointers_made" bson:"opponent_three_pointers_made"`
	PlusMinusPerHundred                  float64 `json:"plus_minus_per_hundred" bson:"plus_minus_per_hundred"`
	Possessions                          int     `json:"possessions" bson:"possessions"`
	PersonalFouls                        int     `json:"personal_fouls" bson:"personal_fouls"`
	PersonalFoulsDrawn                   int     `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	Season                               string  `json:"season" bson:"season"`
	WinOrLoss                            string  `json:"win_or_loss" bson:"win_or_loss"`
}

func (t TeamGame) String() string {
	return Print(t)
}

type PlayersInGame struct {
	TeamPlayers     []*Player `json:"team" bson:"team"`
	OpponentPlayers []*Player `json:"opponent" bson:"opponent"`
}
