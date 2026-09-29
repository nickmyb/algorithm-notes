package leetcode

import (
	"reflect"
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestMergeTwoLists(t *testing.T) {
	qs := []struct {
		name string
		l1   []int
		l2   []int
		want []int
	}{
		{"Example 1", []int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{"Example 2", []int{}, []int{}, []int{}},
		{"Example 3", []int{}, []int{0}, []int{0}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := mergeTwoLists(structures.Ints2List(q.l1), structures.Ints2List(q.l2))

			ret := structures.List2Ints(got)
			if !reflect.DeepEqual(ret, q.want) {
				t.Fatalf("mergeTwoLists(%v, %v) = %v, want %v", structures.FormatInts(q.l1), structures.FormatInts(q.l2), ret, q.want)
			}
		})
	}
}
