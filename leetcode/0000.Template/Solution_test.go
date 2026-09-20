package leetcode

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	t.Skip("骨架还没有题解，写完后删掉这一行")

	// 表驱动测试：一行一个用例，三列分别是用例名、输入、期望输出。
	//
	// name 直接写题目里的 Example 几，t.Run 会逐条报告通过与否，失败时也带上
	// 是哪个用例；还能用 go test -run TestSolve/Example_1 单独跑某一个。
	//
	// 用例照着该题 README 里 ## 题目 那节的 Example 转录，题目给几个就写几个。
	qs := []struct {
		name string
		in   []int
		want []int
	}{
		// {"Example 1", []int{2, 7, 11, 15}, []int{0, 1}},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := solve(q.in)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("solve(%v) = %v, want %v", q.in, got, q.want)
			}
		})
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
