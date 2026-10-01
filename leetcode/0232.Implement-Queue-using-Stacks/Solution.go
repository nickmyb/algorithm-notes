package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

type MyQueue struct {
	stack []int
	queue []int
}

func Constructor() MyQueue {
	return MyQueue{}
}

func (this *MyQueue) Push(x int) {
	this.stack = append(this.stack, x)
}

func (this *MyQueue) Pop() int {
	if len(this.queue) == 0 {
		this.Stack2Queue()
	}
	popped := this.queue[len(this.queue)-1]
	this.queue = this.queue[:len(this.queue)-1]
	return popped
}

func (this *MyQueue) Peek() int {
	if len(this.queue) == 0 {
		this.Stack2Queue()
	}
	return this.queue[len(this.queue)-1]
}

func (this *MyQueue) Empty() bool {
	return len(this.stack)+len(this.queue) == 0
}

func (this *MyQueue) Stack2Queue() {
	for len(this.stack) > 0 {
		popped := this.stack[len(this.stack)-1]
		this.stack = this.stack[:len(this.stack)-1]
		this.queue = append(this.queue, popped)
	}
}
