# [20. Valid Parentheses](https://leetcode.cn/problems/valid-parentheses/)

> 难度：Easy

## 题目

Given a string `s` containing just the characters `'('`, `')'`, `'{'`, `'}'`, `'['` and `']'`, determine if the input string is valid.

An input string is valid if:

- Open brackets must be closed by the same type of brackets.
- Open brackets must be closed in the correct order.
- Every close bracket has a corresponding open bracket of the same type.

**Example 1:**

**Input:** s = "()"

**Output:** true

**Example 2:**

**Input:** s = "()[]{}"

**Output:** true

**Example 3:**

**Input:** s = "(]"

**Output:** false

**Example 4:**

**Input:** s = "([])"

**Output:** true

**Example 5:**

**Input:** s = "([)]"

**Output:** false

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists of parentheses only `'()[]{}'`.

<details>
<summary>官方中文题目翻译</summary>

给定一个只包括 `'('`，`')'`，`'{'`，`'}'`，`'['`，`']'` 的字符串 `s` ，判断字符串是否有效。

有效字符串需满足：

- 左括号必须用相同类型的右括号闭合。
- 左括号必须以正确的顺序闭合。
- 每个右括号都有一个对应的相同类型的左括号。

**示例 1：**

**输入：** s = "()"

**输出：** true

**示例 2：**

**输入：** s = "()[]{}"

**输出：** true

**示例 3：**

**输入：** s = "(]"

**输出：** false

**示例 4：**

**输入：** s = "([])"

**输出：** true

**示例 5：**

**输入：** s = "([)]"

**输出：** false

**提示：**

- `1 <= s.length <= 10^4`
- `s` 仅由括号 `'()[]{}'` 组成

</details>

## 题目大意

判断括号是否符合数学意义上的成对出现

## 解题思路

stack或者deque,特别注意ascii中括号不连续!!!

### 同类题

## 复杂度

- 时间复杂度： O(n)
- 空间复杂度： O(n)
