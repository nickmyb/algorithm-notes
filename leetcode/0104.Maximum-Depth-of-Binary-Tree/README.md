# [104. Maximum Depth of Binary Tree](https://leetcode.cn/problems/maximum-depth-of-binary-tree/)

> 难度：Easy

## 题目

Given the `root` of a binary tree, return *its maximum depth*.

A binary tree's **maximum depth** is the number of nodes along the longest path from the root node down to the farthest leaf node.

**Example 1:**

![](https://assets.leetcode.com/uploads/2020/11/26/tmp-tree.jpg)

```
Input: root = [3,9,20,null,null,15,7]
Output: 3
```

**Example 2:**

```
Input: root = [1,null,2]
Output: 2
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-100 <= Node.val <= 100`

<details>
<summary>官方中文题目翻译</summary>

给定一个二叉树 `root` ，返回其最大深度。

二叉树的 **最大深度** 是指从根节点到最远叶子节点的最长路径上的节点数。

**示例 1：**

![](https://assets.leetcode.com/uploads/2020/11/26/tmp-tree.jpg)

```
输入：root = [3,9,20,null,null,15,7]
输出：3
```

**示例 2：**

```
输入：root = [1,null,2]
输出：2
```

**提示：**

- 树中节点的数量在 `[0, 10^4]` 区间内。
- `-100 <= Node.val <= 100`

</details>

## 题目大意

求二叉树高度。

## 解题思路

简单递归即可。

## 复杂度

- 时间复杂度：O(n)，每个节点恰好访问一次,每次只做一次比较和加法。
- 空间复杂度：O(h)，唯一的开销是递归栈,深度等于树高。题目不保证平衡,退化成链时 h = n,所以最坏是 O(n);平衡树则是 O(log n)。

```
calls = 0   # 累计调用次数——只增不减
cur   = 0   # 当前栈深——进 +1，出 -1
peak  = 0   # cur 的历史最大值

def maxDepth(root):
    global calls, cur, peak
    calls += 1; cur += 1; peak = max(peak, cur)   # 入口
    if root is None:
        cur -= 1; return 0                        # 出口 1
    l = maxDepth(root.left)
    r = maxDepth(root.right)
    cur -= 1                                      # 出口 2
    return max(l, r) + 1

calls 和 peak 的区别全在减不减上:calls 只增,量的是"一共调用过多少次"(时间);cur 有进有出,peak 量的是"同一瞬间最多几帧"(空间)。
```
