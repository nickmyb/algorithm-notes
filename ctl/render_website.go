//go:build ignore

// 这个文件里的代码负责渲染上游 LeetCode-Go 的 Hugo 站点（website/）：
// 第二章「算法专题」页面、书籍左侧目录，以及把题解 README 拷进第四章。
// 本仓库没有引入 website/，所以整个文件用 //go:build ignore 停用，不参与编译。
//
// 要恢复：
//  1. 把 website/ 目录、ctl/meta/ 里的复杂度元数据、ctl/template/ 下各专题的 .md 模板补回来
//  2. 删掉本文件第一行的 //go:build ignore
//  3. 在 ctl/render.go 的 newBuildCommand() 里放开 newBuildChapterTwo() / newBuildMenu()
//  4. 在 ctl/refresh.go 的 refresh() 里放开对应调用
//
// 相关的停用文件还有 ctl/label.go、ctl/pdf.go、ctl/util/website.go。

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"encoding/json"

	m "github.com/nickmyb/algorithm-notes/ctl/models"
	"github.com/nickmyb/algorithm-notes/ctl/util"
	"github.com/spf13/cobra"
)

var (
	chapterTwoList = []string{"Array", "String", "Two Pointers", "Linked List", "Stack", "Tree", "Dynamic Programming", "Backtracking", "Depth First Search", "Breadth First Search",
		"Binary Search", "Math", "Hash Table", "Sorting", "Bit Manipulation", "Union Find", "Sliding Window", "Segment Tree", "Binary Indexed Tree"}
	chapterTwoFileName = []string{"Array", "String", "Two_Pointers", "Linked_List", "Stack", "Tree", "Dynamic_Programming", "Backtracking", "Depth_First_Search", "Breadth_First_Search",
		"Binary_Search", "Math", "Hash_Table", "Sorting", "Bit_Manipulation", "Union_Find", "Sliding_Window", "Segment_Tree", "Binary_Indexed_Tree"}
	chapterTwoSlug = []string{"array", "string", "two-pointers", "linked-list", "stack", "tree", "dynamic-programming", "backtracking", "depth-first-search", "breadth-first-search",
		"binary-search", "math", "hash-table", "sorting", "bit-manipulation", "union-find", "sliding-window", "segment-tree", "binary-indexed-tree"}
)

func newBuildChapterTwo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chapter-two",
		Short: "Build Chapter Two commands",
		Run: func(cmd *cobra.Command, args []string) {
			buildChapterTwo(true)
		},
	}
	return cmd
}

func newBuildMenu() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "menu",
		Short: "Build Menu commands",
		Run: func(cmd *cobra.Command, args []string) {
			buildBookMenu()
		},
	}
	return cmd
}

// internal: true  渲染的链接都是 hugo 内部链接，用户生成 hugo web
//
//	false 渲染的链接是外部 HTTPS 链接，用于生成 PDF
func buildChapterTwo(internal bool) {
	var (
		gr        m.GraphQLResp
		questions []m.Question
		count     int
	)
	for index, tag := range chapterTwoSlug {
		body := getTagProblemList(tag)
		// 返回值校验：接口为空 / 解析失败 / 未取到题目（多见于被限流）时跳过该 tag，
		// 避免后续用空数据继续处理时偶发 panic。跳过的 tag 重新运行即可补齐。
		if len(body) == 0 {
			fmt.Printf("跳过 ChapterTwo[%v]：接口返回为空（可能被限流）\n", tag)
			continue
		}
		err := json.Unmarshal(body, &gr)
		if err != nil {
			fmt.Printf("跳过 ChapterTwo[%v]：JSON 解析失败 %v\n", tag, err)
			continue
		}
		questions = gr.Data.TopicTag.Questions
		if len(questions) == 0 {
			fmt.Printf("跳过 ChapterTwo[%v]：未取到题目数据（可能被限流）\n", tag)
			continue
		}
		mdrows := m.ConvertMdModelFromQuestions(questions)
		sort.Sort(m.SortByQuestionID(mdrows))
		solutions, _ := util.LoadSolutions()
		solutionIds := util.SolutionIDs(solutions)
		tl, err := loadMetaData(fmt.Sprintf("./meta/%v", chapterTwoFileName[index]))
		if err != nil {
			fmt.Printf("err = %v\n", err)
		}
		tls := m.GenerateTagMdRows(solutionIds, tl, mdrows, internal)
		//  按照模板渲染 README
		res, err := renderChapterTwo(fmt.Sprintf("./template/%v.md", chapterTwoFileName[index]), m.TagLists{TagLists: tls})
		if err != nil {
			fmt.Println(err)
			return
		}
		if internal {
			util.WriteFile(fmt.Sprintf("../website/content/ChapterTwo/%v.md", chapterTwoFileName[index]), res)
		} else {
			util.WriteFile(fmt.Sprintf("./pdftemp/ChapterTwo/%v.md", chapterTwoFileName[index]), res)
		}

		count++
	}
	fmt.Printf("write %v files successful", count)
}

