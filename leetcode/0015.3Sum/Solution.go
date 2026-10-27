package leetcode

import "sort"

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func threeSum(nums []int) [][]int {
	sort.Ints(nums)
	n := len(nums)
	ans := [][]int{}
	for i := 0; i < n-2; i++ {
		x := nums[i]
		if i > 0 && x == nums[i-1] { // 第一个数去重
			continue
		}
		if x+nums[i+1]+nums[i+2] > 0 { // 最小的三个都 > 0，后面只会更大
			break
		}
		if x+nums[n-2]+nums[n-1] < 0 { // x 配最大的两个都 < 0，x 太小
			continue
		}
		j, k := i+1, n-1
		for j < k {
			s := x + nums[j] + nums[k]
			if s < 0 {
				j++
			} else if s > 0 {
				k--
			} else {
				ans = append(ans, []int{x, nums[j], nums[k]})
				for j++; j < k && nums[j] == nums[j-1]; j++ { // 第二个数去重
				}
				for k--; k > j && nums[k] == nums[k+1]; k-- { // 第三个数去重
				}
			}
		}
	}
	return ans
}
