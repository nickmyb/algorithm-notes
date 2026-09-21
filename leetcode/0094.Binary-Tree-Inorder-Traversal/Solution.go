package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

// TreeNode 起别名，题解本体里就能写裸的 *TreeNode
type TreeNode = structures.TreeNode

// ===== 以下是题解本体，与提交到 LeetCode 的代码一字不差 =====

func inorderTraversal(root *TreeNode) []int {
	return inorder(root, []int{})
}

func inorder(root *TreeNode, input []int) []int {
	if root == nil {
		return input
	}

	output := inorder(root.Left, input)
	output = append(output, root.Val)
	output = inorder(root.Right, output)

	return output
}
