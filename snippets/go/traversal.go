package snippets

// TreeNode 照抄 LeetCode 给的定义。题解里的 TreeNode 来自 structures，
// 这里只能 import 标准库，所以另定义一份，inorder 才能和题解副本逐字一致。
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// inorder 按中序（左、根、右）把以 root 为根的子树的值追加到 *ret。
// BST 的中序序列就是键从小到大，所以验证 BST 可以检查它是否严格递增。
// 传 *[]int 是因为 append 可能换底层数组，各层递归要共享同一个切片头；
// 调用方传 &ret，ret 原有的元素保留在前面。
func inorder(root *TreeNode, ret *[]int) {
	if root == nil {
		return
	}

	inorder(root.Left, ret)
	*ret = append(*ret, root.Val)
	inorder(root.Right, ret)
}
