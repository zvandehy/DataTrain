package model

import "go.mongodb.org/mongo-driver/bson"

func (f *TeamFilter) MongoFilter() bson.M {
	filter := bson.M{}
	if f.TeamID != nil {
		filter["team_id"] = *f.TeamID
	}
	if f.Name != nil {
		filter["name"] = *f.Name
	}
	if f.Abbreviation != nil {
		filter["abbreviation"] = *f.Abbreviation
	}
	// if f.TeamCity != nil {
	// 	filter["team_city"] = *f.TeamCity
	// }
	// if f.TeamConference != nil {
	// 	filter["team_conference"] = *f.TeamConference
	// }
	// if f.TeamDivision != nil {
	// 	filter["team_division"] = *f.TeamDivision
	// }
	return filter
}
