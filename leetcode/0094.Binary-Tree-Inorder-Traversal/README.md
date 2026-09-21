# [94. Binary Tree Inorder Traversal](https://leetcode.cn/problems/binary-tree-inorder-traversal/)

> 难度：Easy

## 题目

Given the `root` of a binary tree, return *the inorder traversal of its nodes' values*.

**Example 1:**

**Input:** root = [1,null,2,3]

**Output:** [1,3,2]

**Explanation:**

![](https://assets.leetcode.com/uploads/2024/08/29/screenshot-2024-08-29-202743.png)

**Example 2:**

**Input:** root = [1,2,3,4,5,null,8,null,null,6,7,9]

**Output:** [4,2,6,5,7,1,3,9,8]

**Explanation:**

![](https://assets.leetcode.com/uploads/2024/08/29/tree_2.png)

**Example 3:**

**Input:** root = []

**Output:** []

**Example 4:**

**Input:** root = [1]

**Output:** [1]

**Constraints:**

- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

**Follow up:** Recursive solution is trivial, could you do it iteratively?

<details>
<summary>官方中文题目翻译</summary>

给定一个二叉树的根节点 `root` ，返回 *它的 **中序** 遍历* 。

**示例 1：**

![](https://assets.leetcode.com/uploads/2020/09/15/inorder_1.jpg)

```
输入：root = [1,null,2,3]
输出：[1,3,2]
```

**示例 2：**

```
输入：root = []
输出：[]
```

**示例 3：**

```
输入：root = [1]
输出：[1]
```

**提示：**

- 树中节点数目在范围 `[0, 100]` 内
- `-100 <= Node.val <= 100`

**进阶:** 递归算法很简单，你可以通过迭代算法完成吗？

</details>

## 题目大意

二叉树 左-中-右 的遍历方式。

## 解题思路

简单递归即可。

## 复杂度

- 时间复杂度：O(n)，n 为节点数。每个节点恰好访问一次，每次的追加操作在三门语言里都是摊还 O(1)（底层几何增长，n 次追加的总拷贝量是 O(n)，不是 O(n²)）
- 空间复杂度：O(n)。递归栈深度等于树高 h，题目不保证平衡，退化成链时 h = n；平衡树则为 O(log n)。输出列表本身也是 O(n)，但那是题目要求的返回值，按惯例不计入额外空间
