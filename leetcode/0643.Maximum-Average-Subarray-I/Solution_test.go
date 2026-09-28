package leetcode

import (
	"reflect"
	"testing"
)

func TestFindMaxAverage(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		k    int
		want float64
	}{
		{"Example 1", []int{1, 12, -5, -6, 50, 3}, 4, 12.75},
		{"Example 2", []int{5}, 1, 5},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := findMaxAverage(q.nums, q.k)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("findMaxAverage(%v, %v) = %v, want %v", q.nums, q.k, got, q.want)
			}
		})
	}
}
