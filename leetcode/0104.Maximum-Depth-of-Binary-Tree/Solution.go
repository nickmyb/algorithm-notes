package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

type TreeNode = structures.TreeNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	maxLeft := maxDepth(root.Left)
	maxRight := maxDepth(root.Right)
	return max(maxLeft, maxRight) + 1
}
