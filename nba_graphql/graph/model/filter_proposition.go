package model

// TODO: Proposition Type is in format of "PointsAssists" instead of "points_assists", which is inconsistent with AverageStats/Games

func (f *PropositionFilter) UnmarshalJSON(data []byte) error {
	s, err := NewStat(string(data[:]))
	if err != nil {
		return err
	}
	*f.PropositionType = s
	return nil
}

func (f *PropositionFilter) UnmarshalBSON(data []byte) error {
	s, err := NewStat(string(data[:]))
	if err != nil {
		return err
	}
	*f.PropositionType = s
	return nil
}
