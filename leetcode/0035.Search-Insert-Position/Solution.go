package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func searchInsert(nums []int, target int) int {
	lo, hi, mid := 0, len(nums)-1, -1

	for lo <= hi {
		mid = (lo + hi) / 2

		if target == nums[mid] {
			return mid
		} else if target < nums[mid] {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	if target > nums[mid] {
		return mid + 1
	}
	return mid
}
