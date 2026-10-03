# [35. Search Insert Position](https://leetcode.cn/problems/search-insert-position/)

> 难度：Easy

## 题目

Given a sorted array of distinct integers and a target value, return the index if the target is found. If not, return the index where it would be if it were inserted in order.

You must write an algorithm with `O(log n)` runtime complexity.

**Example 1:**

```
Input: nums = [1,3,5,6], target = 5
Output: 2
```

**Example 2:**

```
Input: nums = [1,3,5,6], target = 2
Output: 1
```

**Example 3:**

```
Input: nums = [1,3,5,6], target = 7
Output: 4
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `-10^4 <= nums[i] <= 10^4`
- `nums` contains **distinct** values sorted in **ascending** order.
- `-10^4 <= target <= 10^4`

<details>
<summary>官方中文题目翻译</summary>

给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。

请必须使用时间复杂度为 `O(log n)` 的算法。

**示例 1:**

```
输入: nums = [1,3,5,6], target = 5
输出: 2
```

**示例 2:**

```
输入: nums = [1,3,5,6], target = 2
输出: 1
```

**示例 3:**

```
输入: nums = [1,3,5,6], target = 7
输出: 4
```

**提示:**

- `1 <= nums.length <= 10^4`
- `-10^4 <= nums[i] <= 10^4`
- `nums` 为 **无重复元素** 的 **升序** 排列数组
- `-10^4 <= target <= 10^4`

</details>

## 题目大意

找到元素应该准确插入的index位置

## 解题思路

1. 二分搜索
2. 确定插入的index

```
claude code给出的解释

循环过程中一直有这样的关系：
- lo 左边的元素都小于 target（因为 lo = mid + 1 只在 nums[mid] < target 时发生）；
- hi 右边的元素都大于 target（因为 hi = mid - 1 只在 nums[mid] > target 时发生）。

不变式靠什么成立
- 一开始 lo = 0，左边没有元素，自然成立。
- lo 只在 nums[mid] < target 时变成 mid + 1。数组有序，nums[0..mid] 都 ≤ nums[mid]，所以都 < target。
- hi 这一侧同理：hi + 1 及右边的元素都 > target。

循环结束时 lo == hi + 1，所以 lo 正好是第一个大于 target 的位置，也就是插入位置。两端也成立：比所有元素都小时 lo 停在 0，比所有元素都大时 lo 停在 len(nums)。不用再像现在这样比一次 nums[mid] 来决定是 mid 还是 mid + 1。
```

### 同类题

- 前置：[704. Binary Search](../0704.Binary-Search/)：二分查找
- [34. Find First and Last Position of Element in Sorted Array](../0034.Find-First-and-Last-Position-of-Element-in-Sorted-Array/)

## 复杂度

- 时间复杂度：O(log(n))
- 空间复杂度：O(1)
