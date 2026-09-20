package structures

import (
	"fmt"
	"strings"
)

// FormatInts 把层序数组里的 NULL 哨兵还原成 "null"，用于测试失败时的报错信息。
//
// Go 的 []int 装不了 nil，所以建树用的空节点是哨兵 NULL = -1<<63。直接 %v 打出来
// 是 [1 -9223372036854775808 2 3]，看不出对应题目里的哪个输入；Python 和 Java
// 那边显示的是 None / null，三门语言的报错对不上。
//
//	structures.FormatInts([]int{1, NULL, 2, 3})  →  "[1,null,2,3]"
//
// 输出格式和 LeetCode 题面里的写法一致，可以直接和题目给的 Input 对照。
func FormatInts(ints []int) string {
	parts := make([]string, 0, len(ints))
	for _, v := range ints {
		if v == NULL {
			parts = append(parts, "null")
			continue
		}
		parts = append(parts, fmt.Sprintf("%d", v))
	}
	return "[" + strings.Join(parts, ",") + "]"
}
