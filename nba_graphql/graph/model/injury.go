package model

//TODO: Might eventually want team IDs and/or maybe missed gameIDs
type Injury struct {
	Status     string `json:"status" bson:"status"`
	StartDate  string `json:"start_date" bson:"start_date"`
	ReturnDate string `json:"return_date" bson:"return_date"`
	PlayerID   int    `json:"playerID" bson:"playerID"`
}
