# [733. Flood Fill](https://leetcode.cn/problems/flood-fill/)

> 难度：Easy

## 题目

You are given an image represented by an `m x n` grid of integers `image`, where `image[i][j]` represents the pixel value of the image. You are also given three integers `sr`, `sc`, and `color`. Your task is to perform a **flood fill** on the image starting from the pixel `image[sr][sc]`.

To perform a **flood fill**:

- Begin with the starting pixel and change its color to `color`.
- Perform the same process for each pixel that is **directly adjacent** (pixels that share a side with the original pixel, either horizontally or vertically) and shares the **same color** as the starting pixel.
- Keep **repeating** this process by checking neighboring pixels of the *updated* pixels and modifying their color if it matches the original color of the starting pixel.
- The process **stops** when there are **no more** adjacent pixels of the original color to update.

Return the **modified** image after performing the flood fill.

**Example 1:**

**Input:** image = [[1,1,1],[1,1,0],[1,0,1]], sr = 1, sc = 1, color = 2

**Output:** [[2,2,2],[2,2,0],[2,0,1]]

**Explanation:**

![](https://assets.leetcode.com/uploads/2021/06/01/flood1-grid.jpg)

From the center of the image with position `(sr, sc) = (1, 1)` (i.e., the red pixel), all pixels connected by a path of the same color as the starting pixel (i.e., the blue pixels) are colored with the new color.

Note the bottom corner is **not** colored 2, because it is not horizontally or vertically connected to the starting pixel.

**Example 2:**

**Input:** image = [[0,0,0],[0,0,0]], sr = 0, sc = 0, color = 0

**Output:** [[0,0,0],[0,0,0]]

**Explanation:**

The starting pixel is already colored with 0, which is the same as the target color. Therefore, no changes are made to the image.

**Constraints:**

- `m == image.length`
- `n == image[i].length`
- `1 <= m, n <= 50`
- `0 <= image[i][j], color < 2^16`
- `0 <= sr < m`
- `0 <= sc < n`

<details>
<summary>官方中文题目翻译</summary>

有一幅以 `m x n` 的二维整数数组表示的图画 `image` ，其中 `image[i][j]` 表示该图画的像素值大小。你也被给予三个整数 `sr` ,  `sc` 和 `color` 。你应该从像素 `image[sr][sc]` 开始对图像进行上色 **填充** 。

为了完成 **上色工作**：

- 从初始像素开始，将其颜色改为 `color`。
- 对初始坐标的 **上下左右四个方向上** 相邻且与初始像素的原始颜色同色的像素点执行相同操作。
- 通过检查与初始像素的原始颜色相同的相邻像素并修改其颜色来继续 **重复** 此过程。
- 当 **没有** 其它原始颜色的相邻像素时 **停止** 操作。

最后返回经过上色渲染 **修改** 后的图像 。

**示例 1:**

![](https://assets.leetcode.com/uploads/2021/06/01/flood1-grid.jpg)

**输入：** image = [[1,1,1],[1,1,0],[1,0,1]]，sr = 1, sc = 1, color = 2

**输出：** [[2,2,2],[2,2,0],[2,0,1]]

**解释：** 在图像的正中间，坐标 `(sr,sc)=(1,1)` （即红色像素）,在路径上所有符合条件的像素点的颜色都被更改成相同的新颜色（即蓝色像素）。

注意，右下角的像素 **没有** 更改为2，因为它不是在上下左右四个方向上与初始点相连的像素点。

**示例 2:**

**输入：** image = [[0,0,0],[0,0,0]], sr = 0, sc = 0, color = 0

**输出：** [[0,0,0],[0,0,0]]

**解释：** 初始像素已经用 0 着色，这与目标颜色相同。因此，不会对图像进行任何更改。

**提示:**

- `m == image.length`
- `n == image[i].length`
- `1 <= m, n <= 50`
- `0 <= image[i][j], color < 2^16`
- `0 <= sr < m`
- `0 <= sc < n`

</details>

## 题目大意

从起始点出发,找上下左右数值相同的位置改数字

## 解题思路

1. 确定起点 → 检查边界和访问条件 → 标记当前格 → 遍历四个邻居

### 同类题

## 复杂度

- 时间复杂度： floodFill = O(mn)
- 空间复杂度： floodFill = O(mn)
