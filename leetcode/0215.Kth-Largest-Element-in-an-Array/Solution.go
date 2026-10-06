package leetcode

// ===== 本地接线区 =====
import (
	"cmp"
)

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

// region @snippets/go/heap.go

// MaxPQ 是 algs4 2.4 节用数组实现的最大堆，元素类型 T 由调用方指定。
// 「最大」按 cmp 的顺序算，必须用 NewMaxPQ 创建，零值 MaxPQ 的 cmp 为 nil，不能直接使用。
// 要最小堆就把 cmp 反过来传，例如 func(a, b int) int { return cmp.Compare(b, a) }。
//
// 堆是一棵完全二叉树：节点从上到下、从左到右排满，中间没有空位，
// 所以第 k 个节点的位置完全由 k 决定，不需要指针。pq[0] 空着不用，
// 第 k 个节点存在 pq[k]，它的孩子是 2k、2k+1，父节点是 k/2：
//
//	    1
//	 2     3
//	4 5   6 7
//
// 堆序：每个节点都不小于它的孩子，所以最大值总在 pq[1]。
//
// 之前的版本，在仓库根目录执行：
//
//	git show eddb32d:snippets/go/heap.go    # 用指针链接节点的最小堆 MinPQ
//	git show c0668d1:snippets/go/heap.go    # 数组实现、元素类型固定为 int 的 MaxPQ
type MaxPQ[T any] struct {
	cmp func(T, T) int
	pq  []T // pq[1..n] 存元素，pq[0] 不用，len(pq) 总是 n+1
	n   int // 元素个数
}

// NewMaxPQ 返回一个按 cmp 比较元素的空堆，切片里只有占位的 pq[0]，cmp 不能为 nil。
// cmp(a, b) 小于 0 表示 a < b，等于 0 表示相等，大于 0 表示 a > b，
// 只看符号，不要求恰好是 -1、0、1。有序类型可以直接传 cmp.Compare[T]。
func NewMaxPQ[T any](cmp func(T, T) int) *MaxPQ[T] {
	return &MaxPQ[T]{
		cmp: cmp,
		pq:  make([]T, 1),
		n:   0,
	}
}

// NewMaxPQCap 返回一个空堆，预留 maxN 个元素的容量（pq[0] 另占一格，所以是 maxN+1）。
// maxN 只是容量提示，不是上限：超过之后 append 照常扩容。
func NewMaxPQCap[T any](cmp func(T, T) int, maxN int) *MaxPQ[T] {
	return &MaxPQ[T]{
		cmp: cmp,
		pq:  make([]T, 1, maxN+1),
		n:   0,
	}
}

// NewMaxPQFrom 用 a 中的元素建堆，不修改 a。O(n)。
// 对照 algs4 MaxPQ.java 的 MaxPQ(Key[] keys)，书中 2.4 节堆排序的第一阶段「堆的构造」：
// 先把 a 原样复制到 pq[1..n]，再从最后一个有孩子的节点 n/2 往前逐个 sink。
// 编号大于 n/2 的节点都是叶子，各自已经是合法的堆；sink(k) 时 k 的两棵子树已经是堆，
// 下沉完以 k 为根的子树也成了堆，到 k = 1 时整棵树就是堆。
// 逐个 Insert 结果一样，但要 O(n log n)；这里大部分节点在底层、下沉距离很短，
// 比较次数不超过 2n（algs4 命题 R）。
//
// 必须复制而不是直接拿 a 当底层数组：sink 会打乱调用方的切片，之后 append 还可能和它共用内存。
func NewMaxPQFrom[T any](cmp func(T, T) int, a []T) *MaxPQ[T] {
	pq := &MaxPQ[T]{
		cmp: cmp,
		pq:  make([]T, len(a)+1),
		n:   len(a),
	}
	copy(pq.pq[1:], a)
	// 2k是k的左节点, 2k <= pq.n
	for k := pq.n / 2; k >= 1; k-- {
		pq.sink(k)
	}
	return pq
}

// Insert 把 x 放到最后一个位置，再上浮到满足堆序的位置。O(log n)。
func (pq *MaxPQ[T]) Insert(x T) {
	pq.n += 1
	pq.pq = append(pq.pq, x)
	pq.swim(pq.n)
}

// Max 返回最大值但不删除，空堆时 ok 为 false。O(1)。
func (pq *MaxPQ[T]) Max() (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}
	return pq.pq[1], true
}

// DelMax 删除并返回最大值，空堆时 ok 为 false。O(log n)。
// 把最后一个元素搬到根上、截掉最后一格，再让它下沉。
// 截掉这一格不只是释放空间：下一次 Insert 用 append，不截的话新元素会排到旧值后面。
// 截掉之前先把这一格清零：截掉后它还留在底层数组里，T 含指针时会让 GC 回收不了
// 已经出队的元素，algs4 里 pq[n+1] = null 也是为此。
func (pq *MaxPQ[T]) DelMax() (T, bool) {
	var zero T
	if !pq.IsEmpty() {
		m := pq.pq[1]
		pq.pq[1] = pq.pq[pq.n]
		pq.pq[pq.n] = zero
		pq.pq = pq.pq[:pq.n]
		pq.n -= 1
		pq.sink(1)

		return m, true
	}
	return zero, false
}

// IsEmpty 报告堆是否为空。
func (pq *MaxPQ[T]) IsEmpty() bool {
	return pq.Size() == 0
}

// Size 返回元素个数。
func (pq *MaxPQ[T]) Size() int {
	return pq.n
}

// less 报告 pq[i] 是否小于 pq[j]。
func (pq *MaxPQ[T]) less(i, j int) bool {
	return pq.cmp(pq.pq[i], pq.pq[j]) < 0
}

// exch 交换 pq[i] 和 pq[j]。
func (pq *MaxPQ[T]) exch(i, j int) {
	pq.pq[i], pq.pq[j] = pq.pq[j], pq.pq[i]
}

// swim 让第 k 个节点上浮：比父节点大就和父节点交换，直到到达根或不再比父节点大。
func (pq *MaxPQ[T]) swim(k int) {
	for k > 1 && pq.less(k/2, k) {
		pq.exch(k/2, k)
		k = k / 2
	}
}

// sink 让第 k 个节点下沉：和较大的孩子交换，直到没有孩子或不小于两个孩子。
// 必须和较大的孩子换，换上去的才能同时不小于另一个孩子。
func (pq *MaxPQ[T]) sink(k int) {
	// 只要还存在左子节点就继续循环
	for 2*k <= pq.n {
		j := 2 * k
		// 存在右子节点才需要选更大的节点
		if j+1 <= pq.n && pq.less(j, j+1) {
			j += 1
		}
		if !pq.less(k, j) {
			return
		}
		pq.exch(k, j)

		k = j
	}
}

// endregion @snippets/go/heap.go

func findKthLargest(nums []int, k int) int {
	maxPQ := NewMaxPQFrom(cmp.Compare[int], nums)

	maxK := 0
	for _ = range k {
		maxK, _ = maxPQ.DelMax()
	}
	return maxK
}
