package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

// region @snippets/go/frequency.go

// countAscii 统计字符串中每个 ASCII 字符出现的次数。
// 返回长度为 128 的数组，直接以 ASCII 码作为下标，区分大小写。
// 调用方需保证 str 只包含 ASCII 字节（0~127），否则下标越界会 panic。
func countAscii(str string) [128]int {
	counts := [128]int{}

	for i := range len(str) {
		counts[str[i]] += 1
	}

	return counts
}

// endregion @snippets/go/frequency.go

func containLetters(s, t [128]int) bool {
	for i := range 128 {
		if s[i] < t[i] {
			return false
		}
	}

	return true
}

func minWindow(s string, t string) string {
	need := countAscii(t) // 目标计数减去窗口计数，负数表示富余。
	missing := len(t)     // 尚缺的字符总数，重复字符分别计数。
	left := 0
	bestStart, bestLen := 0, len(s)+1

	for right := 0; right < len(s); right++ {
		c := s[right]
		if need[c] > 0 {
			missing--
		}
		need[c]--

		for missing == 0 {
			if length := right - left + 1; length < bestLen {
				bestStart, bestLen = left, length
			}

			c := s[left]
			need[c]++
			if need[c] > 0 {
				missing++
			}
			left++
		}
	}

	if bestLen > len(s) {
		return ""
	}
	return s[bestStart : bestStart+bestLen]
}

func minWindowScanCounts(s string, t string) string {
	lettersT := countAscii(t)
	letters := [128]int{}
	subS := ""

	i := 0
	for j := 0; j < len(s); j++ {
		cIndex := s[j]
		letters[cIndex] += 1

		if letters[cIndex] < lettersT[cIndex] {
		} else if containLetters(letters, lettersT) {
			if subS == "" {
				subS = s[i : j+1]
			}
			if len(s[i:j+1]) < len(subS) {
				subS = s[i : j+1]
			}

			for ; i <= j && containLetters(letters, lettersT); i++ {
				sub := s[i : j+1]
				if len(sub) < len(subS) {
					subS = sub
				}
				letters[s[i]] -= 1
			}
		}
	}

	return subS
}
