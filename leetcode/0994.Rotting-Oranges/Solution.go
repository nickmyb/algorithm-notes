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
			}
		}
	}

	var bfs func(int, int, int)
	bfs = func(x, y int, m int) {
		if x < 0 || x >= row || y < 0 || y >= col {
			return
		}
		if grid[x][y] == 0 {
			return
		}

		// 发现新橘子就扩散
		grid[x][y] = 0
		maxM = max(maxM, m)

		q = append(q, []int{x - 1, y, m + 1})
		q = append(q, []int{x + 1, y, m + 1})
		q = append(q, []int{x, y - 1, m + 1})
		q = append(q, []int{x, y + 1, m + 1})
	}

	// bfs
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		bfs(cur[0], cur[1], cur[2])
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
