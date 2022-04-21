package model

type Injury struct {
	Injury     string `json:"injury_name" bson:"injury_name"`
	InjuryDate string `json:"injury_date" bson:"injury_date"`
	ReturnDate string `json:"return_date" bson:"return_date"`
	PlayerID   int    `json:"playerID" bson:"playerID"`
}
