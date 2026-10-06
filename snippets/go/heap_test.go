package snippets

import (
	"slices"
	"testing"
)

// checkHeap 检查 MinPQ 的内部结构，每次 Add / RemoveSmallest 之后都调用。
// 只看出队顺序抓不到指针没接好的错误：交换节点时漏改孩子的 parent、
// 漏改 dir、漏更新 pq.root，往往要再操作几次才炸，炸的位置离出错的地方很远。
func checkHeap(t *testing.T, pq *MinPQ) {
	t.Helper()
	if pq.root != nil && pq.root.parent != nil {
		t.Fatalf("root.parent = %p, want nil", pq.root.parent)
	}

	// 层序遍历：完全二叉树里 nil 之后不能再出现节点
	count := 0
	seenNil := false
	queue := []*hNode{pq.root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if n == nil {
			seenNil = true
			continue
		}
		if seenNil {
			t.Fatalf("not a complete tree: node %d after a gap", n.key)
		}
		count++
		for _, c := range []struct {
			child *hNode
			dir   int
		}{{n.left, -1}, {n.right, 1}} {
			if c.child == nil {
				queue = append(queue, nil)
				continue
			}
			if c.child.parent != n {
				t.Fatalf("node %d: parent pointer not pointing to %d", c.child.key, n.key)
			}
			if c.child.dir != c.dir {
				t.Fatalf("node %d: dir = %d, want %d", c.child.key, c.child.dir, c.dir)
			}
			if c.child.key < n.key {
				t.Fatalf("heap order violated: child %d < parent %d", c.child.key, n.key)
			}
			queue = append(queue, c.child)
		}
	}
	if count != pq.size {
		t.Fatalf("tree has %d nodes, size = %d", count, pq.size)
	}
}

// addAll 依次 Add，并在每次之后检查结构
func addAll(t *testing.T, pq *MinPQ, xs []int) {
	t.Helper()
	for _, x := range xs {
		pq.Add(x)
		checkHeap(t, pq)
	}
}

// drain 不断 RemoveSmallest 直到空，返回出队序列
func drain(t *testing.T, pq *MinPQ) []int {
	t.Helper()
	var got []int
	for pq.Size() > 0 {
		x, ok := pq.RemoveSmallest()
		if !ok {
			t.Fatalf("RemoveSmallest ok = false with size %d", pq.Size()+1)
		}
		got = append(got, x)
		checkHeap(t, pq)
	}
	return got
}

func TestMinPQEmpty(t *testing.T) {
	// 空堆上 root 是 nil，直接取 root.key 会空指针 panic
	pq := NewMinPQ()
	if pq.Size() != 0 {
		t.Fatalf("Size = %d, want 0", pq.Size())
	}
	if _, ok := pq.GetSmallest(); ok {
		t.Fatal("GetSmallest on empty: ok = true, want false")
	}
	if _, ok := pq.RemoveSmallest(); ok {
		t.Fatal("RemoveSmallest on empty: ok = true, want false")
	}
	checkHeap(t, pq)
}

