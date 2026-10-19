package leetcode

import "testing"

func TestOrangesRotting(t *testing.T) {
	qs := []struct {
		name string
		grid [][]int
		want int
	}{
		{
			name: "Example 1",
			grid: [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}},
			want: 4,
		},
		{
			name: "Example 2",
			grid: [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}},
			want: -1,
		},
		{
			name: "Example 3",
			grid: [][]int{{0, 2}},
			want: 0,
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := orangesRotting(q.grid)
			if got != q.want {
				t.Fatalf("orangesRotting(%v) = %d, want %d", q.grid, got, q.want)
			}
		})
	}
}
