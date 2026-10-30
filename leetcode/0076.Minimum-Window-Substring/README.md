# [76. Minimum Window Substring](https://leetcode.cn/problems/minimum-window-substring/)

> 难度：Hard

## 题目

Given two strings `s` and `t` of lengths `m` and `n` respectively, return *the **minimum window*** ***substring*** *of* `s` *such that every character in* `t` *(**including duplicates**) is included in the window*. If there is no such substring, return *the empty string* `""`.

The testcases will be generated such that the answer is **unique**.

**Example 1:**

```
Input: s = "ADOBECODEBANC", t = "ABC"
Output: "BANC"
Explanation: The minimum window substring "BANC" includes 'A', 'B', and 'C' from string t.
```

**Example 2:**

```
Input: s = "a", t = "a"
Output: "a"
Explanation: The entire string s is the minimum window.
```

**Example 3:**

```
Input: s = "a", t = "aa"
Output: ""
Explanation: Both 'a's from t must be included in the window.
Since the largest window of s only has one 'a', return empty string.
```

**Constraints:**

- `m == s.length`
- `n == t.length`
- `1 <= m, n <= 10^5`
- `s` and `t` consist of uppercase and lowercase English letters.

**Follow up:** Could you find an algorithm that runs in `O(m + n)` time?

<details>
<summary>官方中文题目翻译</summary>

给定两个字符串 `s` 和 `t`，长度分别是 `m` 和 `n`，返回 s 中的 **最短窗口 子串**，使得该子串包含 `t` 中的每一个字符（**包括重复字符**）。如果没有这样的子串，返回空字符串 `""`。

测试用例保证答案唯一。

**示例 1：**

```
输入：s = "ADOBECODEBANC", t = "ABC"
输出："BANC"
解释：最小覆盖子串 "BANC" 包含来自字符串 t 的 'A'、'B' 和 'C'。
```

**示例 2：**

```
输入：s = "a", t = "a"
输出："a"
解释：整个字符串 s 是最小覆盖子串。
```

**示例 3:**

```
输入: s = "a", t = "aa"
输出: ""
解释: t 中两个字符 'a' 均应包含在 s 的子串中，
因此没有符合条件的子字符串，返回空字符串。
```

**提示：**

- `m == s.length`
- `n == t.length`
- `1 <= m, n <= 10^5`
- `s` 和 `t` 由英文字母组成

**进阶：** 你能设计一个在 `O(m + n)` 时间内解决此问题的算法吗？

</details>

## 题目大意

## 解题思路

1. 覆盖的子串要求每个字母的个数都比目标字符串多; 找到以某个字母结尾的子串满足条件后尝试缩小左边界,找到以某字母结尾的最小子串; 往后查找以其他字母结尾的满足要求的子串;

### 同类题

- [3. Longest Substring Without Repeating Characters](../0003.Longest-Substring-Without-Repeating-Characters/)：滑动窗口，求无重复字符的最长子串
- [438. Find All Anagrams in a String](../0438.Find-All-Anagrams-in-a-String/)：滑动窗口，查找异位词子串

## 复杂度

- 时间复杂度：O(n)
- 空间复杂度：O(1)
