package helpers

import (
	"testing"
)

func Test_Smallest(t *testing.T) {
	type args struct {
		ff []float64
	}
	tests := []struct {
		name string
		ff   []float64
		want float64
	}{
		{name: "1", ff: []float64{0.1, -0.2, 4, 14, -2}, want: -2.0},
		{name: "2", ff: []float64{0.1, -0.2, 4, 14}, want: -0.2},
		{name: "3", ff: []float64{0.1, 0.2, 4, 14, 2}, want: 0.1},
		{name: "4", ff: []float64{2, 3, 4, 5, 1}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Smallest(tt.ff); got != tt.want {
				t.Errorf("smallest() = %v, want %v", got, tt.want)
			}
		})
	}
}
