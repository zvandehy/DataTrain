package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type DBProposition struct {
	PlayerID     int        `db:"playerID"`
	GameID       string     `db:"gameID"`
	OpponentID   int        `db:"opponentID"`
	PlayerName   string     `db:"playerName"`
	StartTime    *time.Time `db:"startTime"`
	StatType     string     `db:"statType"`
	Target       float64    `db:"target"`
	Sportsbook   string     `db:"sportsbook"`
	LastModified *time.Time `db:"lastModified"`
	CreatedAt    *time.Time `db:"CreatedAt"`
	UpdatedAt    *time.Time `db:"UpdatedAt"`
}

type Proposition struct {
	Sportsbook   SportsbookOption `db:"sportsbook" json:"sportsbook" bson:"sportsbook"`
	Target       float64          `db:"target" json:"target" bson:"target"`
	TypeRaw      string           `db:"statType" json:"type" bson:"type"`
	Type         Stat             `json:"propType" bson:"propType"`
	LastModified *time.Time       `json:"lastModified" bson:"lastModified"`
	Outcome      PropOutcome      `json:"outcome" bson:"outcome"`
	ActualResult *float64         `json:"actualResult" bson:"actualResult"`
	Accuracy     float64          `json:"accuracy" bson:"accuracy"`
	Game         *PlayerGame      `json:"game" bson:"game"`
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
