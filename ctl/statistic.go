package main

import (
	"sort"

	m "github.com/nickmyb/algorithm-notes/ctl/models"
	"github.com/nickmyb/algorithm-notes/ctl/util"
)

func statisticalData(problemsMap map[int]m.StatStatusPairs, solutionIds []int) (easyTotal, mediumTotal, hardTotal, optimizingEasy, optimizingMedium, optimizingHard int32, optimizingIds []int) {
	easyTotal, mediumTotal, hardTotal, optimizingEasy, optimizingMedium, optimizingHard, optimizingIds = 0, 0, 0, 0, 0, 0, []int{}
	for _, v := range problemsMap {
		switch m.DifficultyMap[v.Difficulty.Level] {
		case "Easy":
			{
				easyTotal++
				if v.Status == "ac" && util.BinarySearch(solutionIds, v.Stat.FrontendQuestionID.Int()) == -1 {
					optimizingEasy++
					optimizingIds = append(optimizingIds, v.Stat.FrontendQuestionID.Int())
				}
			}
		case "Medium":
			{
				mediumTotal++
				if v.Status == "ac" && util.BinarySearch(solutionIds, v.Stat.FrontendQuestionID.Int()) == -1 {
					optimizingMedium++
					optimizingIds = append(optimizingIds, v.Stat.FrontendQuestionID.Int())
				}
			}
		case "Hard":
			{
				hardTotal++
				if v.Status == "ac" && util.BinarySearch(solutionIds, v.Stat.FrontendQuestionID.Int()) == -1 {
					optimizingHard++
					optimizingIds = append(optimizingIds, v.Stat.FrontendQuestionID.Int())
				}
			}
		}
	}
	sort.Ints(optimizingIds)
	return easyTotal, mediumTotal, hardTotal, optimizingEasy, optimizingMedium, optimizingHard, optimizingIds
}
