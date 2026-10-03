package snippets

// BST 是 algs4 3.2 节的二叉查找树符号表，键为 string，值为 int。
// 零值即为空树，可以直接使用：var t BST。
type BST struct {
	root *node
}

// node 是树的结点。下面的方法都允许 nil 接收者，nil 表示空子树，
// 所以递归时可以直接调用 x.left.get(key)，不必先判空。
type node struct {
	key         string
	value       int
	left, right *node
	n           int // 以该结点为根的子树结点数
}

// get 在以 x 为根的子树中查找 key。
func (x *node) get(key string) (int, bool) {
	if x == nil {
		return 0, false
	}

	switch {
	case key < x.key:
		return x.left.get(key)
	case key > x.key:
		return x.right.get(key)
	default:
		return x.value, true
	}
}

// put 在以 x 为根的子树中插入或覆盖 key，返回插入后的子树根，
// 调用方用它更新指向这棵子树的链接。
func (x *node) put(key string, val int) *node {
	if x == nil {
		return &node{key: key, value: val, n: 1}
	}

	switch {
	case key < x.key:
		x.left = x.left.put(key, val)
	case key > x.key:
		x.right = x.right.put(key, val)
	default:
		x.value = val
	}

	x.n = 1 + x.left.size() + x.right.size()

	return x
}

// size 返回以 x 为根的子树结点数，空子树为 0。
func (x *node) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

// Get 返回 key 对应的值；ok 为 false 表示 key 不存在，此时值为 0。
func (t *BST) Get(key string) (int, bool) {
	return t.root.get(key)
}

// Put 插入键值对，key 已存在时覆盖旧值。
func (t *BST) Put(key string, val int) {
	t.root = t.root.put(key, val)
}

// Size 返回树中键值对的数量。
func (t *BST) Size() int {
	return t.root.size()
}
