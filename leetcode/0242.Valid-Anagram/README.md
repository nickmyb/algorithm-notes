# [242. Valid Anagram](https://leetcode.cn/problems/valid-anagram/)

> 难度：Easy

## 题目

Given two strings `s` and `t`, return `true` if `t` is an anagram of `s`, and `false` otherwise.

**Example 1:**

**Input:** s = "anagram", t = "nagaram"

**Output:** true

**Example 2:**

**Input:** s = "rat", t = "car"

**Output:** false

**Constraints:**

- `1 <= s.length, t.length <= 5 * 10^4`
- `s` and `t` consist of lowercase English letters.

**Follow up:** What if the inputs contain Unicode characters? How would you adapt your solution to such a case?

<details>
<summary>官方中文题目翻译</summary>

给定两个字符串 `s` 和 `t` ，编写一个函数来判断 `t` 是否是 `s` 的 字母异位词。

**示例 1:**

```
输入: s = "anagram", t = "nagaram"
输出: true
```

**示例 2:**

```
输入: s = "rat", t = "car"
输出: false
```

**提示:**

- `1 <= s.length, t.length <= 5 * 10^4`
- `s` 和 `t` 仅包含小写字母

**进阶:** 如果输入字符串包含 unicode 字符怎么办？你能否调整你的解法来应对这种情况？

</details>

## 题目大意

定义字母异位词: 字母异位词是通过重新排列不同单词或短语的字母而形成的单词或短语，并使用所有原字母一次。

## 解题思路

1. 排序str
2. 计算str的字母表签名

## 复杂度

- 时间复杂度： 1. O(n*log(n)); 2. O(n)
- 空间复杂度： O(n)
