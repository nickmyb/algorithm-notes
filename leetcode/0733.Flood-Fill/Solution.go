package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func floodFill(image [][]int, sr int, sc int, color int) [][]int {
	original := image[sr][sc]
	if original == color {
		return image
	}

	rows := len(image)
	cols := len(image[0])

	var dfs func(int, int)
	dfs = func(x, y int) {
		if x < 0 || x >= rows || y < 0 || y >= cols {
			return
		}
		if image[x][y] != original {
			return
		}

		image[x][y] = color
		dfs(x-1, y)
		dfs(x+1, y)
		dfs(x, y-1)
		dfs(x, y+1)
	}

	dfs(sr, sc)
	return image
}
