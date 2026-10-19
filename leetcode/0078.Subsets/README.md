# [78. Subsets](https://leetcode.cn/problems/subsets/)

> 难度：Medium

## 题目

Given an integer array `nums` of **unique** elements, return *all possible* *subsets* *(the power set)*.

The solution set **must not** contain duplicate subsets. Return the solution in **any order**.

**Example 1:**

```
Input: nums = [1,2,3]
Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
```

**Example 2:**

```
Input: nums = [0]
Output: [[],[0]]
```

**Constraints:**

- `1 <= nums.length <= 10`
- `-10 <= nums[i] <= 10`
- All the numbers of `nums` are **unique**.

<details>
<summary>官方中文题目翻译</summary>

给你一个整数数组 `nums` ，数组中的元素 **互不相同** 。返回该数组所有可能的子集（幂集）。

解集 **不能** 包含重复的子集。你可以按 **任意顺序** 返回解集。

**示例 1：**

```
输入：nums = [1,2,3]
输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
```

**示例 2：**

```
输入：nums = [0]
输出：[[],[0]]
```

**提示：**

- `1 <= nums.length <= 10`
- `-10 <= nums[i] <= 10`
- `nums` 中的所有元素 **互不相同**

</details>

## 题目大意

## 解题思路

1. 递归(任一元素都有在subset和不在subset两种状态)

### 同类题

## 复杂度

设 n = len(nums)。每个元素都有“选”或“不选”两种状态，所以一共有 2^n 个子集。

时间 O(n·2^n)：代码中的 append([]int{nums[0]}, subset...) 会复制 subset 的元素。虽然子集有 2^n 个，但生成每个子集的成本并非都为 O(1)。所有子集的元素个数之和是 n·2^(n-1)（每个元素恰好出现在一半的子集中），因此总时间是 O(n·2^n)。

空间 O(n·2^n)：返回值要保存全部子集；其中有 2^n 个切片，元素总数也是 n·2^(n-1)。递归调用栈只有 O(n)，但它不是主要的空间开销，所以 README 中的 O(n) 不能作为这份实现的总空间复杂度。

- 时间复杂度：O(n*2^n)
- 空间复杂度：O(n*2^n)
