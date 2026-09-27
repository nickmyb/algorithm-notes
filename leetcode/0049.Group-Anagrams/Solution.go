package leetcode

import "sort"

func groupAnagrams(strs []string) [][]string {
	m := make(map[[26]int][]string)
	var ret [][]string

	for _, str := range strs {
		k := string2Alphabet(str)
		m[k] = append(m[k], str)
	}

	for _, v := range m {
		ret = append(ret, v)
	}

	return ret
}

func string2Alphabet(str string) [26]int {
	alphabet := [26]int{}

	for _, s := range str {
		alphabet[s-'a'] += 1
	}

	return alphabet
}

func groupAnagramsSorted(strs []string) [][]string {
	m := make(map[string][]string)
	var ret [][]string

	for _, str := range strs {
		k := []rune(str)
		sort.Slice(k, func(i, j int) bool {
			return k[i] < k[j]
		})
		sk := string(k)
		m[sk] = append(m[sk], str)
	}

	for _, v := range m {
		ret = append(ret, v)
	}

	return ret
}
