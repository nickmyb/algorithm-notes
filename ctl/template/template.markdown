# algorithm-notes

我的算法题解仓库，支持 Go / Python / Java 多语言题解。

## 初心

这个仓库是给自己刷题用的：**把刷过的题整理到一起，并且记下自己是怎么做出来的。**

LeetCode 会存提交记录，但存的是结果——一次次提交散在各自的题目里，翻起来是一长串时间戳。当时为什么选这个思路、先写错了哪一版、复杂度怎么算出来的、那几个测试用例各自在防什么，都没留下。过两个月再看自己的 AC 代码，和看别人的没多大区别。

所以这里一道题一个目录，代码和记录放在一起：

- **解题思路和复杂度写在题目 README 里**，和题解代码同一个目录，不用另开一个笔记本，也不用在两处之间对着找。
- **测试用例本身就是记录**。为什么只测这几个、哪个用例在防哪种写错，想清楚了就写在用例旁边，下次再看不用重新推一遍。
- **同一道题可以三门语言并排写**，横着一看就知道各语言的写法差在哪。`TreeNode` / `ListNode` 由仓库统一提供，定义和 LeetCode 给的一字不差。
- **根 README 自动汇总**：做过哪些题、各语言各多少，`make readme` 刷一遍就有，不用手工维护目录。

剩下的功夫都花在让「记录」这件事不费劲上——因为一旦费劲，就不会坚持：

- **题解本体和 LeetCode 提交框逐字一致**，三门语言都是。本地多出来的那几行接线（`package`、import、类型别名）集中在开头，标记之后的部分复制进、复制出都不用改一个字。
- **`make new ID=1` 一条命令开题**：查标题、建目录、抓英文题目和官方中文翻译、按需生成三门语言骨架，打开就能写；写完 `make test ID=1` 直接跑。
- **输出里没有噪音**，只有自己那道题的用例名和结果。用不用 IDE 都行——编辑器加 `make test` 足够，想用 JetBrains 就 `make ide`。

clone 完三条命令就能开始：

```sh
git switch -c my-solutions init
make init
make new ID=1
```

如果你也想要这样一个仓库，不用 fork 我的题解——从 `init` 标签起步，那里是干净的骨架和工具链，没有任何我的个人内容。

## 关于这个仓库

