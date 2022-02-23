package model

type Player struct {
	Name        string   `json:"name" bson:"first_name"`
	PlayerID    int      `json:"playerID" bson:"playerID"`
	Seasons     []string `json:"seasons" bson:"seasons"`
	Position    string   `json:"position" bson:"position"`
	CurrentTeam string   `json:"currentTeam" bson:"teamABR"`
}
