package leetcode

import (
	"reflect"
	"slices"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	qs := []struct {
		name string
		s    string
		p    string
		want []int
	}{
		{"Example 1", "cbaebabacd", "abc", []int{0, 6}},
		{"Example 2", "abab", "ab", []int{0, 1, 2}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := findAnagrams(q.s, q.p)
			// 题目允许任意顺序返回索引；排序后比较，保留重复项。
			gotSorted := slices.Clone(got)
			slices.Sort(gotSorted)
			if !reflect.DeepEqual(gotSorted, q.want) {
				t.Fatalf("findAnagrams(%q, %q) = %v, want %v", q.s, q.p, got, q.want)
			}
		})
	}
}
