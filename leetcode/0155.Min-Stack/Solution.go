package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

const NullValue = -999

type MinStack struct {
	nums []int
	min  int
	// min [0, index - 1] = preMin[index]
	preMin []int
}

func Constructor() MinStack {
	return MinStack{}
}

func (this *MinStack) Push(value int) {
	this.nums = append(this.nums, value)
	length := len(this.nums)

	if length == 1 {
		this.min = value
		this.preMin = append(this.preMin, NullValue)
	} else {
		this.preMin = append(this.preMin, this.min)
		if value < this.min {
			this.min = value
		}
	}
}

func (this *MinStack) Pop() {
	if len(this.nums) > 0 {
		top := this.Top()
		preMin := this.preMin[len(this.nums)-1]
		this.nums = this.nums[:len(this.nums)-1]
		this.preMin = this.preMin[:len(this.preMin)-1]

		if top == this.min {
			// 不记录preMin的话就需要计算当前最小值
			this.min = preMin
		}
	}
}

func (this *MinStack) Top() int {
	if len(this.nums) > 0 {
		return this.nums[len(this.nums)-1]
	}
	return NullValue
}

func (this *MinStack) GetMin() int {
	if len(this.nums) > 0 {
		return this.min
	}
	return NullValue
}
