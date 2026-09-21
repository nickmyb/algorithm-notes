package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const workflowProblems = `{"stat_status_pairs":[{"stat":{"frontend_question_id":"1","question__title":"Two Sum","question__title_slug":"two-sum"},"difficulty":{"level":1}}]}`

// 每个流程都在独立仓库中运行，接口用本地服务替代，不触碰作者题解或真实 Cookie。
func workflowRepo(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"ctl/template", "leetcode/0000.Template"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"Solution.go", "Solution_test.go", "solution.py", "solution_test.py", "Solution.java", "SolutionTest.java"} {
		body, err := os.ReadFile(filepath.Join("../leetcode/0000.Template", name))
		if err != nil {
			t.Fatal(err)
		}
		workflowWrite(t, filepath.Join(root, "leetcode/0000.Template", name), string(body))
	}
	template, err := os.ReadFile(readmeTemplate)
	if err != nil {
		t.Fatal(err)
	}
	workflowWrite(t, filepath.Join(root, "ctl/template/template.markdown"), string(template))
	t.Chdir(filepath.Join(root, "ctl"))
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	oldReq, oldAnonymous, oldCookie := req, anonymous, sentCookie
	oldAll, oldGraphql := AllProblemURL, QraphqlURL
	t.Cleanup(func() {
		req, anonymous, sentCookie = oldReq, oldAnonymous, oldCookie
		AllProblemURL, QraphqlURL = oldAll, oldGraphql
	})
	req, anonymous, sentCookie = nil, true, false
	AllProblemURL, QraphqlURL = server.URL+"/all", server.URL+"/graphql"
	return root
}

func workflowWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func workflowAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/graphql" {
		fmt.Fprint(w, `{"data":{"question":{"content":"<p>English fixture</p>","translatedContent":"<p>中文测试数据</p>"}}}`)
		return
	}
	fmt.Fprint(w, workflowProblems)
}

func TestNewProblemAndAnonymousReadme(t *testing.T) {
	root := workflowRepo(t, workflowAPI)
	cmd := newNewCommand()
	cmd.SetArgs([]string{"1", "--langs", "go,python"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	problem := filepath.Join(root, "leetcode/0001.Two-Sum")
	for _, file := range []string{"Solution.go", "Solution_test.go", "solution.py", "solution_test.py", "README.md"} {
		if _, err := os.Stat(filepath.Join(problem, file)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(problem, "Solution.java")); !os.IsNotExist(err) {
		t.Fatal("--langs go,python 不应创建 Java 文件")
	}
	readme, err := os.ReadFile(filepath.Join(problem, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"English fixture", "中文测试数据", "## 题目大意\n\n## 解题思路"} {
		if !strings.Contains(string(readme), text) {
			t.Errorf("题目 README 缺少 %q", text)
		}
	}
	workflowWrite(t, filepath.Join(problem, "Solution.go"), "作者的代码，不应覆盖")
	if err := scaffold(1, []string{"go"}); err == nil {
		t.Fatal("重复建题应该失败")
	}
	body, _ := os.ReadFile(filepath.Join(problem, "Solution.go"))
	if string(body) != "作者的代码，不应覆盖" {
		t.Fatal("重复建题覆盖了作者的代码")
	}

	build := newBuildREADME()
	build.SetArgs([]string{"--anonymous"})
	if err := build.Execute(); err != nil {
		t.Fatal(err)
	}
	generated, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	if !strings.Contains(text, "|0001|Two Sum|[Go](") || !strings.Contains(text, ") [Python](") {
		t.Fatal("生成的题目表格缺少新题或语言链接")
	}
	if strings.Contains(text, "## 个人数据") || strings.Contains(text, "{{.") {
		t.Fatal("匿名 README 不应出现个人数据或未替换的占位符")
	}
}

func TestReadmeFailureKeepsExistingFile(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"HTTP failure", 503, "unavailable"},
		{"invalid JSON", 200, "not JSON"},
		{"empty catalogue", 200, `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := workflowRepo(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			})
			file := filepath.Join(root, "README.md")
			workflowWrite(t, file, "existing README")
			cmd := newBuildREADME()
			cmd.SetArgs([]string{"--anonymous"})
			if err := cmd.Execute(); err == nil {
				t.Fatal("生成失败必须向命令行返回错误")
			}
			body, _ := os.ReadFile(file)
			if string(body) != "existing README" {
				t.Fatal("请求失败改动了已有 README")
			}
		})
	}
}

func TestReadmeWriteFailureIsReported(t *testing.T) {
	root := workflowRepo(t, workflowAPI)
	if err := os.Mkdir(filepath.Join(root, "README.md"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := buildREADME(); err == nil {
		t.Fatal("README 路径不可写时必须失败")
	}
}
