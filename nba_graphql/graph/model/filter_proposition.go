package model

func (f *PropositionFilter) UnmarshalJSON(data []byte) error {
	*f.PropositionType = NewStat(string(data[:]))
	return nil
}

func (f *PropositionFilter) UnmarshalBSON(data []byte) error {
	*f.PropositionType = NewStat(string(data[:]))
	return nil
}

// // MarshalGQL implements the graphql.Marshaler interface
// func (f PropositionFilter) MarshalGQL(w io.Writer) {

// 	io.WriteString(w, string(f))
// }

// func (f *PropositionFilter) UnmarshalGQL(v interface{}) error {
// 	season, ok := v.(string)
// 	if !ok {
// 		return fmt.Errorf("Propo must be a string")
// 	}

// 	*f.PropositionType = NewStat(season)
// 	return nil
// }
