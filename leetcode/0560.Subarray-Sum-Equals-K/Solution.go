package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func subarraySum(nums []int, k int) int {
	// 前缀和为K出现过v次
	count := map[int]int{0: 1}
	sum := 0
	subs := 0

	for _, n := range nums {
		sum += n
		subs += count[sum-k]
		count[sum] += 1
	}

	return subs
}

func subarraySumTraversals(nums []int, k int) int {
	subs := 0

	for i := 0; i < len(nums); i++ {
		sum := nums[i]
		if sum == k {
			subs += 1
		}

		for j := i + 1; j < len(nums); j++ {
			sum += nums[j]
			if sum == k {
				subs += 1
			}
		}
	}

	return subs
}
