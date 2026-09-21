package structures

// ListNode 是链接节点
// 这个不能复制到*_test.go文件中。会导致Travis失败
type ListNode struct {
	Val  int
	Next *ListNode
}

// List2Ints convert List to []int
func List2Ints(head *ListNode) []int {
	// 按节点身份判环，正常长链表不能因为超过固定长度而失败。
	seen := make(map[*ListNode]struct{})
	res := []int{}
	for head != nil {
		if _, ok := seen[head]; ok {
			panic("链表有环，无法展开")
		}
		seen[head] = struct{}{}

		res = append(res, head.Val)
		head = head.Next
	}

	return res
}

// Ints2List convert []int to List
func Ints2List(nums []int) *ListNode {
	if len(nums) == 0 {
		return nil
	}

	l := &ListNode{}
	t := l
	for _, v := range nums {
		t.Next = &ListNode{Val: v}
		t = t.Next
	}
	return l.Next
}

// GetNodeWith returns the first node with val
func (l *ListNode) GetNodeWith(val int) *ListNode {
	res := l
	for res != nil {
		if res.Val == val {
			break
		}
		res = res.Next
	}
	return res
}

// Ints2ListWithCycle returns a list whose tail point to pos-indexed node
// head's index is 0
// if pos = -1, no cycle
func Ints2ListWithCycle(nums []int, pos int) *ListNode {
	head := Ints2List(nums)
	if pos == -1 {
		return head
	}
	c := head
	for pos > 0 {
		c = c.Next
		pos--
	}
	tail := c
	for tail.Next != nil {
		tail = tail.Next
	}
	tail.Next = c
	return head
}
