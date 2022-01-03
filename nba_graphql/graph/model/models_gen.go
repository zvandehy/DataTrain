// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Player struct {
	ID        string        `json:"id" bson:"_id"`
	PlayerID  int           `json:"playerID" bson:"playerID"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	TeamAbr   string        `json:"teamABR" bson:"teamABR"`
	Seasons   []string      `json:"seasons" bson:"seasons"`
	Games     []*PlayerGame `json:"games" bson:"games"`
}

type PlayerGame struct {
	AssistPercentage             float64 `json:"assist_percentage" bson:"assist_percentage"`
	Assists                      int     `json:"assists" bson:"assists"`
	Date                         string  `json:"date" bson:"date"`
	DefensiveReboundPercentage   float64 `json:"defensive_rebound_percentage" bson:"defensive_rebound_percentage"`
	DefensiveRebounds            int     `json:"defensive_rebounds" bson:"defensive_rebounds"`
	EffectiveFieldGoalPercentage float64 `json:"effective_field_goal_percentage" bson:"effective_field_goal_percentage"`
	FieldGoalPercentage          float64 `json:"field_goal_percentage" bson:"field_goal_percentage"`
	FieldGoalsAttempted          int     `json:"field_goals_attempted" bson:"field_goals_attempted"`
	FieldGoalsMade               int     `json:"field_goals_made" bson:"field_goals_made"`
	FreeThrowsAttempted          int     `json:"free_throws_attempted" bson:"free_throws_attempted"`
	FreeThrowsMade               int     `json:"free_throws_made" bson:"free_throws_made"`
	FreeThrowsPercentage         float64 `json:"free_throws_percentage" bson:"free_throws_percentage"`
	GameID                       string  `json:"gameID" bson:"gameID"`
	HomeOrAway                   string  `json:"home_or_away" bson:"home_or_away"`
	Minutes                      string  `json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   float64 `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            int     `json:"offensive_rebounds" bson:"offensive_rebounds"`
	Opponent                     int     `json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           int     `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                int     `json:"personal_fouls" bson:"personal_fouls"`
	Points                       int     `json:"points" bson:"points"`
	Season                       string  `json:"season" bson:"season"`
	ThreePointPercentage         float64 `json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       int     `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            int     `json:"three_pointers_made" bson:"three_pointers_made"`
	TotalRebounds                int     `json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       float64 `json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    int     `json:"turnovers" bson:"turnovers"`
	Usage                        float64 `json:"usage" bson:"usage"`
}

type Team struct {
	ID           string      `json:"id" bson:"_id"`
	TeamID       int         `json:"teamID" bson:"teamID"`
	Abbreviation string      `json:"abbreviation" bson:"abbreviation"`
	City         string      `json:"city" bson:"city"`
	Name         string      `json:"name" bson:"name"`
	Games        []*TeamGame `json:"games" bson:"games"`
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
	Opponent                             int     `json:"opponent" bson:"opponent"`
	TeamID                               string  `json:"teamID" bson:"teamID"`
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
