package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func numIslands(grid [][]byte) int {
	num := 0

	row := len(grid)
	col := len(grid[0])

	var dfs func(int, int)
	dfs = func(x, y int) {
		if x < 0 || x >= row || y < 0 || y >= col {
			return
		}
		if grid[x][y] == '0' {
			return
		}

		grid[x][y] = '0'

		dfs(x-1, y)
		dfs(x+1, y)
		dfs(x, y-1)
		dfs(x, y+1)
	}

	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == '1' {
				num += 1
				dfs(r, c)
			}
		}
	}

	return num
}
