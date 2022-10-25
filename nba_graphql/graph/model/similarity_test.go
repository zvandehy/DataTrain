package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPlayerSimilarityMatrix_AddPlayer(t *testing.T) {
	type fields struct {
		Matrix map[int]SimilarityVector
	}
	type args struct {
		playerID      int
		playerAverage PlayerAverage
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       error
		wantMatrix map[int]SimilarityVector
	}{
		{
			name: "Add player to empty matrix",
			fields: fields{
				Matrix: make(map[int]SimilarityVector),
			},
			args: args{
				playerID: 1,
				playerAverage: PlayerAverage{
					Points:   1,
					Assists:  1,
					Rebounds: 1,
				},
			},
			want: nil,
			wantMatrix: map[int]SimilarityVector{
				1: {
					Comparisons: make(map[int]PlayerDiff),
					Average: PlayerAverage{
						Points:   1,
						Assists:  1,
						Rebounds: 1,
					},
				},
			},
		},
		{
			name: "Add player to matrix with two players",
			fields: fields{
				Matrix: map[int]SimilarityVector{
					1: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							2: PlayerDiff{
								Points:   1,
								Assists:  1,
								Rebounds: 1,
							},
						},
						Average: PlayerAverage{
							Points:   10,
							Assists:  5,
							Rebounds: 6,
						},
					},
					2: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points:   -1,
								Assists:  -1,
								Rebounds: -1,
							},
						},
						Average: PlayerAverage{
							Points:   9,
							Assists:  4,
							Rebounds: 5,
						},
					},
				},
			},
			args: args{
				playerID: 3,
				playerAverage: PlayerAverage{
					Points:   11,
					Assists:  4,
					Rebounds: 6,
				},
			},
			want: nil,
			wantMatrix: map[int]SimilarityVector{
				// Player `1` should have a new comparison with player `3`
				1: SimilarityVector{
					Comparisons: map[int]PlayerDiff{
						// Player 1 averages +1 points, +1 assists, and +1 rebounds from player 2
						2: PlayerDiff{
							Points:   1,
							Assists:  1,
							Rebounds: 1,
						},
						// Player 1 averages -1 points, +1 assists, and +0 rebounds from player 3
						3: PlayerDiff{
							Points:   -1,
							Assists:  1,
							Rebounds: 0,
						},
					},
					// Player 1 Averages 10 points, 5 assists, and 6 rebounds
					Average: PlayerAverage{
						Points:   10,
						Assists:  5,
						Rebounds: 6,
					},
				},
				2: SimilarityVector{
					Comparisons: map[int]PlayerDiff{
						1: PlayerDiff{
							Points:   -1,
							Assists:  -1,
							Rebounds: -1,
						},
						3: PlayerDiff{
							Points:   -2,
							Assists:  0,
							Rebounds: -1,
						},
					},
					Average: PlayerAverage{
						Points:   9,
						Assists:  4,
						Rebounds: 5,
					},
				},
				3: SimilarityVector{
					Comparisons: map[int]PlayerDiff{
						1: PlayerDiff{
							Points:   1,
							Assists:  -1,
							Rebounds: 0,
						},
						2: PlayerDiff{
							Points:   2,
							Assists:  0,
							Rebounds: 1,
						},
					},
					Average: PlayerAverage{
						Points:   11,
						Assists:  4,
						Rebounds: 6,
					},
				},
			},
		},
		{
			name: "Add player to matrix, but player already exists",
			fields: fields{
				Matrix: map[int]SimilarityVector{
					1: SimilarityVector{
						Comparisons: map[int]PlayerDiff{},
						Average: PlayerAverage{
							Points:   10,
							Assists:  5,
							Rebounds: 6,
						},
					},
				},
			},
			args: args{
				playerID: 1,
				playerAverage: PlayerAverage{
					Points:   11,
					Assists:  4,
					Rebounds: 6,
				},
			},
			want: fmt.Errorf("player 1 already exists in matrix"),
			wantMatrix: map[int]SimilarityVector{
				1: SimilarityVector{
					Comparisons: map[int]PlayerDiff{},
					Average: PlayerAverage{
						Points:   10,
						Assists:  5,
						Rebounds: 6,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PlayerSimilarityMatrix{
				Matrix: tt.fields.Matrix,
			}
			got := m.AddPlayer(tt.args.playerID, tt.args.playerAverage)
			if (tt.want != nil && got != nil && tt.want.Error() != got.Error()) || (tt.want == nil && got != nil) || (tt.want != nil && got == nil) {
				t.Errorf("PlayerSimilarityMatrix.AddPlayer() = \n%v, want \n%v", got, tt.want)
			}
			if !reflect.DeepEqual(m.Matrix, tt.wantMatrix) {
				t.Errorf("PlayerSimilarityMatrix.AddPlayer() = %v, want %v", m.Matrix, tt.wantMatrix)
			}
		})
	}
}

