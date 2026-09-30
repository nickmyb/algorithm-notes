package leetcode

import (
	"testing"
)

func TestMinStack(t *testing.T) {
	// 设计类题目：ops / args / want 三列照抄题目 Input 和 Output 的三行数组，
	// want 里的 null 写成 nil。第一个操作固定是构造函数，从第二个开始逐个回放。
	qs := []struct {
		name string
		ops  []string
		args [][]int
		want []any
	}{
		{
			"Example 1",
			[]string{"MinStack", "push", "push", "push", "getMin", "pop", "top", "getMin"},
			[][]int{{}, {-2}, {0}, {-3}, {}, {}, {}, {}},
			[]any{nil, nil, nil, nil, -3, nil, 0, -2},
		},
		{
			"Leetcode 33",
			[]string{"MinStack", "push", "push", "push", "push", "getMin", "pop", "getMin", "pop", "getMin", "pop", "getMin"},
			[][]int{{}, {2}, {0}, {3}, {0}, {}, {}, {}, {}, {}, {}, {}},
			[]any{nil, nil, nil, nil, nil, 0, nil, 0, nil, 0, nil, 2},
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			ms := Constructor()

			for i := 1; i < len(q.ops); i++ {
				var got any
				switch q.ops[i] {
				case "push":
					ms.Push(q.args[i][0])
				case "pop":
					ms.Pop()
				case "top":
					got = ms.Top()
				case "getMin":
					got = ms.GetMin()
				default:
					t.Fatalf("未知操作 %q", q.ops[i])
				}

				// 后面的操作依赖前面的状态，第一处不一致就停下
				if got != q.want[i] {
					t.Fatalf("第 %d 步 %s(%v) = %v, want %v", i, q.ops[i], q.args[i], got, q.want[i])
				}
			}
		})
	}
}
