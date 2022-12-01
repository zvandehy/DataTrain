package model

import "math"

type StandardizedPlayerStats struct {
	Name                   string  `db:"name"`
	Id                     int     `db:"playerID"`
	NGames                 int     `db:"games"`
	Assists                float64 `db:"ZSCORE_assists"`
	DefensiveRebounds      float64 `db:"ZSCORE_defensiveRebounds"`
	OffensiveRebounds      float64 `db:"ZSCORE_offensiveRebounds"`
	FieldGoalsAttempted    float64 `db:"ZSCORE_fieldGoalsAttempted"`
	FieldGoalsMade         float64 `db:"ZSCORE_fieldGoalsMade"`
	FreeThrowsAttempted    float64 `db:"ZSCORE_freeThrowsAttempted"`
	FreeThrowsMade         float64 `db:"ZSCORE_freeThrowsMade"`
	PersonalFoulsDrawn     float64 `db:"ZSCORE_personalFoulsDrawn"`
	PersonalFouls          float64 `db:"ZSCORE_personalFouls"`
	Points                 float64 `db:"ZSCORE_points"`
	ThreePointersAttempted float64 `db:"ZSCORE_threePointersAttempted"`
	ThreePointersMade      float64 `db:"ZSCORE_threePointersMade"`
	Rebounds               float64 `db:"ZSCORE_rebounds"`
	Turnovers              float64 `db:"ZSCORE_turnovers"`
	Blocks                 float64 `db:"ZSCORE_blocks"`
	Steals                 float64 `db:"ZSCORE_steals"`
	PotentialAssists       float64 `db:"ZSCORE_potentialAssists"`
	Passes                 float64 `db:"ZSCORE_passes"`
	Minutes                float64 `db:"ZSCORE_minutes"`
	HeightInches           float64 `db:"ZSCORE_heightInches"`
	Weight                 float64 `db:"ZSCORE_weight"`
}

func StandardizedPlayerStatsFields() []string {
	return []string{
		"assists",
		"defensiveRebounds",
		"offensiveRebounds",
		"fieldGoalsAttempted",
		"fieldGoalsMade",
		"freeThrowsAttempted",
		"freeThrowsMade",
		"personalFoulsDrawn",
		"personalFouls",
		"points",
		"threePointersAttempted",
		"threePointersMade",
		"rebounds",
		"turnovers",
		"blocks",
		"steals",
		"potentialAssists",
		"passes",
		"minutes",
		"heightInches",
		"weight",
	}
}

func (s StandardizedPlayerStats) Field(field string) float64 {
	switch field {
	case "assists":
		return s.Assists
	case "defensiveRebounds":
		return s.DefensiveRebounds
	case "offensiveRebounds":
		return s.OffensiveRebounds
	case "fieldGoalsAttempted":
		return s.FieldGoalsAttempted
	case "fieldGoalsMade":
		return s.FieldGoalsMade
	case "freeThrowsAttempted":
		return s.FreeThrowsAttempted
	case "freeThrowsMade":
		return s.FreeThrowsMade
	case "personalFoulsDrawn":
		return s.PersonalFoulsDrawn
	case "personalFouls":
		return s.PersonalFouls
	case "points":
		return s.Points
	case "threePointersAttempted":
		return s.ThreePointersAttempted
	case "threePointersMade":
		return s.ThreePointersMade
	case "rebounds":
		return s.Rebounds
	case "turnovers":
		return s.Turnovers
	case "blocks":
		return s.Blocks
	case "steals":
		return s.Steals
	case "potentialAssists":
		return s.PotentialAssists
	case "passes":
		return s.Passes
	case "minutes":
		return s.Minutes
	case "heightInches":
		return s.HeightInches
	case "weight":
		return s.Weight
	}
	return -1
}

func (a *StandardizedPlayerStats) CosineSimilarityTo(b StandardizedPlayerStats) float64 {
	vectorA, vectorB := []float64{}, []float64{}
	for _, field := range StandardizedPlayerStatsFields() {
		vectorA = append(vectorA, a.Field(field))
		vectorB = append(vectorB, b.Field(field))
	}
	return CosineSimilarity(vectorA, vectorB)
}

func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return -1
	}
	var dotProduct, magnitudeA, magnitudeB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += math.Pow(a[i], 2)
		magnitudeB += math.Pow(b[i], 2)
	}
	return math.Round((dotProduct/(math.Sqrt(magnitudeA)*math.Sqrt(magnitudeB)))*100.0) / 100.0
}
