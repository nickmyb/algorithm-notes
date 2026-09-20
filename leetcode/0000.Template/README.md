# 0000. 题解骨架

这个目录不是题解，是**新开一题时的复制源**。题号 `0000` 被 ctl 特殊对待：扫描时跳过，不计入仓库 README 的任何统计。

新开一题不用手动复制，用 `make new ID=1`（等价于 `cd ctl && go run . new 1`）。完整的目录约定和写题流程见[仓库 README](../../README.md#目录结构)。

## 目录里有什么

| 文件 | 说明 |
|:---|:---|
| `Solution.go` / `Solution_test.go` | Go 题解与表驱动测试，`package leetcode` |
| `solution.py` / `solution_test.py` | Python 题解与 pytest 测试 |
| `Solution.java` / `SolutionTest.java` | Java 题解与 main 式测试，default package |
| `README.md` | 题目描述、题目大意、解题思路、复杂度 |

只写其中一两门语言完全可以，`make new` 加 `LANGS=go,python` 就只生成对应文件。仓库 README 表格里的 `Solution` 列按目录里**实际存在**的文件渲染，写了几门就显示几个链接。

## 题解 README 的格式

`make new` 会生成填好标题、链接、难度的 README，照着填即可：

```markdown
# [1. Two Sum](https://leetcode.com/problems/two-sum/)

> 难度：Easy

## 题目

（粘贴 LeetCode 原题描述）

## 题目大意

（一两句话说清楚题目要求）

## 解题思路

（思路推导，为什么这么做）

## 复杂度

- 时间复杂度：O(n)
- 空间复杂度：O(n)
```
