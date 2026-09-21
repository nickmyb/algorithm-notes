package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SolutionsDir 是题解根目录。ctl 的所有命令都在 ctl/ 目录下执行，所以这里是相对路径。
const SolutionsDir = "../leetcode/"

// TemplateID 是骨架目录 leetcode/0000.Template 的题号。它只是新开一题时的复制源，
// 不是真正的题解，所以扫描时会跳过，不计入 README 的任何统计。
const TemplateID = 0

// Language 描述一门题解语言：怎么从文件名认出它，以及在 README 里怎么显示。
// 以后要加 C++、Rust，在 Languages 里加一行就够了，其余代码不用动。
type Language struct {
	Name       string // README 链接上显示的名字
	Ext        string // 题解文件后缀
	TestSuffix string // 测试文件后缀；命中的文件是测试，不算题解
	Entry      string // 题解入口文件名；目录里有多个该语言的文件时，README 链到它
}

// Languages 是仓库支持的语言，README 里的题解链接按这个顺序排列。
var Languages = []Language{
	{Name: "Go", Ext: ".go", TestSuffix: "_test.go", Entry: "Solution.go"},
	{Name: "Python", Ext: ".py", TestSuffix: "_test.py", Entry: "solution.py"},
	{Name: "Java", Ext: ".java", TestSuffix: "Test.java", Entry: "Solution.java"},
}

// SolutionFile 是某道题下某一门语言的题解文件。
type SolutionFile struct {
	Lang string // 语言名，取自 Language.Name
	Name string // 文件名，如 "Solution.go"
}

// Solution 是 leetcode/ 下的一个题目目录。
type Solution struct {
	ID    int            // 题号
	Dir   string         // 目录名，如 "0001.Two-Sum"
	Files []SolutionFile // 各语言的题解文件，按 Languages 的顺序排列
}

// dirPattern 匹配 "0001.Two-Sum" 这样的题目目录名。
// 原版用 name[4] == '.' 判断，遇到短于 5 个字符的文件名（如 .DS_Store）会越界 panic。
var dirPattern = regexp.MustCompile(`^(\d{4})\.(.+)$`)

// LoadSolutions 扫描题解目录，返回写了题解的题目（按题号升序）和
// 建了目录但还没放进任何题解文件的题目数量（即"在尝试中"的题）。
func LoadSolutions() (solutions []Solution, pending int) {
	entries, err := os.ReadDir(SolutionsDir)
	if err != nil {
		fmt.Printf("读取题解目录失败: %v\n", err)
		return nil, 0
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		match := dirPattern.FindStringSubmatch(e.Name())
		if match == nil {
			continue
		}
		id, err := strconv.Atoi(match[1])
		if err != nil {
			fmt.Printf("目录名里的题号解析失败: %v (%v)\n", e.Name(), err)
			continue
		}
		if id == TemplateID {
			continue
		}
		files, err := detectLanguages(filepath.Join(SolutionsDir, e.Name()))
		if err != nil {
			fmt.Printf("读取 %v 失败: %v\n", e.Name(), err)
			continue
		}
		if len(files) == 0 {
			pending++
			continue
		}
		solutions = append(solutions, Solution{ID: id, Dir: e.Name(), Files: files})
	}

	sort.Slice(solutions, func(i, j int) bool { return solutions[i].ID < solutions[j].ID })
	fmt.Printf("读取了 %v 道题的题解，另有 %v 道题建了目录但还没有题解\n", len(solutions), pending)
	return solutions, pending
}

// detectLanguages 找出一个题目目录里每门语言的题解入口文件。
//
// 一道题一门语言通常只有一个文件，但有些题需要单独的辅助数据结构（比如线段树写在
// SegmentTree.go 里），目录里就会出现同语言的多个文件。这时优先认 Language.Entry
// 指定的入口文件；没有入口文件才退回字典序第一个，免得 README 链到辅助文件上去。
func detectLanguages(dir string) ([]SolutionFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	files := []SolutionFile{}
	for _, lang := range Languages {
		candidates := []string{}
		for _, name := range names {
			if strings.HasSuffix(name, lang.TestSuffix) || filepath.Ext(name) != lang.Ext {
				continue
			}
			candidates = append(candidates, name)
		}
		if len(candidates) == 0 {
			continue
		}

		entry := candidates[0]
		for _, name := range candidates {
			if name == lang.Entry {
				entry = name
				break
			}
		}
		files = append(files, SolutionFile{Lang: lang.Name, Name: entry})
	}
	return files, nil
}

// SolutionIDs 取出题号列表，顺序与入参一致（LoadSolutions 已保证升序），
// 供 BinarySearch 和统计代码使用。
func SolutionIDs(solutions []Solution) []int {
	ids := make([]int, 0, len(solutions))
	for _, s := range solutions {
		ids = append(ids, s.ID)
	}
	return ids
}

// CountByLanguage 统计每门语言各写了多少道题的题解。
func CountByLanguage(solutions []Solution) map[string]int {
	counts := map[string]int{}
	for _, s := range solutions {
		for _, f := range s.Files {
			counts[f.Lang]++
		}
	}
	return counts
}

// WriteFile define
func WriteFile(fileName string, content []byte) error {
	// os.WriteFile 会截断旧内容，并把打开、写入、关闭文件的错误传回调用方。
	return os.WriteFile(fileName, content, 0644)
}

// LoadFile define
func LoadFile(filePath string) ([]byte, error) {
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
		output = append(output, line...)
		output = append(output, []byte("\n")...)
	}
}

// CopyFile define
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// BinarySearch define
func BinarySearch(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + (high-low)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}
