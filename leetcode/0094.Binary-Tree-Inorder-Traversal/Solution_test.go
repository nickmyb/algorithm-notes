package leetcode

import (
	"reflect"
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestInorderTraversal(t *testing.T) {
	NULL := structures.NULL
	// 题目英文版给的 4 个示例，原样转录。
	// 注意中文版只有 3 个，缺了下面这棵大树——官方翻译滞后，以英文版为准。
	qs := []struct {
		name string
		in   []int
		want []int
	}{
		{"Example 1", []int{1, NULL, 2, 3}, []int{1, 3, 2}},
		{"Example 2", []int{1, 2, 3, 4, 5, NULL, 8, NULL, NULL, 6, 7, 9}, []int{4, 2, 6, 5, 7, 1, 3, 9, 8}},
		// 题解从 []int{} 起手而不是 var result []int，所以空树返回空切片而非 nil。
		// 这点要留意：LeetCode 判题对 nil 和空数组一视同仁，reflect.DeepEqual 不。
		{"Example 3 空树", []int{}, []int{}},
		{"Example 4 单节点", []int{1}, []int{1}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := inorderTraversal(structures.Ints2TreeNode(q.in))
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("inorderTraversal(%v) = %v, want %v", structures.FormatInts(q.in), got, q.want)
			}
		})
	}
}
