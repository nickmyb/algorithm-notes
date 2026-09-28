package leetcode

import (
	"reflect"
	"testing"
)

func TestMoveZeroes(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		want []int
	}{
		{"Example 1", []int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{"Example 2", []int{0}, []int{0}},
	}

	impls := map[string]func([]int){
		"moveZeroes":           moveZeroes,
		"moveZeroesTwoPointer": moveZeroesTwoPointer,
	}

	for _, q := range qs {
		for implName, impl := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				// 因为多个实现要用同样的数组做测试,需要复制数组用于测试
				nums := append([]int(nil), q.in...)
				impl(nums)

				if !reflect.DeepEqual(nums, q.want) {
					t.Fatalf("%v(%v) = %v, want %v", implName, q.in, nums, q.want)
				}
			})
		}
	}
}
