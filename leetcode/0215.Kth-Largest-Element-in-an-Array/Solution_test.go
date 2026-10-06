package leetcode

import (
	"slices"
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

	impls := map[string]func([]int, int) int{
		"findKthLargest":      findKthLargest,
		"findKthLargestMinPQ": findKthLargestMinPQ,
	}

	for _, q := range qs {
		for implName, solve := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				// 传副本：以后加原地划分的解法（如快速选择）会改动 nums，不能影响别的实现
				got := solve(slices.Clone(q.nums), q.k)
				if got != q.want {
					t.Fatalf("%s(%v, %d) = %v, want %v", implName, q.nums, q.k, got, q.want)
				}
			})
		}
	}
}
