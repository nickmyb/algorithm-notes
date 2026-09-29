package leetcode

import (
	"strconv"
	"testing"
)

func TestNumArray(t *testing.T) {
	qs := []struct {
		name    string
		nums    []int
		queries [][2]int
		want    []int
	}{
		{"Example 1", []int{-2, 0, 3, -5, 2, -1}, [][2]int{{0, 2}, {2, 5}, {0, 5}}, []int{1, -1, -3}},
	}

	for _, q := range qs {
		na := Constructor(q.nums)

		for index, query := range q.queries {
			t.Run(q.name+" ["+strconv.Itoa(index)+"]", func(t *testing.T) {

				got := na.SumRange(query[0], query[1])
				if got != q.want[index] {
					t.Fatalf("NumArray(%v).SumRange(%v) = %v, want %v", q.nums, query, got, q.want[index])
				}
			})
		}
	}
}
