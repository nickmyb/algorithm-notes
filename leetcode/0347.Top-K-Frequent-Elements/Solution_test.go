package leetcode

import (
	"slices"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"Example 1", []int{1, 1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{"Example 2", []int{1}, 1, []int{1}},
		{"Example 3", []int{1, 2, 1, 2, 1, 2, 3, 1, 3, 2}, 2, []int{1, 2}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := topKFrequent(q.nums, q.k)
			// 题目允许任意顺序返回，排序后再比较
			sorted := slices.Clone(got)
			slices.Sort(sorted)
			if !slices.Equal(sorted, q.want) {
				t.Fatalf("topKFrequent(%v, %d) = %v, want %v in any order", q.nums, q.k, got, q.want)
			}
		})
	}
}
