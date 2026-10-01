package leetcode

import (
	"reflect"
	"testing"
)

func TestSearchRange(t *testing.T) {
	qs := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"Example 1", []int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{"Example 2", []int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{"Example 3", []int{}, 0, []int{-1, -1}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := searchRange(q.nums, q.target)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("searchRange(%v, %d) = %v, want %v", q.nums, q.target, got, q.want)
			}
		})
	}
}
