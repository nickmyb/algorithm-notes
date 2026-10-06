package snippets

// 检查器写在 _test.go 里是必须的：本包非测试文件里的每个顶层函数、类型和方法
// 都会被当成 snippet，检查器的函数放进去就会把自己也登记成 snippet。

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

// structuresDir 是 LeetCode 平台类型的本地实现。snippet 里和它同名的类型
// （如 traversal.go 的 TreeNode）只是为了让 snippet 自己能编译的替身，
// 题解里用类型别名接到 structures，不复制，所以不算 snippet。
const structuresDir = "../../structures/go"

// funcDef 是一个 snippet：顶层函数、类型或方法。名字（map 的键）的写法：
//
//	abs           顶层函数
//	type MaxPQ    类型，带 "type " 前缀，不会和函数、方法重名
//	MaxPQ.sink    方法，接收者类型去掉 * 和类型参数
type funcDef struct {
	file string // 所在文件名，关联关系按它分组
	line int
	pos  string // file:line，报错用
	sig  string // 从 func 关键字到返回值结尾，即函数头；类型为空
	body string // 从 func / type 关键字到结尾，不含文档注释
}

// TestSolutionsMatchSnippets 是真正对仓库生效的检查，make snippets 和全量测试都跑它。
// 通过时打印关联关系：同一个 snippet 文件引出的题目就是代码上相关联的题，
// 这份清单只从代码推出，不在任何地方手写。
func TestSolutionsMatchSnippets(t *testing.T) {
	canon, errs := loadSnippets(".", structuresDir)
	uses, more := checkSolutions(canon, solutionRoots)
	for _, e := range append(errs, more...) {
		t.Error(e)
	}
	if !t.Failed() {
		fmt.Print(usageReport(canon, uses))
	}
}

// usageReport 按 snippet 文件 → 函数 → 题目分组，函数、类型、方法按在文件里的先后排列。
func usageReport(canon map[string]funcDef, uses map[string][]string) string {
	names := make([]string, 0, len(canon))
	for name := range canon {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		a, b := canon[names[i]], canon[names[j]]
		if a.file != b.file {
			return a.file < b.file
		}
		return a.line < b.line
	})
	var b strings.Builder
	file := ""
	for _, name := range names {
		if def := canon[name]; def.file != file {
			file = def.file
			fmt.Fprintf(&b, "snippets/go/%s\n", file)
		}
		problems := uses[name]
		if len(problems) == 0 {
			fmt.Fprintf(&b, "  %s: （还没有题目引用）\n", name)
			continue
		}
		fmt.Fprintf(&b, "  %s: %s\n", name, strings.Join(problems, ", "))
	}
	return b.String()
}

// loadSnippets 读取 dir 下非测试文件的顶层函数、类型和方法，并检查只 import 了标准库。
// platformDir 里定义过的类型名不算 snippet，传空串表示没有。
func loadSnippets(dir, platformDir string) (map[string]funcDef, []string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, []string{err.Error()}
	}
	platform, errs := platformTypes(platformDir)
	canon := map[string]funcDef{}
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
		for name, def := range declarations(fset, file, src) {
			if platform[name] {
				continue
			}
			canon[name] = def
		}
	}
	return canon, errs
}

// platformTypes 返回 dir 下非测试文件定义的类型，键的写法同 funcDef。
func platformTypes(dir string) (map[string]bool, []string) {
	types := map[string]bool{}
	if dir == "" {
		return types, nil
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, []string{err.Error()}
	}
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
		for name := range declarations(fset, file, src) {
			if strings.HasPrefix(name, "type ") {
				types[name] = true
			}
		}
	}
	return types, errs
}

