package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestThreeSum(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		want [][]int
	}{
		{"Example 1", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"Example 2", []int{0, 1, 1}, [][]int{}},
		{"Example 3", []int{0, 0, 0}, [][]int{{0, 0, 0}}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := threeSum(slices.Clone(q.in))
			if !reflect.DeepEqual(normalizeTriplets(got), normalizeTriplets(q.want)) {
				t.Fatalf("threeSum(%v) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}

// 题目允许三元组内部和结果列表任意排序；保留重复项，以便发现重复答案。
func normalizeTriplets(triplets [][]int) [][]int {
	normalized := make([][]int, len(triplets))
	for i, triplet := range triplets {
		normalized[i] = slices.Clone(triplet)
		slices.Sort(normalized[i])
	}
	slices.SortFunc(normalized, func(a, b []int) int {
		return slices.Compare(a, b)
	})
	return normalized
}