func loadMetaData(filePath string) (map[int]m.TagList, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, metaMap := bufio.NewReader(f), map[int]m.TagList{}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return metaMap, nil
			}
			return nil, err
		}
		s := strings.Split(string(line), "|")
		// 字段不足的行（空行 / 格式异常）直接跳过，避免下面取 s[1]、s[4..6] 时越界 panic
		if len(s) < 7 {
			continue
		}
		v, _ := strconv.Atoi(strings.Split(s[1], ".")[0])
		// v[0] 是题号，s[4] time, s[5] space, s[6] favorite
		metaMap[v] = m.TagList{
			FrontendQuestionID: int32(v),
			Acceptance:         "",
			Difficulty:         "",
			TimeComplexity:     s[4],
			SpaceComplexity:    s[5],
			Favorite:           s[6],
		}
	}
}

func renderChapterTwo(filePath string, tls m.TagLists) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, output := bufio.NewReader(f), []byte{}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, nil
			}
			return nil, err
		}
		if ok, _ := regexp.Match("{{.AvailableTagTable}}", line); ok {
			reg := regexp.MustCompile("{{.AvailableTagTable}}")
			newByte := reg.ReplaceAll(line, []byte(tls.AvailableTagTable()))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
}

func buildBookMenu() {
	copyLackFile()
	// 按照模板重新渲染 Menu
	res, err := renderBookMenu("./template/menu.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	util.WriteFile("../website/content/menu/index.md", res)
	fmt.Println("generate Menu successful")
}

// 拷贝 leetcode 目录下的题解 README 文件至第四章对应文件夹中
func copyLackFile() {
	solutions, _ := util.LoadSolutions()
	_, ch4Ids := util.LoadChapterFourDir()

	needCopy := []string{}
	for _, s := range solutions {
		if util.BinarySearch(ch4Ids, s.ID) == -1 {
			needCopy = append(needCopy, s.Dir)
		}
	}
	if len(needCopy) > 0 {
		fmt.Printf("有 %v 道题需要拷贝到第四章中\n", len(needCopy))
		for i := 0; i < len(needCopy); i++ {
			if needCopy[i][4] == '.' {
				tmp, err := strconv.Atoi(needCopy[i][:4])
				if err != nil {
					fmt.Println(err)
				}
				err = os.MkdirAll(fmt.Sprintf("../website/content/ChapterFour/%v", util.GetChpaterFourFileNum(tmp)), os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
				util.CopyFile(fmt.Sprintf("../website/content/ChapterFour/%v/%v.md", util.GetChpaterFourFileNum(tmp), needCopy[i]), fmt.Sprintf("../leetcode/%v/README.md", needCopy[i]))
				util.CopyFile(fmt.Sprintf("../website/content/ChapterFour/%v/_index.md", util.GetChpaterFourFileNum(tmp)), "./template/collapseSection.md")
			}
		}
	} else {
		fmt.Printf("【第四章没有需要添加的题解，已经完整了】\n")
	}
}

func generateMenu() string {
	res := ""
	res += menuLine(chapterOneMenuOrder, "ChapterOne")
	res += menuLine(chapterTwoFileOrder, "ChapterTwo")
	res += menuLine(chapterThreeFileOrder, "ChapterThree")
	chapterFourFileOrder, _ := getChapterFourFileOrder()
	res += menuLine(chapterFourFileOrder, "ChapterFour")
	return res
}

func menuLine(order []string, chapter string) string {
	res := ""
	for i := 0; i < len(order); i++ {
		if i == 1 && chapter == "ChapterOne" {
			res += fmt.Sprintf("  - [%v]({{< relref \"/%v/%v\" >}})\n", chapterMap[chapter][order[i]], chapter, order[i])
			continue
		}
		if i == 0 {
			res += fmt.Sprintf("- [%v]({{< relref \"/%v/%v.md\" >}})\n", chapterMap[chapter][order[i]], chapter, order[i])
		} else {
			if chapter == "ChapterFour" {
				res += fmt.Sprintf("    - [%v]({{< relref \"/%v/%v.md\" >}})\n", order[i], chapter, order[i])
			} else {
				res += fmt.Sprintf("  - [%v]({{< relref \"/%v/%v.md\" >}})\n", chapterMap[chapter][order[i]], chapter, order[i])
			}
		}
	}
	return res
}

func renderBookMenu(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, output := bufio.NewReader(f), []byte{}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, nil
			}
			return nil, err
		}
		if ok, _ := regexp.Match("{{.BookMenu}}", line); ok {
			reg := regexp.MustCompile("{{.BookMenu}}")
			newByte := reg.ReplaceAll(line, []byte(generateMenu()))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
}
