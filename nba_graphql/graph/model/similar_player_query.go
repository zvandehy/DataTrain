package model

import "time"

type SimilarPlayerQuery struct {
	SimilarPlayerInput SimilarPlayerInput `json:"similarPlayerInput"`
	ToPlayerID         int                `json:"toPlayerID"`
	EndDate            *time.Time         `json:"endDate"`
}
