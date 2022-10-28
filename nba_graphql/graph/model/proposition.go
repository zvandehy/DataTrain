package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Proposition struct {
	Sportsbook    SportsbookOption    `json:"sportsbook" bson:"sportsbook"`
	Target        float64             `json:"target" bson:"target"`
	TypeRaw       string              `json:"type" bson:"type"`
	Type          Stat                `json:"propType" bson:"propType"`
	LastModified  *time.Time          `json:"lastModified" bson:"lastModified"`
	ProjectionRef *Projection         `json:"projectionRef" bson:"projectionRef"`
	Analysis      *PropositionSummary `json:"analysis" bson:"analysis"`
}

func (p *Proposition) UnmarshalBSON(data []byte) error {
	type Alias Proposition
	bson.Unmarshal(data, (*Alias)(p))
	t, err := NewStat(p.TypeRaw)
	if err != nil {
		return err
	}
	p.Type = t
	return nil
}

func (p *Proposition) UnmarshalJSON(data []byte) error {
	type Alias Proposition
	json.Unmarshal(data, (*Alias)(p))
	t, err := NewStat(p.TypeRaw)
	if err != nil {
		return err
	}
	p.Type = t
	return nil
}

func (p *Proposition) MarshalJSON() ([]byte, error) {
	type Alias Proposition
	p.TypeRaw = string(p.Type)
	return json.Marshal((*Alias)(p))
}

func (p *Proposition) Match(propositionFilter PropositionFilter) bool {
	if propositionFilter.Sportsbook != nil && *propositionFilter.Sportsbook != p.Sportsbook {
		return false
	}
	if propositionFilter.PropositionType != nil && *propositionFilter.PropositionType != p.Type {
		return false
	}
	return true
}
