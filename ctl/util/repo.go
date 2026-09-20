package util

// 仓库地址，用于拼 README 里的题解链接。换仓库或改默认分支时只需要改这里。
const (
	RepoURL       = "https://github.com/nickmyb/algorithm-notes"
	DefaultBranch = "main"
)

// Site 决定 ctl 从哪个 LeetCode 站点取题库数据，也决定 README 里题目链接指向哪边。
//
// 选 leetcode.cn 的理由（2026-09 实测两边的 /api/problems/all/）：
//
//  1. 题库是 leetcode.com 的严格超集。.com 的 4055 道编号题 .cn 一道不缺，
//     另外还多出 388 道 .com 上根本查不到的题：LCR 194、面试题 109、LCP 82、LCS 3。
//  2. 编号题的 question__title / question__title_slug 在 .cn 上仍然是英文
//     （题 1 是 "Two Sum" / "two-sum"，题 15 是 "3Sum" / "3sum"），
//     所以 0001.Two-Sum 这套目录命名在两个站点通用，换站点不用重命名任何目录。
//  3. 会员题数量两边一致（783 道），换站点不会多出或少掉付费内容。
//
// 代价有两个，都已经在代码里处理了：
//
//   - .cn 的 frontend_question_id 是 JSON 字符串而不是数字，而且混着 "LCP 82"、
//     "面试题 17.14" 这种非数字题号。见 models.QuestionID 的解析。
//   - 非数字题号没法放进 leetcode/%04d.Title 这套目录命名，ctl 会跳过它们。
//     这些题以后单独开 lcp/ 这样的同级目录收纳。
//
// 换回 leetcode.com 只要改这一行。注意 Cookie 必须来自同一个站点，
// .com 的 LEETCODE_SESSION 在 .cn 上无效，反之亦然。
const Site = "https://leetcode.cn"
