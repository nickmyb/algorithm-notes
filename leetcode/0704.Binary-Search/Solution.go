package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

func search(nums []int, target int) int {
	return binarySearch(target, nums)
}

// binarySearch 二分查找
func binarySearch(key int, a []int) int {
	lo, hi := 0, len(a)-1

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