// checkSolutions 在题解里找和 snippet 同名的顶层函数、类型和方法，
// 返回每个 snippet 被哪些题目目录引用，以及不一致的报错。
func checkSolutions(canon map[string]funcDef, roots []string) (uses map[string][]string, errs []string) {
	uses = map[string][]string{}
	if len(canon) == 0 {
		return uses, nil
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
			for name, got := range declarations(fset, file, src) {
				want, ok := canon[name]
				if !ok {
					continue
				}
				switch {
				case got.body == want.body:
					uses[name] = append(uses[name], filepath.Base(filepath.Dir(path)))
				case got.sig != want.sig:
					// 签名不同时多半是同名的另一个函数，只报第 1 行不一致看不出该怎么改。
					// 类型的 sig 两边都为空，不会走到这里
					errs = append(errs, fmt.Sprintf("%s: %s 的签名与 snippet %s 不同\n"+
						"  snippet: %s\n  题解:    %s\n"+
						"  Go 没有重载，同名函数（同一类型上的同名方法）一律当作 snippet 的副本检查：\n"+
						"  - 是副本：改回 snippet 的签名\n"+
						"  - 是不同的函数：改个名字",
						got.pos, name, want.pos, want.sig, got.sig))
				default:
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
	for _, problems := range uses {
		sort.Strings(problems)
	}
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

// declarations 返回文件里的顶层函数、类型和方法，键的写法见 funcDef。
// 类型别名（type TreeNode = structures.TreeNode）是接线，不是副本，跳过。
func declarations(fset *token.FileSet, file *ast.File, src []byte) map[string]funcDef {
	defs := map[string]funcDef{}
	at := func(pos token.Pos, sig, body string) funcDef {
		p := fset.Position(pos)
		return funcDef{
			file: filepath.Base(p.Filename),
			line: p.Line,
			pos:  fmt.Sprintf("%s:%d", filepath.ToSlash(p.Filename), p.Line),
			sig:  strings.ReplaceAll(sig, "\r\n", "\n"),
			body: strings.ReplaceAll(body, "\r\n", "\n"),
		}
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			name := d.Name.Name
			if d.Recv != nil {
				name = receiverType(d.Recv.List[0].Type) + "." + name
			}
			// 从 fn.Type.Func（func 关键字，方法则在接收者之前）开始截取，跳过 fn.Doc
			start := fset.Position(d.Type.Func).Offset
			sigEnd := fset.Position(d.Type.End()).Offset
			end := fset.Position(d.End()).Offset
			defs[name] = at(d.Pos(), string(src[start:sigEnd]), string(src[start:end]))
		case *ast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			for _, spec := range d.Specs {
				ts := spec.(*ast.TypeSpec)
				if ts.Assign.IsValid() {
					continue
				}
				// 从类型名截到结尾再补上 type：分组写法 type ( … ) 里每个类型单独比较，
				// 和单独一行的 type X … 形式一致。结构体字段上的注释在范围内，要一致
				start := fset.Position(ts.Pos()).Offset
				end := fset.Position(ts.End()).Offset
				defs["type "+ts.Name.Name] = at(ts.Pos(), "", "type "+string(src[start:end]))
			}
		}
	}
	return defs
}

// receiverType 取接收者的类型名：*MaxPQ[T] → MaxPQ。
func receiverType(expr ast.Expr) string {
	for {
		switch e := expr.(type) {
		case *ast.StarExpr:
			expr = e.X
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		case *ast.ParenExpr:
			expr = e.X
		case *ast.Ident:
			return e.Name
		default:
			return fmt.Sprintf("%T", e)
		}
	}
}

// firstDiff 给出第一处不同的行，函数和类型通常不长，这样足够定位。
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
			return fmt.Sprintf("  第 %d 行（从 func / type 关键字算起）\n  snippet: %q\n  题解:    %q", i+1, wl, gl)
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
	canon, errs := loadSnippets(filepath.Join(root, "snippets"), "")
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
	canon, _ := loadSnippets(filepath.Join(root, "snippets"), "")
	_, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 1 || !strings.Contains(errs[0], "Solution.go:3: abs") || !strings.Contains(errs[0], "x <= 0") {
		t.Fatalf("errs=%v", errs)
	}
}

func TestCheckerHintsOnSignatureMismatch(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/abs.go": sampleSnippet,
		"leetcode/0001.A/Solution.go": `package leetcode

func abs(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
`,
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"), "")
	_, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 1 {
		t.Fatalf("errs=%v", errs)
	}
	for _, want := range []string{
		"Solution.go:3: abs 的签名与 snippet",
		"snippet: func abs(x int) int\n",
		"题解:    func abs(x, y int) int\n",
		"改个名字",
	} {
		if !strings.Contains(errs[0], want) {
			t.Fatalf("缺少 %q:\n%s", want, errs[0])
		}
	}
}

// 题解里和 snippet 函数同名的方法（键是 T.abs）不算副本
func TestCheckerSkipsTestsMethodsAndIDECopies(t *testing.T) {
	stale := "package leetcode\n\nfunc abs(x int) int { return x }\n"
	root := writeFiles(t, map[string]string{
		"snippets/abs.go":                                   sampleSnippet,
		"leetcode/0001.A/Solution_test.go":                  stale,
		"leetcode/0001.A/out/production/0001.A/Solution.go": stale,
		"leetcode/0001.A/Solution.go":                       "package leetcode\n\ntype T struct{}\n\nfunc (T) abs(x int) int { return x }\n",
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"), "")
	uses, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(uses)+len(errs) != 0 {
		t.Fatalf("uses=%v errs=%v", uses, errs)
	}
}

func TestUsageReportGroupsByFile(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/letters.go":         "package snippets\n\nfunc countLetters() {}\n\nfunc anagrammatize() {}\n",
		"snippets/math.go":            "package snippets\n\nfunc abs() {}\n",
		"leetcode/0242.B/Solution.go": "package leetcode\n\nfunc countLetters() {}\n",
		"leetcode/0049.A/Solution.go": "package leetcode\n\nfunc countLetters() {}\n",
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"), "")
	uses, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	want := "snippets/go/letters.go\n" +
		"  countLetters: 0049.A, 0242.B\n" +
		"  anagrammatize: （还没有题目引用）\n" +
		"snippets/go/math.go\n" +
		"  abs: （还没有题目引用）\n"
	if got := usageReport(canon, uses); got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestSnippetsOnlyImportStdlib(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/bad.go": "package snippets\n\nimport \"github.com/x/y\"\n\nfunc f() { y.Z() }\n",
		"snippets/ok.go":  "package snippets\n\nimport \"sort\"\n\nfunc g(a []int) { sort.Ints(a) }\n",
		// 测试文件可以用任何依赖，里面的函数也不算 snippet
		"snippets/ok_test.go": "package snippets\n\nimport \"github.com/x/y\"\n\nfunc h() { y.Z() }\n",
	})
	canon, errs := loadSnippets(filepath.Join(root, "snippets"), "")
	if len(errs) != 1 || !strings.Contains(errs[0], "github.com/x/y") {
		t.Fatalf("errs=%v", errs)
	}
	if _, ok := canon["h"]; ok {
		t.Fatal("测试文件里的函数不应算作 snippet")
	}
}

const samplePQ = `package snippets

// PQ 是一个泛型的堆。
type PQ[T any] struct {
	items []T // 字段上的注释在比较范围内
}

// newPQ 返回空堆。
func newPQ[T any]() *PQ[T] {
	return &PQ[T]{}
}

// push 加入一个元素。
func (pq *PQ[T]) push(x T) {
	pq.items = append(pq.items, x)
}
`

// 类型和方法按 type PQ、PQ.push 识别，接收者里的 * 和 [T] 不影响名字
func TestCheckerChecksTypesAndMethods(t *testing.T) {
	copyWith := func(old, new string) string {
		return strings.Replace(strings.Replace(samplePQ, "package snippets", "package leetcode", 1), old, new, 1)
	}
	tests := []struct {
		name     string
		solution string
		wantUses []string
		wantErr  string
	}{
		{"identical copy with own doc", copyWith("// push 加入一个元素。", "// push 这题里入堆。"),
			[]string{"PQ.push", "newPQ", "type PQ"}, ""},
		{"method body drift", copyWith("append(pq.items, x)", "append(pq.items, x, x)"),
			[]string{"newPQ", "type PQ"}, "PQ.push 与"},
		{"type field drift", copyWith("items []T", "items []*T"),
			[]string{"PQ.push", "newPQ"}, "type PQ 与"},
		{"receiver renamed", copyWith("func (pq *PQ[T]) push(x T) {\n\tpq.items = append(pq.items, x)",
			"func (q *PQ[T]) push(x T) {\n\tq.items = append(q.items, x)"),
			[]string{"newPQ", "type PQ"}, "PQ.push 的签名与 snippet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeFiles(t, map[string]string{
				"snippets/pq.go":              samplePQ,
				"leetcode/0215.A/Solution.go": tt.solution,
			})
			canon, errs := loadSnippets(filepath.Join(root, "snippets"), "")
			uses, more := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
			errs = append(errs, more...)
			var used []string
			for name := range uses {
				used = append(used, name)
			}
			sort.Strings(used)
			if strings.Join(used, ",") != strings.Join(tt.wantUses, ",") {
				t.Fatalf("uses = %v, want %v (errs=%v)", used, tt.wantUses, errs)
			}
			switch {
			case tt.wantErr == "" && len(errs) != 0:
				t.Fatalf("errs = %v, want none", errs)
			case tt.wantErr != "" && (len(errs) != 1 || !strings.Contains(errs[0], tt.wantErr)):
				t.Fatalf("errs = %v, want one containing %q", errs, tt.wantErr)
			}
		})
	}
}

