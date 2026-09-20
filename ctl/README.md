# ctl

题解仓库的命令行工具，从 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 的 ctl 改造而来。所有命令都在 `ctl/` 目录下执行，仓库根的 `Makefile` 只是它的薄封装。

```sh
cd ctl && go run . <command>
```

| 命令 | Makefile 封装 | 作用 |
|:---|:---|:---|
| `new <题号>` | `make new ID=1` | 按题号建题目目录，抓题目描述写进 README，复制各语言骨架 |
| `build readme` | `make readme` | 重新生成仓库根 `README.md` |
| `build readme --anonymous` | `make readme-anon` | 同上，但跳过 `config.toml` 强制匿名，个人数据全为 0 |
| `refresh` | — | 目前等同于 `build readme` |
| `version` | — | 打印版本 |

`new` 支持 `--langs` 只生成部分语言，如 `go run . new 15 --langs go,python`。

## 数据来源

默认走 **leetcode.cn**，理由写在 `util/repo.go` 的 `Site` 常量注释里（一句话：题库是 leetcode.com 的严格超集，而编号题的标题和 slug 仍是英文，目录命名两边通用）。换站点改那一行即可。

- 题目列表：`/api/problems/all/` —— 题号、标题、slug、难度、通过率
- 题目描述：GraphQL `questionData` —— 英文原文和官方中文翻译，两者都是 HTML，由 `html2md.go` 转成 Markdown

两个接口都**不需要登录**。题目列表缓存在 `ctl/.cache/problems.json`，24 小时内不重复请求；断网时用过期缓存兜底。

## 什么时候需要登录

大部分情况都不需要。下面这张表是实测结果：

| 操作 | 需要登录？ | 不登录会怎样 |
|:---|:---|:---|
| `new <免费题>` | 否 | 完全正常，中英文题目描述都能拿到 |
| `new <会员题>` | 否 | 目录、标题、难度都正常，但**题目描述是空的**，要自己登录后复制进去 |
| `build readme` 的题目表格 | 否 | 完全正常 |
| `build readme` 的「个人数据」 | **是** | AC 数全是 0，Perfection Rate 显示 `-` |
| `build readme` 的「已 AC 但未收录」 | **是** | 列表为空 |

会员题（`paid_only`）比较特殊：`new` 会先提示「第 N 题是会员题，题目描述需要登录后自己复制」。而且**光登录还不够，账号得真有 LeetCode 会员**，否则接口照样返回空描述。

## 可选：配置登录

不配也能用，只是上表里标了「是」的那两处没内容。

反过来，**配了 Cookie 之后要生成一份不含个人数据的 README**（比如给 `template` 分支用），加 `--anonymous`：

```sh
make readme-anon
```

它直接跳过 `config.toml`，不读也不发 Cookie。

想要「个人数据」有内容：

1. 在 `ctl/` 下新建 `config.toml`（已在 `.gitignore` 里，不会被提交）
2. 填入下面的内容，把占位符换成自己的

```toml
Username="你的用户名"
Cookie="csrftoken=XXXXXXXX; LEETCODE_SESSION=YYYYYYYY;"
CSRFtoken="XXXXXXXX"
```

3. Cookie 从**已登录的 leetcode.cn** 页面取（浏览器开发者工具 → Application → Cookies），需要 `csrftoken` 和 `LEETCODE_SESSION` 两项

Cookie 不跨站点通用：leetcode.com 的 session 在 leetcode.cn 上无效，反之亦然。`Password` 字段保留是为了兼容原版结构，实际请求不用它，留空即可。

## README 模板

`template/template.markdown` 是仓库根 README 的模板，`build readme` 把里面的占位符替换掉：

| 占位符 | 内容 |
|:---|:---|
| `{{.PersonalData}}` | 个人 AC 数据整节，**含 `## 个人数据` 标题**。未登录时渲染成空串，整节消失 |
| `{{.LanguageTable}}` | 各语言题解数量统计 |
| `{{.TotalNum}}` | 一句话进度说明 |
| `{{.SolvedTable}}` | 已写题解的题目表格 |
| `{{.AvailableTable}}` | LeetCode 全部题目的表格（4000+ 行，默认没用上） |
| `{{.OptimizingTable}}` | 已 AC 但仓库里还没收录题解的题目 |

改 README 的版式直接改这个模板，不要改生成出来的 `README.md`——它每次 `make readme` 都会被覆盖。

注意 `{{.PersonalData}}` 和别的占位符不一样：**它连 `## 个人数据` 这个标题一起渲染**。标题放在 Go 代码里而不是模板里，是为了让它能跟着内容一起消失——未登录时接口返回的 AC 计数全是 0，与其展示一张全 0 的表（看起来像"一道题都没做"而不是"没拿到数据"），不如整节不渲染；模板里要是留着标题，就会变成一个后面什么都没有的空标题。

也因此，**模板里不要在注释或正文里写出 `{{.` 开头的占位符字面量** —— 替换是全文本匹配的，会把注释里那个也一起换掉。

## 多语言探测

`util/util.go` 里的 `Languages` 决定 README 的 `Solution` 列怎么渲染：

```go
{Name: "Go",     Ext: ".go",   TestSuffix: "_test.go", Entry: "Solution.go"},
{Name: "Python", Ext: ".py",   TestSuffix: "_test.py", Entry: "solution.py"},
{Name: "Java",   Ext: ".java", TestSuffix: "Test.java", Entry: "Solution.java"},
```

`Entry` 是入口文件名：一个目录里同语言有多个文件时（比如题解之外还有个 `SegmentTree.go`）优先链到它，避免 README 链到辅助文件上。

加一门语言（C++、Rust）在这里加一行，**ctl 这一侧**就完事了——`Solution` 列会自动多出该语言的链接，语言统计表自动多一行，`ctl new --langs` 也会认。但 ctl 之外还有五步：骨架文件、测试脚本、`structures/<lang>/` 的共享数据结构、把共享结构接进该语言的编译、共享结构自己的测试。完整清单见 [AGENTS.md 的「加一门新语言」](../AGENTS.md#加一门新语言)。

题号 `0000` 保留给 `leetcode/0000.Template` 骨架目录，扫描时跳过，不计入任何统计。

## 已停用的功能

上游 ctl 还能渲染 Hugo 站点（website/）和导出 PDF 电子书。本仓库没有引入 `website/`，这些代码用 `//go:build ignore` 停用，源码保留着以便将来恢复：

| 文件 | 原本的作用 |
|:---|:---|
| `render_website.go` | 渲染站点第二章「算法专题」、书籍目录、拷题解进第四章 |
| `label.go` | 站点各章节的顺序表和标题映射 |
| `pdf.go` | 把整本书合并成单个 Markdown 以导出 PDF |
| `rangking.go` | 抓个人主页排名（上游的解析逻辑已被 LeetCode 改版弄失效） |
| `template_render.go` | 旧的 README 渲染实现，已被 `render.go` 取代 |
| `models/tagproblem.go` | 「算法专题」的 GraphQL 模型 |
| `util/website.go` | 站点目录相关的工具函数 |

每个文件开头都写了恢复步骤。注意恢复时还要补回 `ctl/meta/`（复杂度元数据）和 `ctl/template/` 下的各专题模板——那些是原作者针对他自己 798 道题解的数据，没有一起搬过来。
