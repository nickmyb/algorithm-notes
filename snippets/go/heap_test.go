package snippets

import (
	"cmp"
	"slices"
	"testing"
)

// checkHeap 检查 MaxPQ 的内部状态，每次 Insert / DelMax 之后都调用。
// 只看出队顺序时，堆序被破坏往往要再操作几次才表现出来，离出错的那一步很远。
func checkHeap(t *testing.T, pq *MaxPQ[int]) {
	t.Helper()
	// pq[0] 空着不用，所以切片长度总是 n+1；
	// DelMax 没截掉最后一格的话，下一次 Insert 会 append 到 n+2，中间留下旧值
	if len(pq.pq) != pq.n+1 {
		t.Fatalf("len(pq) = %d, want n+1 = %d", len(pq.pq), pq.n+1)
	}
	for i := 2; i <= pq.n; i++ {
		if pq.cmp(pq.pq[i/2], pq.pq[i]) < 0 {
			t.Fatalf("heap order violated: pq[%d] = %d < pq[%d] = %d", i/2, pq.pq[i/2], i, pq.pq[i])
		}
	}
}

// insertAll 依次 Insert，并在每次之后检查结构
func insertAll(t *testing.T, pq *MaxPQ[int], xs []int) {
	t.Helper()
	for _, x := range xs {
		pq.Insert(x)
		checkHeap(t, pq)
	}
}

// drain 不断 DelMax 直到空，返回出队序列
func drain(t *testing.T, pq *MaxPQ[int]) []int {
	t.Helper()
	var got []int
	for !pq.IsEmpty() {
		x, ok := pq.DelMax()
		if !ok {
			t.Fatalf("DelMax ok = false with size %d", pq.Size()+1)
		}
		got = append(got, x)
		checkHeap(t, pq)
	}
	return got
}

func TestMaxPQEmpty(t *testing.T) {
	// 空堆上 Max / DelMax 要报 ok = false，而不是去读 pq[1] 越界 panic
	pq := NewMaxPQ(cmp.Compare[int])
	if !pq.IsEmpty() || pq.Size() != 0 {
		t.Fatalf("IsEmpty = %v, Size = %d, want true, 0", pq.IsEmpty(), pq.Size())
	}
	if _, ok := pq.Max(); ok {
		t.Fatal("Max on empty: ok = true, want false")
	}
	if _, ok := pq.DelMax(); ok {
		t.Fatal("DelMax on empty: ok = true, want false")
	}
	checkHeap(t, pq)

	// 删空之后同理：切片截回只剩 pq[0]，Max 不能返回刚删掉的 7
	pq.Insert(7)
	pq.DelMax()
	if x, ok := pq.Max(); ok {
		t.Fatalf("Max after emptying = (%d, true), want ok = false", x)
	}
}

func TestMaxPQDrainSorted(t *testing.T) {
	// n 从 1 取到 33：跨过 1、2、4、8、16、32 这些层的边界，
	// 最后一个节点分别落在左孩子和右孩子上，sink 的两个边界判断都会走到。
	ascending := func(n int) []int {
		xs := make([]int, n)
		for i := range xs {
			xs[i] = i + 1
		}
		return xs
	}
	descending := func(n int) []int {
		xs := ascending(n)
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
		// 每个新元素都是最大的，swim 必须一路换到根：抓「只换一层」、循环条件写成 k > 2
		{"ascending", ascending},
		// Insert 时一次都不用换；出队时把最后一个小元素放到根，sink 必须一路换到底
		{"descending", descending},
		// 前两种只走最左或最右的一条路径，乱序才会走到中间的子树
		{"shuffled", shuffled},
	}
	for _, in := range inputs {
		t.Run(in.name, func(t *testing.T) {
			for n := 1; n <= 33; n++ {
				xs := in.gen(n)
				pq := NewMaxPQ(cmp.Compare[int])
				insertAll(t, pq, xs)
				if pq.Size() != n {
					t.Fatalf("n=%d: Size = %d after inserts", n, pq.Size())
				}
				want := slices.Clone(xs)
				slices.Sort(want)
				slices.Reverse(want)
				if got, ok := pq.Max(); !ok || got != want[0] {
					t.Fatalf("n=%d: Max = (%d, %v), want (%d, true)", n, got, ok, want[0])
				}
				if got := drain(t, pq); !slices.Equal(got, want) {
					t.Fatalf("n=%d: drain(%v) = %v, want %v", n, xs, got, want)
				}
			}
		})
	}
}

