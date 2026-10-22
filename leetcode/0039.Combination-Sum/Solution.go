package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func combinationSum(candidates []int, target int) [][]int {
	var result [][]int
	var path []int

	var dfs func(start, remain int)
	dfs = func(start, remain int) {
		if remain == 0 {
			r := append([]int(nil), path...)
			result = append(result, r)
			return
		}

		for i := start; i < len(candidates); i++ {
			if candidates[i] > remain {
				continue
			}
			path = append(path, candidates[i])
			dfs(i, remain-candidates[i])
			path = path[:len(path)-1]
		}
	}

	dfs(0, target)
	return result
}

func combinationSumNumOrNot(candidates []int, target int) [][]int {
	var ret [][]int

	if target == 0 {
		return [][]int{}
	}

	if len(candidates) == 1 {
		if target%candidates[0] != 0 {
			return [][]int{}
		}

		var r []int
		for range target / candidates[0] {
			r = append(r, candidates[0])
		}
		return [][]int{r}
	}

	c := candidates[0]
	ret = append(ret, combinationSumNumOrNot([]int{c}, target)...)

	maxC := target / c
	var partCandidates []int
	partCandidates = candidates[1:]

	for numC := range maxC + 1 {
		partCombinationSum := combinationSumNumOrNot(partCandidates, target-numC*c)

		var r []int
		for range numC {
			r = append(r, c)
		}

		for _, pc := range partCombinationSum {
			// 要新创建slice
			combination := append(append([]int(nil), r...), pc...)
			ret = append(ret, combination)
		}
	}

	return ret
}
