package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

type TreeNode = structures.TreeNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func levelOrder(root *TreeNode) [][]int {
	ret := [][]int{}
	r := []*TreeNode{root}

	for len(r) != 0 {
		next, curNums := traverseNode(r)
		if curNums != nil {
			ret = append(ret, curNums)
		}
		r = next
	}

	return ret
}

func traverseNode(cur []*TreeNode) ([]*TreeNode, []int) {
	var next []*TreeNode
	var curNums []int

	for _, n := range cur {
		if n != nil {
			curNums = append(curNums, n.Val)
			next = append(next, n.Left, n.Right)
		}
	}

	return next, curNums
}

func levelOrderRecursive(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	return levelNode([]*TreeNode{root}, nil)
}

func levelNode(root []*TreeNode, prev [][]int) [][]int {
	if len(root) == 0 {
		return prev
	}

	var curLevel []int
	preLength := len(root)

	// range 在开始时就固定了要遍历的长度，新加进去的孩子要到下一层才会被处理
	for _, n := range root {
		curLevel = append(curLevel, n.Val)
		if n.Left != nil {
			root = append(root, n.Left)
		}
		if n.Right != nil {
			root = append(root, n.Right)
		}
	}
	root = root[preLength:]
	prev = append(prev, curLevel)

	return levelNode(root, prev)
}
