package snippets

// BST 是 algs4 3.2 节的二叉查找树符号表，键类型 K、值类型 V 由调用方指定。
// 键的顺序由 cmp 决定，必须用 NewBST 创建，零值 BST 的 cmp 为 nil，不能直接使用。
//
// 泛型写法参照 Go 官方博客 When To Use Generics 的 General purpose data structures 一节：
// https://go.dev/blog/when-generics
//
// 改成泛型之前是 string 键、int 值的版本，在仓库根目录执行：
//
//	git show 3caa89e:snippets/go/bst.go                  # 查看旧版本
//	git diff 3caa89e 250ca21 -- snippets/go/bst.go       # 查看旧版本改成泛型的差异
type BST[K, V any] struct {
	cmp  func(K, K) int
	root *node[K, V]
}

// NewBST 返回一棵按 cmp 比较键的空树，cmp 不能为 nil。
// cmp(a, b) 小于 0 表示 a < b，等于 0 表示相等，大于 0 表示 a > b，
// 只看符号，不要求恰好是 -1、0、1。有序类型可以直接传 cmp.Compare[K]，
// string 也可以传 strings.Compare。
func NewBST[K, V any](cmp func(K, K) int) *BST[K, V] {
	return &BST[K, V]{cmp: cmp}
}

// node 是树的结点。nil 表示空子树：get、put 遇到 nil 的 x 直接处理，
// size 允许 nil 接收者，所以递归时可以直接传 x.left、调用 x.left.size()，不必先判空。
// get、put 要用 cmp，写在 BST 上，x 作参数传入。
type node[K, V any] struct {
	key         K
	value       V
	left, right *node[K, V]
	n           int // 以该结点为根的子树结点数
}

// get 在以 x 为根的子树中查找 key。
func (t *BST[K, V]) get(x *node[K, V], key K) (V, bool) {
	if x == nil {
		var zero V
		return zero, false
	}

	switch c := t.cmp(key, x.key); {
	case c < 0:
		return t.get(x.left, key)
	case c > 0:
		return t.get(x.right, key)
	default:
		return x.value, true
	}
}

// put 在以 x 为根的子树中插入或覆盖 key，返回插入后的子树根，
// 调用方用它更新指向这棵子树的链接。
func (t *BST[K, V]) put(x *node[K, V], key K, val V) *node[K, V] {
	if x == nil {
		return &node[K, V]{key: key, value: val, n: 1}
	}

	switch c := t.cmp(key, x.key); {
	case c < 0:
		x.left = t.put(x.left, key, val)
	case c > 0:
		x.right = t.put(x.right, key, val)
	default:
		x.value = val
	}

	x.n = 1 + x.left.size() + x.right.size()

	return x
}

// size 返回以 x 为根的子树结点数，空子树为 0。
func (x *node[K, V]) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

// Get 返回 key 对应的值；ok 为 false 表示 key 不存在，此时值为 V 的零值。
func (t *BST[K, V]) Get(key K) (V, bool) {
	return t.get(t.root, key)
}

// Put 插入键值对，key 已存在时覆盖旧值。
func (t *BST[K, V]) Put(key K, val V) {
	t.root = t.put(t.root, key, val)
}

// Size 返回树中键值对的数量。
func (t *BST[K, V]) Size() int {
	return t.root.size()
}
