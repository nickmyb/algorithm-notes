package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func countWindowLetters(str string) [128]int {
	letters := [128]int{}

	for _, s := range str {
		letters[s-'A'] += 1
	}

	return letters
}

func containLetters(s, t [128]int) bool {
	for i := range 128 {
		if s[i] < t[i] {
			return false
		}
	}

	return true
}

func minWindow(s string, t string) string {
	lettersT := countWindowLetters(t)
	letters := [128]int{}
	subS := ""

	i := 0
	for j := 0; j < len(s); j++ {
		cIndex := s[j] - 'A'
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
				letters[s[i]-'A'] -= 1
			}
		}
	}

	return subS
}
