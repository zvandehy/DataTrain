package resolver

import "testing"

func TestPoolVariance(t *testing.T) {
	type args struct {
		datasets [][]float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Example from paper",
			args: args{
				datasets: [][]float64{
					{32, 36, 27, 28, 30, 31},
					{32, 34, 30, 33, 29, 36, 24},
					{39, 40, 42},
				},
			},
			want: 22.84,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PoolVariance(tt.args.datasets); got != tt.want {
				t.Errorf("PoolVariance() = %v, want %v", got, tt.want)
			}
		})
	}
}
