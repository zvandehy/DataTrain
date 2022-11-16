package model

import "go.mongodb.org/mongo-driver/bson"

func (input *ProjectionFilter) MongoFilter() bson.M {
	filter := input.Period.MongoFilter()
	if input.PropositionFilter != nil {
		if input.PropositionFilter.Sportsbook != nil {
			filter["propositions.sportsbook"] = input.PropositionFilter.Sportsbook
		}
		// TODO: Implement TypeFilter
		// if input.PropositionFilter.PropositionType != nil {
		// 	filter["propositions.type"] = input.PropositionFilter.PropositionType
		// }
	}
	return filter
}
