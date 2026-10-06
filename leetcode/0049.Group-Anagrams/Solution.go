package leetcode

import "sort"

func groupAnagrams(strs []string) [][]string {
	m := make(map[[26]int][]string)
	var ret [][]string

	for _, str := range strs {
		k := countLetters(str)
		m[k] = append(m[k], str)
	}

	for _, v := range m {
		ret = append(ret, v)
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

func groupAnagramsSorted(strs []string) [][]string {
	m := make(map[string][]string)
	var ret [][]string

	for _, str := range strs {
		anagram := anagrammatize(str)
		m[anagram] = append(m[anagram], str)
	}

	for _, v := range m {
		ret = append(ret, v)
	}

	return ret
}

// region @snippets/go/anagram.go

// anagrammatize 将字符串按字符升序排序，返回其规范形式。
// 互为变位词的字符串排序后结果相同，因此可用作变位词分组的 key。
// 按 rune 排序，支持多字节字符（如中文）。
func anagrammatize(str string) string {
	r := []rune(str)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

// endregion @snippets/go/anagram.go
