package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

// searchRange 分别计算左右界
func searchRange(nums []int, target int) []int {
	lo := searchRangeLo(nums, target)
	hi := searchRangeHi(nums, target)

	return []int{lo, hi}
}

func searchRangeLo(nums []int, target int) int {
	lo, hi, ret := 0, len(nums)-1, -1

	for lo <= hi {
		mid := (lo + hi) / 2

		if nums[mid] == target {
			ret = mid
			hi = mid - 1
		} else if nums[mid] > target {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	return ret
}

func searchRangeHi(nums []int, target int) int {
	lo, hi, ret := 0, len(nums)-1, -1

	for lo <= hi {
		mid := (lo + hi) / 2

		if nums[mid] == target {
			ret = mid
			lo = mid + 1
		} else if nums[mid] > target {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	return ret
}

// searchRangeRecursive 递归
func searchRangeRecursive(nums []int, target int, left int, right int) []int {
	index := binarySearchRange(target, nums, left, right)
	if index == -1 {
		return []int{-1, -1}
	}

	return []int{betterIndex(index, searchRangeRecursive(nums, target, left, index-1)[0], -1), betterIndex(index, searchRangeRecursive(nums, target, index+1, right)[1], 1)}
}

func binarySearchRange(key int, a []int, left int, right int) int {
	lo, hi := left, right

	for lo <= hi {
		mid := (lo + hi) / 2

		if key == a[mid] {
			return mid
		} else if key < a[mid] {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return -1
}

func betterIndex(index int, bIndex int, t int) int {
	if bIndex == -1 {
		return index
	}
	if t == 1 {
		return max(index, bIndex)
	}
	if t == -1 {
		return min(index, bIndex)
	}

	return -1
}
