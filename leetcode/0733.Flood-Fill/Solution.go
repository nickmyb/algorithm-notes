package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func floodFill(image [][]int, sr int, sc int, color int) [][]int {
	graph := NewGraph(image)

	fill(graph, sr, sc, image[sr][sc], color)

	return graph.filled
}

func fill(graph *Graph, sr int, sc int, colorB int, colorA int) {
	// sr, sc
	if 0 <= sr && sr < graph.row && 0 <= sc && sc < graph.column && graph.marked[sr][sc] == false {
		graph.marked[sr][sc] = true

		if graph.image[sr][sc] == colorB {
			graph.filled[sr][sc] = colorA
		} else {
			return
		}

		fill(graph, sr-1, sc, colorB, colorA)
		fill(graph, sr+1, sc, colorB, colorA)
		fill(graph, sr, sc-1, colorB, colorA)
		fill(graph, sr, sc+1, colorB, colorA)
	} else {
		return
	}
}

type Graph struct {
	image  [][]int
	marked [][]bool
	filled [][]int
	row    int
	column int
}

func NewGraph(image [][]int) *Graph {
	g := &Graph{
		image:  image,
		row:    len(image),
		column: len(image[0]),
	}

	for _, row := range image {
		var markedRow []bool
		var filledRow []int

		for _, block := range row {
			markedRow = append(markedRow, false)
			filledRow = append(filledRow, block)
		}

		g.marked = append(g.marked, markedRow)
		g.filled = append(g.filled, filledRow)
	}

	return g
}
