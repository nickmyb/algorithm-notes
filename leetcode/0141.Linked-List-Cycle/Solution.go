package leetcode

// ===== 本地接线区 =====

import structures "github.com/nickmyb/algorithm-notes/structures/go"

type ListNode = structures.ListNode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func hasCycle(head *ListNode) bool {
	faster := head

	for faster != nil && faster.Next != nil {
		faster = faster.Next.Next
		head = head.Next
		if faster == head {
			return true
		}
	}

	return false
}

func hasCycleMap(head *ListNode) bool {
	m := map[*ListNode]struct{}{}

	for head != nil {
		_, ok := m[head]

		if ok {
			return true
		}

		m[head] = struct{}{}
		head = head.Next
	}

	return false
}
