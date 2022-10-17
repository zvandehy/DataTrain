package model

// TODO: Proposition Type is in format of "PointsAssists" instead of "points_assists", which is inconsistent with AverageStats/Games

func (f *PropositionFilter) UnmarshalJSON(data []byte) error {
	*f.PropositionType = NewStat(string(data[:]))
	return nil
}

func (f *PropositionFilter) UnmarshalBSON(data []byte) error {
	*f.PropositionType = NewStat(string(data[:]))
	return nil
}
