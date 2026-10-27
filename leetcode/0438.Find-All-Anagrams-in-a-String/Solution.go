package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func findAnagrams(s string, p string) []int {
	var ret []int

	lettersP := countLetters(p)

	var letters [26]int
	i := -1

	for j := 0; j < len(s); j++ {
		cIndex := s[j] - 'a'
		letters[cIndex] += 1

		for letters[cIndex] > lettersP[cIndex] {
			i += 1
			letters[s[i]-'a'] -= 1
		}

		if letters == lettersP {
			ret = append(ret, i+1)
		}
	}

	return ret
}

// region @snippets/go/anagram.go

// countLetters 统计字符串中每个小写字母出现的次数。
// 返回长度为 26 的数组，下标 0~25 依次对应 'a'~'z'。
// 调用方需保证 str 只包含小写字母 a-z，否则下标越界会 panic。
func countLetters(str string) [26]int {
	letters := [26]int{}

	for _, s := range str {
		letters[s-'a'] += 1
	}

	return letters
}

// endregion @snippets/go/anagram.go
