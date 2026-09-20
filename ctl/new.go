package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	m "github.com/nickmyb/algorithm-notes/ctl/models"
	"github.com/nickmyb/algorithm-notes/ctl/util"
	"github.com/spf13/cobra"
)

const (
	// templateDir 是多语言骨架目录，new 命令从这里复制题解文件。
	templateDir = util.SolutionsDir + "0000.Template"
	// problemsCache 缓存 LeetCode 题库，免得每开一题都打一次接口（顺便躲开限流）。
	problemsCache = "./.cache/problems.json"
	// problemsCacheTTL 是缓存有效期。题库变动很慢，一天刷一次足够。
	problemsCacheTTL = 24 * time.Hour
)

var newLangs string

func newNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <题号>",
		Short: "Scaffold a solution directory from leetcode/0000.Template",
		Long: "按题号去 LeetCode 查题目标题，在 leetcode/ 下建好目录，\n" +
			"再从 0000.Template 复制各语言的题解骨架和 README。",
		Example: "  algorithm-notes new 1\n  algorithm-notes new 15 --langs go,python",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 参数不对是用法错误，这时打出用法说明有帮助
			id, err := strconv.Atoi(args[0])
			if err != nil || id <= 0 {
				return fmt.Errorf("题号必须是正整数，收到 %q", args[0])
			}
			// 往下都是运行时错误（目录已存在、拉不到题库…），跟用法无关，
			// 再打一整篇 usage 只会把真正的错误信息淹掉
			cmd.SilenceUsage = true
			return scaffold(id, strings.Split(newLangs, ","))
		},
	}
	cmd.Flags().StringVar(&newLangs, "langs", allLangNames(), "要生成哪几门语言的骨架，逗号分隔")
	return cmd
}

func allLangNames() string {
	names := make([]string, 0, len(util.Languages))
	for _, l := range util.Languages {
		names = append(names, strings.ToLower(l.Name))
	}
	return strings.Join(names, ",")
}

// resolveLanguages 把 --langs 的取值翻译成 util.Languages 里的条目，保持声明顺序。
func resolveLanguages(names []string) ([]util.Language, error) {
	wanted := map[string]bool{}
	for _, n := range names {
		n = strings.ToLower(strings.TrimSpace(n))
		if n != "" {
			wanted[n] = true
		}
	}
	if len(wanted) == 0 {
		return nil, errors.New("--langs 不能为空")
	}

	langs := []util.Language{}
	for _, l := range util.Languages {
		if wanted[strings.ToLower(l.Name)] {
			delete(wanted, strings.ToLower(l.Name))
			langs = append(langs, l)
		}
	}
	if len(wanted) > 0 {
		unknown := make([]string, 0, len(wanted))
		for n := range wanted {
			unknown = append(unknown, n)
		}
		return nil, fmt.Errorf("不认识的语言 %v，可选：%v", unknown, allLangNames())
	}
	return langs, nil
}

func scaffold(id int, langNames []string) error {
	langs, err := resolveLanguages(langNames)
	if err != nil {
		return err
	}

	problem, err := findProblem(id)
	if err != nil {
		return err
	}

	title := strings.TrimSpace(problem.Stat.QuestionTitle)
	dir := filepath.Join(util.SolutionsDir, fmt.Sprintf("%04d.%s", id, m.StandardizedTitle(title)))
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("目录已存在，没有改动它：%v", dir)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	created := []string{"README.md"}
	slug := strings.TrimSpace(problem.Stat.QuestionTitleSlug)
	// 题目描述单独走一次 GraphQL。拉不到就留空小节，不影响目录创建。
	detail := getQuestionDetail(slug)
	if detail == nil {
		fmt.Println("【没取到题目描述，README 的题目小节留空，自己补一下】")
	}
	util.WriteFile(filepath.Join(dir, "README.md"), []byte(problemReadme(
		id, title, slug, m.DifficultyMap[problem.Difficulty.Level], detail)))

	for _, lang := range langs {
		files, err := templateFiles(lang)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Printf("【%v 在 %v 里没有骨架文件，跳过】\n", lang.Name, templateDir)
			continue
		}
		for _, name := range files {
			if _, err := util.CopyFile(filepath.Join(dir, name), filepath.Join(templateDir, name)); err != nil {
				return err
			}
			created = append(created, name)
		}
	}

	fmt.Printf("已创建 %v\n", dir)
	for _, name := range created {
		fmt.Printf("  %v\n", name)
	}
	return nil
}

