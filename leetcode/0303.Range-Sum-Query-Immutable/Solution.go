package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

type NumArray struct {
	nums      []int
	prefixSum []int
}

func Constructor(nums []int) NumArray {
	na := NumArray{
		nums:      nums,
		prefixSum: make([]int, len(nums)+1),
	}

	for i, num := range nums {
		na.prefixSum[i+1] = na.prefixSum[i] + num
	}

	return na

}

func (this *NumArray) SumRange(left int, right int) int {
	return this.prefixSum[right+1] - this.prefixSum[left]
}
