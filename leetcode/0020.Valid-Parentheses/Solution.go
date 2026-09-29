package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
// [] {} 在 ascii 不是连续的!!!
func isValid(s string) bool {
	var slice []int32

	for _, ch := range s {
		if ch == '(' || ch == '[' || ch == '{' {
			slice = append(slice, ch)
		}

		if ch == ')' || ch == ']' || ch == '}' {
			if len(slice) < 1 {
				return false
			}
			lastCh := slice[len(slice)-1]

			if ch == ')' {
				if lastCh != '(' {
					return false
				}
			}
			if ch == ']' {
				if lastCh != '[' {
					return false
				}
			}
			if ch == '}' {
				if lastCh != '{' {
					return false
				}
			}

			slice = slice[:len(slice)-1]
		}
	}

	return len(slice) == 0
}
