package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func lengthOfLongestSubstring(s string) int {
	ascii := [128]int{}
	start, end, maxLength := 0, 0, 0

	for i, ch := range s {
		end = i
		ascii[ch] += 1

		for ascii[ch] > 1 {
			ascii[s[start]] -= 1
			start += 1
		}

		maxLength = max(maxLength, end-start+1)
	}

	return maxLength
}

func lengthOfLongestSubstringMap(s string) int {
	m := make(map[uint8]int) // for i, ch(int 32) := range s
	start := 0
	end := 0
	maxLength := 0

	for i := range len(s) {
		ch := s[i]
		index, ok := m[ch]
		// 记录ch的最后一次出现
		// 非重复子串只能是ch上次出现的坐标后一个字符到当前字符
		m[ch] = i
		end = i

		if ok {
			// 重复字符出现时可能需要start右移到index+1
			// start本来就在index右侧的时候不需要移动
			start = max(start, index+1)
		}

		currentLength := end - start + 1
		maxLength = max(maxLength, currentLength)
	}

	return maxLength
}

// 重复出现字符的下一个字符到当前字符都记录一遍
func lengthOfLongestSubstringTimeout(s string) int {
	if len(s) == 0 {
		return 0
	}

	maxLength := 1
	// 注意考虑空字符串的情况
	m := map[uint8]int{s[0]: 0}

	for i := 1; i < len(s); i++ {
		ch := s[i]
		_, ok := m[ch]

		if ok {
			j := m[ch] + 1
			m = make(map[uint8]int)
			for j <= i {
				m[s[j]] = j
				j++
			}
			// 有重复字符的时候最长子串肯定没变,不需要计算maxLength
			continue
		}

		m[s[i]] = i
		maxLength = max(len(m), maxLength)
	}

	return maxLength
}
