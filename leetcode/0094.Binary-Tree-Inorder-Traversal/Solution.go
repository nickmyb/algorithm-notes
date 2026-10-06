package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

// TreeNode 起别名，题解本体里就能写裸的 *TreeNode
type TreeNode = structures.TreeNode

// ===== 以下是题解本体，与提交到 LeetCode 的代码一字不差 =====

func inorderTraversal(root *TreeNode) []int {
	ret := []int{}
	inorder(root, &ret)
	return ret
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
