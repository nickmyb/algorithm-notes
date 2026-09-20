package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

// QuestionID 是题目在站点上显示的题号，也就是 frontend_question_id。
//
// 两个站点返回的 JSON 类型不一样，所以不能直接用 int32 接：
//
//	leetcode.com  "frontend_question_id": 1        ← 数字
//	leetcode.cn   "frontend_question_id": "1"      ← 字符串
//	leetcode.cn   "frontend_question_id": "LCP 82" ← 字符串，而且不是数字
//
// 本仓库默认走 leetcode.cn（理由见 util.Site），所以必须两种形态都能解析，
// 否则 json.Unmarshal 会在第一条 LCP 记录上直接失败。
//
// 非数字题号（LCR / 面试题 / LCP / LCS 共 388 道）的 Num 为 0，原文留在 Raw 里。
// 这类题没法放进 leetcode/%04d.Title 这套目录命名，调用方用 Numbered() 跳过。
type QuestionID struct {
	Num int32  // 数字题号；非数字题号为 0
	Raw string // 站点上显示的原始题号，如 "1" 或 "LCP 82"
}

// UnmarshalJSON 同时接受 JSON 数字和 JSON 字符串。
func (q *QuestionID) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	// 字符串形态先剥掉引号；用 json.Unmarshal 而不是手工 Trim，
	// 免得题号里出现转义字符时解析错。
	if strings.HasPrefix(s, `"`) {
		var unquoted string
		if err := json.Unmarshal(b, &unquoted); err != nil {
			return err
		}
		s = strings.TrimSpace(unquoted)
	}
	q.Raw = s

	n, err := strconv.Atoi(s)
	if err != nil {
		// 非数字题号不算错误，交给调用方按 Numbered() 过滤
		q.Num = 0
		return nil
	}
	q.Num = int32(n)
	return nil
}

// MarshalJSON 统一输出成字符串，保证写回缓存的 JSON 还能被自己读回来。
func (q QuestionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.Raw)
}

// Numbered 表示这是一道普通编号题，能放进 leetcode/%04d.Title 目录。
func (q QuestionID) Numbered() bool {
	return q.Num > 0
}

// Int 返回数字题号，方便和 map[int] 之类打交道。
func (q QuestionID) Int() int {
	return int(q.Num)
}
