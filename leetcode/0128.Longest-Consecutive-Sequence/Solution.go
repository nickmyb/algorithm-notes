package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func longestConsecutive(nums []int) int {
	m := make(map[int]struct{})

	for _, num := range nums {
		m[num] = struct{}{}
	}

	longest := 0
	count := 0

	for k, _ := range m {
		if _, ok := m[k-1]; ok {
			continue
		}

		for _, ok := m[k]; ok; {
			count += 1
			k += 1
			_, ok = m[k]
		}
		longest = max(longest, count)
		count = 0
	}

	return longest
}
