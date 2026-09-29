package leetcode

// ===== 本地接线区 =====
import structures "github.com/nickmyb/algorithm-notes/structures/go"

type ListNode = structures.ListNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	var prev *ListNode

	for head != nil {
		next := head.Next
		head.Next = prev

		prev = head
		head = next
	}

	return prev
}

func reverseListStack(head *ListNode) *ListNode {
	var ret *ListNode
	stack := []*ListNode{nil}

	for head != nil {
		stack = append(stack, head)
		head = head.Next
	}

	length := len(stack)

	if length == 0 {
		return nil
	}

	ret = stack[length-1]
	i := 1

	for length-i-1 >= 0 {
		stack[length-i].Next = stack[length-i-1]
		i++
	}

	return ret
}
