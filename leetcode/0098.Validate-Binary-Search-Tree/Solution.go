package leetcode

// ===== 本地接线区 =====
import (
	"math"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

type TreeNode = structures.TreeNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func isValidBST(root *TreeNode) bool {
	return isValidBSTLoHi(root, math.MinInt, math.MaxInt)
}

func isValidBSTLoHi(root *TreeNode, lo int, hi int) bool {
	if root == nil {
		return true
	}

	return lo < root.Val && root.Val < hi && isValidBSTLoHi(root.Left, lo, root.Val) && isValidBSTLoHi(root.Right, root.Val, hi)
}

func isValidBSTInorder(root *TreeNode) bool {
	var ret []int

	inorder(root, &ret)

	for i := range len(ret) - 1 {
		if !(ret[i] < ret[i+1]) {
			return false
		}
	}

	return true
}

// region @snippets/go/traversal.go

func inorder(root *TreeNode, ret *[]int) {
	if root == nil {
		return
	}

	inorder(root.Left, ret)
	*ret = append(*ret, root.Val)
	inorder(root.Right, ret)
}

// endregion @snippets/go/traversal.go
