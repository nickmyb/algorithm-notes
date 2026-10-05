package leetcode

import (
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestIsValidBST(t *testing.T) {
	NULL := structures.NULL

	qs := []struct {
		name string
		in   []int
		want bool
	}{
		{"Example 1", []int{2, 1, 3}, true},
		{"Example 2", []int{5, 1, 4, NULL, NULL, 3, 6}, false},
		// 以下题目没给，各自抓一种两个 Example 都放过的错误实现（已对拍）
		// 只和直接孩子比较：3 比父节点 6 小，但在根 5 的右子树里
		{"Grandchild violates root", []int{5, 4, 6, NULL, NULL, 3, 7}, false},
		// 边界写成 < 而不是 <=：题目要求严格大于/小于，相等不合法
		{"Duplicate values", []int{2, 2, 2}, false},
		// 用 int32 极值当初始上下界：节点值正好取到 2^31-1 / -2^31 时被误判
		{"Max int32", []int{2147483647}, true},
		{"Min int32", []int{-2147483648}, true},
	}

	impls := map[string]func(*TreeNode) bool{
		"isValidBST":        isValidBST,
		"isValidBSTInorder": isValidBSTInorder,
	}

	for _, q := range qs {
		for implName, solve := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := solve(structures.Ints2TreeNode(q.in))
				if got != q.want {
					t.Fatalf("%s(%v) = %v, want %v", implName, structures.FormatInts(q.in), got, q.want)
				}
			})
		}
	}
}
