# [1. Two Sum](https://leetcode.cn/problems/two-sum/)

> 难度：Easy

## 题目

You are given an array of integers `nums` and an integer `target`, return *indices of the two numbers such that they add up to `target`*.

You may assume that each input would have ***exactly* one solution**, and you may not use the *same* element twice.

You can return the answer in any order.

**Example 1:**

```
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].
```

**Example 2:**

```
Input: nums = [3,2,4], target = 6
Output: [1,2]
```

**Example 3:**

```
Input: nums = [3,3], target = 6
Output: [0,1]
```

**Constraints:**

- `2 <= nums.length <= 10^4`
- `-10^9 <= nums[i] <= 10^9`
- `-10^9 <= target <= 10^9`
- **Only one valid answer exists.**

**Follow-up:** Can you come up with an algorithm that is less than `O(n^2)` time complexity?

<details>
<summary>官方中文题目翻译</summary>

给定一个整数数组 `nums` 和一个整数目标值 `target`，请你在该数组中找出 **和为目标值** *`target`*  的那 **两个** 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。

**示例 1：**

```
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
```

**示例 2：**

```
输入：nums = [3,2,4], target = 6
输出：[1,2]
```

**示例 3：**

```
输入：nums = [3,3], target = 6
输出：[0,1]
```

**提示：**

- `2 <= nums.length <= 10^4`
- `-10^9 <= nums[i] <= 10^9`
- `-10^9 <= target <= 10^9`
- **只会存在一个有效答案**

**进阶：** 你可以想出一个时间复杂度小于 `O(n^2)` 的算法吗？

</details>

## 题目大意

计算数组中求和为目标结果的组合

## 解题思路

1. 通过遍历检查剩余值是否已经存在历史记录
2. twoSumIndexList 记录 对应数字 的 所有index list

## 复杂度

- 时间复杂度： twoSum = O(n), 一遍遍历即可; twoSumIndexList = O(n), 看起来是O(n^2), 但第二遍遍历nums的时候, 题目保证的唯一性可以保证l最大检索长度为2, 最坏情况为O(n) + O(2*n)
- 空间复杂度： Θ(n)
