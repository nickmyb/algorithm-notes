package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func numIslands(grid [][]byte) int {
	var marked [][]bool
	num := 0

	row := len(grid)
	col := len(grid[0])

	for _, _row := range grid {
		var r []bool
		for range _row {
			r = append(r, false)
		}
		marked = append(marked, r)
	}

	var dfs func(int, int)
	dfs = func(x, y int) {
		if x < 0 || x >= row || y < 0 || y >= col {
			return
		}
		if marked[x][y] {
			return
		}

		marked[x][y] = true

		if grid[x][y] == '0' {
			return
		}
		dfs(x-1, y)
		dfs(x+1, y)
		dfs(x, y-1)
		dfs(x, y+1)
	}

	for i, _row := range grid {
		for j, _ := range _row {
			if !marked[i][j] && grid[i][j] == '1' {
				num += 1
				dfs(i, j)
			}
		}
	}

	return num
}
