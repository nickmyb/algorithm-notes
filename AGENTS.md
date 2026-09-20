# AGENTS.md

给 AI 编码助手（Claude Code / ChatGPT / Codex 等）看的项目说明。人类看 [README.md](./README.md)。

`CLAUDE.md` 通过 `@AGENTS.md` 导入本文件，两边共用同一份说明，改这里就够了。

## 项目概述

个人算法题解仓库，一道题一个目录，同时支持 Go / Python / Java 三门语言的题解。

仓库结构和 `ctl` 工具链以 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 为模板，改造点见下文「相对上游的修改」。题解代码由仓库作者本人编写，初始化代码由 Claude Code 生成。

## 常用命令

```sh
make help                        # 列出所有命令
make init                        # 初始化：检查工具链、装依赖、跑测试、生成 README
make new ID=1                    # 新开一题（三门语言）
make new ID=15 LANGS=go          # 只要某一门或某几门语言
make test                        # 三门语言全跑
make test-go / test-python / test-java
make readme                      # 重新生成根 README.md
make fmt / vet / tidy / clean
```

`make` 只是 `ctl` 的薄封装，`make new ID=1` 等价于 `cd ctl && go run . new 1`。**ctl 的所有命令都必须在 `ctl/` 目录下执行**，它用的是相对路径（`../leetcode/`、`./template/`）。

## 环境版本

只在这三个版本上测试过，各有唯一声明来源，改版本要改对应的那处：

| 语言 | 版本 | 声明位置 |
|:---|:---|:---|
| Go | 1.26.4 | `go.mod` 的 `go` 指令 |
| Python | 3.12 | `scripts/init.sh` 的 `PYTHON_MIN_MAJOR/MINOR` |
| Java | 17 | `javatest.sh` 的 `JAVA_RELEASE`（传给 `javac --release`） |

CI（`.github/workflows/test.yml`）跟着这三处走，不要在 CI 里另写一套版本。

## 目录结构

```
leetcode/<NNNN>.<英文标题>/    # 一题一目录，NNNN 是四位题号
├── README.md                 # 题目 / 题目大意 / 解题思路 / 复杂度
├── Solution.go  Solution_test.go
├── solution.py  solution_test.py
└── Solution.java SolutionTest.java
ctl/                          # 命令行工具
structures/                   # Go 题解共用数据结构
conftest.py pytest.ini        # Python 测试基建
gotest.sh javatest.sh         # Go / Java 测试脚本
scripts/init.sh Makefile
```

## 不要破坏的约束

下面每一条都是实测得出的，看起来像"可以简化的怪写法"，实际不能改：

### 1. `README.md` 是生成产物，不要直接编辑

根目录的 `README.md` 每次 `make readme` 都会被整体覆盖。要改版式或正文，改 **`ctl/template/template.markdown`**，里面的 `{{.Xxx}}` 占位符由 `ctl/render.go` 替换。

### 2. Python 题解不能用普通 import

几百道题都叫 `solution.py` / `solution_test.py`。实测：

- pytest 默认的 prepend 模式 → `import file mismatch`，**整个收集阶段中断**
- 换 `--import-mode=importlib` → 能收集，但测试里 `from solution import ...` 变成 `ModuleNotFoundError`（importlib 模式不把题目目录加进 `sys.path`）

所以根目录 `conftest.py` 提供 `solution` fixture，按**文件路径**加载同目录的 `solution.py`，模块名带上题目目录名保证唯一。测试必须写成：

```python
def test_two_sum(solution):
    assert solution.twoSum([2, 7, 11, 15], 9) == [0, 1]
```

不要"顺手改成" `from solution import twoSum`，那样一定会坏。

### 3. Java 必须逐目录编译

每道题的题解类都叫 `Solution`（default package）。`javac leetcode/*/*.java` 一次性编译会报 `duplicate class: Solution`。`javatest.sh` 按目录分别编译到各自的 `out/java/<题目>/`，不要合并成一次编译。

Java 测试没有引入 JUnit，就是一个 `main` + 断言抛 `AssertionError`。

### 4. 题号可能不是数字

数据源是 **leetcode.cn**（理由在 `ctl/util/repo.go` 的 `Site` 注释里）。它的 `frontend_question_id` 是 JSON **字符串**，而且混着 `"LCP 82"`、`"面试题 17.14"` 这类非数字题号（共 388 道）。

所以有 `models.QuestionID` 自定义 `UnmarshalJSON`，数字和字符串都吃，非数字的 `Num` 为 0。**不要把它改回 `int32`**，那会让整个 `json.Unmarshal` 在第一条 LCP 记录上失败。

非数字题号的题会被跳过（`%04d.标题` 表达不了），将来收这类题要新开 `lcp/` 之类的同级目录。

### 5. 中文强调后面的空格是必需的

`html2md.go` 的 `padCJKEmphasis` 会在以中文标点收尾的强调后补一个空格。这不是排版洁癖：CommonMark 规定闭合标记前是标点、后面又不是空白或标点时不算"右侧贴合"，`**进阶：**你可以…` 会整段渲染成**字面星号**。实测验证过。

