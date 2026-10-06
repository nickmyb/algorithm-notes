package snippets

// MaxPQ 是 algs4 2.4 节用数组实现的最大堆，元素类型为 int。
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
// 改成数组之前是用指针链接节点的最小堆 MinPQ，在仓库根目录执行：
//
//	git show eddb32d:snippets/go/heap.go    # 查看指针版本
type MaxPQ struct {
	pq []int // pq[1..n] 存元素，pq[0] 不用，len(pq) 总是 n+1
	n  int   // 元素个数
}

// NewMaxPQ 返回一个空堆，切片里只有占位的 pq[0]。
func NewMaxPQ() *MaxPQ {
	return &MaxPQ{
		pq: []int{0},
		n:  0,
	}
}

// Insert 把 x 放到最后一个位置，再上浮到满足堆序的位置。O(log n)。
func (pq *MaxPQ) Insert(x int) {
	pq.n += 1
	pq.pq = append(pq.pq, x)
	pq.swim(pq.n)
}

// Max 返回最大值但不删除，空堆时 ok 为 false。O(1)。
func (pq *MaxPQ) Max() (int, bool) {
	if pq.IsEmpty() {
		return 0, false
	}
	return pq.pq[1], true
}

// DelMax 删除并返回最大值，空堆时 ok 为 false。O(log n)。
// 把最后一个元素搬到根上、截掉最后一格，再让它下沉。
// 截掉这一格不只是释放空间：下一次 Insert 用 append，不截的话新元素会排到旧值后面。
func (pq *MaxPQ) DelMax() (int, bool) {
	if !pq.IsEmpty() {
		m := pq.pq[1]
		pq.pq[1] = pq.pq[pq.n]
		pq.pq = pq.pq[:pq.n]
		pq.n -= 1
		pq.sink(1)

		return m, true
	}
	return 0, false
}

// IsEmpty 报告堆是否为空。
func (pq *MaxPQ) IsEmpty() bool {
	return pq.Size() == 0
}

// Size 返回元素个数。
func (pq *MaxPQ) Size() int {
	return pq.n
}

// exch 交换 pq[i] 和 pq[j]。
func (pq *MaxPQ) exch(i, j int) {
	pq.pq[i], pq.pq[j] = pq.pq[j], pq.pq[i]
}

// swim 让第 k 个节点上浮：比父节点大就和父节点交换，直到到达根或不再比父节点大。
func (pq *MaxPQ) swim(k int) {
	for k > 1 && pq.pq[k/2] < pq.pq[k] {
		pq.exch(k/2, k)
		k = k / 2
	}
}

// sink 让第 k 个节点下沉：和较大的孩子交换，直到没有孩子或不小于两个孩子。
// 必须和较大的孩子换，换上去的才能同时不小于另一个孩子。
func (pq *MaxPQ) sink(k int) {
	// 只要还存在左子节点就继续循环
	for 2*k <= pq.n {
		j := 2 * k
		// 存在右子节点才需要选更大的节点
		if j+1 <= pq.n && pq.pq[j+1] > pq.pq[j] {
			j += 1
		}
		if pq.pq[k] >= pq.pq[j] {
			return
		}
		pq.exch(k, j)

		k = j
	}
}
