package leetcode

import (
	"reflect"
	"testing"
)

// 表驱动测试：每个 question 是一组「入参 + 期望结果」，往 qs 里加一行就是加一个用例。
type question struct {
	para para
	ans  ans
}

// para 是入参
type para struct {
	nums []int
}

// ans 是期望结果
type ans struct {
	one []int
}

func TestSolve(t *testing.T) {
	t.Skip("骨架还没有题解，写完后删掉这一行")

	qs := []question{
		// {para{[]int{2, 7, 11, 15}}, ans{[]int{0, 1}}},
	}

	// 一题有多种解法时，把每种实现登记进来，同一组用例跑全部实现，
	// 顺带保证它们结果一致。
	impls := map[string]func([]int) []int{
		"solve": solve,
		// "solveBruteForce": solveBruteForce,
	}

	for name, solve := range impls {
		t.Run(name, func(t *testing.T) {
			for _, q := range qs {
				got := solve(q.para.nums)
				if !reflect.DeepEqual(got, q.ans.one) {
					t.Fatalf("%v(%v) = %v, want %v", name, q.para.nums, got, q.ans.one)
				}
			}
		})
	}
}
