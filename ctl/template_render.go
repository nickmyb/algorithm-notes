//go:build ignore

// 用 html/template 渲染 README 的旧实现，已被 ctl/render.go 里的 renderReadme 取代
// （后者支持 {{.LanguageTable}} 等多语言占位符）。保留作参考，用 //go:build ignore 停用。

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	m "github.com/nickmyb/algorithm-notes/ctl/models"
	"github.com/nickmyb/algorithm-notes/ctl/util"
)

func makeReadmeFile(mdrows m.Mdrows) {
	file := "./README.md"
	os.Remove(file)
	var b bytes.Buffer
	tmpl := template.Must(template.New("readme").Parse(readTMPL("template.markdown")))
	err := tmpl.Execute(&b, mdrows)
	if err != nil {
		fmt.Println(err)
	}
	// 保存 README.md 文件
	util.WriteFile(file, b.Bytes())
}

func readTMPL(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}
