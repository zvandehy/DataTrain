package model

import (
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

func (f TeamFilter) String() string {
	return util.Print(f)
}

func (o *Operator) Evaluate(left, right float64) bool {
	switch *o {
	case OperatorEq:
		return left == right
	case OperatorNeq:
		return left != right
	case OperatorGt:
		return left > right
	case OperatorLt:
		return left < right
	case OperatorGte:
		return left >= right
	case OperatorLte:
		return left <= right
	}
	return false
}
