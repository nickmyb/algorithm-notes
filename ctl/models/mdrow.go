package models

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/nickmyb/algorithm-notes/ctl/util"
)

// dashRuns 匹配连续的多个 -。题目标题里本来就带 - 时（如 "K-diff Pairs in an Array"），
// 空格替换后会出现 -- 甚至 ---，这里统一收敛成一个 -。
var dashRuns = regexp.MustCompile(`-{2,}`)

// StandardizedTitle 把 LeetCode 的题目标题转成目录名：空格换成 -，
// 去掉在路径里不方便出现的标点，再把连续的 - 合成一个。
//
//	"Two Sum"                  → "Two-Sum"
//	"Best Time to Buy a Stock" → "Best-Time-to-Buy-a-Stock"
//	"K-diff Pairs in an Array" → "K-diff-Pairs-in-an-Array"
func StandardizedTitle(orig string) string {
	s := strings.TrimSpace(orig)
	// 去掉 / 尤其重要：留着会让题目目录凭空多一层。
	s = strings.NewReplacer(
		" ", "-", "'", "", "%", "", "(", "", ")", "",
		",", "", "?", "", "/", "", ":", "", `"`, "", "!", "",
	).Replace(s)
	s = dashRuns.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// Mdrow define
type Mdrow struct {
	FrontendQuestionID int32  `json:"question_id"`
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
	SolutionPath       string `json:"solution_path"`
	Acceptance         string `json:"acceptance"`
	Difficulty         string `json:"difficulty"`
	Frequency          string `json:"frequency"`
}

// GenerateMdRows 把扫描到的题解填进 README 的表格行。
// 有题解的题目，Solution 列渲染成各语言题解文件的链接；没题解的题目留空。
func GenerateMdRows(solutions []util.Solution, mdrows []Mdrow) {
	byID := make(map[int32]util.Solution, len(solutions))
	for _, s := range solutions {
		byID[int32(s.ID)] = s
	}

	matched := 0
	for i := range mdrows {
		mdrows[i].QuestionTitle = strings.TrimSpace(mdrows[i].QuestionTitle)
		s, ok := byID[mdrows[i].FrontendQuestionID]
		if !ok {
			mdrows[i].SolutionPath = ""
			continue
		}
		mdrows[i].SolutionPath = SolutionLinks(s)
		matched++
	}

	// 目录里有、但 LeetCode 题库里查不到的题号，多半是目录名写错了，提示一下而不是静默丢掉。
	if matched != len(solutions) {
		for _, s := range solutions {
			if _, ok := byID[int32(s.ID)]; !ok {
				fmt.Printf("【题号 %v 在 LeetCode 题库里不存在，检查目录名：%v】\n", s.ID, s.Dir)
			}
		}
	}
}

// SolutionLinks 把一道题的多语言题解渲染成 "[Go](…) [Python](…) [Java](…)"。
// 每个链接直接指向该语言的题解文件，而不是题目目录，这样三门语言的链接不会互相重复。
func SolutionLinks(s util.Solution) string {
	links := make([]string, 0, len(s.Files))
	for _, f := range s.Files {
		links = append(links, fmt.Sprintf("[%s](%s/blob/%s/leetcode/%s/%s)",
			f.Lang, util.RepoURL, util.DefaultBranch,
			url.PathEscape(s.Dir), url.PathEscape(f.Name)))
	}
	return strings.Join(links, " ")
}

// LanguageTable 渲染各语言的题解数量统计表。
func LanguageTable(solutions []util.Solution) string {
	counts := util.CountByLanguage(solutions)

	res := "|  语言  |  题解数  |  占比  |\n"
	res += "|:--------:|:--------:|:--------:|\n"
	for _, lang := range util.Languages {
		// 占比是「这门语言写了多少题 / 一共做了多少题」。一道题写了三门语言就三行都算一次，
		// 所以几列加起来可以超过 100%。
		res += fmt.Sprintf("|%v|%v|%v|\n",
			lang.Name, counts[lang.Name], percent(counts[lang.Name], len(solutions)))
	}
	// 合计行的占比恒等于 100%，是噪音，留空。
	res += fmt.Sprintf("|**合计题数**|**%v**|—|\n", len(solutions))
	return res
}

// percent 算百分比，分母为 0 时返回 "-"。
//
// 原版几处百分比都是直接相除：一道题都没收录时、或者没配 config.toml 导致 AC 数为 0 时，
// 0/0 会算出 NaN 并原样写进 README，表格里就是一片 "NaN%"。
func percent[T int | int32](num, den T) string {
	if den == 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f%%", float64(num)/float64(den)*100)
}

// | 0001 | Two Sum | [Go](…/leetcode/0001.Two-Sum/Solution.go) [Python](…) | 45.6% | Easy |
func (m Mdrow) tableLine() string {
	return fmt.Sprintf("|%04d|%v|%v|%v|%v|\n",
		m.FrontendQuestionID, m.QuestionTitle, m.SolutionPath, m.Acceptance, m.Difficulty)
}

// SortByQuestionID define
type SortByQuestionID []Mdrow

func (a SortByQuestionID) Len() int      { return len(a) }
func (a SortByQuestionID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByQuestionID) Less(i, j int) bool {
	return a[i].FrontendQuestionID < a[j].FrontendQuestionID
}

// Mdrows define
type Mdrows struct {
	Mdrows []Mdrow
}

// | No. | Title | Solution | Acceptance | Difficulty |
// |:---:|:---|:---:|:---:|:---:|
//
// 上游表头还有一列 Frequency，但 tableLine 从来没填过值——接口返回的 frequency
// 字段对未登录用户恒为 0，所以每一行末尾都是空的。这里直接去掉那一列。
func (mds Mdrows) table() string {
	res := "| No. |  Title  |  Solution  |  Acceptance |  Difficulty |\n"
	res += "|:--------:|:--------------------------------------------------------------|:--------:|:--------:|:--------:|\n"
	for _, p := range mds.Mdrows {
		res += p.tableLine()
	}
	return res
}

// AvailableTable define
func (mds Mdrows) AvailableTable() string {
	return mds.table()
}

// SolvedTable 只渲染已经写了题解的题目。题库有 4000+ 道题，全量表格太长，
// README 里默认只列自己做过的。
func (mds Mdrows) SolvedTable() string {
	solved := Mdrows{}
	for _, row := range mds.Mdrows {
		if row.SolutionPath != "" {
			solved.Mdrows = append(solved.Mdrows, row)
		}
	}
	sort.Sort(SortByQuestionID(solved.Mdrows))
	return solved.table()
}
