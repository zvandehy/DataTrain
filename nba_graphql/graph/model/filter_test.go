package model

import (
	"reflect"
	"testing"

	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

func TestPrint(t *testing.T) {
	type args struct {
		Str            string
		StrPtr         *string
		Int            int
		IntPtr         *int
		StrSlice       []string
		StrPtrSlice    []*string
		StrSlicePtr    *[]string
		StrPtrSlicePtr *[]*string
		Obj            struct {
			//TODO: need to test unexported fields
			Foo string
		}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				Str: "hello world",
			},
			want: "{Str:hello world}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.Print(tt.args); got != tt.want {
				t.Errorf("Print() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerFilter_FilterPlayerStats(t *testing.T) {
	players := []*Player{
		{
			PlayerID:   1,
			GamesCache: []*PlayerGame{{Points: 10, Minutes: 10}, {Points: 20, Minutes: 20}},
		},
		{
			PlayerID:   2,
			GamesCache: []*PlayerGame{{Points: 10, Minutes: 5}, {Points: 10, Minutes: 5}},
		},
		{
			PlayerID:   3,
			GamesCache: []*PlayerGame{{Points: 10, Minutes: 5}, {Points: 20, Minutes: 10}},
		},
	}
	tests := []struct {
		name        string
		statFilters []*StatFilter
		wantIDs     []int
	}{
		{
			name: "PerGame OperatorGt",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorGt, Value: 10, Mode: StatModePerGame},
			},
			wantIDs: []int{1, 3},
		},
		{
			name: "PerGame OperatorGte",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorGte, Value: 10, Mode: StatModePerGame},
			},
			wantIDs: []int{1, 2, 3},
		},
		{
			name: "PerGame OperatorLt",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorLt, Value: 15, Mode: StatModePerGame},
			},
			wantIDs: []int{2},
		},
		{
			name: "PerGame OperatorLte",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorLte, Value: 15, Mode: StatModePerGame},
			},
			wantIDs: []int{1, 2, 3},
		},
		{
			name: "PerGame OperatorEq",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorEq, Value: 15, Mode: StatModePerGame},
			},
			wantIDs: []int{1, 3},
		},
		{
			name: "PerGame OperatorNeq",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorNeq, Value: 15, Mode: StatModePerGame},
			},
			wantIDs: []int{2},
		},
		{
			name: "PerMinute Gte",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorGte, Value: 2, Mode: StatModePerMinute},
			},
			wantIDs: []int{2, 3},
		},
		{
			name: "Total Gt",
			statFilters: []*StatFilter{
				{Stat: Points, Operator: OperatorGt, Value: 20, Mode: StatModeTotal},
			},
			wantIDs: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			filter := PlayerFilter{StatFilters: &tt.statFilters}
			got := filter.FilterPlayerStats(players, nil)
			gotIDs := []int{}
			for _, player := range got {
				gotIDs = append(gotIDs, player.PlayerID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Errorf("%v => PlayerFilter.FilterPlayerStats() = %v, want %v", tt.name, gotIDs, tt.wantIDs)
			}
		})
	}

}