func TestPlayerSimilarityMatrix_AddPlayers(t *testing.T) {
	type fields struct {
		Matrix map[int]SimilarityVector
	}
	type args struct {
		players []PlayerAverage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PlayerSimilarityMatrix{
				Matrix: tt.fields.Matrix,
			}
			if err := m.AddPlayers(tt.args.players); (err != nil) != tt.wantErr {
				t.Errorf("PlayerSimilarityMatrix.AddPlayers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStatOfInterest_Add(t *testing.T) {
	type fields struct {
		Scores []float64
		Stat   Stat
		Mean   float64
		StdDev float64
	}
	type args struct {
		playerAverage PlayerAverage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatOfInterest{
				Scores: tt.fields.Scores,
				Stat:   tt.fields.Stat,
				Mean:   tt.fields.Mean,
				StdDev: tt.fields.StdDev,
			}
			s.Add(tt.args.playerAverage)
		})
	}
}

func TestStatOfInterest_ZScore(t *testing.T) {
	type fields struct {
		Scores []float64
		Stat   Stat
		Mean   float64
		StdDev float64
	}
	type args struct {
		playerAverage PlayerAverage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatOfInterest{
				Scores: tt.fields.Scores,
				Stat:   tt.fields.Stat,
				Mean:   tt.fields.Mean,
				StdDev: tt.fields.StdDev,
			}
			if got := s.ZScore(tt.args.playerAverage); got != tt.want {
				t.Errorf("StatOfInterest.ZScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerSimilarityMatrix_AddNormalizedPlayers(t *testing.T) {
	type fields struct {
		Matrix map[int]SimilarityVector
	}
	type args struct {
		players         []PlayerAverage
		statsOfInterest []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   PlayerSimilarityMatrix
	}{
		{name: "Range of Points",

			fields: fields{
				Matrix: map[int]SimilarityVector{},
			},
			args: args{
				players: []PlayerAverage{
					{
						Player: Player{PlayerID: 1},
						Points: 10,
					},
					{
						Player: Player{PlayerID: 2},
						Points: 5,
					},
					{
						Player: Player{PlayerID: 3},
						Points: 10,
					},
					{
						Player: Player{PlayerID: 4},
						Points: 11,
					},
					{
						Player: Player{PlayerID: 5},
						Points: 6,
					},
				},
				statsOfInterest: []string{"points"},
			},
			want: PlayerSimilarityMatrix{
				Matrix: map[int]SimilarityVector{
					1: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							2: PlayerDiff{
								Points: 2.06,
							},
							3: PlayerDiff{
								Points: 0,
							},
							4: PlayerDiff{
								Points: -0.41,
							},
							5: PlayerDiff{
								Points: 1.65,
							},
						},
						Average: PlayerAverage{
							Player: Player{PlayerID: 1},
							Points: 0.66,
						},
					},
					2: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points: -2.06,
							},
							3: PlayerDiff{
								Points: -2.06,
							},
							4: PlayerDiff{
								Points: -2.47,
							},
							5: PlayerDiff{
								Points: -0.41,
							},
						},
						Average: PlayerAverage{
							Player: Player{PlayerID: 2},
							Points: -1.40,
						},
					},
					3: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points: 0,
							},
							2: PlayerDiff{
								Points: 2.06,
							},
							4: PlayerDiff{
								Points: -0.41,
							},
							5: PlayerDiff{
								Points: 1.65,
							},
						},
						Average: PlayerAverage{
							Player: Player{PlayerID: 3},
							Points: 0.66,
						},
					},
					4: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points: 0.41,
							},
							2: PlayerDiff{
								Points: 2.47,
							},
							3: PlayerDiff{
								Points: 0.41,
							},
							5: PlayerDiff{
								Points: 2.06,
							},
						},
						Average: PlayerAverage{
							Player: Player{PlayerID: 4},
							Points: 1.07,
						},
					},
					5: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points: -1.65,
							},
							2: PlayerDiff{
								Points: 0.41,
							},
							3: PlayerDiff{
								Points: -1.65,
							},
							4: PlayerDiff{
								Points: -2.06,
							},
						},
						Average: PlayerAverage{
							Player: Player{PlayerID: 5},
							Points: -0.99,
						},
					},
				},
			},
		},
		{
			name: "All Stats",
			fields: fields{
				Matrix: map[int]SimilarityVector{},
			},
			args: args{
				players: []PlayerAverage{
					{
						Player:                 Player{PlayerID: 1},
						Points:                 1,
						Rebounds:               1,
						Assists:                1,
						Steals:                 1,
						Blocks:                 1,
						Turnovers:              1,
						DefensiveRebounds:      1,
						OffensiveRebounds:      1,
						ThreePointersMade:      1,
						ThreePointersAttempted: 1,
						FieldGoalsMade:         1,
						FieldGoalsAttempted:    1,
						FreeThrowsMade:         1,
						FreeThrowsAttempted:    1,
						Height:                 1,
						Weight:                 1,
						GamesPlayed:            1,
						PersonalFoulsDrawn:     1,
						PersonalFouls:          1,
						Minutes:                1,
					},
					{
						Player:                 Player{PlayerID: 2},
						Points:                 0,
						Rebounds:               0,
						Assists:                0,
						Steals:                 0,
						Blocks:                 0,
						Turnovers:              0,
						DefensiveRebounds:      0,
						OffensiveRebounds:      0,
						ThreePointersMade:      0,
						ThreePointersAttempted: 0,
						FieldGoalsMade:         0,
						FieldGoalsAttempted:    0,
						FreeThrowsMade:         0,
						FreeThrowsAttempted:    0,
						Height:                 0,
						Weight:                 0,
						GamesPlayed:            0,
						PersonalFoulsDrawn:     0,
						PersonalFouls:          0,
						Minutes:                0,
					},
				},
				statsOfInterest: []string{
					"games_played",
					"height",
					"weight",
					"minutes",
					"points",
					"assists",
					"rebounds",
					"blocks",
					"defensive_rebounds",
					"field_goals_attempted",
					"field_goals_made",
					"free_throws_attempted",
					"free_throws_made",
					"offensive_rebounds",
					"personal_fouls_drawn",
					"personal_fouls",
					"steals",
					"three_pointers_attempted",
					"three_pointers_made",
					"turnovers",
				},
			},
			want: PlayerSimilarityMatrix{
				Matrix: map[int]SimilarityVector{
					1: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							2: PlayerDiff{
								Points:                 2,
								Rebounds:               2,
								Assists:                2,
								Steals:                 2,
								Blocks:                 2,
								Turnovers:              2,
								DefensiveRebounds:      2,
								OffensiveRebounds:      2,
								ThreePointersMade:      2,
								ThreePointersAttempted: 2,
								FieldGoalsMade:         2,
								FieldGoalsAttempted:    2,
								FreeThrowsMade:         2,
								FreeThrowsAttempted:    2,
								Height:                 2,
								Weight:                 2,
								GamesPlayed:            2,
								PersonalFoulsDrawn:     2,
								PersonalFouls:          2,
								Minutes:                2,
							},
						},
						Average: PlayerAverage{
							Player:                 Player{PlayerID: 1},
							Points:                 1,
							Rebounds:               1,
							Assists:                1,
							Steals:                 1,
							Blocks:                 1,
							Turnovers:              1,
							DefensiveRebounds:      1,
							OffensiveRebounds:      1,
							ThreePointersMade:      1,
							ThreePointersAttempted: 1,
							FieldGoalsMade:         1,
							FieldGoalsAttempted:    1,
							FreeThrowsMade:         1,
							FreeThrowsAttempted:    1,
							Height:                 1,
							Weight:                 1,
							GamesPlayed:            1,
							PersonalFoulsDrawn:     1,
							PersonalFouls:          1,
							Minutes:                1,
						},
					},
					2: SimilarityVector{
						Comparisons: map[int]PlayerDiff{
							1: PlayerDiff{
								Points:                 -2,
								Rebounds:               -2,
								Assists:                -2,
								Steals:                 -2,
								Blocks:                 -2,
								Turnovers:              -2,
								DefensiveRebounds:      -2,
								OffensiveRebounds:      -2,
								ThreePointersMade:      -2,
								ThreePointersAttempted: -2,
								FieldGoalsMade:         -2,
								FieldGoalsAttempted:    -2,
								FreeThrowsMade:         -2,
								FreeThrowsAttempted:    -2,
								Height:                 -2,
								Weight:                 -2,
								GamesPlayed:            -2,
								PersonalFoulsDrawn:     -2,
								PersonalFouls:          -2,
								Minutes:                -2,
							},
						},
						Average: PlayerAverage{
							Player:                 Player{PlayerID: 2},
							Points:                 -1,
							Rebounds:               -1,
							Assists:                -1,
							Steals:                 -1,
							Blocks:                 -1,
							Turnovers:              -1,
							DefensiveRebounds:      -1,
							OffensiveRebounds:      -1,
							ThreePointersMade:      -1,
							ThreePointersAttempted: -1,
							FieldGoalsMade:         -1,
							FieldGoalsAttempted:    -1,
							FreeThrowsMade:         -1,
							FreeThrowsAttempted:    -1,
							Height:                 -1,
							Weight:                 -1,
							GamesPlayed:            -1,
							PersonalFoulsDrawn:     -1,
							PersonalFouls:          -1,
							Minutes:                -1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PlayerSimilarityMatrix{
				Matrix: tt.fields.Matrix,
			}
			m.AddNormalizedPlayers(tt.args.players)
			for i, row := range m.Matrix {
				if !reflect.DeepEqual(row, tt.want.Matrix[i]) {
					t.Errorf("Got Player Similarity Vector [%d] = \n%v, but want \n%v", i, m.Matrix[i], tt.want.Matrix[i])
				}
			}

		})
	}
}

func TestPlayerSimilarityMatrix_CompareAverages(t *testing.T) {
	type fields struct {
		Matrix map[int]SimilarityVector
	}
	type args struct {
		in        int
		averageIn PlayerAverage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int]PlayerDiff
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PlayerSimilarityMatrix{
				Matrix: tt.fields.Matrix,
			}
			if got := m.CompareAverages(tt.args.in, tt.args.averageIn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlayerSimilarityMatrix.CompareAverages() = %v, want %v", got, tt.want)
			}
		})
	}
}
