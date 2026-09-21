package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	m "github.com/nickmyb/algorithm-notes/ctl/models"
	"github.com/nickmyb/algorithm-notes/ctl/util"
	"github.com/spf13/cobra"
)

// readmeTemplate 是 README 的模板，里面的 {{.Xxx}} 占位符由 renderReadme 替换。
const readmeTemplate = "./template/template.markdown"

func newBuildCommand() *cobra.Command {
	mc := &cobra.Command{
		Use:   "build <subcommand>",
		Short: "Build doc related commands",
	}
	mc.AddCommand(
		newBuildREADME(),
		// 下面两个子命令负责渲染 Hugo 站点的第二章（算法专题）和书籍目录，
		// 本仓库没有引入 website/，已随 ctl/render_website.go 一起停用。
		// newBuildChapterTwo(),
		// newBuildMenu(),
	)
	return mc
}

func newBuildREADME() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "readme",
		Short:        "Build readme.md commands",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildREADME()
		},
	}
	cmd.Flags().BoolVar(&anonymous, "anonymous", false,
		"忽略 config.toml，以未登录身份请求，不渲染个人数据（生成 template 分支的 README 时用）")
	return cmd
}

func buildREADME() error {
	var lpa m.LeetCodeProblemAll

	// 请求所有题目信息
	body := getProblemAllList()
	if len(body) == 0 {
		return fmt.Errorf("拉取 LeetCode 题库失败（网络不通或被限流），README 未改动")
	}
	if err := json.Unmarshal(body, &lpa); err != nil {
		return fmt.Errorf("解析题库数据失败，README 未改动: %w", err)
	}
	if len(lpa.StatStatusPairs) == 0 {
		return fmt.Errorf("题库数据为空，README 未改动")
	}

	// 拼凑 README 需要渲染的数据
	problems := lpa.StatStatusPairs
	info := m.ConvertUserInfoModel(lpa)
	// 非数字题号（leetcode.cn 独有的 LCR / 面试题 / LCP / LCS）的 Int() 都是 0，
	// 不跳过的话会全挤进 key 0 并污染难度统计。
	problemsMap := make(map[int]m.StatStatusPairs, len(problems))
	for _, v := range problems {
		if !v.Stat.FrontendQuestionID.Numbered() {
			continue
		}
		problemsMap[v.Stat.FrontendQuestionID.Int()] = v
	}

	mdrows := m.ConvertMdModelFromSsp(problems)
	sort.Sort(m.SortByQuestionID(mdrows))

	solutions, pending := util.LoadSolutions()
	m.GenerateMdRows(solutions, mdrows)

	var optimizingIds []int
	info.EasyTotal, info.MediumTotal, info.HardTotal,
		info.OptimizingEasy, info.OptimizingMedium, info.OptimizingHard,
		optimizingIds = statisticalData(problemsMap, util.SolutionIDs(solutions))
	omdrows := m.ConvertMdModelFromIds(problemsMap, optimizingIds)
	sort.Sort(m.SortByQuestionID(omdrows))

	// 按照模板渲染 README
	res, err := renderReadme(readmeTemplate, solutions, pending,
		m.Mdrows{Mdrows: mdrows}, m.Mdrows{Mdrows: omdrows}, info)
	if err != nil {
		return err
	}
	if err := util.WriteFile("../README.md", res); err != nil {
		return fmt.Errorf("写入 README 失败: %w", err)
	}
	fmt.Println("write file successful")
	reportPersonalData(info)
	return nil
}

// reportPersonalData 在写完 README 之后说明个人数据是否可用。
//
// 未登录时接口不报错，只是所有 AC 计数都返回 0，所以不渲染个人数据，
// 避免把"没拿到数据"误读成"一道题都没做"。
//
// 用两个信号区分三种状态：sentCookie 是本地事实（发没发凭据），info.UserName
// 是服务器的判断（认不认）。发了但服务器不认，就是 Cookie 过期——这种情况接口
// 照常返回 200，只是数据是未登录的，光看任何一个信号都会误判。
func reportPersonalData(info m.UserInfo) {
	switch {
	case info.UserName != "":
		solved := info.AcEasy + info.AcMedium + info.AcHard
		optimizing := info.OptimizingEasy + info.OptimizingMedium + info.OptimizingHard
		fmt.Printf("个人数据（%v）：已 AC %v 题，其中 %v 题还没在仓库里写题解\n",
			info.UserName, solved, optimizing)

	case sentCookie:
		fmt.Println("【带了 Cookie 但服务器按未登录处理，多半是 Cookie 过期了】")
		fmt.Println("  README 的「个人数据」一节因此没有渲染，这不代表你没做过题。")
		fmt.Println("  重新登录 leetcode.cn 取一次 csrftoken 和 LEETCODE_SESSION，更新 ctl/config.toml。")

	default:
		fmt.Println("【未登录，README 不渲染「个人数据」一节】")
		fmt.Println("  接口对未登录请求返回的 AC 计数全是 0，与其展示一张全 0 的表，不如不展示。")
		fmt.Println("  配好 ctl/config.toml 之后这一节会自动出现，见 ctl/README.md。")
	}
}

// renderReadme 读模板，把 {{.Xxx}} 占位符替换成渲染好的内容。
//
// 支持的占位符：
//
//	{{.PersonalData}}    个人 AC 数据表（需要 config.toml 里的 Cookie 才有数据）
//	{{.LanguageTable}}   各语言题解数量统计表
//	{{.TotalNum}}        一句话进度说明
//	{{.SolvedTable}}     已写题解的题目表格
//	{{.AvailableTable}}  LeetCode 全部题目的表格（4000+ 行，默认不用）
//	{{.OptimizingTable}} 已经 AC 但仓库里还没收录题解的题目表格
func renderReadme(filePath string, solutions []util.Solution, pending int,
	mdrows, omdrows m.Mdrows, user m.UserInfo) ([]byte, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	optimizing := user.OptimizingEasy + user.OptimizingMedium + user.OptimizingHard
	replacements := map[string]string{
		"{{.PersonalData}}":  user.PersonalData(),
		"{{.LanguageTable}}": m.LanguageTable(solutions),
		"{{.TotalNum}}": fmt.Sprintf("以下已经收录了 %v 道题的题解，还有 %v 道题在尝试中",
			len(solutions), pending),
		"{{.SolvedTable}}":    mdrows.SolvedTable(),
		"{{.AvailableTable}}": mdrows.AvailableTable(),
		"{{.OptimizingTable}}": fmt.Sprintf("以下 %v 道题已经 AC 但仓库里还没有收录题解\n\n%v",
			optimizing, omdrows.AvailableTable()),
	}

	out := string(src)
	for placeholder, value := range replacements {
		out = strings.ReplaceAll(out, placeholder, value)
	}
	return []byte(out), nil
}
