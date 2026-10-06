package leetcode

import (
	"testing"
)

func TestFindKthLargest(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"Example 1", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"Example 2", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := findKthLargest(q.nums, q.k)
			if got != q.want {
				t.Fatalf("findKthLargest(%v, %d) = %v, want %v", q.nums, q.k, got, q.want)
			}
		})
	}
}
