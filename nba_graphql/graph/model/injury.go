package model

type Injury struct {
	Description string `json:"description" bson:"description"`
	StartDate   string `json:"start_date" bson:"start_date"`
	ReturnDate  string `json:"return_date" bson:"return_date"`
	PlayerID    int    `json:"playerID" bson:"playerID"`
}

var Injury2 = Injury{
	Description: "Ankle",
	StartDate:   "4/10/2022",
	ReturnDate:  "4/25/2022",
	PlayerID:    200746,
}
var Injury3 = Injury{
	Description: "Elbow",
	StartDate:   "4/08/2022",
	ReturnDate:  "4/31/2022",
	PlayerID:    2546,
}
var Injury4 = Injury{
	Description: "Elbow",
	StartDate:   "4/08/2022",
	ReturnDate:  "4/31/2022",
	PlayerID:    203507,
}
var Injury5 = Injury{
	Description: "Wrist",
	StartDate:   "4/01/2022",
	ReturnDate:  "5/01/2022",
	PlayerID:    1630173,
}
