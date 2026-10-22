package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestCombinationSum(t *testing.T) {
	qs := []struct {
		name       string
		candidates []int
		target     int
		want       [][]int
	}{
		{
			name:       "Example 1",
			candidates: []int{2, 3, 6, 7},
			target:     7,
			want:       [][]int{{2, 2, 3}, {7}},
		},
		{
			name:       "Example 2",
			candidates: []int{2, 3, 5},
			target:     8,
			want:       [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}},
		},
		{
			name:       "Example 3",
			candidates: []int{2},
			target:     1,
			want:       [][]int{},
		},
		{
			name:       "Example 4",
			candidates: []int{8, 7, 4, 3},
			target:     11,
			want:       [][]int{{4, 4, 3}, {7, 4}, {8, 3}},
		},
	}
	impls := []struct {
		name  string
		solve func([]int, int) [][]int
	}{
		{"combinationSum", combinationSum},
		{"combinationSumNumOrNot", combinationSumNumOrNot},
	}

	for _, q := range qs {
		for _, impl := range impls {
			t.Run(q.name+"/"+impl.name, func(t *testing.T) {
				got := impl.solve(q.candidates, q.target)
				if !reflect.DeepEqual(normalizeCombinations(got), normalizeCombinations(q.want)) {
					t.Fatalf("%s(%v, %d) = %v, want %v", impl.name, q.candidates, q.target, got, q.want)
				}
			})
		}
	}
}

func normalizeCombinations(combinations [][]int) [][]int {
	normalized := make([][]int, len(combinations))
	for i, combination := range combinations {
		normalized[i] = append([]int(nil), combination...)
		slices.Sort(normalized[i])
	}
	slices.SortFunc(normalized, func(a, b []int) int {
		return slices.Compare(a, b)
	})
	return normalized
}
