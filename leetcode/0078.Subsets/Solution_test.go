package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestSubsets(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "Example 1",
			nums: []int{1, 2, 3},
			want: [][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		},
		{
			name: "Example 2",
			nums: []int{0},
			want: [][]int{{}, {0}},
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := subsets(q.nums)
			if !reflect.DeepEqual(normalizeSubsets(got), normalizeSubsets(q.want)) {
				t.Fatalf("subsets(%v) = %v, want %v", q.nums, got, q.want)
			}
		})
	}
}

func normalizeSubsets(subsets [][]int) [][]int {
	normalized := make([][]int, len(subsets))
	for i, subset := range subsets {
		normalized[i] = append([]int(nil), subset...)
		slices.Sort(normalized[i])
	}
	slices.SortFunc(normalized, func(a, b []int) int {
		return slices.Compare(a, b)
	})
	return normalized
}
