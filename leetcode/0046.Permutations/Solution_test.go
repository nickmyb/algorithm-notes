package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestPermute(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "Example 1",
			nums: []int{1, 2, 3},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
		{
			name: "Example 2",
			nums: []int{0, 1},
			want: [][]int{{0, 1}, {1, 0}},
		},
		{
			name: "Example 3",
			nums: []int{1},
			want: [][]int{{1}},
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := permute(q.nums)
			if !reflect.DeepEqual(normalizePermutations(got), normalizePermutations(q.want)) {
				t.Fatalf("permute(%v) = %v, want %v", q.nums, got, q.want)
			}
		})
	}
}

func normalizePermutations(permutations [][]int) [][]int {
	normalized := append([][]int(nil), permutations...)
	slices.SortFunc(normalized, func(a, b []int) int {
		return slices.Compare(a, b)
	})
	return normalized
}
