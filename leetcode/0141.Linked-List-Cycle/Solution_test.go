package leetcode

import (
	"testing"

	structures "github.com/nickmyb/algorithm-notes/structures/go"
)

func TestHasCycle(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		pos  int
		want bool
	}{
		{"Example 1", []int{3, 2, 0, -4}, 1, true},
		{"Example 2", []int{1, 2}, 0, true},
		{"Example 3", []int{1}, -1, false},
	}

	impls := map[string]func(*ListNode) bool{
		"hasCycle":    hasCycle,
		"hasCycleMap": hasCycleMap,
	}

	for _, q := range qs {
		for implName, impl := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := impl(structures.Ints2ListWithCycle(q.in, q.pos))
				if got != q.want {
					t.Fatalf("%v(%v, pos=%v) = %v, want %v", implName, structures.FormatInts(q.in), q.pos, got, q.want)
				}
			})
		}
	}
}

// 输入是树或链表时，用 structures.Ints2TreeNode / Ints2List 建结构，
// 报错里的输入用 structures.FormatInts 打印——它把空节点的哨兵还原成 null，
// 否则 %v 会打出 -9223372036854775808，和题面里的 [1,null,2,3] 对不上：
//
//	got := maxDepth(structures.Ints2TreeNode(q.in))
//	t.Fatalf("maxDepth(%v) = %v, want %v", structures.FormatInts(q.in), got, q.want)
//
// 一题写了多种解法时，把它们登记进一个 map，同一组用例跑全部实现，
// 顺带保证几种解法结果一致。子测试名会变成 Example_1/solveBruteForce：
//
//	impls := map[string]func([]int) []int{
//		"solve":           solve,
//		"solveBruteForce": solveBruteForce,
//	}
//	for _, q := range qs {
//		for implName, solve := range impls {
//			t.Run(q.name+"/"+implName, func(t *testing.T) { ... })
//		}
//	}
