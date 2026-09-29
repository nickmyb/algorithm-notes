# [3. Longest Substring Without Repeating Characters](https://leetcode.cn/problems/longest-substring-without-repeating-characters/)

> 难度：Medium

## 题目

Given a string `s`, find the length of the **longest** **substring** without duplicate characters.

**Example 1:**

```
Input: s = "abcabcbb"
Output: 3
Explanation: The answer is "abc", with the length of 3. Note that "bca" and "cab" are also correct answers.
```

**Example 2:**

```
Input: s = "bbbbb"
Output: 1
Explanation: The answer is "b", with the length of 1.
```

**Example 3:**

```
Input: s = "pwwkew"
Output: 3
Explanation: The answer is "wke", with the length of 3.
Notice that the answer must be a substring, "pwke" is a subsequence and not a substring.
```

**Constraints:**

- `0 <= s.length <= 10^5`
- `s` consists of English letters, digits, symbols and spaces.

<details>
<summary>官方中文题目翻译</summary>

给定一个字符串 `s` ，请你找出其中不含有重复字符的 **最长 子串** 的长度。

**示例 1:**

```
输入: s = "abcabcbb"
输出: 3
解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。注意 "bca" 和 "cab" 也是正确答案。
```

**示例 2:**

```
输入: s = "bbbbb"
输出: 1
解释: 因为无重复字符的最长子串是 "b"，所以其长度为 1。
```

**示例 3:**

```
输入: s = "pwwkew"
输出: 3
解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
     请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。
```

**提示：**

- `0 <= s.length <= 10^5`
- `s` 由英文字母、数字、符号和空格组成

</details>

## 题目大意

最长非重复字符子串

## 解题思路

- 确定字符串的左右边界
  - 左边界直接从左向右依次扫描
  - 当前左边界和字符上次出现的右侧比较,选择更右侧的边界
- 注意考虑一些边界情况,可能导致异常

### 同类题

- [643. Maximum Average Subarray I](../0643.Maximum-Average-Subarray-I/)：固定窗口

## 复杂度

- 时间复杂度： lengthOfLongestSubstring = O(n), start, end 都遍历一次字符串 = O(2*n) = O(n); lengthOfLongestSubstringMap = O(n), 一次遍历 = O(n), start跳跃前进 = O(1) * n = O(n); lengthOfLongestSubstringTimeout = O(n^2);
- 空间复杂度： lengthOfLongestSubstring = O(n), 如果固定字符集的情况可以视为 = O(1); lengthOfLongestSubstringMap = O(n), 需要记录一个字符的map; lengthOfLongestSubstringTimeout = O(n);
