# [39. Combination Sum](https://leetcode.cn/problems/combination-sum/)

> 难度：Medium

## 题目

Given an array of **distinct** integers `candidates` and a target integer `target`, return *a list of all **unique combinations** of* `candidates` *where the chosen numbers sum to* `target`*.* You may return the combinations in **any order**.

The **same** number may be chosen from `candidates` an **unlimited number of times**. Two combinations are unique if the frequency of at least one of the chosen numbers is different.

The test cases are generated such that the number of unique combinations that sum up to `target` is less than `150` combinations for the given input.

**Example 1:**

```
Input: candidates = [2,3,6,7], target = 7
Output: [[2,2,3],[7]]
Explanation:
2 and 3 are candidates, and 2 + 2 + 3 = 7. Note that 2 can be used multiple times.
7 is a candidate, and 7 = 7.
These are the only two combinations.
```

**Example 2:**

```
Input: candidates = [2,3,5], target = 8
Output: [[2,2,2,2],[2,3,3],[3,5]]
```

**Example 3:**

```
Input: candidates = [2], target = 1
Output: []
```

**Constraints:**

- `1 <= candidates.length <= 30`
- `2 <= candidates[i] <= 40`
- All elements of `candidates` are **distinct**.
- `1 <= target <= 40`

<details>
<summary>官方中文题目翻译</summary>

给你一个 **无重复元素** 的整数数组 `candidates` 和一个目标整数 `target` ，找出 `candidates` 中可以使数字和为目标数 `target` 的 所有 **不同组合** ，并以列表形式返回。你可以按 **任意顺序** 返回这些组合。

`candidates` 中的 **同一个** 数字可以 **无限制重复被选取** 。如果至少一个数字的被选数量不同，则两种组合是不同的。

对于给定的输入，保证和为 `target` 的不同组合数少于 `150` 个。

**示例 1：**

```
输入：candidates = [2,3,6,7], target = 7
输出：[[2,2,3],[7]]
解释：
2 和 3 可以形成一组候选，2 + 2 + 3 = 7 。注意 2 可以使用多次。
7 也是一个候选， 7 = 7 。
仅有这两种组合。
```

**示例 2：**

```
输入: candidates = [2,3,5], target = 8
输出: [[2,2,2,2],[2,3,3],[3,5]]
```

**示例 3：**

```
输入: candidates = [2], target = 1
输出: []
```

**提示：**

- `1 <= candidates.length <= 30`
- `2 <= candidates[i] <= 40`
- `candidates` 的所有元素 **互不相同**
- `1 <= target <= 40`

</details>

## 题目大意

## 解题思路

1. combinationSum: 当前已选的数放在 path 里，还差 remain 才达到目标；接下来只能从下标 start 开始选
2. combinationSumNumOrNot(递归):
   - 只有一个数字num1的时候,有解 n*num1 / 无解
   - 多一个数字num2, (0...n)num2解 + 情况1

### 同类题

## 复杂度

```
先不用急着把时间写成只含 n 和 target 的式子；按“递归走了多少次 × 每次做多少事”来看更清楚。

设 n 是候选数个数，L = ⌊target / 最小候选数⌋ 是最多能选几个数，K 是答案数。

回溯版 combinationSum：

- 一次 dfs 最多检查 n 个候选数。
- 找到答案时，复制长度最多为 L 的 path，这件事共发生 K 次。
- 若实际调用了 R 次 dfs，时间就是 O(nR + KL)。
- 递归栈和 path 最多长 L；答案占 KL，所以总空间 O(L + KL)。

枚举次数版 combinationSumNumOrNot：

- 当前数最多选 L 次。代码对 numC=0…L 都构造一次前缀；长度相加为 0+1+…+L，即每个递归节点最多 O(L²)。
- 一个答案可能在返回途中被逐层复制，最多经过 n 层，每次最多复制 L 个数，总计 O(nKL)。
- 若它实际递归调用了 R₂ 次，时间可写为 O(R₂L² + nKL)。各层还会保存子问题结果，所以空间上界是 O(nL + nKL)。

这里的 R、R₂ 表示实际搜索过的节点数，包含最终没找到答案的分支。因此不能只看“题目保证答案少于 150 个”就认为搜索很快。先掌握这套逐项计数的方法，比记住一个很松的指数上界更有用。
```

- 时间复杂度： combinationSum = O(target/min(candidates)^2); combinationSumNumOrNot = O(n*target/min(candidates));
- 空间复杂度： combinationSum = O(n*target/min(candidates)); combinationSumNumOrNot = O(n*target/min(candidates));
