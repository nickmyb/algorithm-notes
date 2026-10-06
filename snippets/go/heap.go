package snippets

type MaxPQ struct {
	pq []int // 1, ..., n
	n  int
}

func NewMaxPQ() *MaxPQ {
	return &MaxPQ{
		pq: []int{0},
		n:  0,
	}
}

func (pq *MaxPQ) Insert(x int) {
	pq.n += 1
	pq.pq = append(pq.pq, x)
	pq.swim(pq.n)
}
func (pq *MaxPQ) Max() (int, bool) {
	if pq.IsEmpty() {
		return 0, false
	}
	return pq.pq[1], true
}

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

func (pq *MaxPQ) IsEmpty() bool {
	return pq.Size() == 0
}

func (pq *MaxPQ) Size() int {
	return pq.n
}

func (pq *MaxPQ) exch(i, j int) {
	pq.pq[i], pq.pq[j] = pq.pq[j], pq.pq[i]
}

func (pq *MaxPQ) swim(k int) {
	for k > 1 && pq.pq[k/2] < pq.pq[k] {
		pq.exch(k/2, k)
		k = k / 2
	}
}

func (pq *MaxPQ) sink(k int) {
	for 2*k <= pq.n {
		j := 2 * k
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
