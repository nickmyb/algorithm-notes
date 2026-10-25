# [11. Container With Most Water](https://leetcode.cn/problems/container-with-most-water/)

> 难度：Medium

## 题目

Given an array `height`, choose two vertical lines to form a container with the x-axis. Return the greatest amount of water it can hold. The container cannot be tilted.

**Example 1:**

**Input:** height = [1,8,6,2,5,4,8,3,7]

**Output:** 49

**Example 2:**

**Input:** height = [1,1]

**Output:** 1

**Constraints:**

- `n == height.length`
- `2 <= n <= 10^5`
- `0 <= height[i] <= 10^4`

## 题目大意

## 解题思路

```
当前左右端点是 l、r，假设 height[l] < height[r]，当前面积就是：

height[l] × (r-l)

如果保留左端点 l，把右端点换成内部任意位置 k：

新面积 = min(height[l], height[k]) × (k-l)
       ≤ height[l] × (k-l)
       ≤ height[l] × (r-l)
```

### 同类题

## 复杂度

- 时间复杂度：O(n)
- 空间复杂度：O(1)
