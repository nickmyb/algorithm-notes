# [704. Binary Search](https://leetcode.cn/problems/binary-search/)

> 难度：Easy

## 题目

Given an array of integers `nums` which is sorted in ascending order, and an integer `target`, write a function to search `target` in `nums`. If `target` exists, then return its index. Otherwise, return `-1`.

You must write an algorithm with `O(log n)` runtime complexity.

**Example 1:**

```
Input: nums = [-1,0,3,5,9,12], target = 9
Output: 4
Explanation: 9 exists in nums and its index is 4
```

**Example 2:**

```
Input: nums = [-1,0,3,5,9,12], target = 2
Output: -1
Explanation: 2 does not exist in nums so return -1
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `-10^4 < nums[i], target < 10^4`
- All the integers in `nums` are **unique**.
- `nums` is sorted in ascending order.

<details>
<summary>官方中文题目翻译</summary>

给定一个 `n` 个元素有序的（升序）整型数组 `nums` 和一个目标值 `target`  ，写一个函数搜索 `nums` 中的 `target`，如果 `target` 存在返回下标，否则返回 `-1`。

你必须编写一个具有 `O(log n)` 时间复杂度的算法。

**示例 1:**

```
输入: nums = [-1,0,3,5,9,12], target = 9
输出: 4
解释: 9 出现在 nums 中并且下标为 4
```

**示例 2:**

```
输入: nums = [-1,0,3,5,9,12], target = 2
输出: -1
解释: 2 不存在 nums 中因此返回 -1
```

**提示：**

- 你可以假设 `nums` 中的所有元素是不重复的。
- `n` 将在 `[1, 10000]`之间。
- `nums` 的每个元素都将在 `[-9999, 9999]`之间。

</details>

## 题目大意

二分查找

## 解题思路

有序序列才可以使用二分查找

### 同类题

- [35. Search Insert Position](../0035.Search-Insert-Position/)

## 复杂度

- 时间复杂度：O(log(n))
- 空间复杂度：O(1)
