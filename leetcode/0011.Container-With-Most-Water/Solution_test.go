package leetcode

import "testing"

func TestMaxArea(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		want int
	}{
		{"Example 1", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"Example 2", []int{1, 1}, 1},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := maxArea(q.in)
			if got != q.want {
				t.Fatalf("maxArea(%v) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}
