# [643. Maximum Average Subarray I](https://leetcode.cn/problems/maximum-average-subarray-i/)

> 难度：Easy

## 题目

You are given an integer array `nums` consisting of `n` elements, and an integer `k`.

Find a contiguous subarray whose **length is equal to** `k` that has the maximum average value and return *this value*. Any answer with a calculation error less than `10^-5` will be accepted.

**Example 1:**

```
Input: nums = [1,12,-5,-6,50,3], k = 4
Output: 12.75000
Explanation: Maximum average is (12 - 5 - 6 + 50) / 4 = 51 / 4 = 12.75
```

**Example 2:**

```
Input: nums = [5], k = 1
Output: 5.00000
```

**Constraints:**

- `n == nums.length`
- `1 <= k <= n <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

<details>
<summary>官方中文题目翻译</summary>

给你一个由 `n` 个元素组成的整数数组 `nums` 和一个整数 `k` 。

请你找出平均数最大且 **长度为 `k`** 的连续子数组，并输出该最大平均数。

任何误差小于 `10^-5` 的答案都将被视为正确答案。

**示例 1：**

```
输入：nums = [1,12,-5,-6,50,3], k = 4
输出：12.75
解释：最大平均数 (12-5-6+50)/4 = 51/4 = 12.75
```

**示例 2：**

```
输入：nums = [5], k = 1
输出：5.00000
```

**提示：**

- `n == nums.length`
- `1 <= k <= n <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

</details>

## 题目大意

= 求最大子串和

## 解题思路

每次遍历求和超时,修改成记录历史和,利用历史和求新和减少O(nk)时间复杂度

### 同类题

- [3. Longest Substring Without Repeating Characters](../0003.Longest-Substring-Without-Repeating-Characters/)：可变窗口

## 复杂度

- 时间复杂度： O(n)
- 空间复杂度： O(1)
