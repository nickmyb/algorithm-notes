# [167. Two Sum II - Input Array Is Sorted](https://leetcode.cn/problems/two-sum-ii-input-array-is-sorted/)

> 难度：Medium

## 题目

You are given a **1-indexed** array of integers `numbers` that is already **sorted in non-decreasing order**.

Find **two** numbers such that they add up to a specific `target` number. Let these two numbers be `numbers[index_1]` and `numbers[index_2]` where `1 <= index_1 < index_2 <= numbers.length`.

Return the indices of the two numbers `index_1` and `index_2` as an integer array `[index_1, index_2]` of length 2.

The tests are generated such that there is **exactly one solution**. You **may not** use the same element twice.

Your solution must use only constant extra space.

**Example 1:**

```
Input: numbers = [2,7,11,15], target = 9
Output: [1,2]
Explanation: The sum of 2 and 7 is 9. Therefore, index1 = 1, index2 = 2. We return [1, 2].
```

**Example 2:**

```
Input: numbers = [2,3,4], target = 6
Output: [1,3]
Explanation: The sum of 2 and 4 is 6. Therefore index1 = 1, index2 = 3. We return [1, 3].
```

**Example 3:**

```
Input: numbers = [-1,0], target = -1
Output: [1,2]
Explanation: The sum of -1 and 0 is -1. Therefore index1 = 1, index2 = 2. We return [1, 2].
```

**Constraints:**

- `2 <= numbers.length <= 3 * 10^4`
- `-1000 <= numbers[i] <= 1000`
- `numbers` is sorted in **non-decreasing order**.
- `-1000 <= target <= 1000`
- The tests are generated such that there is **exactly one solution**.

<details>
<summary>官方中文题目翻译</summary>

给你一个下标从 **1** 开始的整数数组 `numbers` ，该数组已按** 非递减顺序排列**。

请你从数组中找出满足相加之和等于目标数 `target` 的 **两个** 数。令这两个数分别是 `numbers[index_1]` 和 `numbers[index_2]` ，其中 `1 <= index_1 < index_2 <= numbers.length` 。

以长度为 2 的整数数组 `[index_1, index_2]` 的形式返回这两个整数的下标 `index_1` 和 `index_2`。

你可以假设每个输入 **只对应唯一的答案** ，而且你 **不可以** 重复使用相同的元素。

你所设计的解决方案必须只使用常数级的额外空间。

**示例 1：**

```
输入：numbers = [2,7,11,15], target = 9
输出：[1,2]
解释：2 与 7 之和等于目标数 9 。因此 index1 = 1, index2 = 2 。返回 [1, 2] 。
```

**示例 2：**

```
输入：numbers = [2,3,4], target = 6
输出：[1,3]
解释：2 与 4 之和等于目标数 6 。因此 index1 = 1, index2 = 3 。返回 [1, 3] 。
```

**示例 3：**

```
输入：numbers = [-1,0], target = -1
输出：[1,2]
解释：-1 与 0 之和等于目标数 -1 。因此 index1 = 1, index2 = 2 。返回 [1, 2] 。
```

**提示：**

- `2 <= numbers.length <= 3 * 10^4`
- `-1000 <= numbers[i] <= 1000`
- `numbers` 按 **非递减顺序** 排列
- `-1000 <= target <= 1000`
- **仅存在一个有效答案**

</details>

## 题目大意

## 解题思路

- 要求空间复杂度O(1),那就不能使用map

1. 假设答案是`[x...(i, j)...y]`, `x,y要翻越答案前必须先到目标点,并且有题目的答案唯一性保证`。任意时刻, 假设y已经到了j, 因为x还没到i, 所以 x+y<target了,y不可能到y-1,对应x也一样
2. 暴力的遍历

### 同类题

- [125. Valid Palindrome](../0125.Valid-Palindrome/)：双指针
- [283. Move Zeroes](../0283.Move-Zeroes/)：双指针

## 复杂度

- 时间复杂度： 1. O(n); 2. O(n^2)
- 空间复杂度： O(1)
