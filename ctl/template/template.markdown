# algorithm-notes

我的算法题解仓库，支持 Go / Python / Java 多语言题解。

## 关于这个仓库

仓库结构和 `ctl` 工具链以 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 为模板构建，在此感谢原作者 [@halfrost](https://github.com/halfrost) 开源的这套题解组织方式和自动化工具。本仓库在其基础上做了三件事：改造 `ctl` 让 README 支持多语言题解链接、去掉 Hugo 站点与 PDF 电子书相关的部分、补上 Python 和 Java 的测试链路。

代码归属：

- **题解代码**（`leetcode/` 下各题的 `Solution.*`）由我本人编写。卡住的题会参考 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 的实现，思路借鉴自那里。
- **仓库初始化代码**（`ctl/` 工具链改造、`Makefile`、`conftest.py`、各测试脚本、CI 配置和本 README 模板）由 [Claude Code](https://claude.com/claude-code) 生成。

## 快速开始

```sh
git clone git@github.com:nickmyb/algorithm-notes.git
cd algorithm-notes
git checkout init     # 干净的起点：只有工具链和题解骨架
make init             # 检查工具链、装依赖、跑通测试、生成 README
make new ID=1         # 开始第一题
```

`make init` 需要 Go（`ctl` 是 Go 写的，必需）；Python 和 Java 缺了只会警告并跳过对应的测试链路，不写那门语言就不受影响。

## 环境要求

**项目只在下面这三个版本上测试过**，其他版本不保证可用：

| 语言 | 版本 | 版本来源 | 说明 |
|:---|:---|:---|:---|
| Go | **1.26.4** | `go.mod` 的 `go` 指令 | Go 没有 LTS，官方只维护最近两个大版本。CI 用 `go-version-file: go.mod`，和本地同一个来源 |
| Python | **3.12** | `scripts/init.sh` 的版本检查 | Python 没有 LTS，每个小版本一律 5 年支持期。低于 3.12 时 `make init` 会直接报错退出 |
| Java | **17**（LTS） | `javatest.sh` 的 `--release` | Java 的 LTS 每两年一个：8 / 11 / 17 / 21 / 25 |

三者都只有一个声明来源，CI 跟着本地走，不会各说各的。

Java 的 `--release 17` 同时约束语言特性和可用 API：即使本地装的是 JDK 21，写了 record 或 switch 模式匹配也会在本地就编译失败，而不是推上去才被 CI 拦下。要换目标版本改 `javatest.sh` 里的 `JAVA_RELEASE`，或临时 `JAVA_RELEASE=21 make test-java`。

## 目录结构

```
algorithm-notes/
├── leetcode/                    # 题解，一题一目录
│   ├── 0000.Template/           # 新开一题的复制源，不计入统计
│   └── 0001.Two-Sum/
│       ├── README.md            # 题目、题目大意、解题思路、复杂度
│       ├── Solution.go          # Go 题解
│       ├── Solution_test.go     # Go 表驱动测试
│       ├── solution.py          # Python 题解
│       ├── solution_test.py     # Python 测试
│       ├── Solution.java        # Java 题解
│       └── SolutionTest.java    # Java 测试
├── ctl/                         # 命令行工具：建题目录、生成 README
├── structures/                  # 各语言共用的数据结构（TreeNode、ListNode 等）
│   ├── go/                      # package structures，带 halfrost 原有的 9 个结构
│   ├── java/                    # TreeNode/ListNode + 辅助，test/ 下是它们的测试
│   └── python/                  # tree_node.py / list_node.py + 各自的测试
├── conftest.py                  # pytest 的 solution fixture，见下文
├── pytest.ini                   # pytest 配置
├── gotest.sh                    # Go 测试 + 覆盖率
├── javatest.sh                  # Java 逐目录编译并测试
├── scripts/init.sh              # make init 的实现
├── AGENTS.md                    # 给 AI 编码助手看的项目说明
├── CLAUDE.md                    # 导入 AGENTS.md，Claude Code 和 ChatGPT 共用一份
└── Makefile                     # 所有日常命令的入口
```

题目目录名是 `%04d.英文标题`：标题里的空格转 `-`，`' % ( ) , ? / : " !` 去掉，连续的 `-` 收敛成一个。这个名字由 `make new` 自动生成，不用手写。

以后要收 LCP、剑指 Offer 这类非编号题，会新开 `lcp/` 这样的同级目录，因为 `%04d` 表达不了 `LCP 82` 这种题号。

## IDE 配置

测试脚本内部会切到仓库根，所以 Run Configuration 的 working directory 设成哪里都不影响结果。

### GoLand —— 打开仓库根

| 配置项 | 值 |
|:---|:---|
| GOROOT | 指向 `go.mod` 要求的版本（当前 1.26.4） |
| Go Modules | 自动从 `go.mod` 识别，不用动 |

`ctl/` 下 7 个带 `//go:build ignore` 的文件会被标成排除在构建外、灰掉，这是故意的（见 [ctl/README.md](./ctl/README.md) 的「已停用的功能」），不要去"修"。

### PyCharm —— 打开仓库根

| 配置项 | 值 |
|:---|:---|
| Python Interpreter | `<仓库根>/.venv/bin/python`（`make init` 建的） |
| **`structures/python`** | 目录树右键 → **Mark Directory as → Sources Root** |
| 测试框架 | pytest，自动读 `pytest.ini` |

标记 Sources Root 那步是关键：运行时是 `conftest.py` 把这个目录加进 `sys.path` 的，属于运行期行为，IDE 的静态分析看不到，不标记的话 `from tree_node import TreeNode` 会一直标红。

### IntelliJ IDEA

Java 侧的约束：每道题都是 default package 里的 `class Solution`，**同一个编译范围里只能有一个**，否则撞成一片 `Duplicate class Solution`（`javac` 同理，所以 `javatest.sh` 才逐目录编译）。

两种配法，都是让同一时刻只有一个题目目录进入 Java 的编译范围。Project SDK / Module SDK 都选 **17**（和 `javatest.sh` 的 `JAVA_RELEASE` 一致）。

**推荐：打开仓库根 + 标记 Sources Root**

1. 打开仓库根，设好 Project SDK 17
2. 目录树上把**当前在写的那道题的目录**和 `structures/java` 都标记为 **Mark Directory as → Sources Root**
3. 换题时取消上一题的标记，标到新的题目目录上

一个窗口搞定，换题只要改两次标记，不用建模块。代价是同一时刻只有一道题被索引。

推荐它还有一个更实际的理由：**IDEA 的编译输出会落在仓库根的 `out/`**，而 `go list ./leetcode/...` 和 pytest 的 `testpaths` 都够不到那里，不会污染测试。下面那种配法就没这个好处。

**备选：按题打开 + 模块依赖**

1. 直接把题目目录作为项目打开，如 `leetcode/0094.Binary-Tree-Inorder-Traversal`
2. `Project Structure` → `Modules` → `+` 添加 `structures/java`，命名为 `structures`
3. 选中题目模块 → **`Dependencies`** 标签 → `+` → **Module Dependency** → 选 `structures`，Scope 为 `Compile`

模块关系明确，但换题要重开窗口重配一次。

> **这种配法有个坑**：IDEA 会把编译输出建在**题目目录里面**（`leetcode/0094.xxx/out/`），而且把所有非 `.java` 文件**当资源复制进去**——`solution.py`、`solution_test.py`、`Solution.go`、`Solution_test.go` 全都有一份副本。这些副本会随着你改代码而过时，而**测试工具不看 `.gitignore`**：
>
> - `go list ./leetcode/...` 把 `out/production/<题目>/` 当成一个真实的包编译测试
> - pytest 把 `out/` 里那份旧 `solution_test.py` 也收集了
>
> 旧题解配旧测试永远是绿的，比测试失败更危险。仓库已经挡住了这一层（`gotest.sh` 过滤 `/out/`、`pytest.ini` 的 `norecursedirs` 含 `out`），所以现在用是安全的；但如果你想从根上避免，在 `Project Structure` → `Project` → **Compiler output** 里把输出路径改到题目目录之外。

配好之后 `TreeNode` / `TreeNodes` 在题解里能正常补全和跳转。嫌配置麻烦的话，IDE 只当编辑器用，验证一律走 `make test-java ID=94`。

## 命令

```sh
make help                        # 列出所有命令
make new ID=1                    # 新开一题，自动抓题目描述填进 README
make new ID=15 LANGS=go,python   # 只生成指定语言的骨架
make test                        # 跑三门语言的全部测试
make test ID=94                  # 只测一道题（那题没写的语言会安静跳过）
make test-go / test-python / test-java   # 同样支持 ID=94
make readme                      # 重新生成本 README 的题目表格
make fmt / vet / tidy            # Go 格式化、静态检查、依赖整理
make clean                       # 清构建产物，不碰题解
```

`make` 只是 `ctl` 的薄封装，`make new ID=1` 等价于 `cd ctl && go run . new 1`。细节见 [ctl/README.md](./ctl/README.md)。

`LANGS` 可以只指定一门或几门语言，没写的语言就不生成骨架文件：

```sh
make new ID=15 LANGS=go          # 只有 Solution.go / Solution_test.go
make new ID=20 LANGS=python      # 只有 solution.py / solution_test.py
make new ID=21 LANGS=go,java     # Go 和 Java
```

下面的题目表格按目录里**实际存在**的文件渲染，写了几门语言就显示几个链接，事后补另一门语言只要把文件放进去再 `make readme`。

## 写一道题

```sh
make new ID=1     # 建出 leetcode/0001.Two-Sum/，题目描述已经填好
# 写题解，把测试骨架里的 t.Skip / pytest.skip 删掉
make test ID=1    # 只测这一道，反复改的时候快得多
make test         # 提交前全量跑一遍
make readme       # 刷新下面的题目表格
```

`make test ID=1` 会按题号找到目录（`94` → `leetcode/0094.*`），那道题没写的语言安静跳过 —— 比如只写了 Java 的题，Go 和 Python 不会因为"没有测试"而报错。题号不存在时直接报错退出，不会静默退化成全量。

题解 README 的四个小节：

| 小节 | 内容 |
|:---|:---|
| `## 题目` | `make new` 自动填：英文原文 + 折叠的官方中文翻译 |
| `## 题目大意` | **自己写**，一两句话复述题意。这一步的价值就在于自己写，所以没有自动填 |
| `## 解题思路` | 思路推导，为什么这么做 |
| `## 复杂度` | 时间 / 空间复杂度 |

前三节沿用 halfrost 原仓库的格式，`## 复杂度` 是本仓库加的——原仓库把复杂度记在 `ctl/meta/` 里用于渲染站点，那套机制没有搬过来，所以在 README 里留个固定位置。

> **统计的口径**：ctl 判断一道题"有没有题解"，看的是目录里**有没有对应语言的文件，不看内容**。所以 `make new` 一建完目录，那道题就计入下面的统计了，哪怕里面还是骨架。
>
> 由此有个实际建议：**只打算写一门语言时，`make new` 记得带 `LANGS=`**。三门全生成却只写一门的话，下面的表格会显示三个链接、语言统计三门都算 1，而点进去另外两个是空骨架。事后想补语言，把文件放进目录再 `make readme` 即可，不用重建。

## 多语言约定

### 题解本体可以原样复制到 LeetCode

题解文件的结构固定为两段：

```
[本地接线区]   ← LeetCode 上不需要、本地才要的几行
[题解本体]     ← 逐字等于 LeetCode 提交框里的内容
```

从 LeetCode 复制进来只需**补接线**，从仓库复制出去只需**删接线**，本体一个字都不改。

**什么算接线**：凡是 LeetCode 提交框里不存在的代码行——Go 的 `package leetcode`、三门语言的 import、Go 的类型别名。所以文件里那行 `===== 以下是题解本体 =====` 的标记在**所有 import 之后**，不在文件最开头。

| 语言 | 本地接线区 | 题解本体 |
|:---|:---|:---|
| Go | `package leetcode` + `import structures` + `type TreeNode = structures.TreeNode` | `func inorderTraversal(root *TreeNode) []int {…}` |
| Java | `import java.util.List;` 等（LeetCode 预置了 `java.util.*`） | `class Solution {…}` |
| Python | `from tree_node import TreeNode`（LeetCode 预置了节点类型） | `class Solution:` + 带 `self` 的方法 |

由此有三条细节是刻意的，不要"顺手改掉"：

- **Java 是 `class Solution`，不带 `public`** —— LeetCode 模板里带 `public` 的是 `TreeNode`，`Solution` 没有。文件名仍是 `Solution.java`，非 public 类不要求匹配文件名。
- **Python 是 `class Solution:` 里带 `self` 的方法**，不是顶层函数。Python 靠缩进表达结构，改成顶层函数意味着每次复制都要删类声明、去 `self`、整段反缩进。测试里写 `solution.Solution().方法名(...)`。
- **Go 用类型别名** `type TreeNode = structures.TreeNode`（是别名 `=`，不是新类型），这样本体里能写裸的 `*TreeNode`。这是上游 halfrost 的做法。

### 为什么三门语言的文件名不完全一致

| 语言 | 题解 | 测试 | 原因 |
|:---|:---|:---|:---|
| Go | `Solution.go` | `Solution_test.go` | 包名来自文件里的 `package` 声明，文件名对编译器没有意义，全仓库同名零代价 |
| Python | `solution.py` | `solution_test.py` | 模块名就是文件名，同名会在 `sys.modules` 里互相顶替，靠 `conftest.py` 的 fixture 化解 |
| Java | `Solution.java` | `SolutionTest.java` | `public class` 名必须和文件名一致；同名类不能一次编译，靠 `javatest.sh` 逐目录编译化解 |

小写的 `solution.py` 和大写的 `Solution.java` 是各语言自己的命名惯例，不是随手写的。

### Python 的 solution fixture

几百道题都叫 `solution_test.py`，pytest 默认的 prepend 模式会直接报 `import file mismatch` 中断收集；换成 importlib 模式虽然能收集，但测试里 `from solution import ...` 又会 `ModuleNotFoundError`。

仓库根的 `conftest.py` 提供了一个 `solution` fixture，按**文件路径**加载同目录的 `solution.py`，模块名带上题目目录名保证唯一。所以测试这样写：

```python
def test_two_sum(solution):
    assert solution.twoSum([2, 7, 11, 15], 9) == [0, 1]
```

### 同一语言的多种解法

写在**同一个题解文件里的多个函数**。入口函数保持 LeetCode 给的签名名，其余解法加后缀说明算法：

```
twoSum              # 入口，哈希表解法
twoSumBruteForce
twoSumTwoPointers
```

测试不为每种解法单写一份，而是**同一组用例跑所有实现**，顺带保证它们结果一致（Go 用 `t.Run` 子测试，Python 用 `parametrize`）。

### 共享数据结构

`TreeNode`、`ListNode` 这些在 LeetCode 上是平台提供的，本仓库由 `structures/` 扮演这个角色，**不要在题目目录里重复定义**。

| 语言 | 题解里怎么用 |
|:---|:---|
| Go | `import structures "github.com/nickmyb/algorithm-notes/structures/go"` |
| Java | 直接用 `TreeNode` / `TreeNodes`，`javatest.sh` 编译时会带上 |
| Python | `from tree_node import TreeNode, build_tree` |

目前三门语言都有 `TreeNode` 和 `ListNode`（含建树、建链表、遍历、成环等测试辅助）。Go 那边额外带着 halfrost 的 `Stack` / `Queue` / `Heap` / `PriorityQueue`，那是为了补 Go 标准库的缺口——Java 有 `ArrayDeque` / `PriorityQueue`，Python 有 `deque` / `heapq`，不需要也不应该手写一份。`NestedInteger` / `Interval` / `Point` 用到再补。

这些辅助是所有树/链表题的地基，写错了会让题解测试给出**假结果**，所以三门语言各自都有测试，跟着 `make test` 全量跑（单题模式不跑）。

节点类型的定义和 LeetCode 给的一字不差（已用官方 `codeSnippets` 逐字比对验证），所以题解可以在本地和提交框之间原样复制。建树、建链表这类**测试辅助**单独放（Java 在 `TreeNodes` / `ListNodes`，Python 在 `tree_node.py` / `list_node.py` 的模块级函数），不挂到节点类上，免得破坏这个性质：

```java
// 题目给的输入是 [1,null,2,3]，测试里就这么写
TreeNode root = TreeNodes.build(1, null, 2, 3);
```

Go 的目录叫 `go` 而包名是 `structures`（`go` 是关键字不能当包名），import 路径末段和包名不一致是刻意的，为了和 `java/`、`python/` 对称；实测 `build` / `vet` / `gofmt` / `test` 都正常。

**新增一个共享结构**：往对应的 `structures/<lang>/` 里加文件就生效，不用改测试脚本。三点注意：

- **Java** 的每个类都会进入**每一个**题目目录的编译单元，类名要足够独特，和某道题里的辅助类重名会直接编译失败
- **Python** 的文件名就是模块名，别起 `queue.py`、`heapq.py` 这种和标准库重名的
- 节点类型照搬 LeetCode 的定义不要改，测试辅助另外放（`TreeNode` / `TreeNodes` 分开就是为此）

只在真用到时加。halfrost 的 Go 版有 13 个树相关函数，Java / Python 侧目前只移植了建树和层序展开两个。

### 什么时候才另开文件

只有需要独立的辅助数据结构时，比如第 307 题的线段树：

```
0307.Range-Sum-Query-Mutable/
├── SegmentTree.go     ← 辅助结构，按结构命名
├── Solution.go        ← 题解入口，调用 SegmentTree
└── solution.py        ← Python 没有这个约束，类直接写在同一个文件里
```

文件名按**结构**命名，不要叫 `Solution2.go`——那会读成"第二种解法"，和上面的约定冲突。Java 的辅助类声明成 `public` 就必须单独一个文件，不写 `public` 可以塞进 `Solution.java`；Go 纯属风格选择，不拆也能编译。

下面表格里的 `Solution` 列按目录里**实际存在**的文件渲染，写了几门语言就显示几个链接。

## 分支与标签

| 引用 | 内容 |
|:---|:---|
| `main` | 我的题解 |
| `template` | 干净的骨架 + 工具链，不含任何题解，随工具链演进往前走 |
| `init` 标签 | 始终指向 `template` 的最新提交，新用户从这里起步 |

工具链更新后用 `git tag -f init template && git push -f origin init` 把标签移到最新。

## 数据来源

`ctl` 默认从 **leetcode.cn** 取数据，理由写在 `ctl/util/repo.go` 的注释里：题库是 leetcode.com 的严格超集（.com 的 4055 道编号题一道不缺，另有 388 道 LCR / 面试题 / LCP / LCS），而编号题的标题和 slug 仍是英文，目录命名在两个站点通用。换站点改那一个常量即可。

**日常用不着登录**：`make new` 建目录、抓题目描述，`make readme` 生成题目表格，都是匿名请求。

需要登录的只有两处，都在 README 的统计部分：

| 需要登录 | 不登录的表现 |
|:---|:---|
| 下面的「个人数据」表格 | AC 数全是 0，Perfection Rate 显示 `-` |
| 「已 AC 但还没写题解」的列表 | 为空（模板里默认没启用这块） |

另有一个例外：**会员题**未登录时拿不到题目描述，`make new` 会建好目录但「题目」小节是空的，需要自己登录后复制；而且账号得真有 LeetCode 会员，光登录不够。

配置方法见 [ctl/README.md](./ctl/README.md)，配置文件 `ctl/config.toml` 已在 `.gitignore` 里。

---
{{.PersonalData}}
## 题解统计

{{.LanguageTable}}

## 题目列表

{{.TotalNum}}

{{.SolvedTable}}

---

## License

题解代码以 [MIT](./LICENSE) 授权。仓库结构与 `ctl` 工具链改编自 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go)，原项目同为 MIT 授权。
