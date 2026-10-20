# [46. Permutations](https://leetcode.cn/problems/permutations/)

> 难度：Medium

## 题目

Given an array `nums` of distinct integers, return all the possible permutations. You can return the answer in any order.

**Example 1:**

```
Input: nums = [1,2,3]
Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

**Example 2:**

```
Input: nums = [0,1]
Output: [[0,1],[1,0]]
```

**Example 3:**

```
Input: nums = [1]
Output: [[1]]
```

**Constraints:**

- `1 <= nums.length <= 6`
- `-10 <= nums[i] <= 10`
- All the integers of `nums` are unique.

<details>
<summary>官方中文题目翻译</summary>

给定一个不含重复数字的数组 `nums` ，返回其 **所有可能的全排列** 。你可以 **按任意顺序** 返回答案。

**示例 1：**

```
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

**示例 2：**

```
输入：nums = [0,1]
输出：[[0,1],[1,0]]
```

**示例 3：**

```
输入：nums = [1]
输出：[[1]]
```

**提示：**

- `1 <= nums.length <= 6`
- `-10 <= nums[i] <= 10`
- `nums` 中的所有整数 **互不相同**

</details>

## 题目大意

## 解题思路

1. 递归,新元素可以插入任意n个元素组成的一个排列的n+1个位置

### 同类题

- [78. Subsets](../0078.Subsets/)

## 复杂度

长度为 n 的数组有 n! 个排列；代码为每个排列复制 n 个元素，返回值也要存下这 n! 个长度为 n 的切片

- 时间复杂度：O(n*n!)
- 空间复杂度：O(n*n!)
