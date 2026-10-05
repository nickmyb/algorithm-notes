package leetcode

import (
	"reflect"
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestLevelOrder(t *testing.T) {
	NULL := structures.NULL

	qs := []struct {
		name string
		in   []int
		want [][]int
	}{
		{"Example 1", []int{3, 9, 20, NULL, NULL, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},
		{"Example 2", []int{1}, [][]int{{1}}},
		{"Example 3", []int{}, [][]int{}},
	}

	impls := map[string]func(*TreeNode) [][]int{
		"levelOrder":          levelOrder,
		"levelOrderRecursive": levelOrderRecursive,
	}

	for _, q := range qs {
		for implName, solve := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := solve(structures.Ints2TreeNode(q.in))
				// LeetCode 把 nil 和空切片都序列化成 []，这里同样视为相等
				if len(got) == 0 && len(q.want) == 0 {
					return
				}
				if !reflect.DeepEqual(got, q.want) {
					t.Fatalf("%s(%v) = %v, want %v", implName, structures.FormatInts(q.in), got, q.want)
				}
			})
		}
	}
}
