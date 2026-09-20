package models

import (
	"fmt"
	"strings"
)

// LeetCodeProblemAll define
type LeetCodeProblemAll struct {
	UserName        string            `json:"user_name"`
	NumSolved       int32             `json:"num_solved"`
	NumTotal        int32             `json:"num_total"`
	AcEasy          int32             `json:"ac_easy"`
	AcMedium        int32             `json:"ac_medium"`
	AcHard          int32             `json:"ac_hard"`
	StatStatusPairs []StatStatusPairs `json:"stat_status_pairs"`
	FrequencyHigh   float64           `json:"frequency_high"`
	FrequencyMid    float64           `json:"frequency_mid"`
	CategorySlug    string            `json:"category_slug"`
	AcEasyTotal     int32
	AcMediumTotal   int32
	AcHardTotal     int32
}

// ConvertUserInfoModel define
func ConvertUserInfoModel(lpa LeetCodeProblemAll) UserInfo {
	info := UserInfo{}
	info.UserName = lpa.UserName
	info.NumSolved = lpa.NumSolved
	info.NumTotal = lpa.NumTotal
	info.AcEasy = lpa.AcEasy
	info.AcMedium = lpa.AcMedium
	info.AcHard = lpa.AcHard
	info.FrequencyHigh = lpa.FrequencyHigh
	info.FrequencyMid = lpa.FrequencyMid
	info.CategorySlug = lpa.CategorySlug
	return info
}

// StatStatusPairs define
type StatStatusPairs struct {
	Stat       Stat       `json:"stat"`
	Status     string     `json:"status"`
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
	IsFavor    bool       `json:"is_favor"`
	Frequency  float64    `json:"frequency"`
	Progress   float64    `json:"progress"`
}

// ConvertMdModelFromSsp define
//
// 非数字题号（leetcode.cn 独有的 LCR / 面试题 / LCP / LCS）放不进
// leetcode/%04d.Title 这套目录命名，直接跳过，不进 README 表格。
func ConvertMdModelFromSsp(problems []StatStatusPairs) []Mdrow {
	mdrows := []Mdrow{}
	for _, problem := range problems {
		if !problem.Stat.FrontendQuestionID.Numbered() {
			continue
		}
		mdrows = append(mdrows, toMdrow(problem))
	}
	return mdrows
}

// ConvertMdModelFromIds define
func ConvertMdModelFromIds(problemsMap map[int]StatStatusPairs, ids []int) []Mdrow {
	mdrows := []Mdrow{}
	for _, v := range ids {
		mdrows = append(mdrows, toMdrow(problemsMap[v]))
	}
	return mdrows
}

func toMdrow(problem StatStatusPairs) Mdrow {
	return Mdrow{
		FrontendQuestionID: problem.Stat.FrontendQuestionID.Num,
		QuestionTitle:      strings.TrimSpace(problem.Stat.QuestionTitle),
		QuestionTitleSlug:  strings.TrimSpace(problem.Stat.QuestionTitleSlug),
		Acceptance:         acceptance(problem.Stat),
		Difficulty:         DifficultyMap[problem.Difficulty.Level],
		Frequency:          fmt.Sprintf("%f", problem.Frequency),
	}
}

// acceptance 算通过率。原版直接相除，提交数为 0 的新题会算出 NaN 写进 README，这里挡一下。
func acceptance(s Stat) string {
	if s.TotalSubmitted == 0 {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", s.TotalAcs/s.TotalSubmitted*100)
}

// Stat define
type Stat struct {
	QuestionTitle     string  `json:"question__title"`
	QuestionTitleSlug string  `json:"question__title_slug"`
	TotalAcs          float64 `json:"total_acs"`
	TotalSubmitted    float64 `json:"total_submitted"`
	Acceptance        string
	Difficulty        string
	// 原版是 int32，但 leetcode.cn 返回的是字符串题号，而且混着 "LCP 82" 这种
	// 非数字题号，用 int32 会让整个 Unmarshal 失败。见 QuestionID。
	FrontendQuestionID QuestionID `json:"frontend_question_id"`
}

// Difficulty define
type Difficulty struct {
	Level int32 `json:"level"`
}

// DifficultyMap define
var DifficultyMap = map[int32]string{
	1: "Easy",
	2: "Medium",
	3: "Hard",
}
