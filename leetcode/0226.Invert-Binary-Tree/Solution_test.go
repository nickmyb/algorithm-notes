package leetcode

import (
	"reflect"
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestInvertTree(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		want []int
	}{
		{"Example 1", []int{4, 2, 7, 1, 3, 6, 9}, []int{4, 7, 2, 9, 6, 3, 1}},
		{"Example 2", []int{2, 1, 3}, []int{2, 3, 1}},
		{"Example 3", []int{}, []int{}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := structures.Tree2ints(invertTree(structures.Ints2TreeNode(q.in)))
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("invertTree(%v) = %v, want %v",
					structures.FormatInts(q.in), structures.FormatInts(got), structures.FormatInts(q.want))
			}
		})
	}
}
