package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func findMaxAverage(nums []int, k int) float64 {
	prevSumK := sumK(nums, k, 0)
	maxSumK := prevSumK

	for i := 1; i < len(nums)-k+1; i++ {
		sum := deltaSumK(nums, k, i, prevSumK)

		if sum > maxSumK {
			maxSumK = sum
		}

		prevSumK = sum
	}

	return float64(maxSumK) / float64(k)
}

func sumK(nums []int, k int, start int) int {
	ret := 0

	for i := start; i < start+k; i++ {
		ret += nums[i]
	}

	return ret
}

func deltaSumK(nums []int, k int, start int, prevSum int) int {
	return prevSum - nums[start-1] + nums[start+k-1]
}
