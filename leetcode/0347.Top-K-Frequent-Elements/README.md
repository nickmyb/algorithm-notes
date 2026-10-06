# [347. Top K Frequent Elements](https://leetcode.cn/problems/top-k-frequent-elements/)

> 难度：Medium

## 题目

Given an integer array `nums` and an integer `k`, return *the* `k` *most frequent elements*. You may return the answer in **any order**.

**Example 1:**

**Input:** nums = [1,1,1,2,2,3], k = 2

**Output:** [1,2]

**Example 2:**

**Input:** nums = [1], k = 1

**Output:** [1]

**Example 3:**

**Input:** nums = [1,2,1,2,1,2,3,1,3,2], k = 2

**Output:** [1,2]

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `k` is in the range `[1, the number of unique elements in the array]`.
- It is **guaranteed** that the answer is **unique**.

**Follow up:** Your algorithm's time complexity must be better than `O(n log n)`, where n is the array's size.

<details>
<summary>官方中文题目翻译</summary>

给你一个整数数组 `nums` 和一个整数 `k` ，请你返回其中出现频率前 `k` 高的元素。你可以按 **任意顺序** 返回答案。

**示例 1：**

**输入：** nums = [1,1,1,2,2,3], k = 2

**输出：** [1,2]

**示例 2：**

**输入：** nums = [1], k = 1

**输出：** [1]

**示例 3：**

**输入：** nums = [1,2,1,2,1,2,3,1,3,2], k = 2

**输出：** [1,2]

**提示：**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `k` 的取值范围是 `[1, 数组中不相同的元素的个数]`
- 题目数据保证答案唯一，换句话说，数组中前 `k` 个高频元素的集合是唯一的

**进阶：** 你所设计算法的时间复杂度 **必须** 优于 `O(n log n)` ，其中 `n` 是数组大小。

</details>

## 题目大意

## 解题思路

1. Top-K问题

### 同类题

- 前置：[215. Kth Largest Element in an Array](../0215.Kth-Largest-Element-in-an-Array/)：Top-K问题

## 复杂度

- 时间复杂度：topKFrequent = O(n + m log k); n 是数组长度，m 是不同元素的个数（m ≤ n）
- 空间复杂度：topKFrequent = O(m);
