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
make test ID=94                  # 只测一道题，那题没写的语言安静跳过
make test-go / test-python / test-java   # 同样支持 ID=94
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
structures/                   # 各语言共用的数据结构，按语言分子目录
├── go/                       # package structures，import 路径 .../structures/go
├── java/                     # TreeNode.java 等，javac 编译每个题目目录时一并带上
└── python/                   # tree_node.py / list_node.py，conftest.py 把它加进 sys.path
conftest.py pytest.ini        # Python 测试基建
gotest.sh pytest.sh javatest.sh   # 三门语言的测试脚本，都接受可选的目录参数
scripts/init.sh               # make init 的实现
scripts/problem-dir.sh        # 题号 -> 目录，供 Makefile 的 ID= 用
Makefile
```

### 单题测试

`make test ID=94` 只测一道题。链路是 `Makefile` 用 `scripts/problem-dir.sh` 把题号解析成目录，再把目录作为参数传给三个测试脚本。

三个脚本都必须把「这道题没写这门语言」当成**跳过而不是失败**，改动它们时不要破坏这一点：

- `gotest.sh`：`compgen -G "$dir/*.go"` 判断有没有 Go 文件，没有就跳过——否则 `go test` 报 `no Go files` 让整条命令变红
- `pytest.sh`：pytest 在「一个测试都没收集到」时返回 **5**，单题模式下要把它当成 0。这正是 `pytest.sh` 存在的唯一理由，不要把它简化成直接调 pytest
- `javatest.sh`：本来就会跳过没有 `.java` 的目录，另外在 `total == 0` 时提前返回 0

题号解析失败时 Makefile 用 `$(error ...)` 直接中止，**不能**静默退化成全量测试——那会让你以为测了单题，实际跑了整个仓库。

## 共享数据结构

`TreeNode`、`ListNode` 这些在 LeetCode 上是平台提供的，本仓库由 `structures/` 扮演这个角色，**不要在题目目录里重复定义**（Java 会直接撞类名编译失败）。

| 语言 | 位置 | 题解里怎么用 | 接入机制 |
|:---|:---|:---|:---|
| Go | `structures/go/` | `import structures "github.com/nickmyb/algorithm-notes/structures/go"` | 普通 Go 包 |
| Java | `structures/java/` | 直接用 `TreeNode` / `TreeNodes` | `javatest.sh` 把 `structures/java/*.java` 加进每次 `javac` 的源文件列表 |
| Python | `structures/python/` | `from tree_node import TreeNode, build_tree` | `conftest.py` 把该目录加进 `sys.path` |

### 现在有哪些，以及为什么只有这些

Go 那边带着 halfrost 的 9 个文件，但它们不是一类东西，**不要照数量去补齐 Java/Python**：

| 类别 | Go 里的 | Java / Python 要不要 |
|:---|:---|:---|
| LeetCode 题目里真正出现的节点类型 | `TreeNode`、`ListNode`、`NestedInteger`、`Interval`、`Point` | 要。目前移植了 `TreeNode` 和 `ListNode` |
| 补 Go 标准库缺口的容器 | `Stack`、`Queue`、`Heap`、`PriorityQueue` | **不要**。Java 有 `ArrayDeque`/`PriorityQueue`，Python 有 `deque`/`heapq`，手写一份反而没人会用 |
| 测试辅助函数 | `Ints2TreeNode`、`Tree2ints`… | 要，跟着对应的节点类型一起移植 |

`NestedInteger`（4 道题）、`Interval` / `Point`（早期题型，现在多用 `int[][]`）用到再补。Java 的 `NestedInteger` 在 LeetCode 上是 `interface`，本地测试还要配一个实现类，比另外两个复杂。

### 共享结构必须有测试

`structures/` 里的辅助方法是所有树/链表题的地基，**写错了会让题解测试给出假结果**——初版 Java `toLevelOrder` 就因为往 `ArrayDeque` 塞 null 子节点而 NPE（Python 的 `deque` 允许 None，同样的写法在那边不会暴露）。所以三门语言各自都有测试，且都跟着全量测试跑：

| 语言 | 测试位置 | 由谁跑 |
|:---|:---|:---|
| Go | `structures/go/*_test.go` | `gotest.sh` 末尾的 `go test ./structures/...` |
| Java | `structures/java/test/StructuresTest.java` | `javatest.sh` 全量模式下单独编译运行 |
| Python | `structures/python/*_test.py` | `pytest.ini` 的 `testpaths` 包含 `structures/python` |

Java 的测试放在 `test/` **子目录**是必须的：`javatest.sh` 取共享结构用的是 `structures/java/*.java` 这个非递归 glob，放同级的话这个测试类会被编进每一道题的产物里。

单题模式（`make test ID=94`）**不跑**共享结构的测试——那个模式是为了快速迭代一道题。改了 `structures/` 记得跑一次全量 `make test`。

Go 的目录叫 `go` 但包名是 `structures`（`go` 是关键字，不能当包名），所以 import 路径末段和包名不一致——这是刻意的，为了和 `java/`、`python/` 对称。实测 `build`/`vet`/`gofmt`/`test` 全部正常。题解里用显式别名 `import structures ".../structures/go"` 让它更好读。

**节点类型的定义要和 LeetCode 给的一字不差**（`TreeNode.java` / `tree_node.py` 里的 `TreeNode` 都是照搬的），这样题解在本地和在提交框里可以原样来回复制。测试辅助方法另外放（Java 放 `TreeNodes`，Python 放同文件的函数），不要挂到节点类上。

### 怎么新增一个共享结构

三门语言的接入都是自动的（Java 按 glob 取整个目录、Python 靠 `sys.path`、Go 是普通包），**加文件就生效，不用改测试脚本**。但各语言有各自的正确性约束：

**Go** —— 在 `structures/go/` 加 `.go` 文件，包名必须是 `structures`。没有别的限制。

**Java** —— 在 `structures/java/` 加 `.java` 文件，public 类名和文件名一致。关键约束：

> 这里的每个类都会被加进**每一个**题目目录的编译单元。所以类名要足够独特，一旦和某道题里的辅助类重名，那道题直接编译失败（`duplicate class`）。反过来说，题目目录里也不要定义和这里同名的类。

**都要配测试**——见上一节，`structures/` 的东西写错了会让所有相关题解的测试给出假结果。

**Python** —— 在 `structures/python/` 加 `.py` 文件。关键约束：

> 这个目录在 `sys.path` 上，**文件名就是模块名**。不要起 `queue.py`、`heapq.py`、`collections.py` 这类和标准库重名的名字。`conftest.py` 用的是 `sys.path.append` 而不是 `insert(0)`，所以标准库优先，重名时是你的模块取不到（会立刻报错），而不是标准库被静默遮蔽——但仍然别这么命名。

**三门语言共同的约定：**

1. **节点类型的定义照搬 LeetCode 给的那段注释，一字不改**（字段名、构造器都不动）。这样题解在本地和在提交框之间可以原样复制。
2. **测试辅助另外放**，不要挂到节点类上：Java 放 `XxxNodes` 工具类（`TreeNode` → `TreeNodes`），Python 放同文件的模块级函数，Go 放同包的函数。挂上去就破坏了第 1 条。
3. **命名跨语言可对应**：`Ints2TreeNode`（Go）/ `TreeNodes.build`（Java）/ `build_tree`（Python）是同一件事，各自用本语言的惯例，但读者要能一眼对上。
4. 只在真需要时加。halfrost 的 Go 版树相关函数有 13 个，Java / Python 侧移植了建树、层序展开、三种遍历、查找、比较，其余用到哪个补哪个。

加新语言时，`structures/<lang>/` 只是其中一步，完整的六步清单见下文「加一门新语言」。

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
    assert solution.Solution().twoSum([2, 7, 11, 15], 9) == [0, 1]
```

题解本体是 `class Solution`（见约束 8），所以要先 `solution.Solution()` 实例化。
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

### 8. 题解本体必须能原样复制到 LeetCode

**这是本仓库最硬的一条约束。** 题解文件的结构固定为：

```
[本地接线区]   ← LeetCode 上不需要、本地才要的几行
[题解本体]     ← 逐字等于 LeetCode 提交框里的内容
```

两个方向都要成立：从 LeetCode 复制进来只需**补接线**，从仓库复制出去只需**删接线**，本体一个字都不改。

| 语言 | 本地接线区 | 题解本体 |
|:---|:---|:---|
| Go | `package leetcode` + `import structures ".../structures/go"` + `type TreeNode = structures.TreeNode` | `func inorderTraversal(root *TreeNode) []int {…}` |
| Java | `import java.util.List;` 等（LeetCode 预置了 `java.util.*`） | `class Solution {…}` |
| Python | `from tree_node import TreeNode`（LeetCode 预置了节点类型） | `class Solution:` + 带 `self` 的方法 |

**接线区的判据**：凡是 LeetCode 提交框里**不存在**的代码行，都属于接线区。这包括 Go 的 `package leetcode`、三门语言的 import 语句、Go 的类型别名。所以「以下是题解本体」这个标记必须放在**所有 import 之后**，而不是文件最开头——放错位置的话，按标记复制出去会连 import 一起带走。

各语言的接线内容不同，标记位置也不同：

| 语言 | 接线区包含 | 标记放在 |
|:---|:---|:---|
| Go | `package leetcode` + `import` + `type TreeNode = structures.TreeNode` | 类型别名之后 |
| Java | `import java.util.…` | import 之后 |
| Python | `from typing import …`、`from tree_node import …` | import 之后 |

由此有三条不能改的细节：

- **Java 的类声明是 `class Solution`，不带 `public`**。LeetCode 模板里带 `public` 的是 `TreeNode`，`Solution` 没有。文件名仍是 `Solution.java`，非 public 类不要求和文件名匹配，测试在同一个 default package 里照样访问得到。
- **Python 的题解是 `class Solution:` 里带 `self` 的方法**，不是顶层函数。Python 靠缩进表达结构，改成顶层函数意味着每道题复制时都要删类声明、去 `self`、整段反缩进——三件事，不是"差一个 self"。测试里用 `solution.Solution().方法名(...)` 取。
- **Go 用类型别名而不是直接写 `*structures.TreeNode`**：`type TreeNode = structures.TreeNode`（注意是别名 `=`，不是定义新类型），这样本体里的签名和 LeetCode 一字不差。这是上游 halfrost 的做法，照搬。

### 9. 统计看文件存不存在，不看内容

`util.LoadSolutions` 判断一道题"有没有题解"，看的是目录里有没有对应语言的文件，**不看文件内容**。所以 `make new` 一建完目录，那道题就立刻计入 README 的统计，哪怕里面还是骨架。

这是有意保持简单：正常流程是「建目录 → 写题解 → 跑测试 → `make readme`」，刷新 README 时题解已经写完了。但如果先 `make new` 一批题占坑，统计会虚高。

同一个判据也决定了单题测试跳不跳过某门语言，于是有个实际后果——**只打算写一门语言时，`make new` 要带 `LANGS=`**：

| | `make new ID=94 LANGS=java` | `make new ID=94`（三门都生成，只写了 Java） |
|:---|:---|:---|
| 目录里 | 只有 `Solution.java` | `Solution.go` / `solution.py` 也在，是空骨架 |
| `make test ID=94` | Go/Python 跳过 | Go/Python 照跑，靠骨架里的 `t.Skip` / `pytest.skip` 通过 |
| README | `[Java]` | `[Go] [Python] [Java]`，点进去两个是空骨架 |
| 语言统计 | Java 1 / Go 0 / Python 0 | 三门都算 1 ⚠️ |

事后补语言只要把文件放进目录再 `make readme`，不用重建目录。

由此，`LoadSolutions` 返回的 `pending`（"建了目录但还没有题解"）在实际中几乎恒为 0，因为 `make new` 总会创建语言文件。它只在手工建了空目录时才会非零。**不要因为这个计数总是 0 就以为逻辑坏了。**

### 10. 题解 README 的「题目大意」故意留空

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

六步，缺一不可。前三步只让 README 认得这门语言，**第四、五步才让树和链表的题能写**——漏掉的话，加完语言会发现第一道 `TreeNode` 的题就卡住。

**1. 登记语言** —— `ctl/util/util.go` 的 `Languages` 表加一行：

```go
{Name: "C++", Ext: ".cpp", TestSuffix: "_test.cpp", Entry: "Solution.cpp"},
```

README 的渲染逻辑不用动，`Solution` 列会自动多出这门语言的链接，语言统计表也自动多一行。

**2. 骨架文件** —— `leetcode/0000.Template/` 里加 `Solution.cpp` 和 `SolutionTest.cpp`。`ctl new --langs cpp` 靠扩展名从这里复制，不需要改 `new.go`。骨架要遵守约束 8：分「本地接线区」和「题解本体」两段，本体形态和 LeetCode 该语言的模板一致。

**3. 测试脚本** —— 照 `javatest.sh` 抄一个 `cpptest.sh`。必须做到：

- 接受可选的目录参数（单题模式），不给参数则遍历 `ROOTS`
- 「这道题没写这门语言」按**跳过**处理，不是失败（见「单题测试」一节）
- 失败时退出码非零

**4. 共享数据结构** —— 新建 `structures/cpp/`，至少把 `TreeNode` 和 `ListNode` 照 LeetCode 该语言的定义移植过去，外加建树/建链表的测试辅助。不要照搬 Go 的 `Stack`/`Queue`/`Heap`/`PriorityQueue`，先看这门语言的标准库有没有。

**5. 把共享结构接进编译** —— 这是最容易漏的一步，各语言机制不同：

| 语言 | 接法 |
|:---|:---|
| Go | 普通包，题解 `import` + 类型别名即可，脚本不用改 |
| Java | `javatest.sh` 把 `structures/java/*.java` 加进每次 `javac` 的源文件列表 |
| Python | `conftest.py` 把 `structures/python` 加进 `sys.path` |
| C / C++ | `cpptest.sh` 编译时加 `-I structures/cpp`（已验证可行） |
| Rust | 见下 |

**6. 共享结构的测试** —— `structures/<lang>/` 里的辅助方法写错了会让所有相关题解给出假结果（见「共享结构必须有测试」），所以要跟着全量测试跑。注意 Java 那种「共享目录按 glob 全量编译」的语言，测试要放在 `test/` 子目录，否则测试类会被编进每一道题。

最后在 `Makefile` 加 `test-<lang>` 目标（记得支持 `ID=`）、并入 `test`，再在 `.github/workflows/test.yml` 加一个 job。

**Rust 要单独决策**：Rust 的模块名必须是合法标识符，而 `0094.Binary-Tree-Inorder-Traversal` 数字开头、带 `.` 和 `-`，三条全犯，常规 `mod` 声明指不到题目目录。逃生口是 `#[path = "..."] mod xxx;` 加 `rustc --test <任意路径>.rs`，或者用 Cargo workspace + 一个 `src/lib.rs` 逐题 `#[path]` 登记（`ctl new` 可以自动追加）。**这两条都没有实测过**，选哪条要先验证。注意问题出在题目目录的命名上，和 `structures/rust/` 无关。

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

这个文件是可选的，绝大部分操作不需要登录：

| 操作 | 需要登录？ |
|:---|:---|
| `new <免费题>`、`build readme` 的题目表格 | 否 |
| `new <会员题>` 的题目描述 | 是，且账号要有 LeetCode 会员，否则接口返回空 |
| `build readme` 的「个人数据」和「已 AC 但未收录」 | 是 |

所以不要因为「没有 config.toml」就认为工具坏了或去加什么回退逻辑——未登录是正常路径。