func TestMaxPQOnlyLeftChild(t *testing.T) {
	// 出队 9 后 n = 2，把 3 放到根，它只有左孩子 5，而且 5 就是最后一个节点。
	// sink 的循环条件写成 2*k < n 会直接退出，3 停在根上，下一次出队得到 3。
	//
	//	  9            3
	//	 / \    →     /
	//	5   3        5
	pq := NewMaxPQ(cmp.Compare[int])
	insertAll(t, pq, []int{9, 5, 3})
	if got := drain(t, pq); !slices.Equal(got, []int{9, 5, 3}) {
		t.Fatalf("drain = %v, want [9 5 3]", got)
	}
}

func TestMaxPQEqualChildren(t *testing.T) {
	// 出队 9 后，把 1 放到根，它的两个孩子都是 2。
	// sink 选较大孩子时若写成 `if l > r {…} else if r > l {…}`，
	// 平局两个分支都不进，1 就停在根上。
	//
	//	    9
	//	   / \
	//	  2   2
	//	 /
	//	1
	pq := NewMaxPQ(cmp.Compare[int])
	insertAll(t, pq, []int{9, 2, 2, 1})
	if got := drain(t, pq); !slices.Equal(got, []int{9, 2, 2, 1}) {
		t.Fatalf("drain = %v, want [9 2 2 1]", got)
	}
}

func TestMaxPQInterleaved(t *testing.T) {
	// 前面的用例都是先全部 Insert 再全部出队。交替进行时，
	// Insert 放的位置依赖 DelMax 正确地截掉了最后一格、n 减对了；
	// 否则新元素会排在旧值后面，或者 n 和切片长度对不上。
	pq := NewMaxPQ(cmp.Compare[int])
	ops := []struct {
		add  int  // pop 为 false 时 Insert 这个值
		pop  bool // true 表示 DelMax
		want int  // pop 时期望出队的值
	}{
		{add: 5}, {add: 7}, {add: 2},
		{pop: true, want: 7},
		{add: 9}, {add: 6},
		{pop: true, want: 9},
		{pop: true, want: 6},
		{add: 8}, {add: 3}, {add: 4},
		{pop: true, want: 8},
		{pop: true, want: 5},
		{pop: true, want: 4},
		{pop: true, want: 3},
		{pop: true, want: 2},
	}
	for i, op := range ops {
		if op.pop {
			got, ok := pq.DelMax()
			if !ok || got != op.want {
				t.Fatalf("op %d: DelMax = (%d, %v), want (%d, true)", i, got, ok, op.want)
			}
		} else {
			pq.Insert(op.add)
		}
		checkHeap(t, pq)
	}
	if !pq.IsEmpty() {
		t.Fatalf("Size = %d at end, want 0", pq.Size())
	}
}

func TestMaxPQReuseAfterEmpty(t *testing.T) {
	// 取走最后一个元素时切片要截回只剩 pq[0]。
	// 漏了的话下一次 Insert 会 append 到下标 2，而 n = 1 指向的 pq[1] 还是旧值。
	pq := NewMaxPQ(cmp.Compare[int])
	insertAll(t, pq, []int{4})
	if got := drain(t, pq); !slices.Equal(got, []int{4}) {
		t.Fatalf("drain = %v, want [4]", got)
	}
	insertAll(t, pq, []int{6, 9})
	if got := drain(t, pq); !slices.Equal(got, []int{9, 6}) {
		t.Fatalf("drain = %v, want [9 6]", got)
	}
}

func TestMaxPQCustomCmp(t *testing.T) {
	// 反过来的 cmp 得到最小堆。cmp 用减法，返回值不只是 -1、0、1：
	// less 里写成 cmp(a, b) == -1 的话，除了相差 1 的两个数，其余都被当成「不小于」，
	// 堆就不动了。cmp.Compare 恰好只返回 -1、0、1，前面的用例抓不到这种写法。
	pq := NewMaxPQ(func(a, b int) int { return b - a })
	insertAll(t, pq, []int{5, 1, 9, 3, 7})
	if got := drain(t, pq); !slices.Equal(got, []int{1, 3, 5, 7, 9}) {
		t.Fatalf("drain = %v, want [1 3 5 7 9]", got)
	}
}

