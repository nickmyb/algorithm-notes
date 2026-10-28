package leetcode

import "testing"

func TestSubarraySum(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"Example 1", []int{1, 1, 1}, 2, 2},
		{"Example 2", []int{1, 2, 3}, 3, 2},
	}

	impls := map[string]func([]int, int) int{
		"subarraySum":           subarraySum,
		"subarraySumTraversals": subarraySumTraversals,
	}

	for _, q := range qs {
		for implName, impl := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := impl(q.nums, q.k)
				if got != q.want {
					t.Fatalf("%v(%v, %v) = %v, want %v", implName, q.nums, q.k, got, q.want)
				}
			})
		}
	}
}
