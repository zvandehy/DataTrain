package database

// func TestCreateFilter(t *testing.T) {
// 	type args struct {
// 		filter model.GameFilter
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bson.M
// 	}{
// 		{
// 			name: "complete filter",
// 			args: args{
// 				filter: model.GameFilter{
// 					TeamID:     &[]int{1}[0],
// 					OpponentID: &[]int{2}[0],
// 					PlayerID:   &[]int{3}[0],
// 					GameID:     &[]string{"4"}[0],
// 					Seasons:    &[][]model.SeasonOption{{model.SEASON_2022_23}}[0],
// 					StartDate:  &[]string{"01-01-2022"}[0],
// 					EndDate:    &[]string{"12-12-2023"}[0],
// 				},
// 			},
// 			want: bson.M{
// 				"teamID":   1,
// 				"opponent": 2,
// 				"playerID": 3,
// 				"gameID":   4,
// 				"season":   "2022-23",
// 				"date":     bson.M{"$gte": "01-01-2022", "$lt": "12-12-2023"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := createGameFilter(tt.args.filter)
// 			// TODO: compare with deep equal
// 			gotStr := fmt.Sprintf("%v", got)
// 			wantStr := fmt.Sprintf("%v", tt.want)
// 			if gotStr != wantStr {
// 				t.Errorf("got:\n%v\nbut want:\n%v", gotStr, wantStr)
// 			}
// 		})
// 	}
// }
