package leetcode

// ===== 本地接线区 =====
//
// 用到 TreeNode / ListNode 时，在这里 import 共享结构并起一个类型别名，
// 这样题解本体里就能写裸的 *TreeNode，和 LeetCode 给的签名一字不差：
//
//	import structures "github.com/nickmyb/algorithm-notes/structures/go"
//
//	type TreeNode = structures.TreeNode
//	type ListNode = structures.ListNode
//
// 注意是类型别名（=），不是定义新类型，两者在 Go 里是同一个类型。
//
// package 声明本身也属于接线，LeetCode 的提交框里没有这一行。

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

// 题解骨架。把函数名和签名改成本题要求的，例如 Two Sum 是：
//
//	func twoSum(nums []int, target int) []int
//
// 同一题的多种解法写在本文件里的多个函数，入口函数保持 LeetCode 给的签名名，
// 其余解法加后缀说明算法：twoSum / twoSumBruteForce / twoSumTwoPointers。
// 需要独立辅助数据结构时才另开文件，文件名按结构命名（SegmentTree.go），
// 不要叫 Solution2.go —— 那会和"第二种解法"混淆。
//
// Go 的包名来自 package 声明，跟文件名无关，所以每道题的文件都可以
// 叫 Solution.go 而不冲突。
func solve(nums []int) []int {
	return nil
}
