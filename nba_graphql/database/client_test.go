package database

import (
	"fmt"
	"testing"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateFilter(t *testing.T) {
	type args struct {
		filters []model.GameFilter
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "single complete filter",
			args: args{
				filters: []model.GameFilter{
					{
						TeamID:     &[]int{1}[0],
						OpponentID: &[]int{2}[0],
						PlayerID:   &[]int{3}[0],
						GameID:     &[]string{"4"}[0],
						Season:     &[]string{"2022-23"}[0],
						StartDate:  &[]string{"01-01-2022"}[0],
						EndDate:    &[]string{"12-12-2023"}[0],
					},
				},
			},
			want: bson.M{
				"$or": bson.A{bson.M{
					"teamID":   1,
					"opponent": 2,
					"playerID": 3,
					"gameID":   4,
					"season":   "2022-23",
					"date":     bson.M{"$gte": "01-01-2022", "$lt": "12-12-2023"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createGameFilter(tt.args.filters)
			// TODO: compare with deep equal
			gotStr := fmt.Sprintf("%v", got)
			wantStr := fmt.Sprintf("%v", tt.want)
			if gotStr != wantStr {
				t.Errorf("got:\n%v\nbut want:\n%v", gotStr, wantStr)
			}
		})
	}
}
