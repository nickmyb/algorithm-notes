package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func maxSlidingWindow(nums []int, k int) []int {
	var ret []int
	var candidates []int

	for i := 0; i < len(nums); i++ {
		if i-k >= 0 && candidates[0] == i-k {
			candidates = candidates[1:]
		}

		for len(candidates) > 0 && nums[candidates[len(candidates)-1]] <= nums[i] {
			candidates = candidates[:len(candidates)-1]
		}
		candidates = append(candidates, i)

		if i >= k-1 {
			ret = append(ret, nums[candidates[0]])
		}
	}

	return ret
}
