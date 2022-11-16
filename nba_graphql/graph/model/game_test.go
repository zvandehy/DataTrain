package model

import "testing"

func TestNewStat(t *testing.T) {
	type args struct {
		stat string
	}
	tests := []struct {
		name string
		args args
		want Stat
	}{
		{name: "lowercase", args: args{stat: "points"}, want: Points},
		{name: "uppercase", args: args{stat: "POINTS"}, want: Points},
		{name: "invalid", args: args{stat: "invalid"}, want: ""},
		{name: "Spaces", args: args{stat: "Fantasy Score"}, want: FantasyScore},
		{name: "3-PT Made", args: args{stat: "3-PT Made"}, want: ThreePointersMade},
		{name: "Blks+Stls", args: args{stat: "Blks+Stls"}, want: BlocksSteals},
		{name: "blks_stls", args: args{stat: "blks_stls"}, want: BlocksSteals},
		{name: "Pts+Rebs+Asts", args: args{stat: "Pts+Rebs+Asts"}, want: PointsReboundsAssists},
		{name: "Pts+Rebs", args: args{stat: "Pts+Rebs"}, want: PointsRebounds},
		{name: "pts_rebs", args: args{stat: "pts_rebs"}, want: PointsRebounds},
		{name: "Pts+Asts", args: args{stat: "Pts+Asts"}, want: PointsAssists},
		{name: "pts_asts", args: args{stat: "pts_asts"}, want: PointsAssists},
		{name: "Rebs+Asts", args: args{stat: "Rebs+Asts"}, want: ReboundsAssists},
		{name: "rebs_asts", args: args{stat: "rebs_asts"}, want: ReboundsAssists},
		{name: "Fantasy Score", args: args{stat: "Fantasy Score"}, want: FantasyScore},
		{name: "Free Throws Made", args: args{stat: "Free Throws Made"}, want: FreeThrowsMade},
		{name: "Turnovers", args: args{stat: "Turnovers"}, want: Turnovers},
		{name: "turnovers", args: args{stat: "turnovers"}, want: Turnovers},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewStat(tt.args.stat); got != tt.want {
				t.Errorf("NewStat(%v) = %v, want %v", tt.args.stat, got, tt.want)
			}
		})
	}
}
