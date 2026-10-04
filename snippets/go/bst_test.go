package snippets

import (
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
