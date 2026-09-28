# [217. Contains Duplicate](https://leetcode.cn/problems/contains-duplicate/)

> 难度：Easy

## 题目

Given an integer array `nums`, return `true` if any value appears **at least twice** in the array, and return `false` if every element is distinct.

**Example 1:**

**Input:** nums = [1,2,3,1]

**Output:** true

**Explanation:**

The element 1 occurs at the indices 0 and 3.

**Example 2:**

**Input:** nums = [1,2,3,4]

**Output:** false

**Explanation:**

All elements are distinct.

**Example 3:**

**Input:** nums = [1,1,1,3,3,4,3,2,4,2]

**Output:** true

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

<details>
<summary>官方中文题目翻译</summary>

给你一个整数数组 `nums` 。如果任一值在数组中出现 **至少两次** ，返回 `true` ；如果数组中每个元素互不相同，返回 `false` 。

**示例 1：**

**输入：** nums = [1,2,3,1]

**输出：** true

**解释：**

元素 1 在下标 0 和 3 出现。

**示例 2：**

**输入：** nums = [1,2,3,4]

**输出：** false

**解释：**

所有元素都不同。

**示例 3：**

**输入：** nums = [1,1,1,3,3,4,3,2,4,2]

**输出：** true

**提示：**

- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

</details>

## 题目大意

数组中是否存在相同数字

## 解题思路

1. 实现一个简单的set
2. 排序后比较是否存在相同的相邻数字

## 复杂度

- 时间复杂度： O(n)
- 空间复杂度： O(n)
