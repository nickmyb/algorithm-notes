package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func subsets(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	if len(nums) == 1 {
		return [][]int{[]int{}, []int{nums[0]}}
	}

	var ret [][]int
	for _, subset := range subsets(nums[1:]) {
		ret = append(ret, subset)
		ret = append(ret, append([]int{nums[0]}, subset...))
	}

	return ret
}
