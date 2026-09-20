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
├── structures/                  # Go 题解共用的数据结构（ListNode、TreeNode 等）
├── conftest.py                  # pytest 的 solution fixture，见下文
├── pytest.ini                   # pytest 配置
├── gotest.sh                    # Go 测试 + 覆盖率
├── javatest.sh                  # Java 逐目录编译并测试
├── scripts/init.sh              # make init 的实现
└── Makefile                     # 所有日常命令的入口
```

题目目录名是 `%04d.英文标题`：标题里的空格转 `-`，`' % ( ) , ? / : " !` 去掉，连续的 `-` 收敛成一个。这个名字由 `make new` 自动生成，不用手写。

以后要收 LCP、剑指 Offer 这类非编号题，会新开 `lcp/` 这样的同级目录，因为 `%04d` 表达不了 `LCP 82` 这种题号。

## 命令

```sh
make help                        # 列出所有命令
make new ID=1                    # 新开一题，自动抓题目描述填进 README
make new ID=15 LANGS=go,python   # 只生成指定语言的骨架
make test                        # 跑三门语言的全部测试
make test-go / test-python / test-java
make readme                      # 重新生成本 README 的题目表格
make fmt / vet / tidy            # Go 格式化、静态检查、依赖整理
make clean                       # 清构建产物，不碰题解
```

`make` 只是 `ctl` 的薄封装，`make new ID=1` 等价于 `cd ctl && go run . new 1`。细节见 [ctl/README.md](./ctl/README.md)。

## 写一道题

```sh
make new ID=1     # 建出 leetcode/0001.Two-Sum/，题目描述已经填好
# 写题解，把测试骨架里的 t.Skip / pytest.skip 删掉
make test         # 三门语言一起验
make readme       # 刷新下面的题目表格
```

题解 README 的四个小节：

| 小节 | 内容 |
|:---|:---|
| `## 题目` | `make new` 自动填：英文原文 + 折叠的官方中文翻译 |
| `## 题目大意` | **自己写**，一两句话复述题意。这一步的价值就在于自己写，所以没有自动填 |
| `## 解题思路` | 思路推导，为什么这么做 |
| `## 复杂度` | 时间 / 空间复杂度 |

前三节沿用 halfrost 原仓库的格式，`## 复杂度` 是本仓库加的——原仓库把复杂度记在 `ctl/meta/` 里用于渲染站点，那套机制没有搬过来，所以在 README 里留个固定位置。

## 多语言约定

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

拉题目列表和题目描述都不需要登录。想让下面的「个人数据」表格有内容，按 [ctl/README.md](./ctl/README.md) 配好 `ctl/config.toml`（已在 `.gitignore` 里）。

---

## 个人数据

{{.PersonalData}}

## 题解统计

{{.LanguageTable}}

## 题目列表

{{.TotalNum}}

{{.SolvedTable}}

---

## License

题解代码以 [MIT](./LICENSE) 授权。仓库结构与 `ctl` 工具链改编自 [halfrost/LeetCode-Go](https://github.com/halfrost/LeetCode-Go)，原项目同为 MIT 授权。
