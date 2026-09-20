package models

import (
	"fmt"
)

// UserInfo define
type UserInfo struct {
	UserName         string `json:"user_name"`
	NumSolved        int32  `json:"num_solved"`
	NumTotal         int32  `json:"num_total"`
	AcEasy           int32  `json:"ac_easy"`
	AcMedium         int32  `json:"ac_medium"`
	AcHard           int32  `json:"ac_hard"`
	EasyTotal        int32
	MediumTotal      int32
	HardTotal        int32
	OptimizingEasy   int32
	OptimizingMedium int32
	OptimizingHard   int32
	FrequencyHigh    float64 `json:"frequency_high"`
	FrequencyMid     float64 `json:"frequency_mid"`
	CategorySlug     string  `json:"category_slug"`
}

// |    |  Easy  |  Medium  |  Hard |  Total |  optimizing |
// |:--------:|:--------------------------------------------------------------|:--------:|:--------:|:--------:|:--------:|
func (ui UserInfo) table() string {
	res := "|    |  Easy  |  Medium  |  Hard |  Total |\n"
	res += "|:--------:|:--------:|:--------:|:--------:|:--------:|\n"
	res += fmt.Sprintf("|Optimizing|%v|%v|%v|%v|\n", ui.OptimizingEasy, ui.OptimizingMedium, ui.OptimizingHard, ui.OptimizingEasy+ui.OptimizingMedium+ui.OptimizingHard)
	res += fmt.Sprintf("|Accepted|**%v**|**%v**|**%v**|**%v**|\n", ui.AcEasy, ui.AcMedium, ui.AcHard, ui.AcEasy+ui.AcMedium+ui.AcHard)
	res += fmt.Sprintf("|Total|%v|%v|%v|%v|\n", ui.EasyTotal, ui.MediumTotal, ui.HardTotal, ui.EasyTotal+ui.MediumTotal+ui.HardTotal)
	res += fmt.Sprintf("|Perfection Rate|%v|%v|%v|%v|\n",
		perfection(ui.OptimizingEasy, ui.AcEasy),
		perfection(ui.OptimizingMedium, ui.AcMedium),
		perfection(ui.OptimizingHard, ui.AcHard),
		perfection(ui.OptimizingEasy+ui.OptimizingMedium+ui.OptimizingHard, ui.AcEasy+ui.AcMedium+ui.AcHard))
	res += fmt.Sprintf("|Completion Rate|%v|%v|%v|%v|\n",
		percent(ui.AcEasy, ui.EasyTotal),
		percent(ui.AcMedium, ui.MediumTotal),
		percent(ui.AcHard, ui.HardTotal),
		percent(ui.AcEasy+ui.AcMedium+ui.AcHard, ui.EasyTotal+ui.MediumTotal+ui.HardTotal))
	return res
}

// perfection 算「已收录题解占已 AC 题目的比例」，即 1 - 待补题解数 / AC 数。
// 没配 config.toml 时 AC 数是 0，返回 "-" 而不是 NaN，见 percent。
func perfection(optimizing, accepted int32) string {
	if accepted == 0 {
		return "-"
	}
	return percent(accepted-optimizing, accepted)
}

// PersonalData define
func (ui UserInfo) PersonalData() string {
	return ui.table()
}
