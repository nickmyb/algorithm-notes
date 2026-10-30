package leetcode

import (
	"reflect"
	"testing"
)

func TestMaxSlidingWindow(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"Example 1", []int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{"Example 2", []int{1}, 1, []int{1}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := maxSlidingWindow(q.nums, q.k)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("maxSlidingWindow(%v, %v) = %v, want %v", q.nums, q.k, got, q.want)
			}
		})
	}
}
