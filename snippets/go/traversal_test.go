package snippets

import (
	"slices"
	"testing"
)

func TestInorder(t *testing.T) {
	// 三层的树才能抓到只处理一层、不往下递归的写法（那样得到 [2 4 5]）：
	//
	//	    4
	//	   / \
	//	  2   5
	//	 / \
	//	1   3
	tree := &TreeNode{Val: 4,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 5},
	}
	tests := []struct {
		name string
		root *TreeNode
		want []int
	}{
		// 写成前序得到 [4 2 1 3 5]，后序 [1 3 2 5 4]，左右写反 [5 4 3 2 1]
		{"three levels", tree, []int{1, 2, 3, 4, 5}},
		// 漏了 nil 判断会在 root.Left 上空指针 panic
		{"empty tree", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			inorder(tt.root, &got)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("inorder = %v, want %v", got, tt.want)
			}
		})
	}
}
