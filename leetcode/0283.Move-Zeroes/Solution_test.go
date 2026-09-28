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

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			moveZeroes(q.in)
			got := q.in
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("moveZeroes(%v) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}
