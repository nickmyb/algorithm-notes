package snippets

import (
	"slices"
	"strings"
	"testing"
)

// 每个用例都注明它能抓到的错误写法，说不出来的不加。

func TestBSTGet(t *testing.T) {
	// 值取插入序号。按这个顺序插入后树的形状：
	//
	//	        S
	//	      /   \
	//	     E     X
	//	    / \
	//	   A   R
	//	    \  /
	//	    C H
	bst := NewBST[string, int](strings.Compare)
	for i, k := range []string{"S", "E", "X", "A", "R", "C", "H"} {
		bst.Put(k, i)
	}
	tests := []struct {
		name   string
		key    string
		want   int
		wantOK bool
	}{
		// 根：抓相等分支写错（比如漏了 default，落到某一侧继续找）
		{"根", "S", 0, true},
		// 左子树深处：抓 put/get 递归时传了 x 而不是 x.left，会无限递归
		{"左子树深处", "C", 5, true},
		// 右子树：抓 get 里左右分支写反，X > S 却去左边找
		{"右子树", "X", 2, true},
		// 不存在：抓用 -1 之类的哨兵值代替 ok，0 和 false 必须一起返回
		{"不存在", "B", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := bst.Get(tt.key)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("Get(%q) = %d, %v, want %d, %v", tt.key, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestBSTPutOverwrite(t *testing.T) {
	// 重复 key：抓相等时没有覆盖而是又插了一个结点（Get 取到旧值、Size 多算）
	bst := NewBST[string, int](strings.Compare)
	bst.Put("S", 1)
	bst.Put("E", 2)
	bst.Put("E", 3)
	if got, _ := bst.Get("E"); got != 3 {
		t.Errorf("Get(E) = %d, want 3", got)
	}
	if got := bst.Size(); got != 2 {
		t.Errorf("Size() = %d, want 2", got)
	}
}

func TestBSTCmpSign(t *testing.T) {
	// cmp 返回 -4 这类非 ±1 的值：抓 switch 写成 case -1 / case 1，
	// 其余值落进 default 被当成相等，覆盖掉别的 key。strings.Compare 只返回 -1、0、1，
	// 上面的用例都抓不到这种写法。
	bst := NewBST[int, string](func(a, b int) int { return a - b })
	bst.Put(5, "five")
	bst.Put(1, "one")  // cmp(1, 5) = -4
	bst.Put(9, "nine") // cmp(9, 5) = 4
	for _, tt := range []struct {
		key  int
		want string
	}{{5, "five"}, {1, "one"}, {9, "nine"}} {
		if got, ok := bst.Get(tt.key); got != tt.want || !ok {
			t.Errorf("Get(%d) = %q, %v, want %q, true", tt.key, got, ok, tt.want)
		}
	}
	if got := bst.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
}

// bstKeys 用 Select 按排名取出全部键。Select 靠每个结点的 n 定位，
// 所以 n 没维护好、或者树里丢了或多出结点，取出来的序列都会不对。
func bstKeys(bst *BST[string, int]) []string {
	keys := make([]string, 0, bst.Size())
	for k := range bst.Size() {
		key, _ := bst.Select(k)
		keys = append(keys, key)
	}
	return keys
}

// newDeleteTestBST 在 TestBSTMinMax 那棵树上多插一个 Q，
// 让根 M 的后继 Q 不是 M 的直接右孩子，而在 T 的左子树里：
//
//	     M
//	   /   \
//	  E     T
//	 /     / \
//	A     Q   Z
//	 \       /
//	  C     W
func newDeleteTestBST() *BST[string, int] {
	bst := NewBST[string, int](strings.Compare)
	for _, k := range []string{"M", "E", "T", "A", "Z", "C", "W", "Q"} {
		bst.Put(k, 0)
	}
	return bst
}

func TestBSTDelete(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want []string
	}{
		// 叶子：抓删除后没有更新路径上的 n，Size 仍是 8，Select 取出一个空键
		{"叶子", "C", []string{"A", "E", "M", "Q", "T", "W", "Z"}},
		// 只有右孩子：抓两个单孩子分支写反（返回 x.left），C 跟着丢掉
		{"只有右孩子", "A", []string{"C", "E", "M", "Q", "T", "W", "Z"}},
		// 只有左孩子：同上，左右对称，W 跟着丢掉
		{"只有左孩子", "Z", []string{"A", "C", "E", "M", "Q", "T", "W"}},
		// 两个孩子，后继在更深处：抓先接左子树再 deleteMin
		// （后继有了左孩子，deleteMin 会越过它去删 A）、
		// 没把后继从右子树摘掉（Q 出现两次）、deleteMin 没更新 n
		{"两个孩子", "M", []string{"A", "C", "E", "Q", "T", "W", "Z"}},
		// 不存在：抓没有判 nil，走到空子树后 x.key 直接 panic
		{"不存在", "D", []string{"A", "C", "E", "M", "Q", "T", "W", "Z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := newDeleteTestBST()
			bst.Delete(tt.key)
			if got := bstKeys(bst); !slices.Equal(got, tt.want) {
				t.Errorf("Delete(%q) 后 = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestBSTSize(t *testing.T) {
	bst := NewBST[string, int](strings.Compare)
	// 空树：抓 size 没处理 nil 接收者，root 为 nil 时直接 panic
	if got := bst.Size(); got != 0 {
		t.Errorf("空树 Size() = %d, want 0", got)
	}
	// 根的两侧都有子树：抓 n 只累加了一侧（algs4 的 size(left) + size(right) + 1 漏项）
	for _, k := range []string{"B", "A", "C"} {
		bst.Put(k, 0)
	}
	if got := bst.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
}

func TestBSTMinMax(t *testing.T) {
	bst := NewBST[string, int](strings.Compare)
	// 空树：抓 min/max 没处理 nil 接收者，root 为 nil 时直接 panic；
	// 也抓 Min/Max 漏了 ok，空树必须返回 false
	if got, ok := bst.Min(); got != "" || ok {
		t.Errorf("空树 Min() = %q, %v, want \"\", false", got, ok)
	}
	if got, ok := bst.Max(); got != "" || ok {
		t.Errorf("空树 Max() = %q, %v, want \"\", false", got, ok)
	}

	// 按这个顺序插入后树的形状，最小、最大都在第三层，而且各有一个朝内的孩子：
	//
	//	        M
	//	      /   \
	//	     E     T
	//	    /       \
	//	   A         Z
	//	    \       /
	//	     C     W
	for _, k := range []string{"M", "E", "T", "A", "Z", "C", "W"} {
		bst.Put(k, 0)
	}
	tests := []struct {
		name string
		f    func() (string, bool)
		want string
	}{
		// 抓递归方向写错（判断 x.left 却递归 x.right，M 转到 T 后 T 没有左孩子，得到 T）、
		// 少递归一层（直接返回 x.left 得到 E），以及把「没有左孩子」写成「是叶子」
		// （A 有右孩子 C，那种写法会越过 A 走到 nil）
		{"Min", bst.Min, "A"},
		// 同上，左右对称：方向写反得到 E，少一层得到 T，判叶子会越过 Z（它有左孩子 W）走到 nil
		{"Max", bst.Max, "Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := tt.f(); got != tt.want || !ok {
				t.Errorf("%s() = %q, %v, want %q, true", tt.name, got, ok, tt.want)
			}
		})
	}
}

func TestBSTDeleteMin(t *testing.T) {
	// 空树：抓 deleteMin 开头没有判 nil，空树上直接 panic
	empty := NewBST[string, int](strings.Compare)
	empty.DeleteMin()
	if got := empty.Size(); got != 0 {
		t.Errorf("空树 DeleteMin 后 Size() = %d, want 0", got)
	}

	// 最小值 A 有右孩子 C：抓返回 nil 而不是 x.right（C 跟着丢掉）、没更新 n
	bst := newDeleteTestBST()
	bst.DeleteMin()
	want := []string{"C", "E", "M", "Q", "T", "W", "Z"}
	if got := bstKeys(bst); !slices.Equal(got, want) {
		t.Errorf("DeleteMin 后 = %v, want %v", got, want)
	}
}

func TestBSTDeleteMax(t *testing.T) {
	// 空树：抓 deleteMax 开头没有判 nil，空树上直接 panic
	empty := NewBST[string, int](strings.Compare)
	empty.DeleteMax()
	if got := empty.Size(); got != 0 {
		t.Errorf("空树 DeleteMax 后 Size() = %d, want 0", got)
	}

	// 最大值 Z 有左孩子 W：抓返回 nil 而不是 x.left（W 跟着丢掉）、没更新 n
	bst := newDeleteTestBST()
	bst.DeleteMax()
	want := []string{"A", "C", "E", "M", "Q", "T", "W"}
	if got := bstKeys(bst); !slices.Equal(got, want) {
		t.Errorf("DeleteMax 后 = %v, want %v", got, want)
	}
}

func TestBSTFloorCeiling(t *testing.T) {
	// 和 TestBSTMinMax 同一棵树，排好序是 A C E M T W Z：
	//
	//	        M
	//	      /   \
	//	     E     T
	//	    /       \
	//	   A         Z
	//	    \       /
	//	     C     W
	bst := NewBST[string, int](strings.Compare)
	for _, k := range []string{"M", "E", "T", "A", "Z", "C", "W"} {
		bst.Put(k, 0)
	}
	tests := []struct {
		name   string
		f      func(string) (string, bool)
		key    string
		want   string
		wantOK bool
	}{
		// key 在树中：抓相等分支写错，比如把相等并进 c < 0 去左子树找，得到 C
		{"Floor 命中", bst.Floor, "E", "E", true},
		// 答案在右子树里：抓 key > x.key 时直接返回 x、不去右子树找，在根停下得到 M。
		// 特意选根的右半边：其余 Floor 用例都在 M 的左子树里，只靠它们的话，
		// Floor 从错误的子树出发（比如传了 t.root.left）也照样全过
		{"Floor 在右子树", bst.Floor, "X", "W", true},
		// 右子树非空但全都比 key 大，答案是 x 自己：抓漏了「右子树找不到就返回 x」，得到 nil
		{"Floor 回退到 x", bst.Floor, "B", "A", true},
		// 比所有键都小（ASCII 里 '0' < 'A'）：抓没有 floor 时仍返回某个结点
		{"Floor 不存在", bst.Floor, "0", "", false},

		// 以下和 Floor 左右对称
		{"Ceiling 命中", bst.Ceiling, "T", "T", true},
		// 在根停下会得到 M；同样特意选根的另一半
		{"Ceiling 在左子树", bst.Ceiling, "B", "C", true},
		// Z 的左子树 W 比 Y 小
		{"Ceiling 回退到 x", bst.Ceiling, "Y", "Z", true},
		// "ZZ" 比 "Z" 大
		{"Ceiling 不存在", bst.Ceiling, "ZZ", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := tt.f(tt.key); got != tt.want || ok != tt.wantOK {
				t.Errorf("(%q) = %q, %v, want %q, %v", tt.key, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestBSTSelectRank(t *testing.T) {
	// 和 TestBSTMinMax 同一棵树，排好序是 A C E M T W Z，排名依次是 0 到 6：
	//
	//	        M
	//	      /   \
	//	     E     T
	//	    /       \
	//	   A         Z
	//	    \       /
	//	     C     W
	bst := NewBST[string, int](strings.Compare)
	for _, k := range []string{"M", "E", "T", "A", "Z", "C", "W"} {
		bst.Put(k, 0)
	}

	selectTests := []struct {
		name   string
		k      int
		want   string
		wantOK bool
	}{
		// 在左子树里要先向左再向右（M → E → A → C）：抓两个分支写反
		{"左子树", 1, "C", true},
		// 连续向右两次再向左（M → T → Z → W）：抓去右子树时排名减 c 而不是 c+1，
		// 每向右一次多算一位，最后走进空子树返回 false
		{"右子树", 5, "W", true},
		// 越界：抓没有判 nil，走到空子树后 x.left.size() 直接 panic
		{"越界", 7, "", false},
	}
	for _, tt := range selectTests {
		t.Run("Select "+tt.name, func(t *testing.T) {
			if got, ok := bst.Select(tt.k); got != tt.want || ok != tt.wantOK {
				t.Errorf("Select(%d) = %q, %v, want %q, %v", tt.k, got, ok, tt.want, tt.wantOK)
			}
		})
	}

	rankTests := []struct {
		name string
		key  string
		want int
	}{
		// 根：抓相等时多加了 x 自己（返回 left.size()+1，得到 4）
		{"根", "M", 3},
		// 向右两次：抓向右时漏了 +1 或漏了左子树的结点数
		{"右子树", "W", 5},
		// 不在树中：抓没有判 nil 直接 panic，或者以为 key 一定在树里
		{"不在树中", "D", 2},
	}
	for _, tt := range rankTests {
		t.Run("Rank "+tt.name, func(t *testing.T) {
			if got := bst.Rank(tt.key); got != tt.want {
				t.Errorf("Rank(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}
