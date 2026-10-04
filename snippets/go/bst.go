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

// floor 返回以 x 为根的子树中不大于 key 的最大键所在的结点，没有返回 nil。
// key 小于 x.key 时，x 和它的右子树都太大，答案只能在左子树；
// key 大于 x.key 时，x 是一个候选，但右子树里可能有更接近 key 的，
// 先去右子树找，找不到才是 x。
func (t *BST[K, V]) floor(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}

	switch c := t.cmp(key, x.key); {
	case c < 0:
		return t.floor(x.left, key)
	case c > 0:
		fr := t.floor(x.right, key)
		if fr != nil {
			return fr
		}
		return x
	default:
		return x
	}
}

// ceiling 返回以 x 为根的子树中不小于 key 的最小键所在的结点，没有返回 nil，和 floor 左右对称。
func (t *BST[K, V]) ceiling(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}

	switch c := t.cmp(key, x.key); {
	case c < 0:
		cl := t.ceiling(x.left, key)
		if cl != nil {
			return cl
		}
		return x
	case c > 0:
		return t.ceiling(x.right, key)
	default:
		return x
	}
}

// rank 返回以 x 为根的子树中小于 key 的键的个数，key 不必在树中。
// key 小于 x.key 时，x 和右子树都不算，只数左子树；
// key 大于 x.key 时，左子树和 x 都小于 key，再加上右子树里比 key 小的；
// 相等时正好是左子树的结点数。
func (t *BST[K, V]) rank(x *node[K, V], key K) int {
	if x == nil {
		return 0
	}

	switch c := t.cmp(key, x.key); {
	case c < 0:
		return t.rank(x.left, key)
	case c > 0:
		return x.left.size() + 1 + t.rank(x.right, key)
	default:
		return x.left.size()
	}
}

// node 是树的结点。nil 表示空子树：get、put、floor、ceiling、rank 遇到 nil 的 x 直接处理，
// size、min、max、nSelect 允许 nil 接收者，所以递归时可以直接传 x.left、调用 x.left.size()，不必先判空。
// get、put、floor、ceiling、rank 要用 cmp，写在 BST 上，x 作参数传入；
// size、min、max、nSelect 不用 cmp，写在 node 上。
type node[K, V any] struct {
	key         K
	value       V
	left, right *node[K, V]
	n           int // 以该结点为根的子树结点数
}

// size 返回以 x 为根的子树结点数，空子树为 0。
func (x *node[K, V]) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

// min 返回以 x 为根的子树中键最小的结点，空子树返回 nil。
// 比 x 小的键都在左子树，所以一路向左，走到没有左孩子的结点就是最小，
// 它可以有右孩子，右子树里的键都比它大。
func (x *node[K, V]) min() *node[K, V] {
	if x == nil {
		return nil
	}

	if x.left == nil {
		return x
	}
	return x.left.min()
}

// max 返回以 x 为根的子树中键最大的结点，空子树返回 nil，和 min 左右对称。
func (x *node[K, V]) max() *node[K, V] {
	if x == nil {
		return nil
	}

	if x.right == nil {
		return x
	}
	return x.right.max()
}

// nSelect 返回以 x 为根的子树中排名为 k 的结点（从 0 开始数，即恰好有 k 个键比它小），
// k 越界（小于 0 或不小于子树结点数）时返回 nil。select 是 Go 的关键字，所以叫 nSelect。
// 设左子树有 c 个结点：k < c 时答案在左子树，排名不变；k > c 时答案在右子树，
// 左子树和 x 共 c+1 个键排在前面，在右子树里的排名是 k-c-1；k == c 时就是 x。
func (x *node[K, V]) nSelect(k int) *node[K, V] {
	if x == nil {
		return nil
	}

	switch c := x.left.size(); {
	case c < k:
		return x.right.nSelect(k - c - 1)
	case c > k:
		return x.left.nSelect(k)
	default:
		return x
	}
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

// Min 返回最小的键；ok 为 false 表示树为空，此时键为 K 的零值。
func (t *BST[K, V]) Min() (K, bool) {
	x := t.root.min()
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

// Max 返回最大的键；ok 为 false 表示树为空，此时键为 K 的零值。
func (t *BST[K, V]) Max() (K, bool) {
	x := t.root.max()
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

// Floor 返回不大于 key 的最大键，key 在树中时就是它自己；
// ok 为 false 表示所有键都比 key 大（包括树为空），此时键为 K 的零值。
func (t *BST[K, V]) Floor(key K) (K, bool) {
	x := t.floor(t.root, key)
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

// Ceiling 返回不小于 key 的最小键，key 在树中时就是它自己；
// ok 为 false 表示所有键都比 key 小（包括树为空），此时键为 K 的零值。
func (t *BST[K, V]) Ceiling(key K) (K, bool) {
	x := t.ceiling(t.root, key)
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

// Select 返回排名为 k 的键，即从小到大第 k 个（从 0 开始数），Select(0) 就是 Min；
// ok 为 false 表示 k 越界（小于 0 或不小于 Size），此时键为 K 的零值。
func (t *BST[K, V]) Select(k int) (K, bool) {
	x := t.root.nSelect(k)
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

// Rank 返回树中小于 key 的键的个数，key 不必在树中。
// key 在树中时，Select(Rank(key)) 就是 key 自己。
func (t *BST[K, V]) Rank(key K) int {
	return t.rank(t.root, key)
}
