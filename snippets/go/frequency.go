package snippets

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
