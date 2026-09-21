# AGENTS.md

给 AI 编码助手（Claude Code / ChatGPT / Codex 等）看的项目说明。人类看 [README.md](./README.md)。

`CLAUDE.md` 通过 `@AGENTS.md` 导入本文件，两边共用同一份说明，改这里就够了。

**修改或提交前，先检查当前分支、工作区和暂存区，并遵守下文[「分支与标签」](#分支与标签)。工具链默认在 `feature/*` 开发；当前处于 `main` 不代表允许直接提交，「提交」也不等于授权发布或推送。**

## 项目概述

个人算法题解仓库，一道题一个目录，同时支持 Go / Python / Java 三门语言的题解。

仓库结构和 `ctl` 工具链以 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 为模板，改造点见下文「相对上游的修改」。题解代码由仓库作者本人编写，初始化代码由 Claude Code 生成；ChatGPT（Codex）参与后续工具链检查、测试修复、IDE 配置与文档完善，方案和变更由作者确认。

## 常用命令

```sh
make help                        # 列出所有命令
make init                        # 初始化：检查工具链、装依赖、跑测试、生成 README
make ide [IDE=idea]              # 首次创建独立 IDE 项目，无须题号；默认三个 IDE
make new ID=1                    # 新开一题（三门语言）
make new ID=15 LANGS=go          # 只要某一门或某几门语言
make test                        # 三门语言全跑
make test ID=94                  # 只测一道题，那题没写的语言安静跳过
make test-go / test-python / test-java   # 同样支持 ID=94
make readme                      # 重新生成根 README.md
make fmt / vet / tidy / clean
```

`make` 只是 `ctl` 的薄封装，`make new ID=1` 等价于 `cd ctl && go run . new 1`。**ctl 的所有命令都必须在 `ctl/` 目录下执行**，它用的是相对路径（`../leetcode/`、`./template/`）。

## IDE

三门语言的 IDE 配置写在根 README 的「IDE 配置」一节。容易被问到的点：

- **首次创建用 `make ide`，不要求 ID，也不手工只复制运行 XML**。`scripts/ide-init.py` 创建完整的独立项目（`.idea`、`.iml`、运行配置），`IDE=idea/pycharm/goland` 可只建一个。用 IDE 打开命令打印的 `.ide/<IDE>` 文件夹，不打开 `current_problem.xml` 单文件。已有目录一律拒绝覆盖；不要新增默认覆盖或自动删除旧配置的行为。回归测试在 `scripts/ide_init_test.py`，在临时仓库验证首次创建和防覆盖。
- **初始化与选题分开**：Go/Python 的 `__PROBLEM_DIR__` 保留为待配置目标，不能留空后退化成全量测试；Java 初始只接入共享结构，由用户标记当前一道题的 Sources Root。不要把第 94 题或 `0000.Template` 当默认目标。SDK/解释器在 IDE 中由用户选择，不写入机器特定的 SDK 名或绝对路径；Java 语言级别从 `javatest.sh` 读取。
- **多个 JetBrains IDE 不能共用同一个项目目录的 `.idea`**。Project SDK 和 Module SDK 会互相覆盖，表现为 IDEA 跑完后 PyCharm 丢失 interpreter，反过来也一样。各自在 `.ide/idea`、`.ide/pycharm`、`.ide/goland` 保存独立项目配置，再把同一个仓库根加为 Content Root。`.ide/` 不提交。不要通过反复重建 `.venv` 来处理这个问题。
- **每个 IDE 只有一个可编辑的当前题目运行配置，项目只初始化一次**。固定模板在 `ide-templates/<IDE>/current_problem.xml.tpl`，当前配置在 `.ide/<IDE>/.idea/runConfigurations/current_problem.xml`，名称固定为 `Go - Current Problem`、`Python - Current Problem`、`Java - Current Problem`，不带题号。换题修改实际目标，不改名称和模板、不重建项目、不让 `make new` 自动复制一批配置。SDK 和必要的项目文件保留；迁移备份放在被忽略的 `.ide-backups/`。

- **PyCharm 必须把 `structures/python` 标记为 Sources Root** —— 运行时是 `conftest.py` 把它加进 `sys.path` 的，属于运行期行为，IDE 静态分析看不到，不标记则 `from tree_node import ...` 一直标红。这不是配置错误，也别为了"修"它去改 `conftest.py`。
- **PyCharm 单题 pytest 配置使用 Script path，不使用同名模块目标**。每题都叫 `solution_test.py`，`solution_test.test_inorder_traversal` 可能被 IDE 解析到 `0000.Template`。以控制台 `Launching pytest with arguments ...` 的路径为准，而不是运行配置名称；工作目录设为仓库根。模板显示 `test_solve[NOTSET] SKIPPED` 是预期行为，不能靠删除模板的 skip 来修复选错文件。
- **IDEA 同一时刻只能让一个题目目录进入 Java 编译范围** —— 每题都是 default package 的 `class Solution`，多个一起索引会撞成 `Duplicate class`。统一使用独立项目引用仓库根 + 逐题标记 Sources Root，编译输出设为仓库根的 `out/idea`。换题只改运行配置名称不够，还须切换 Sources Root 并 Rebuild；入口始终是 `SolutionTest`。旧的按题打开配置可能仍有题内 `out/`，下面的过滤必须保留。

### 不要去掉测试脚本里对 `out/` 的过滤

IDEA 按「按题打开」那种配法工作时，会在题目目录里建 `out/`，并且把所有非 `.java` 文件**当资源复制进去**——`solution.py`、`solution_test.py`、`Solution.go`、`Solution_test.go` 都有一份副本。这些副本随着改代码而过时，而**测试工具不看 `.gitignore`**。实测两条链路都会把它们当真：

- `go list ./leetcode/...` 把 `out/production/<题目>/` 当成一个**真实的包**编译测试
- pytest 把 `out/` 里那份旧 `solution_test.py` 也收集了（第 94 题当时是 6 个幽灵用例 + 4 个真实用例 = 10）

旧题解配旧测试永远是绿的，**这种"测试说谎"比测试失败危险得多**。所以：

- `gotest.sh` 先单独检查 `go list` 是否成功，再过滤带 `/out/` 的包路径。不能把 `go list` 放回进程替换：其中的失败不会触发外层 `set -e`，曾导致没有 Go 时测试仍报成功
- `pytest.ini` 的 `norecursedirs` 里有 `out`。**一旦设了 `norecursedirs` 就会整个覆盖 pytest 的默认值**，所以那一行把 `*.egg .* build dist` 这些默认项也显式写了出来，删掉任何一个都等于把对应的默认排除关掉

改这两处时不要"简化"掉过滤。

三个测试脚本内部都会切到仓库根（`cd "$(dirname "${BASH_SOURCE[0]}")"`），所以从任何工作目录调用结果都一致。**不要去掉这个 cd**：`gotest.sh` 用的是 `./leetcode/...` 相对路径，从子目录跑会直接报 `no such file or directory`；`pytest.sh` 更隐蔽，pytest 的 rootdir 跟着 cwd 走，从子目录跑会**静默缩小测试范围**（只跑当前目录那几个）却仍然显示绿色。

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

### 日常测试不跑 `0000.Template`

骨架是 `ctl new` 的复制源，属于工具链而不是题解。日常 `make test` 里它只贡献一条 SKIP 和几行噪音，却夹在每道真实题目的输出中间，所以三个测试脚本默认跳过它。

但**不能完全不验**——骨架坏了，之后每次 `make new` 都会产出坏目录。所以：

| 场景 | 是否测骨架 | 怎么控制 |
|:---|:---|:---|
| 日常 `make test` / `make test-go` 等 | 否 | 默认 `TEST_TEMPLATE=0` |
| `make init` | 是 | `scripts/init.sh` 里三处设 `TEST_TEMPLATE=1` |
| CI | 是 | `.github/workflows/test.yml` 三个 job 都设 |
| 显式 `make test ID=0` | 是 | 指定了目录就尊重用户，不受默认排除影响 |

改这些脚本时不要把默认值反过来，也不要在 `make init` / CI 里去掉这个变量——那等于骨架再没人验。回归在 `scripts/tooling_test.py` 的 `test_template_*` 三组。

题号解析失败时 Makefile 用 `$(error ...)` 直接中止，**不能**静默退化成全量测试——那会让你以为测了单题，实际跑了整个仓库。

题号用十进制文本解析并补零，`0094` 和 `94` 必须等价。不要直接用 `printf '%04d' "$id"`：它把前导零当八进制，`0094` 报错后曾误匹配到 `0000.Template`。

测试脚本通过 `scripts/test-common.sh` 统一输出 `PASS / FAIL / SKIP` 汇总；`make test` 由 `scripts/test.sh` 跑完三门语言再汇总，任一失败都必须非零退出。脚本回归测试在 `scripts/tooling_test.py`，跟随全量 Python 测试运行。单题模式仍不跑共享结构和工具链回归测试。

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

链表展开按节点身份检测环。不要恢复上游的 100 节点限制：它会把正常长链表误判成可能有环。三门语言的测试都覆盖 101 个重复值节点，防止长度限制和按值判环两类错误。

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

每道题的题解类都叫 `Solution`（default package）。`javac leetcode/*/*.java` 一次性编译会报 `duplicate class: Solution`。`javatest.sh` 每次在 `out/java/` 下建立独立的临时目录，再按题目分别编译，结束后清理。不要合并编译，也不要复用旧 `.class`：删除测试源码后，旧 `SolutionTest.class` 仍可能被当成当前测试运行。

Java 测试没有引入 JUnit，就是一个 `main` + 断言抛 `AssertionError`。

题目测试用 `structures/java/ExampleTests.java` 记录用例，`check(name, () -> 题解调用, want)` 接受延迟求值，题解异常也会计入失败。末尾的 `finish()` 必须保留，它输出总计，并在失败时抛异常。骨架只调用 `skip` 和 `finish`，仍须保持跳过状态。报告器的测试放在已有的 `structures/java/test/StructuresTest.java` 中。

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

## 写题解测试的分工

`ctl new` 只复制骨架，**不生成任何测试用例**——新建出来的测试文件里一条断言都没有，只有一句 `SKIP`。用例是人写的，或者由 AI 助手按下面的规则代写。

### 默认只转录题目示例

`ctl new` 已经把题目描述（含全部 Example）抓进了题目 README 的 `## 题目` 一节。**代写测试时，把那里的示例机械翻译成断言，题目给几个就是几个，不要自行添加。**

这是转录不是判断：数据来自 LeetCode 官方，不会编错。

**以英文版为准，不要用折叠里的中文翻译。** 同一道题的中英文版本内容会不一致，官方中文翻译滞后。第 94 题实测：

| | 英文版（`content`） | 中文版（`translatedContent`） |
|:---|:---|:---|
| 示例数 | **4 个** | 3 个 |
| `[1,2,3,4,5,null,8,null,null,6,7,9]` | Example 2 | **整个缺失** |

中文版缺的恰好是那棵大树——按下一节的对拍结果，它是唯一能覆盖全部顺序错误的用例。在 leetcode.cn 的网页上默认看到的是中文版，所以会觉得只有 3 个示例。

README 里 `## 题目` 一节本来就是英文版（中文折叠在 `<details>` 里），照着那里转录即可。

### 额外用例必须给出依据

想加题目没给的边界（空输入、单元素、退化结构……），**必须说明它能抓到题目示例抓不到的什么错误**，说不出来就不加。

这条是有教训的。第 94 题最初的测试里，按「中序遍历最容易错在左右顺序，链式结构最能暴露」的惯例加了「纯左链」「纯右链」两个用例。事后写了几个典型的错误实现来对拍，逐个看每个用例能不能把错误实现判失败：

| 用例 | 左右写反 | 写成前序 | 写成后序 | 只判一边的空 |
|:---|:---|:---|:---|:---|
| 题目自带的大树 | 能发现 | 能发现 | 能发现 | 发现不了 |
| 题目自带的空树 | 发现不了 | 发现不了 | 发现不了 | **能发现（仅此一个）** |
| 自行添加的纯左链 | 能发现 | 能发现 | **发现不了** | 发现不了 |
| 自行添加的纯右链 | 能发现 | **发现不了** | 能发现 | 发现不了 |

（「能发现」= 拿有 bug 的实现跑这个用例，测试会失败，说明用例尽到了职责。）

题目自带的大树已经覆盖全部顺序错误，自行添加的两条链各自还**漏掉一种**，严格更弱——给测试增重却没增强。而题目自带的空树是唯一能发现「只判了一边的空」的用例，那种写法在所有非空树上都正确。

所以：**通用的"补边界"直觉不能当依据**。要么能具体说出抓什么错，要么像上面这样实测对拍，否则就只转录题目示例。

### 偏离骨架时，改骨架而不是另写一套

写题解或测试时如果觉得 `leetcode/0000.Template` 的骨架不好用、于是自己另写了一种形式，**那是骨架的问题，去改骨架**，不要让两套形式并存。

骨架的唯一作用是给出可照抄的示范。**没人照抄的骨架是坏的**——它不但没省事，还制造了"仓库里有两种写法"的分裂，后来的人不知道该跟哪个。

这条是有教训的。第 94 题的测试当初没按骨架写，而是另起了一套（匿名 struct + 用例名，而不是骨架里的具名 `question`/`para`/`ans`）。后来第 104 题照着 94 写，于是两道真实题目都偏离了骨架，只有骨架自己是孤例。复盘下来骨架那套确实更差：用例无名，失败时报不出是哪个 Example；Python 更把所有用例塞在一个函数里，第一个 `assert` 挂了后面就不跑。

**判据很简单：实际写出来的东西和骨架不一样时，先问"哪个更好"，然后让两者一致。** 两道题都自发偏离同一个方向，基本可以断定是骨架错了。

### 不要"随机生成"用例

每条断言都要有来源：抄自题目示例，或者针对某个具体错误的推理。测试领域真正的随机化是 property-based testing——随机造输入、拿一份**参考实现**对拍。本仓库没有参考实现，刷题场景也不适合，不要往这个方向走。

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

| 引用 | 职责 |
|:---|:---|
| `main` | 作者的个人题解及个人 README，不是工具链开发分支 |
| `dev` | 工具链改动的集成与验证分支 |
| `feature/<任务名>` | 从 `dev` 派生的单项开发分支，如 `feature/toolchain` |
| `template` | 已验证、可发布的骨架 + 工具链，不含个人题解、个人统计或凭据 |
| `init` 标签 | 用户初始化自己题解仓库的稳定起点，指向已确认发布的 `template` 提交 |

### 开发与发布流程

1. 工具链、公共模板和维护文档默认在从 `dev` 派生的 `feature/<任务名>` 上修改和提交；用户明确指定时才直接在 `dev` 开发。不要直接在 `main` 或 `template` 上开发。
2. 按改动范围完成验证，经用户确认后将开发分支合入 `dev`。本地提交完成不代表已合入其他分支。
3. 只有用户明确要求发布或同步骨架时，才把 `dev` 中已验证的纯工具链改动同步到 `template`，并遵守下文的[骨架 README 生成约束](#骨架-readme-的生成约束)。若来源混有个人内容，只挑选相关提交或差异，不能直接整分支合入。
4. 发布内容检查、测试通过后，才按用户确认的发布范围更新 `init`；开发提交不得自动移动标签。`main` 也只有在用户明确要求时才接收工具链更新。不要为了「同步」强行把所有分支指向同一提交。

### AI 助手操作 Git 的硬性要求

- 修改和提交前检查 `git branch --show-current`、`git status --short --branch` 和暂存区。若当前在 `main` / `template`，先按上述流程进入开发分支；分支或范围不明确时先询问，不能以当前 HEAD 作为授权依据。
- 用户说「提交」只授权约定开发分支上的本地提交，不等于授权合并其他分支、修改标签、重写历史或推送。发布和推送必须有对应的明确授权。
- 使用显式文件路径或分块暂存，先检查 `git diff --cached --name-only`，再审查差异。不要用 `git add .` / `git add -A` 笼统带入个人题解、README 个人数据、`.ai/`、本地 IDE 配置、备份或凭据。
- 切换分支时保留未提交修改和未跟踪文件，不能用 `reset --hard`、`clean` 或覆盖文件来消除阻碍。存在冲突时先询问，或在不改动原工作区的独立 worktree 中处理。
- 若误提交到 `main` / `template`，先核实推送情况，将待保留提交保存到开发分支，再经用户确认修正相关引用；只处理本次误操作，不回退无关历史，不擅自强制推送。
- 只有得到推送授权后才执行 `git push`，并明确指定远端及目标分支或标签，不使用 `--all`、`--mirror`、`--tags` 批量扩大范围。移动已发布的 `init` 或其他历史需要单独说明影响并获得确认。
- 完成后报告当前分支、提交号、实际发生的合并或标签变动，以及是否推送。明确区分「本地提交」「本地合并」「更新标签」和「远端推送」。

## README 生成与个人数据

### 未登录时不渲染「个人数据」一节

`{{.PersonalData}}` 和别的占位符不一样：**它连 `## 个人数据` 这个标题一起渲染**，未登录时整个返回空串，那一节在 README 里直接消失。

标题放在 `models/user.go` 而不是 `template.markdown` 里，就是为了让它跟着内容一起消失——模板里要是留着标题，未登录时会剩一个后面什么都没有的空标题。**不要把标题挪回模板。**

理由：未登录时接口不报错，只是所有 AC 计数返回 0，渲染出来是一张全 0 的表，看起来像"一道题都没做"而不是"没拿到数据"，反而误导人。

另外，**模板里不要在注释或正文中写出 `{{.` 开头的占位符字面量**——替换是全文本匹配的，会把注释里那个也一起换掉（踩过：写了个解释性 HTML 注释，结果输出成"的 ## 标题在　里面"）。

### 骨架 README 的生成约束

`build readme` 只要 `ctl/config.toml` 存在就会带上 Cookie，拿回来的是**本人已登录的数据**，README 的「个人数据」表格里会出现真实的 AC 计数。放在 `main` 上是对的（那张表就是为此存在的），但 `template` / `init` 是给任何人当起点用的，出现别人的刷题数据没有意义。

所以发布时应在不含个人题解的独立 `template` 工作区生成 README，并使用：

```sh
make readme-anon        # 等价于 cd ctl && go run . build readme --anonymous
```

`--anonymous` 会直接跳过 `config.toml`，强制匿名请求，不渲染个人数据一节。

它不会排除磁盘上的题解文件，所以不能在带个人题解的工作区生成后直接用于骨架发布，也不能覆盖作者工作区的个人 README。

这条是踩过坑才加的：作者配好 Cookie 之后，有两个提交把 `Accepted|**1**|...` 写进了本该通用的 README，后来用 rebase 重写历史才清掉。

## 凭据

`ctl/config.toml` 存 LeetCode 的 Cookie，已在 `.gitignore` 里（`config.toml` 不带斜杠，任意层级生效）。**不要提交它，不要把 Cookie 写进任何其他文件或日志。**

这个文件是可选的，绝大部分操作不需要登录：

| 操作 | 需要登录？ |
|:---|:---|
| `new <免费题>`、`build readme` 的题目表格 | 否 |
| `new <会员题>` 的题目描述 | 是，且账号要有 LeetCode 会员，否则接口返回空 |
| `build readme` 的「个人数据」和「已 AC 但未收录」 | 是 |

所以不要因为「没有 config.toml」就认为工具坏了或去加什么回退逻辑——未登录是正常路径。