func TestMinPQDrainSorted(t *testing.T) {
	// n 从 1 取到 33：跨过 1、2、4、8、16、32 这些层的边界。
	// position 把下标换成路径时进制或起点（size 还是 size+1）写错，
	// 小 n 可能碰巧对，一过某个层边界就挂到错误的父节点上。
	descending := func(n int) []int {
		xs := make([]int, n)
		for i := range xs {
			xs[i] = n - i
		}
		return xs
	}
	ascending := func(n int) []int {
		xs := descending(n)
		slices.Reverse(xs)
		return xs
	}
	// 固定的乱序：用线性同余生成，保证每次运行一样
	shuffled := func(n int) []int {
		xs := make([]int, n)
		seed := 7
		for i := range xs {
			seed = (seed*1103515245 + 12345) % (1 << 31)
			xs[i] = seed % 50
		}
		return xs
	}

	inputs := []struct {
		name string
		gen  func(int) []int
	}{
		// 每个新元素都是最小的，swapUp 必须一路换到根：
		// 抓「只换一层」、「到根时没判 parent == nil」、「换到根后没更新 pq.root」
		{"descending", descending},
		// Add 时一次都不用换；出队时把最后一个大元素放到根，swapDown 必须一路换到底
		{"ascending", ascending},
		// 前两种只走最左或最右的一条路径，乱序才会走到中间的子树
		{"shuffled", shuffled},
	}
	for _, in := range inputs {
		t.Run(in.name, func(t *testing.T) {
			for n := 1; n <= 33; n++ {
				xs := in.gen(n)
				pq := NewMinPQ()
				addAll(t, pq, xs)
				if pq.Size() != n {
					t.Fatalf("n=%d: Size = %d after adds", n, pq.Size())
				}
				want := slices.Clone(xs)
				slices.Sort(want)
				if got, _ := pq.GetSmallest(); got != want[0] {
					t.Fatalf("n=%d: GetSmallest = %d, want %d", n, got, want[0])
				}
				if got := drain(t, pq); !slices.Equal(got, want) {
					t.Fatalf("n=%d: drain(%v) = %v, want %v", n, xs, got, want)
				}
			}
		})
	}
}

func TestMinPQEqualChildren(t *testing.T) {
	// 出队 1 后，把 9 放到根，它的两个孩子都是 2。
	// swapDown 选较小孩子时若写成 `if l < r {…} else if r < l {…}`，
	// 平局两个分支都不进，9 就停在根上。
	//
	//	    1
	//	   / \
	//	  2   2
	//	 /
	//	9
	pq := NewMinPQ()
	addAll(t, pq, []int{1, 2, 2, 9})
	if got := drain(t, pq); !slices.Equal(got, []int{1, 2, 2, 9}) {
		t.Fatalf("drain = %v, want [1 2 2 9]", got)
	}
}

func TestMinPQInterleaved(t *testing.T) {
	// 前面的用例都是先全部 Add 再全部出队。交替进行时，
	// Add 用的「下一个空位」依赖 RemoveSmallest 正确地摘掉了最后一个节点、size 减对了；
	// 否则新节点会挂到已被摘走的节点下面，或者覆盖一个还在的节点。
	pq := NewMinPQ()
	ops := []struct {
		add  int  // pop 为 false 时 Add 这个值
		pop  bool // true 表示 RemoveSmallest
		want int  // pop 时期望出队的值
	}{
		{add: 5}, {add: 3}, {add: 8},
		{pop: true, want: 3},
		{add: 1}, {add: 4},
		{pop: true, want: 1},
		{pop: true, want: 4},
		{add: 2}, {add: 7}, {add: 6},
		{pop: true, want: 2},
		{pop: true, want: 5},
		{pop: true, want: 6},
		{pop: true, want: 7},
		{pop: true, want: 8},
	}
	for i, op := range ops {
		if op.pop {
			got, ok := pq.RemoveSmallest()
			if !ok || got != op.want {
				t.Fatalf("op %d: RemoveSmallest = (%d, %v), want (%d, true)", i, got, ok, op.want)
			}
		} else {
			pq.Add(op.add)
		}
		checkHeap(t, pq)
	}
	if pq.Size() != 0 {
		t.Fatalf("Size = %d at end, want 0", pq.Size())
	}
}

func TestMinPQReuseAfterEmpty(t *testing.T) {
	// 取走最后一个元素时要把 pq.root 置回 nil。
	// 漏了的话 GetSmallest 还会返回旧值，下一次 Add 会挂到一个已删除的节点下面。
	pq := NewMinPQ()
	addAll(t, pq, []int{4})
	if got := drain(t, pq); !slices.Equal(got, []int{4}) {
		t.Fatalf("drain = %v, want [4]", got)
	}
	if x, ok := pq.GetSmallest(); ok {
		t.Fatalf("GetSmallest after emptying = (%d, true), want ok = false", x)
	}
	addAll(t, pq, []int{9, 6})
	if got := drain(t, pq); !slices.Equal(got, []int{6, 9}) {
		t.Fatalf("drain = %v, want [6 9]", got)
	}
}