仓库结构和 `ctl` 工具链以 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 为模板构建，在此感谢原作者 [@halfrost](https://github.com/halfrost) 开源的这套题解组织方式和自动化工具。本仓库在其基础上做了三件事：改造 `ctl` 让 README 支持多语言题解链接、去掉 Hugo 站点与 PDF 电子书相关的部分、补上 Python 和 Java 的测试链路。

代码归属：

- **题解代码**（`leetcode/` 下各题的 `Solution.*`）由我本人编写。卡住的题会参考 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go) 的实现，思路借鉴自那里。
- **仓库初始化代码**（`ctl/` 工具链改造、`Makefile`、`conftest.py`、各测试脚本、CI 配置和本 README 模板）由 [Claude Code](https://claude.com/claude-code) 生成。
- **后续工程协作**：ChatGPT（Codex）参与工具链问题排查、多语言测试汇总与回归检查、共享结构修复、IDE 项目隔离、运行配置模板和文档完善。

感谢 ChatGPT 和 Claude 在仓库搭建、排错和维护中的协助。题解代码仍由我本人编写；项目方案与变更由我确认和维护。

## 快速开始

```sh
git clone git@github.com:nickmyb/algorithm-notes.git
cd algorithm-notes
git switch -c my-solutions init  # 从干净的骨架创建自己的分支
make init             # 检查工具链、装依赖、跑通测试、生成 README
make ide              # 可选：用 JetBrains IDE 写题才需要
make new ID=1         # 开始第一题
```

**不用先折腾环境变量。** 装了但没进 `PATH` 是最常见的卡点，所以 `make init` 会先自己找一遍：Go 翻 `/usr/local/go`、`~/go/go*`、`~/sdk/go*`、`/usr/lib/go-*`、`/opt/homebrew`，Python 翻 `/usr/bin`、pyenv，JDK 翻 `$JAVA_HOME`、`/usr/lib/jvm`、SDKMAN、macOS 的 `JavaVirtualMachines`。找到就记进 `local.mk`（本机配置，不进版本库），后面的 `make new` / `make test` 自动跟着用。

只有真的没装才会停下来，这时按它给的链接装完重跑即可。装在冷门位置的话仍可以手动指定：

```sh
make init GO=/path/to/go/bin/go       # PYTHON= / JAVAC= 同理，优先级高于 local.mk
```

Go 是必需的（`ctl` 用 Go 写的），版本下限跟着 `go.mod` 走；Python 和 Java 缺了只会警告并跳过对应的测试链路，不写那门语言就不受影响。

`make ide` 不是必须的——用编辑器写题、靠 `make test` 验证完全可以。只有要在 IntelliJ IDEA / PyCharm / GoLand 里跑题解时才需要，它会为三个 IDE 各建一个独立项目。跑完还要在 IDE 里选一次 SDK/解释器并指定当前题目，命令结束时会把每个 IDE 的具体步骤和示例路径打出来，详见 [ide-templates/README.md](./ide-templates/README.md)。

## 写一道题

从开题到提交的完整流程：

```sh
make new ID=1                    # 1. 建目录，题目描述自动填好
# 或：make new ID=1 LANGS=go,python  # 只选一种建题方式，不要重复创建同一题
                                 # 2. 写题解（见下），删掉测试骨架里的 skip
                                 # 3. 照着题目 README 里的 Example 写测试用例
make test ID=1                   # 4. 只测这一道，反复改的时候快得多
make readme                      # 5. 刷新本 README 的题目表格
make test                        # 6. 提交前全量跑一遍
git add . && git commit          # 7. 提交
```

### 1. 建目录

`make new ID=1` 会去 LeetCode 按题号查标题，建出 `leetcode/0001.Two-Sum/`，把**英文题目描述和折叠的官方中文翻译**写进该题的 `README.md`，再从 `leetcode/0000.Template/` 复制各语言的骨架。

使用 IDE 时，项目只需初始化一次。换题后按 [IDE 模板与换题说明](./ide-templates/README.md)修改已有的 `Current Problem` 目标，不重建 `.ide` 项目，也不用新增运行配置。

**只打算写一门语言时记得带 `LANGS=`。** 三门全生成却只写一门的话，下面的题目表格会显示三个链接、语言统计三门都算 1，而点进去另外两个是空骨架——ctl 判断"有没有题解"看的是**文件存不存在，不看内容**。事后想补语言，把文件放进目录再 `make readme` 即可，不用重建。

### 2. 写题解

题解文件分两段，详见下文[多语言约定](#题解本体可以原样复制到-leetcode)：

```
[本地接线区]   ← import、Go 的 package 和类型别名
[题解本体]     ← 逐字等于 LeetCode 提交框里的内容
```

从 LeetCode 复制进来只需补接线，复制出去只需删接线，本体一个字都不改。`TreeNode`、`ListNode` 由 `structures/` 提供，不用自己定义。

### 3. 写测试

**照着该题 README 里 `## 题目` 那节的 Example 转录**，题目给几个就写几个。注意以**英文版**为准——官方中文翻译滞后，示例可能更少（第 94 题英文版 4 个、中文版 3 个，中文版缺的恰好是覆盖面最广的那个）。

想加题目没给的边界用例，先想清楚它能抓到题目示例抓不到的什么错误；说不出来就别加，那只会给测试增重而不增强。

`make new` 不会替你生成用例。它已经把题目描述抓下来了，看上去只差解析 Example 块，但 LeetCode 的示例结构自己就不统一——老题用 `<pre>`，新题用 `<div class="example-block">`，而且**同一道题的中英文版本可能用不同结构**（第 94 题英文版是 `example-block`，中文版还是 `<pre>`）。再叠加多参数输入、浮点比较、任意顺序等情况，自动解析只能做到「大部分对」。而**测试用例错了比没有更糟**：它会让错误的题解显示成绿色。所以这一步保持手工转录。

Java 用共享的 `ExampleTests` 报告每个用例和最终汇总。删除 `tests.skip(...)`，按下面的形式转录示例，并保留最后的 `tests.finish()`：

```java
ExampleTests tests = new ExampleTests();
Solution s = new Solution();
tests.check("Example 1", () -> s.twoSum(new int[] {2, 7, 11, 15}, 9), new int[] {0, 1});
tests.finish();
```

`check` 会记录断言失败和题解抛出的异常，继续跑后面的用例；`finish` 汇总后在有失败时抛出 `AssertionError`，IDE 和命令行都会正确报错。新建的骨架仍然是跳过状态。

### 4. 跑测试

`make test ID=1` 按题号找目录（`94` → `leetcode/0094.*`），**那道题没写的语言会安静跳过**——只写了 Java 的题，Go 和 Python 不会因为"没有测试"而报错。题号不存在时直接报错退出，不会静默退化成全量。

单题模式不跑 `structures/` 的测试，改了共享结构记得跑一次全量 `make test`。

全量模式也不跑 `0000.Template`——骨架是 `make new` 的复制源，不是题解，放在日常输出里只会添噪音。它由 `make init` 和 CI 负责验证；想单独测用 `make test ID=0`。

三门语言保留各自的用例明细，脚本最后统一输出 `===== Go/Python/Java: PASS/FAIL/SKIP =====`。`make test` 会跑完三门语言，再输出 `===== All: PASS =====` 或 `FAIL`；任意一门失败，整条命令返回非零。`PASS` 表示本次运行没有失败，骨架的跳过数仍应查看用例明细。

在 IDE 的原生运行按钮里，Go 和 Python 的格式由 IDE/测试框架决定；Java 直接运行 `SolutionTest.main` 也会输出用例数及总结果。需要相同的语言汇总格式时，从 IDE 终端运行 `make test ID=1`。

### 5. 填题解 README

`make new` 生成的那份 README 有四个小节：

| 小节 | 谁来填 |
|:---|:---|
| `## 题目` | **自动**：英文原文 + 折叠的官方中文翻译 |
| `## 题目大意` | **自己写**，一两句话复述题意。这一步的价值就在于自己写，所以没有自动填 |
| `## 解题思路` | 思路推导，为什么这么做 |
| `## 复杂度` | 时间 / 空间复杂度 |

前三节沿用 halfrost 原仓库的格式，`## 复杂度` 是本仓库加的——原仓库把复杂度记在 `ctl/meta/` 里用于渲染站点，那套机制没有搬过来，所以在 README 里留个固定位置。

## 命令

```sh
make help                        # 列出所有命令
make init                        # 初始化：检查工具链、装依赖、跑测试、生成 README
make ide [IDE=idea]              # 首次创建独立 IDE 项目，不需要 ID；默认三门语言
make new ID=1 [LANGS=go,python]  # 新开一题，自动抓题目描述填进 README
make test [ID=94]                # 跑测试，带 ID 只测一道
make test-go / test-python / test-java   # 单语言，同样支持 ID=94
make readme                      # 重新生成本 README
make readme-anon                 # 同上但不含个人数据，给 template 分支用
make fmt / vet / tidy            # Go 格式化、静态检查、依赖整理
make clean                       # 清构建产物，不碰题解
```

`make` 只是 `ctl` 的薄封装，`make new ID=1` 等价于 `cd ctl && go run . new 1`。细节见 [ctl/README.md](./ctl/README.md)。

`LANGS` 的取值是 `go`、`python`、`java` 的任意组合，默认三门全生成。下面的题目表格按目录里**实际存在**的文件渲染，写了几门语言就显示几个链接。

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
    assert solution.Solution().twoSum([2, 7, 11, 15], 9) == [0, 1]
```

### 同一语言的多种解法

写在**同一个题解文件里的多个函数**。入口函数保持 LeetCode 给的签名名，其余解法加后缀说明算法：

```
twoSum              # 入口，哈希表解法
twoSumBruteForce
twoSumTwoPointers
```

测试不为每种解法单写一份，而是**同一组用例跑所有实现**，顺带保证它们结果一致（Go 用 `t.Run` 子测试，Python 用 `parametrize`）。

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

只在真用到时加。Java / Python 侧目前提供建树、层序展开、三种遍历、查找和比较，其他辅助用到再补。Java 另外提供测试报告器 `ExampleTests`，供题目测试使用。

## 环境要求

开发和验证都在 **Ubuntu 24.04** 上进行，**语言版本也只测过下面这三个**，别的系统和版本不保证可用：

| 语言 | 版本 | 版本来源 | 说明 |
|:---|:---|:---|:---|
| Go | **1.26.4** | `go.mod` 的 `go` 指令 | Go 没有 LTS，官方只维护最近两个大版本。CI 用 `go-version-file: go.mod`，和本地同一个来源 |
| Python | **3.12** | `scripts/init.sh` 的版本检查 | Python 没有 LTS，每个小版本一律 5 年支持期。低于 3.12 且找不到更高版本时，`make init` 会跳过 Python 并给出提醒 |
| Java | **17**（LTS） | `javatest.sh` 的 `--release` | Java 的 LTS 每两年一个：8 / 11 / 17 / 21 / 25 |

三者都只有一个声明来源，CI 跟着本地走，不会各说各的。

三者都不在 `PATH` 里也没关系，`make init` 会去常见安装位置找（见[快速开始](#快速开始)），
探测结果写进 `local.mk` 供后续命令复用。装在冷门位置时手动指定，优先级高于 `local.mk`：

```sh
make init GO=/path/to/go/bin/go PYTHON=/path/to/python3 JAVAC=/path/to/javac
make test GO=/path/to/go/bin/go        # 其他 make 目标同样接受
JAVA_RELEASE=21 make test-java         # 临时换 Java 目标版本
```

JDK 只需指定 `JAVAC`，配套的 `java` 取同目录的那个——两者必须来自同一个 JDK，否则编译产物跑起来会报 `class file has wrong version`。

Java 的 `--release 17` 同时约束语言特性和可用 API：即使本地装的是 JDK 21，使用 Java 21 的 switch 模式匹配也会在本地编译失败；record 已在 Java 16 正式支持，可以使用。要换目标版本改 `javatest.sh` 里的 `JAVA_RELEASE`，或临时 `JAVA_RELEASE=21 make test-java`。

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
├── conftest.py                  # pytest 的 solution fixture，见「多语言约定」
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

首次配置只需：**执行 `make ide` → 用 IDE 打开指定文件夹 → 选择 SDK 和当前题目**。

在当前仓库根目录运行（不需要题号）：

```sh
make ide
```

默认创建三个独立项目；只用一个 IDE 时，改用 `make ide IDE=idea`、`make ide IDE=pycharm` 或 `make ide IDE=goland`。命令会创建完整项目文件，不需要手工新建空项目或复制 XML。已有目标目录会拒绝覆盖。

然后使用 **File → Open** 打开命令打印的文件夹：

| IDE | 打开的文件夹 | 首次选择的 SDK / 解释器 |
|:---|:---|:---|
| IntelliJ IDEA | `<当前仓库>/.ide/idea` | Project SDK：JDK 17；模块继承项目 SDK |
| PyCharm | `<当前仓库>/.ide/pycharm` | 已有解释器 `<当前仓库>/.venv/bin/python` |
| GoLand | `<当前仓库>/.ide/goland` | GOROOT：`go.mod` 要求的版本 |

**不要打开 `current_problem.xml`、`.iml`、内部 `.idea` 或仓库根。** 如果窗口里只有一个 XML 文件，就不是完整项目。初始化命令不安装 SDK；`.venv` 由 `make init` 创建。

项目文件已经接好源码目录、共享结构、排除项和 Java 输出目录，但**还没有选择题目**。首次运行前：

- GoLand：编辑 `Go - Current Problem`，Test kind 为 Package，填写当前题目的完整 Package path，Pattern 留空。
- PyCharm：编辑 `Python - Current Problem`，Target 选 Script path，选择当前题目的 `solution_test.py`，不要选同名模块。
- IDEA：只将当前一道题标为 Sources Root，保持 `structures/java` 的标记，然后 Rebuild，运行 `Java - Current Problem`。

`__PROBLEM_DIR__` 是待选题标记，不能原样运行；不会默认跑第 94 题或 `0000.Template`。以后换题只改现有目标（Java 切换 Sources Root），不改配置名称、SDK，也不重新执行 `make ide`。

逐项操作、固定模板及安全恢复命令见 **[IDE 首次配置与换题](./ide-templates/README.md)**。

<details>
<summary>为什么要独立项目，以及编译输出注意事项</summary>

IDEA、PyCharm、GoLand 共用仓库根的 `.idea` 时，Project SDK / Module SDK 会相互覆盖。现在各自的项目文件保存在 `.ide/<IDE>`，Content Root 引用同一份仓库源码，互不改写 SDK。`.ide/` 不提交；迁移时无需删除或重建 `.venv`。

编译输出集中放在仓库根的 `out/`，避免在题目目录里出现复制的 Go/Python 源码和测试。

`out/` 按生成工具分开，IDEA 和命令行测试不共用 classpath：

| 路径 | 来源和生命周期 |
|:---|:---|
| `out/idea/production/<模块名>/` | 当前 IDEA 项目的编译结果，供 IDE 运行；可保留直到下次清理 |
| `out/java/run.XXXXXX/` | `javatest.sh` 每次新建的临时目录，结束时自动删除，不复用旧 `.class` |
| `out/production/`、`out/java/leetcode/`、`out/java/structures/` | 旧 IDEA 配置或旧脚本的遗留目录；当前配置不使用，可随构建产物一起清理 |

IDEA 的 [Resource patterns](https://www.jetbrains.com/help/idea/compiler.html) 决定哪些非 Java 文件作为资源复制，所以 `out/idea` 中出现 `.go`、`.py`、README 或 `.iml` 不代表源文件放错了位置。它们不是题解来源，不要编辑或运行这些副本；三个 IDE 都应排除 `out/`，测试脚本里的过滤也必须保留。需要减少资源副本时，可在 Settings → Build, Execution, Deployment → Compiler → Resource patterns 中追加 `!*.go;!*.py;!*.pyc;!*.md;!*.iml`，保留原有排除项；之后在没有测试运行时清理并重新构建。

> **旧的按题打开配置可能留下题内 `out/`**：IDEA 会把编译输出建在题目目录里面（`leetcode/0094.xxx/out/`），而且把非 `.java` 文件作为资源复制进去。这些副本会随着改代码而过时，而测试工具不看 `.gitignore`：
>
> - `go list ./leetcode/...` 把 `out/production/<题目>/` 当成一个真实的包编译测试
> - pytest 把 `out/` 里那份旧 `solution_test.py` 也收集了
>
> 旧题解配旧测试永远是绿的，比测试失败更危险。仓库已经挡住了这一层（`gotest.sh` 过滤 `/out/`、`pytest.ini` 的 `norecursedirs` 含 `out`），所以现在用是安全的；但如果你想从根上避免，在 `Project Structure` → `Project` → **Compiler output** 里把输出路径改到题目目录之外。

配好之后 `TreeNode` / `TreeNodes` 在题解里能正常补全和跳转。嫌配置麻烦的话，IDE 只当编辑器用，验证一律走 `make test-java ID=94`。

</details>

## 分支与标签

| 引用 | 内容 |
|:---|:---|
| `main` | 我的个人题解及个人 README |
| `feature/*` | 从 `dev` 派生的工具链开发分支 |
| `dev` | 工具链集成与验证 |
| `template` | 经确认发布的干净骨架和工具链，不含个人题解或数据 |
| `init` 标签 | 用户初始化仓库的稳定起点，指向已确认发布的 `template` 提交 |

开发、集成、发布和推送分别确认：不直接在 `main` 开发工具链，不因开发提交自动更新 `template` / `init`，不自动推送。完整规范见 [AGENTS.md](./AGENTS.md#分支与标签)。

## 数据来源

`ctl` 默认从 **leetcode.cn** 取数据，理由写在 `ctl/util/repo.go` 的注释里：题库是 leetcode.com 的严格超集（.com 的 4055 道编号题一道不缺，另有 388 道 LCR / 面试题 / LCP / LCS），而编号题的标题和 slug 仍是英文，目录命名在两个站点通用。换站点改那一个常量即可。

**日常用不着登录**：`make new` 建目录、抓题目描述，`make readme` 生成题目表格，都是匿名请求。

需要登录的只有两处，都在 README 的统计部分：

| 需要登录 | 不登录的表现 |
|:---|:---|
| 下面的「个人数据」表格 | 整节不显示；配置有效 Cookie 后自动出现 |
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
