# [125. Valid Palindrome](https://leetcode.cn/problems/valid-palindrome/)

> 难度：Easy

## 题目

A phrase is a **palindrome** if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.

Given a string `s`, return `true` *if it is a **palindrome**, or* `false` *otherwise*.

**Example 1:**

```
Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome.
```

**Example 2:**

```
Input: s = "race a car"
Output: false
Explanation: "raceacar" is not a palindrome.
```

**Example 3:**

```
Input: s = " "
Output: true
Explanation: s is an empty string "" after removing non-alphanumeric characters.
Since an empty string reads the same forward and backward, it is a palindrome.
```

**Constraints:**

- `1 <= s.length <= 2 * 10^5`
- `s` consists only of printable ASCII characters.

<details>
<summary>官方中文题目翻译</summary>

如果在将所有大写字符转换为小写字符、并移除所有非字母数字字符之后，短语正着读和反着读都一样。则可以认为该短语是一个 **回文串** 。

字母和数字都属于字母数字字符。

给你一个字符串 `s`，如果它是 **回文串** ，返回 `true` ；否则，返回 `false` 。

**示例 1：**

```
输入: s = "A man, a plan, a canal: Panama"
输出：true
解释："amanaplanacanalpanama" 是回文串。
```

**示例 2：**

```
输入：s = "race a car"
输出：false
解释："raceacar" 不是回文串。
```

**示例 3：**

```
输入：s = " "
输出：true
解释：在移除非字母数字字符之后，s 是一个空字符串 "" 。
由于空字符串正着反着读都一样，所以是回文串。
```

**提示：**

- `1 <= s.length <= 2 * 10^5`
- `s` 仅由可打印的 ASCII 字符组成

</details>

## 题目大意

回文串: 将所有大写字符转换为小写字符、并移除所有非字母数字字符之后，短语正着读和反着读都一样

## 解题思路

双指针,注意回文的定义**字母数字**,对于特殊字符的额外处理,单个字母转小写可以避免整个字符串的复制造成的空间复杂度

### 同类题

- [283. Move Zeroes](../0283.Move-Zeroes/)：双指针
- [167. Two Sum II - Input Array Is Sorted](../0167.Two-Sum-II-Input-Array-Is-Sorted/)：双指针

## 复杂度

- 时间复杂度： O(n), 遍历1次字符串, 指针移动
- 空间复杂度： O(1), 任何时刻都只存在两个转换后的字符, 读取字符没有空间复杂度
