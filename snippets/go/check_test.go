package snippets

// 检查器写在 _test.go 里是必须的：本包非测试文件里的每个顶层函数都会被当成
// snippet，检查器的函数放进去就会把自己也登记成 snippet。

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// solutionRoots 是题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
var solutionRoots = []string{"../../leetcode"}

type funcDef struct {
	pos  string // file:line，报错用
	body string // 从 func 关键字到函数结尾，不含文档注释
}

// TestSolutionsMatchSnippets 是真正对仓库生效的检查，make snippets 和全量测试都跑它。
func TestSolutionsMatchSnippets(t *testing.T) {
	canon, errs := loadSnippets(".")
	uses, more := checkSolutions(canon, solutionRoots)
	for _, u := range uses {
		t.Log(u)
	}
	for _, e := range append(errs, more...) {
		t.Error(e)
	}
}

// loadSnippets 读取 dir 下非测试文件的顶层函数，并检查只 import 了标准库。
func loadSnippets(dir string) (map[string]funcDef, []string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, []string{err.Error()}
	}
	canon := map[string]funcDef{}
	var errs []string
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		fset, file, src, err := parseFile(path)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		for _, imp := range file.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			// 标准库的第一段路径不带点，第三方模块都以域名开头
			if strings.Contains(strings.SplitN(p, "/", 2)[0], ".") {
				errs = append(errs, fmt.Sprintf("%s: snippet 只能 import 标准库，%q 复制到 LeetCode 编译不过",
					fset.Position(imp.Pos()), p))
			}
		}
		for name, def := range topLevelFuncs(fset, file, src) {
			canon[name] = def
		}
	}
	return canon, errs
}

// checkSolutions 在题解里找和 snippet 同名的顶层函数，返回引用记录和不一致的报错。
func checkSolutions(canon map[string]funcDef, roots []string) (uses, errs []string) {
	if len(canon) == 0 {
		return nil, nil
	}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// out/ 是 IDEA 复制进题目目录的旧副本，理由同 gotest.sh 的过滤
				if path != root && (d.Name() == "out" || strings.HasPrefix(d.Name(), ".")) {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			fset, file, src, err := parseFile(path)
			if err != nil {
				errs = append(errs, err.Error())
				return nil
			}
			for name, got := range topLevelFuncs(fset, file, src) {
				want, ok := canon[name]
				if !ok {
					continue
				}
				if got.body == want.body {
					uses = append(uses, fmt.Sprintf("%s: %s 与 snippet 一致", got.pos, name))
				} else {
					errs = append(errs, fmt.Sprintf("%s: %s 与 %s 不一致\n%s",
						got.pos, name, want.pos, firstDiff(want.body, got.body)))
				}
			}
			return nil
		})
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	sort.Strings(uses)
	sort.Strings(errs)
	return uses, errs
}

func parseFile(path string) (*token.FileSet, *ast.File, []byte, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	return fset, file, src, err
}

func topLevelFuncs(fset *token.FileSet, file *ast.File, src []byte) map[string]funcDef {
	defs := map[string]funcDef{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil {
			continue
		}
		// 从 fn.Type.Func（func 关键字）开始截取，跳过 fn.Doc
		start := fset.Position(fn.Type.Func).Offset
		end := fset.Position(fn.End()).Offset
		defs[fn.Name.Name] = funcDef{
			pos:  fmt.Sprintf("%s:%d", filepath.ToSlash(fset.Position(fn.Pos()).Filename), fset.Position(fn.Pos()).Line),
			body: strings.ReplaceAll(string(src[start:end]), "\r\n", "\n"),
		}
	}
	return defs
}

// firstDiff 给出第一处不同的行，函数通常不长，这样足够定位。
func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			return fmt.Sprintf("  函数第 %d 行\n  snippet: %q\n  题解:    %q", i+1, wl, gl)
		}
	}
	return ""
}

// ---------- 检查器自身的测试，在临时目录里造样本 ----------

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const sampleSnippet = `package snippets

// abs 返回绝对值。
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
`

func TestCheckerAcceptsIdenticalCopyWithOwnDoc(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/abs.go": sampleSnippet,
		"leetcode/0001.A/Solution.go": `package leetcode

// abs 这题里用来算距离，文档注释可以和 snippet 不同。
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
`,
	})
	canon, errs := loadSnippets(filepath.Join(root, "snippets"))
	uses, more := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs)+len(more) != 0 || len(uses) != 1 {
		t.Fatalf("uses=%v errs=%v %v", uses, errs, more)
	}
}

func TestCheckerReportsDrift(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/abs.go": sampleSnippet,
		"leetcode/0001.A/Solution.go": `package leetcode

func abs(x int) int {
	if x <= 0 {
		return -x
	}
	return x
}
`,
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"))
	_, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 1 || !strings.Contains(errs[0], "Solution.go:3: abs") || !strings.Contains(errs[0], "x <= 0") {
		t.Fatalf("errs=%v", errs)
	}
}

func TestCheckerSkipsTestsMethodsAndIDECopies(t *testing.T) {
	stale := "package leetcode\n\nfunc abs(x int) int { return x }\n"
	root := writeFiles(t, map[string]string{
		"snippets/abs.go":                                   sampleSnippet,
		"leetcode/0001.A/Solution_test.go":                  stale,
		"leetcode/0001.A/out/production/0001.A/Solution.go": stale,
		"leetcode/0001.A/Solution.go":                       "package leetcode\n\ntype T struct{}\n\nfunc (T) abs(x int) int { return x }\n",
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"))
	uses, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(uses)+len(errs) != 0 {
		t.Fatalf("uses=%v errs=%v", uses, errs)
	}
}

func TestSnippetsOnlyImportStdlib(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/bad.go": "package snippets\n\nimport \"github.com/x/y\"\n\nfunc f() { y.Z() }\n",
		"snippets/ok.go":  "package snippets\n\nimport \"sort\"\n\nfunc g(a []int) { sort.Ints(a) }\n",
		// 测试文件可以用任何依赖，里面的函数也不算 snippet
		"snippets/ok_test.go": "package snippets\n\nimport \"github.com/x/y\"\n\nfunc h() { y.Z() }\n",
	})
	canon, errs := loadSnippets(filepath.Join(root, "snippets"))
	if len(errs) != 1 || !strings.Contains(errs[0], "github.com/x/y") {
		t.Fatalf("errs=%v", errs)
	}
	if _, ok := canon["h"]; ok {
		t.Fatal("测试文件里的函数不应算作 snippet")
	}
}
