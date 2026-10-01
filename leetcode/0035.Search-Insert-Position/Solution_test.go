package leetcode

import (
	"testing"
)

func TestSearchInsert(t *testing.T) {
	qs := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"Example 1", []int{1, 3, 5, 6}, 5, 2},
		{"Example 2", []int{1, 3, 5, 6}, 2, 1},
		{"Example 3", []int{1, 3, 5, 6}, 7, 4},
		{"target 小于所有元素", []int{1, 3, 5, 6}, 0, 0},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := searchInsert(q.nums, q.target)
			if got != q.want {
				t.Fatalf("searchInsert(%v, %d) = %d, want %d", q.nums, q.target, got, q.want)
			}
		})
	}
}
