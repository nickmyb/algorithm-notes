//go:build ignore

// 这个文件里的工具函数只服务于上游 LeetCode-Go 的 Hugo 站点（website/）和 PDF 电子书，
// 本仓库没有引入 website/，所以整个文件用 //go:build ignore 停用，不参与编译。
//
// 要恢复：把仓库顶上的 website/ 目录补回来，然后删掉本文件第一行的 //go:build ignore。
// 相关的停用文件还有 ctl/render_website.go、ctl/label.go、ctl/pdf.go。

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// GetAllFile define
func GetAllFile(pathname string, fileList *[]string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+fi.Name()+"/", fileList)
		} else {
			*fileList = append(*fileList, fi.Name())
		}
	}
	return *fileList, err
}

// LoadChapterFourDir define
func LoadChapterFourDir() ([]string, []int) {
	files, err := GetAllFile("../website/content/ChapterFour/", &[]string{})
	if err != nil {
		fmt.Println(err)
	}
	solutions, solutionIds, solutionsMap := []string{}, []int{}, map[int]string{}
	for _, f := range files {
		if f[4] == '.' {
			tmp, err := strconv.Atoi(f[:4])
			if err != nil {
				fmt.Println(err)
			}
			solutionIds = append(solutionIds, tmp)
			// len(f.Name())-3 = 文件名去掉 .md 后缀
			solutionsMap[tmp] = f[:len(f)-3]
		}
	}
	sort.Ints(solutionIds)
	fmt.Printf("读取了第四章的 %v 道题的题解\n", len(solutionIds))
	for _, v := range solutionIds {
		if name, ok := solutionsMap[v]; ok {
			solutions = append(solutions, name)
		}
	}
	return solutions, solutionIds
}

// DestoryDir define
func DestoryDir(path string) {
	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		name := fi.Name()
		if strings.Contains(name, "temp") {
			fmt.Println("temp file name:", path)
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delet dir error:", err)
			}
		}
		return nil
	})
}

// GetChpaterFourFileNum define
func GetChpaterFourFileNum(num int) string {
	if num < 100 {
		return fmt.Sprintf("%04d~%04d", (num/100)*100+1, (num/100)*100+99)
	}
	return fmt.Sprintf("%04d~%04d", (num/100)*100, (num/100)*100+99)
}
