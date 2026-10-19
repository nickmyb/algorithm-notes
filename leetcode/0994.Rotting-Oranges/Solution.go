package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func orangesRotting(grid [][]int) int {
	row := len(grid)
	col := len(grid[0])

	var q [][]int
	maxM := 0

	// 找所以烂橘子
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == 2 {
				q = append(q, []int{r, c, 0})
				grid[r][c] = 0
			}
		}
	}

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// BFS 骨架（单源；本题把所有腐烂橘子同时作为起点入队）：
	// q := []int{s}
	// marked[s] = true
	// for len(q) > 0 {
	//     v := q[0]
	//     q = q[1:]
	//     for _, w := range adj(v) {
	//         if !marked[w] {
	//             q = append(q, w)
	//             marked[w] = true
	//             edgeTo[w] = v // 需要还原路径时才记录
	//         }
	//     }
	// }
	// 本题用 grid 的值代替 marked，队列中的 m 记录腐烂所需分钟数，不需要 edgeTo。
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		x, y, m := cur[0], cur[1], cur[2]
		maxM = max(maxM, m)

		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx < 0 || nx >= row || ny < 0 || ny >= col || grid[nx][ny] == 0 {
				continue
			}

			q = append(q, []int{nx, ny, m + 1})
			grid[nx][ny] = 0
		}
	}

	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == 1 {
				return -1
			}
		}
	}

	return maxM
}
