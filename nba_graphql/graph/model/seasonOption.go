package model

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

type SeasonOption string

func (s *SeasonOption) UnmarshalJSON(data []byte) error {
	v := string(data[:])
	*s = NewSeasonOption(v)
	return nil
}

func (s *SeasonOption) UnmarshalBSON(data []byte) error {
	v := string(data[:])
	*s = NewSeasonOption(v)
	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (o *SeasonOption) UnmarshalGQL(v interface{}) error {
	season, ok := v.(string)
	if !ok {
		return fmt.Errorf("SeasonOption must be a string")
	}

	*o = NewSeasonOption(season)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (o SeasonOption) MarshalGQL(w io.Writer) {
	switch o {
	case SEASON_2022_23:
		io.WriteString(w, `"2022-23"`)
	case SEASON_2021_22:
		io.WriteString(w, `"2021-22"`)
	case SEASON_2020_21:
		io.WriteString(w, `"2020-21"`)
	default:
		logrus.Warnf("SeasonOption.MarshalGQL: unknown season option %s", o)
		io.WriteString(w, string(o))
	}
}

func NewSeasonOption(s string) SeasonOption {
	s = util.ClearString(s)
	switch s[len(s)-2:] {
	case "23":
		return SEASON_2022_23
	case "22":
		return SEASON_2021_22
	case "21":
		return SEASON_2020_21
	default:
		logrus.Warnf("Invalid season option: \"%s\" => \"%s\"", s, s[len(s)-2:])
		return SeasonOption(s)
	}
}

const (
	SEASON_2022_23 SeasonOption = "2022-23"
	SEASON_2021_22 SeasonOption = "2021-22"
	SEASON_2020_21 SeasonOption = "2020-21"
)
