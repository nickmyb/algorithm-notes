package leetcode

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func isPalindrome(s string) bool {
	length := len(s)
	start := 0
	end := length - 1

	for start < length && end >= 0 && start <= end {
		if !isAlnum(s[start]) {
			start += 1
			continue
		}
		if !isAlnum(s[end]) {
			end -= 1
			continue
		}

		if chToLower(s[start]) != chToLower(s[end]) {
			return false
		}

		start += 1
		end -= 1
	}

	return true
}

func isAlnum(ch uint8) bool {
	if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') {
		return false
	}

	return true
}

func chToLower(ch uint8) uint8 {
	if ch >= 'A' && ch <= 'Z' {
		return ch + 'a' - 'A'
	}

	return ch
}
