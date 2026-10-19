package leetcode

import "testing"

func TestNumIslands(t *testing.T) {
	qs := []struct {
		name string
		grid [][]byte
		want int
	}{
		{
			name: "Example 1",
			grid: [][]byte{
				[]byte("11110"),
				[]byte("11010"),
				[]byte("11000"),
				[]byte("00000"),
			},
			want: 1,
		},
		{
			name: "Example 2",
			grid: [][]byte{
				[]byte("11000"),
				[]byte("11000"),
				[]byte("00100"),
				[]byte("00011"),
			},
			want: 3,
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := numIslands(q.grid)
			if got != q.want {
				t.Fatalf("numIslands(%q) = %d, want %d", q.grid, got, q.want)
			}
		})
	}
}