### 6. `//go:build ignore` 的文件是有意停用的

`ctl/` 下这些文件不参与编译，**不要"修复"或删除**：

`render_website.go`、`label.go`、`pdf.go`、`rangking.go`、`template_render.go`、`models/tagproblem.go`、`util/website.go`

它们是上游用于渲染 Hugo 站点和导出 PDF 电子书的代码，本仓库没有引入 `website/`，保留源码是为了将来可能恢复。每个文件开头都写了恢复步骤。

### 7. `leetcode/0000.Template` 是骨架，不是题解

题号 `0` 被 `util.LoadSolutions` 跳过，不计入任何统计。`ctl new` 从这里复制各语言文件。它的测试里有 `t.Skip` / `pytest.skip`，**保持跳过状态**。

### 8. 题解 README 的「题目大意」故意留空

`ctl new` 会自动填「题目」（英文原文 + 折叠的官方中文翻译），但**不填「题目大意」**。那是留给作者用一两句话自己复述题意的地方，自动填就失去意义了。不要"补全"它。

上游 halfrost 的三个小节是「题目 / 题目大意 / 解题思路」，「复杂度」是本仓库加的。

## 代码约定

### 多语言文件命名

| 语言 | 题解 | 测试 | 为什么这么命名 |
|:---|:---|:---|:---|
| Go | `Solution.go` | `Solution_test.go` | 包名来自 `package` 声明，文件名对编译器无意义 |
| Python | `solution.py` | `solution_test.py` | 小写是 PEP8 模块命名惯例 |
| Java | `Solution.java` | `SolutionTest.java` | `public class` 名必须和文件名一致 |

大小写不一致是各语言自己的惯例，不是笔误，不要"统一"。

### 同一语言的多种解法

写在**同一个题解文件里的多个函数**：入口函数保持 LeetCode 给的签名名，其余加后缀说明算法（`twoSum` / `twoSumBruteForce` / `twoSumTwoPointers`）。

测试用**同一组用例跑所有实现**（Go 用 `t.Run` 子测试，Python 用 `parametrize`），不要每种解法单写一份测试。

### 什么时候另开文件

只有需要独立辅助数据结构时（如第 307 题的线段树），文件名按**结构**命名：`SegmentTree.go`，不要叫 `Solution2.go`。

`util.Languages` 的 `Entry` 字段就是为此存在的：一个目录里同语言有多个文件时，README 优先链到 `Solution.go` 而不是字典序第一的 `SegmentTree.go`。

### 加一门新语言

改 `ctl/util/util.go` 的 `Languages` 表加一行，再写一个对应的测试脚本（照 `javatest.sh` 抄）、在 `Makefile` 和 CI 里各加一项。README 渲染逻辑不用动。

## 相对上游 halfrost/LeetCode-Go 的修改

**功能改造**

- README 的 `Solution` 列按目录里实际存在的文件渲染多语言链接（上游硬编码 `[Go](...)`），并新增语言统计表
- 新增 `ctl new <题号>`：查标题建目录、抓题目描述填 README、按 `--langs` 复制骨架
- 数据源从 leetcode.com 换到 leetcode.cn，新增 `models.QuestionID` 兼容两种题号形态
- 新增 `html2md.go` 把题目描述 HTML 转 Markdown（带单测）
- 新增 Python / Java 测试链路和对应 CI job
- 删掉 `template/`、`website/`、`topic/`、`note/`；站点与 PDF 代码用 `//go:build ignore` 停用
- 合并成单个 Go module（上游 `structures/`、`ctl/util/`、`ctl/models/` 各有 `go.mod` 靠 `replace` 串联）

**修掉的上游缺陷**

- `util.loadFile` 用 `f.Name()[4] == '.'` 判断目录名，遇到短于 5 字符的文件名越界 panic → 改用正则
- 百分比直接相除，分母为 0 时算出 `NaN%` 写进 README → `percent()` 返回 `-`
- `util.WriteFile` 缺 `O_TRUNC`，写入比原文件短的内容会残留旧尾部
- `config.String()` 把明文密码格式化进字符串 → 改为 `<redacted>`
- `getConfig()` 在 `config.toml` 缺失时 `log.Panic` → 改为可选，未登录也能生成 README

## 分支与标签

| 引用 | 内容 |
|:---|:---|
| `main` | 作者的题解 |
| `template` | 干净的骨架 + 工具链，不含任何题解 |
| `init` 标签 | 始终指向 `template` 的最新提交 |

工具链改动应同时反映到 `template` 分支，并 `git tag -f init template` 把标签移到最新。

## 凭据

`ctl/config.toml` 存 LeetCode 的 Cookie，已在 `.gitignore` 里（`config.toml` 不带斜杠，任意层级生效）。**不要提交它，不要把 Cookie 写进任何其他文件或日志。**

没有这个文件也能正常工作，只是 README 的「个人数据」表格会全是 0。
