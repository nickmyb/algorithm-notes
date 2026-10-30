# [239. Sliding Window Maximum](https://leetcode.cn/problems/sliding-window-maximum/)

> 难度：Hard

## 题目

You are given an array of integers `nums`, there is a sliding window of size `k` which is moving from the very left of the array to the very right. You can only see the `k` numbers in the window. Each time the sliding window moves right by one position.

Return *the max sliding window*.

**Example 1:**

```
Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
Output: [3,3,5,5,6,7]
Explanation:
Window position                Max
---------------               -----
[1  3  -1] -3  5  3  6  7       3
 1 [3  -1  -3] 5  3  6  7       3
 1  3 [-1  -3  5] 3  6  7       5
 1  3  -1 [-3  5  3] 6  7       5
 1  3  -1  -3 [5  3  6] 7       6
 1  3  -1  -3  5 [3  6  7]      7
```

**Example 2:**

```
Input: nums = [1], k = 1
Output: [1]
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `1 <= k <= nums.length`

<details>
<summary>官方中文题目翻译</summary>

给你一个整数数组 `nums`，有一个大小为 `k` 的滑动窗口从数组的最左侧移动到数组的最右侧。你只可以看到在滑动窗口内的 `k` 个数字。滑动窗口每次只向右移动一位。

返回 *滑动窗口中的最大值* 。

**示例 1：**

```
输入：nums = [1,3,-1,-3,5,3,6,7], k = 3
输出：[3,3,5,5,6,7]
解释：
滑动窗口的位置                最大值
---------------               -----
[1  3  -1] -3  5  3  6  7       3
 1 [3  -1  -3] 5  3  6  7       3
 1  3 [-1  -3  5] 3  6  7       5
 1  3  -1 [-3  5  3] 6  7       5
 1  3  -1  -3 [5  3  6] 7       6
 1  3  -1  -3  5 [3  6  7]      7
```

**示例 2：**

```
输入：nums = [1], k = 1
输出：[1]
```

**提示：**

- `1 <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `1 <= k <= nums.length`

</details>

## 题目大意

## 解题思路

1. 单调队列（双端队列），维护可能成为最大值的候选元素; 新元素可以淘汰前面比它小的元素: 因为新元素更大，而且更晚离开窗口，那些旧元素以后都没有机会成为最大值;

```
数组: 3、1、2

加入 3：队列 [3]
加入 1：队列 [3, 1]
加入 2：淘汰 1，队列 [3, 2]

这里为方便展示写的是值，实际存下标。保留 2，是为了等 3 离开窗口后接替它。

遍历到下标 i 时，依次做：

1. 从队首移除已经离开窗口的下标，即 下标 <= i-k。
2. 从队尾移除值小于等于 nums[i] 的下标。
3. 将 i 加入队尾。
4. 窗口凑满 k 个元素后，把队首对应的值加入答案。
```

### 同类题

## 复杂度

- 时间复杂度：O(n)
- 空间复杂度：O(1)
