package leetcode

import (
	"reflect"
	"testing"
)

func TestFloodFill(t *testing.T) {
	qs := []struct {
		name  string
		image [][]int
		sr    int
		sc    int
		color int
		want  [][]int
	}{
		{
			name:  "Example 1",
			image: [][]int{{1, 1, 1}, {1, 1, 0}, {1, 0, 1}},
			sr:    1,
			sc:    1,
			color: 2,
			want:  [][]int{{2, 2, 2}, {2, 2, 0}, {2, 0, 1}},
		},
		{
			name:  "Example 2",
			image: [][]int{{0, 0, 0}, {0, 0, 0}},
			sr:    0,
			sc:    0,
			color: 0,
			want:  [][]int{{0, 0, 0}, {0, 0, 0}},
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := floodFill(q.image, q.sr, q.sc, q.color)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("floodFill(%v, %d, %d, %d) = %v, want %v", q.image, q.sr, q.sc, q.color, got, q.want)
			}
		})
	}
}