// templateFiles 找出骨架目录里属于某门语言的文件（题解 + 测试都要）。
func templateFiles(lang util.Language) ([]string, error) {
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		return nil, fmt.Errorf("读取骨架目录失败: %w", err)
	}
	files := []string{}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == lang.Ext {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

// problemReadme 生成题解 README。
//
// 「题目」放英文原文，官方中文翻译折叠在下面 —— 两种语言都留着，但中文折叠起来
// 不会把 README 撑成两倍长。
//
// 「题目大意」故意留空：halfrost 原仓库里这一节是他自己写的一两句话概要
// （不含示例和约束），不是翻译。自己复述一遍题意是理解题目的一环，自动填就没意义了。
func problemReadme(id int, title, slug, difficulty string, detail *questionDetail) string {
	english, chinese := "", ""
	if detail != nil {
		english = htmlToMarkdown(detail.Content)
		chinese = htmlToMarkdown(detail.TranslatedContent)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# [%d. %s](%s/problems/%s/)\n\n", id, title, util.Site, slug)
	fmt.Fprintf(&b, "> 难度：%s\n\n", difficulty)

	b.WriteString("## 题目\n\n")
	if english != "" {
		b.WriteString(english + "\n\n")
	}
	if chinese != "" {
		b.WriteString("<details>\n<summary>官方中文题目翻译</summary>\n\n")
		b.WriteString(chinese + "\n\n")
		b.WriteString("</details>\n\n")
	}

	b.WriteString("## 题目大意\n\n")
	b.WriteString("## 解题思路\n\n")
	b.WriteString("## 复杂度\n\n")
	b.WriteString("- 时间复杂度：\n")
	b.WriteString("- 空间复杂度：\n")
	return b.String()
}

func findProblem(id int) (m.StatStatusPairs, error) {
	problems, err := loadProblems()
	if err != nil {
		return m.StatStatusPairs{}, err
	}
	for _, p := range problems {
		if p.Stat.FrontendQuestionID.Int() == id {
			if p.PaidOnly {
				fmt.Printf("【第 %v 题是会员题，题目描述需要登录后自己复制】\n", id)
			}
			return p, nil
		}
	}
	return m.StatStatusPairs{}, fmt.Errorf("LeetCode 题库里没有第 %v 题", id)
}

// loadProblems 取题库，优先用本地缓存。缓存过期或损坏时回源，回源失败但有旧缓存则用旧的。
func loadProblems() ([]m.StatStatusPairs, error) {
	if problems, ok := readProblemsCache(); ok {
		return problems, nil
	}

	body := getProblemAllList()
	if len(body) == 0 {
		// 网络不通时退回过期缓存，总比直接失败强。
		if problems, ok := readProblemsCacheIgnoringTTL(); ok {
			fmt.Println("拉取题库失败，使用过期的本地缓存")
			return problems, nil
		}
		return nil, errors.New("拉取 LeetCode 题库失败（网络不通或被限流）")
	}

	var lpa m.LeetCodeProblemAll
	if err := json.Unmarshal(body, &lpa); err != nil {
		return nil, fmt.Errorf("解析题库数据失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(problemsCache), 0755); err == nil {
		util.WriteFile(problemsCache, body)
	}
	return lpa.StatStatusPairs, nil
}

func readProblemsCache() ([]m.StatStatusPairs, bool) {
	info, err := os.Stat(problemsCache)
	if err != nil || time.Since(info.ModTime()) >= problemsCacheTTL {
		return nil, false
	}
	return readProblemsCacheIgnoringTTL()
}

func readProblemsCacheIgnoringTTL() ([]m.StatStatusPairs, bool) {
	body, err := os.ReadFile(problemsCache)
	if err != nil {
		return nil, false
	}
	var lpa m.LeetCodeProblemAll
	if err := json.Unmarshal(body, &lpa); err != nil || len(lpa.StatStatusPairs) == 0 {
		return nil, false
	}
	return lpa.StatStatusPairs, true
}
