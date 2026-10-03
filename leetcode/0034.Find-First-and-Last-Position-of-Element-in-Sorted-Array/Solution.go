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

// 红蓝染色法解释:
// 红蓝的定义: images/0034.Find-First-and-Last-Position-of-Element-in-Sorted-Array/red-blue-definition.png
// 三种区间的完整过程: images/0034.Find-First-and-Last-Position-of-Element-in-Sorted-Array/three-intervals.png

func searchRangeRedBlue(nums []int, target int) []int {
	lo := lowerBoundClosed(nums, target)
	if lo == len(nums) || nums[lo] != target {
		return []int{-1, -1}
	}

	hi := lowerBoundClosed(nums, target+1) - 1

	return []int{lo, hi}
}

// lowerBoundClosed 二分搜索 红蓝染色法 闭区间 [lo, hi], 返回第一个>=target的index
func lowerBoundClosed(nums []int, target int) int {
	lo, hi := 0, len(nums)-1

	// [lo, hi]是没有确定颜色的区间,区间要求非空
	// 区间选择决定循环退出条件
	for lo <= hi {
		mid := (lo + hi) / 2

		// 在有序数组上，这个条件必须是
		// 前面一段都成立，后面一段都不成立，即「红红红…蓝蓝蓝」。这样看一个 mid 才能确定一整段的颜色。< 和 <= 都满足这个前提
		// 且不能是 前面一段不成立 后面成立
		// 要求: 左红右蓝
		// 第一个蓝 = 第一个 if条件 不成立 的位置
		// 最后一个红 = 最后一个 if条件 成立 的位置
		if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	// 不变量
	// 区间右边界外的第一格: 一定是蓝
	// 区间左边界外的第一格: 一定是红
	// 因为不变量的存在,所以可以根据初始值确定红蓝的位置
	// [lo, hi]: 红 = lo - 1; 蓝 = hi + 1
	// [lo, hi): 红 = lo - 1; 蓝 = hi
	// (lo, hi): 红 = lo    ; 蓝 = hi

	// 返回值 可以根据需要返回 第一个蓝或者最后一个红 都可以
	return hi + 1
}

// lowerBoundHalfOpen [lo, hi)
func lowerBoundHalfOpen(nums []int, target int) int {
	lo, hi := 0, len(nums)

	for lo < hi {
		mid := (lo + hi) / 2

		if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return hi
}

// lowerBoundOpen (lo, hi)
func lowerBoundOpen(nums []int, target int) int {
	lo, hi := -1, len(nums)

	for lo+1 < hi {
		mid := (lo + hi) / 2

		if nums[mid] < target {
			lo = mid
		} else {
			hi = mid
		}
	}

	return hi
}
