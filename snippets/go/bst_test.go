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
