package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func moveZeroes(nums []int) {
	zeroCount := 0

	for i, num := range nums {
		if num == 0 {
			zeroCount += 1
		} else {
			nums[i-zeroCount] = num
		}
	}

	end := len(nums)
	for zeroCount > 0 {
		end -= 1
		nums[end] = 0
		zeroCount -= 1
	}
}
