package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"gonum.org/v1/gonum/stat/distuv"
)

type DBProposition struct {
	PlayerID     int        `db:"playerID"`
	GameID       string     `db:"gameID"`
	OpponentID   int        `db:"opponentID"`
	PlayerName   string     `db:"playerName"`
	StatType     string     `db:"statType"`
	Target       float64    `db:"target"`
	Sportsbook   string     `db:"sportsbook"`
	LastModified *time.Time `db:"lastModified"`
	CreatedAt    *time.Time `db:"CreatedAt"`
	UpdatedAt    *time.Time `db:"UpdatedAt"`
}

type Proposition struct {
	PlayerID     int              `db:"playerID"`
	GameID       string           `db:"gameID"`
	OpponentID   int              `db:"opponentID"`
	PlayerName   string           `db:"playerName"`
	Sportsbook   SportsbookOption `db:"sportsbook" json:"sportsbook" bson:"sportsbook"`
	Target       float64          `db:"target" json:"target" bson:"target"`
	TypeRaw      string           `db:"statType" json:"type" bson:"type"`
	Type         Stat             `json:"propType" bson:"propType"`
	LastModified *time.Time       `json:"lastModified" bson:"lastModified" db:"lastModified"`
	Outcome      PropOutcome      `json:"outcome" bson:"outcome"`
	ActualResult *float64         `json:"actualResult" bson:"actualResult"`
	Accuracy     float64          `json:"accuracy" bson:"accuracy"`
	Game         *PlayerGame      `json:"game" bson:"game"`
	CreatedAt    *time.Time       `json:"createdAt" bson:"createdAt" db:"CreatedAt"`
	UpdatedAt    *time.Time       `json:"updatedAt" bson:"updatedAt" db:"UpdatedAt"`
}

type StatDistribution struct {
	StatType Stat    `json:"statType" bson:"statType"`
	Mean     float64 `db:"mean" json:"mean" bson:"mean"`
	StdDev   float64 `db:"stdDev" json:"stdDev" bson:"stdDev"`
}

func PValue(estimation float64, stddev float64, target float64) (float64, Wager) {
	normal := distuv.Normal{Mu: estimation, Sigma: stddev, Src: nil}
	p := normal.CDF(target)
	if p > 0.5 {
		return p, WagerUnder
	}
	return 1 - p, WagerOver
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
