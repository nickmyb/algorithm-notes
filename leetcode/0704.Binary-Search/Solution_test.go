package leetcode

import (
	"testing"
)

func TestSearch(t *testing.T) {
	qs := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"Example 1", []int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{"Example 2", []int{-1, 0, 3, 5, 9, 12}, 2, -1},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := search(q.nums, q.target)
			if got != q.want {
				t.Fatalf("search(%v, %d) = %d, want %d", q.nums, q.target, got, q.want)
			}
		})
	}
}