// 分组写法 type ( … ) 里的类型和单独一行的写法是同一个副本
func TestCheckerMatchesGroupedTypeDecl(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"snippets/a.go":               "package snippets\n\ntype A struct{ x int }\n",
		"leetcode/0001.A/Solution.go": "package leetcode\n\ntype (\n\tA struct{ x int }\n\tB int\n)\n",
	})
	canon, _ := loadSnippets(filepath.Join(root, "snippets"), "")
	uses, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 0 || len(uses["type A"]) != 1 {
		t.Fatalf("uses=%v errs=%v", uses, errs)
	}
}

// snippet 里和 structures 同名的类型是平台类型的替身，题解用别名接到 structures：
// 两边都不能报不一致，替身也不该出现在引用清单里
func TestCheckerSkipsAliasesAndPlatformTypes(t *testing.T) {
	treeNode := "type TreeNode struct {\n\tVal         int\n\tLeft, Right *TreeNode\n}\n"
	inorder := "func inorder(root *TreeNode, ret *[]int) {}\n"
	root := writeFiles(t, map[string]string{
		"structures/TreeNode.go":      "package structures\n\n" + treeNode,
		"snippets/traversal.go":       "package snippets\n\n" + treeNode + "\n" + inorder,
		"leetcode/0094.A/Solution.go": "package leetcode\n\nimport structures \"x/structures/go\"\n\ntype TreeNode = structures.TreeNode\n\n" + inorder,
	})
	canon, errs := loadSnippets(filepath.Join(root, "snippets"), filepath.Join(root, "structures"))
	if _, ok := canon["type TreeNode"]; ok || len(errs) != 0 {
		t.Fatalf("canon has type TreeNode = %v, errs = %v", ok, errs)
	}
	uses, errs := checkSolutions(canon, []string{filepath.Join(root, "leetcode")})
	if len(errs) != 0 || len(uses) != 1 || len(uses["inorder"]) != 1 {
		t.Fatalf("uses=%v errs=%v", uses, errs)
	}

	// 不传 structures 时替身就是普通 snippet，题解里的别名仍然跳过，不报不一致
	canon, _ = loadSnippets(filepath.Join(root, "snippets"), "")
	if _, errs = checkSolutions(canon, []string{filepath.Join(root, "leetcode")}); len(errs) != 0 {
		t.Fatalf("alias reported: %v", errs)
	}
}
