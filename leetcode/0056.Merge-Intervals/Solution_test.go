package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestMerge(t *testing.T) {
	qs := []struct {
		name      string
		intervals [][]int
		want      [][]int
	}{
		{"Example 1", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{"Example 2", [][]int{{1, 4}, {4, 5}}, [][]int{{1, 5}}},
		{"Example 3", [][]int{{4, 7}, {1, 4}}, [][]int{{1, 7}}},
		{"Example 4", [][]int{{2, 3}, {5, 5}, {2, 2}, {3, 4}, {3, 4}}, [][]int{{2, 4}, {5, 5}}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := merge(q.intervals)
			slices.SortFunc(got, func(a, b []int) int { return slices.Compare(a, b) })
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("merge(%v) = %v, want %v", q.intervals, got, q.want)
			}
		})
	}
}
