package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func permute(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}
	if len(nums) == 1 {
		return [][]int{{nums[0]}}
	}

	var ret [][]int
	for _, p := range permute(nums[1:]) {
		for i := range len(p) + 1 {
			var pi []int
			pi = append(pi, p[:i]...)
			pi = append(pi, nums[0])
			pi = append(pi, p[i:]...)
			ret = append(ret, pi)
		}
	}

	return ret
}
