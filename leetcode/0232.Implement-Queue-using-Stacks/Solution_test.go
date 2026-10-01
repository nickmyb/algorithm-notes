package leetcode

import (
	"testing"
)

func TestMyQueue(t *testing.T) {
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
			[]string{"MyQueue", "push", "push", "peek", "pop", "empty"},
			[][]int{{}, {1}, {2}, {}, {}, {}},
			[]any{nil, nil, nil, 1, 1, false},
		},
		{
			// 出队后 queue 里还剩 2 时再 push 3：queue 非空也倒栈的实现会让 3 排到 2 前面
			"交替push pop操作",
			[]string{"MyQueue", "push", "push", "pop", "push", "pop", "pop", "empty"},
			[][]int{{}, {1}, {2}, {}, {3}, {}, {}, {}},
			[]any{nil, nil, nil, 1, nil, 2, 3, true},
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			mq := Constructor()

			for i := 1; i < len(q.ops); i++ {
				var got any
				switch q.ops[i] {
				case "push":
					mq.Push(q.args[i][0])
				case "pop":
					got = mq.Pop()
				case "peek":
					got = mq.Peek()
				case "empty":
					got = mq.Empty()
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
