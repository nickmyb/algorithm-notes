package leetcode

import (
	"reflect"
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestMaxDepth(t *testing.T) {
	NULL := structures.NULL

	qs := []struct {
		name string
		in   []int
		want int
	}{
		{"Example 1", []int{3, 9, 20, NULL, NULL, 15, 7}, 3},
		{"Example 2", []int{1, NULL, 2}, 2},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := maxDepth(structures.Ints2TreeNode(q.in))
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("maxDepth(%v) = %v, want %v", structures.FormatInts(q.in), got, q.want)
			}
		})
	}
}
