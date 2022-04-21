package model

type Player struct {
	FirstName   string   `json:"first_name" bson:"first_name"`
	LastName    string   `json:"last_name" bson:"last_name"`
	PlayerID    int      `json:"playerID" bson:"playerID"`
	Seasons     []string `json:"seasons" bson:"seasons"`
	Position    string   `json:"position" bson:"position"`
	CurrentTeam string   `json:"currentTeam" bson:"teamABR"`
	Height      string   `json:"height" bson:"height"`
	Weight      int      `json:"weight" bson:"weight"`
}

func (p Player) String() string {
	return Print(p)
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
	Margin                       int     `json:"margin" bson:"margin"`
	Minutes                      string  `json:"minutes" bson:"minutes"`
	OffensiveReboundPercentage   float64 `json:"offensive_rebound_percentage" bson:"offensive_rebound_percentage"`
	OffensiveRebounds            int     `json:"offensive_rebounds" bson:"offensive_rebounds"`
	TeamID                       int     `json:"team" bson:"team"`
	OpponentID                   int     `json:"opponent" bson:"opponent"`
	PersonalFoulsDrawn           int     `json:"personal_fouls_drawn" bson:"personal_fouls_drawn"`
	PersonalFouls                int     `json:"personal_fouls" bson:"personal_fouls"`
	Points                       int     `json:"points" bson:"points"`
	PlayerID                     int     `json:"player" bson:"player"`
	Playoffs                     bool    `json:"playoffs" bson:"playoffs"`
	Season                       string  `json:"season" bson:"season"`
	ThreePointPercentage         float64 `json:"three_point_percentage" bson:"three_point_percentage"`
	ThreePointersAttempted       int     `json:"three_pointers_attempted" bson:"three_pointers_attempted"`
	ThreePointersMade            int     `json:"three_pointers_made" bson:"three_pointers_made"`
	TotalRebounds                int     `json:"total_rebounds" bson:"total_rebounds"`
	TrueShootingPercentage       float64 `json:"true_shooting_percentage" bson:"true_shooting_percentage"`
	Turnovers                    int     `json:"turnovers" bson:"turnovers"`
	Blocks                       int     `json:"blocks" bson:"blocks"`
	Steals                       int     `json:"steals" bson:"steals"`
	Usage                        float64 `json:"usage" bson:"usage"`
	WinOrLoss                    string  `json:"win_or_loss" bson:"win_or_loss"`
}

func (p PlayerGame) String() string {
	return Print(p)
}
