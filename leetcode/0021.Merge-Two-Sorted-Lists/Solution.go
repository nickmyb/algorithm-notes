package leetcode

// ===== 本地接线区 =====

import structures "github.com/nickmyb/algorithm-notes/structures/go"

type ListNode = structures.ListNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	sentinel := &ListNode{Val: 0, Next: nil}
	node := sentinel
	node1 := list1
	node2 := list2

	for node1 != nil || node2 != nil {
		if node1 == nil {
			node.Next = node2
			node2 = node2.Next
		} else if node2 == nil {
			node.Next = node1
			node1 = node1.Next
		} else if node1.Val <= node2.Val {
			node.Next = node1
			node1 = node1.Next
		} else {
			node.Next = node2
			node2 = node2.Next
		}

		node.Next.Next = nil
		node = node.Next
	}

	return sentinel.Next
}
