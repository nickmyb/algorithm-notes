package snippets

import (
	"strconv"
)

type hNode struct {
	parent, left, right *hNode
	dir                 int

	key int
}

type MinPQ struct {
	root *hNode
	size int
}

func NewMinPQ() *MinPQ {
	return &MinPQ{
		root: nil,
		size: 0,
	}
}

func (pq *MinPQ) position(i int) (p *hNode, d int) {
	// 第 i 个节点（从 1 开始编号）的父节点 p 和方向 d: -1 left; 1 right
	// i 的二进制去掉开头的 1 就是从根往下的路径，最后一位是方向
	//          1
	//       10   11
	//     100 101 110 111
	bstring := strconv.FormatInt(int64(i), 2)
	p = pq.root
	for _, ch := range bstring[1 : len(bstring)-1] {
		if ch == '0' {
			p = p.left
		}
		if ch == '1' {
			p = p.right
		}
	}
	if bstring[len(bstring)-1] == '0' {
		d = -1
	} else {
		d = 1
	}

	return p, d
}

// 交换 key 而不是交换节点：树的形状不变，parent / left / right / dir 都不用动
func (pq *MinPQ) swapUp(c *hNode) {
	for c.parent != nil && c.parent.key > c.key {
		c.parent.key, c.key = c.key, c.parent.key
		c = c.parent
	}
}

func (pq *MinPQ) swapDown(c *hNode) {
	for {
		// 和较小的孩子交换；两个孩子相等时取左边
		m := c
		if c.left != nil && c.left.key < m.key {
			m = c.left
		}
		if c.right != nil && c.right.key < m.key {
			m = c.right
		}
		if m == c {
			return
		}
		m.key, c.key = c.key, m.key
		c = m
	}
}

func (pq *MinPQ) Add(x int) {
	n := &hNode{key: x}
	pq.size += 1
	if pq.size == 1 {
		pq.root = n
		return
	}

	parent, d := pq.position(pq.size)
	n.parent, n.dir = parent, d
	if d == -1 {
		parent.left = n
	} else {
		parent.right = n
	}
	pq.swapUp(n)
}

func (pq *MinPQ) GetSmallest() (int, bool) {
	if pq.root == nil {
		return 0, false
	}
	return pq.root.key, true
}

func (pq *MinPQ) RemoveSmallest() (int, bool) {
	smallest, ok := pq.GetSmallest()
	if !ok {
		return 0, false
	}
	if pq.size == 1 {
		pq.root = nil
		pq.size = 0
		return smallest, true
	}

	// 摘下最下层最右侧的节点，把它的 key 放到根上，再往下沉
	parent, d := pq.position(pq.size)
	var last *hNode
	if d == -1 {
		last, parent.left = parent.left, nil
	} else {
		last, parent.right = parent.right, nil
	}
	pq.size -= 1

	pq.root.key = last.key
	pq.swapDown(pq.root)
	return smallest, true
}

func (pq *MinPQ) Size() int {
	return pq.size
}
