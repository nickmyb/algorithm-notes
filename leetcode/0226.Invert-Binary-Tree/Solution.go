package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

type TreeNode = structures.TreeNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	root.Left, root.Right = root.Right, root.Left
	invertTree(root.Left)
	invertTree(root.Right)

	return root
}
