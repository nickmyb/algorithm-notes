package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func maxArea(height []int) int {
	left, right := 0, len(height)-1
	m := area(height, left, right)

	l, r := left, right
	for l < r && 0 <= r-1 && l+1 <= len(height)-1 {
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
		m = max(area(height, l, r), m)
	}

	return m
}

func area(height []int, left int, right int) int {
	h := min(height[left], height[right])
	return h * (right - left)
}
