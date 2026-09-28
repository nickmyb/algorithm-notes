# [283. Move Zeroes](https://leetcode.cn/problems/move-zeroes/)

> 难度：Easy

## 题目

Given an integer array `nums`, move all `0`'s to the end of it while maintaining the relative order of the non-zero elements.

**Note** that you must do this in-place without making a copy of the array.

**Example 1:**

```
Input: nums = [0,1,0,3,12]
Output: [1,3,12,0,0]
```

**Example 2:**

```
Input: nums = [0]
Output: [0]
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`

**Follow up:** Could you minimize the total number of operations done?

<details>
<summary>官方中文题目翻译</summary>

给定一个数组 `nums`，编写一个函数将所有 `0` 移动到数组的末尾，同时保持非零元素的相对顺序。

**请注意** ，必须在不复制数组的情况下原地对数组进行操作。

**示例 1:**

```
输入: nums = [0,1,0,3,12]
输出: [1,3,12,0,0]
```

**示例 2:**

```
输入: nums = [0]
输出: [0]
```

**提示**:

- `1 <= nums.length <= 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`

**进阶：** 你能尽量减少完成的操作次数吗？

</details>

## 题目大意

把数组中的0移动到末尾并保持其余数字顺序不变

## 解题思路

1. 对于任意index的数字,他前面有几个已知的0就需要向前移动几个数字,最后要把末尾的0补齐
2. 依次交换 0 和 non-0, 双指针

### 同类题

- [125. Valid Palindrome](../0125.Valid-Palindrome/)：双指针
- [167. Two Sum II - Input Array Is Sorted](../0167.Two-Sum-II-Input-Array-Is-Sorted/)：双指针

## 复杂度

- 时间复杂度： O(n)
- 空间复杂度： O(1)