func TestMaxPQReleasesRemoved(t *testing.T) {
	// DelMax 截掉最后一格后，那一格还在底层数组里。
	// 不清零的话，T 是指针时出队的元素一直被引用着，GC 回收不了，出队顺序上看不出来。
	a, b := 1, 2
	pq := NewMaxPQ(func(x, y *int) int { return cmp.Compare(*x, *y) })
	pq.Insert(&a)
	pq.Insert(&b)
	for pq.Size() > 0 {
		pq.DelMax()
		if p := pq.pq[:pq.n+2][pq.n+1]; p != nil {
			t.Fatalf("slot %d after DelMax = &%d, want nil", pq.n+1, *p)
		}
	}
}

func TestMaxPQCap(t *testing.T) {
	// 对照 algs4 的 MaxPQ(int max)：max 是元素个数，pq[0] 还占一格，所以容量要 max+1。
	// 写成 make([]T, 1, max) 的话，第 max 次 Insert 就会扩容，底层数组换掉。
	for _, max := range []int{0, 1, 5} {
		pq := NewMaxPQCap(cmp.Compare[int], max)
		checkHeap(t, pq)
		first, c := &pq.pq[0], cap(pq.pq)
		for i := range max {
			pq.Insert(i)
			checkHeap(t, pq)
		}
		if &pq.pq[0] != first {
			t.Fatalf("max=%d: reallocated within %d inserts, initial cap = %d, want >= %d", max, max, c, max+1)
		}
		// max 只是预留容量，不是上限：超过之后照常 append 扩容
		insertAll(t, pq, []int{100, 200})
		if got, _ := pq.Max(); got != 200 {
			t.Fatalf("max=%d: Max after exceeding = %d, want 200", max, got)
		}
	}
}

func TestMaxPQFrom(t *testing.T) {
	// 自底向上建堆：从 n/2 往前逐个 sink。n 取 0..33 跨过各层边界，
	// 抓起点写成 n/2-1（最后一个有孩子的节点没处理）、循环到 k > 1 漏掉根。
	for n := 0; n <= 33; n++ {
		a := make([]int, n)
		seed := 7
		for i := range a {
			seed = (seed*1103515245 + 12345) % (1 << 31)
			a[i] = seed % 50
		}
		orig := slices.Clone(a)

		pq := NewMaxPQFrom(cmp.Compare[int], a)
		checkHeap(t, pq)
		// 直接拿 a 当底层数组的话，sink 会把调用方的切片打乱
		if !slices.Equal(a, orig) {
			t.Fatalf("n=%d: a modified to %v, want %v", n, a, orig)
		}
		want := slices.Clone(a)
		slices.Sort(want)
		slices.Reverse(want)
		if got := drain(t, pq); len(want) > 0 && !slices.Equal(got, want) {
			t.Fatalf("n=%d: drain = %v, want %v", n, got, want)
		}
		// 建出来的堆还能继续用
		insertAll(t, pq, []int{3, 8})
		if got, _ := pq.Max(); got != 8 {
			t.Fatalf("n=%d: Max after Insert = %d, want 8", n, got)
		}
	}
}

func TestMaxPQFromLinear(t *testing.T) {
	// From 存在的理由是 O(n) 建堆。逐个 Insert 结果完全一样、只是 O(n log n)，
	// 上面的用例看不出区别，只能数比较次数。
	// 自底向上建堆的比较次数不超过 2n（algs4 命题 R）；升序输入逐个 Insert 时
	// 每个元素都要 swim 到根，n = 1023 时约 8000 次。
	n := 1023
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	count := 0
	NewMaxPQFrom(func(x, y int) int {
		count++
		return cmp.Compare(x, y)
	}, a)
	if count > 2*n {
		t.Fatalf("NewMaxPQFrom used %d comparisons for n=%d, want <= %d", count, n, 2*n)
	}
}
