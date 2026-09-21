# IDE 首次配置与换题

只需三步：**运行 `make ide` → 打开项目文件夹 → 在 IDE 中选 SDK 和题目**。项目只初始化一次，换题不重建。

## 1. 用命令创建项目文件

在你实际要使用的仓库根目录执行（包含 `Makefile` 的目录；测试备份仓库时就在备份内执行）：

```sh
make ide
```

这会一次创建 IDEA、PyCharm、GoLand 三个独立项目。不需要题号，不要求已有题解；需要本机有 Python 3，命令本身不联网、不安装 SDK、不运行测试。

只使用一个 IDE 时，用 `make ide IDE=idea` **代替**上面的命令；另外两个取值为 `pycharm`、`goland`。不要先执行全部初始化，再重复执行单 IDE 初始化。

命令会生成完整的 `.idea`、`.iml` 和一份 `Current Problem` 运行配置；Content Root、共享结构、排除项及 Java 输出路径也已设置。目标目录只要已存在就停止，**不会覆盖已有 SDK 或配置**。

## 2. 在 IDE 中打开文件夹

用欢迎页的 **Open** 或 **File → Open**，选择命令打印的绝对路径：

| 使用的 IDE | 应打开的文件夹 | 运行配置名称 |
|:---|:---|:---|
| IntelliJ IDEA | `<当前仓库>/.ide/idea` | `Java - Current Problem` |
| PyCharm | `<当前仓库>/.ide/pycharm` | `Python - Current Problem` |
| GoLand | `<当前仓库>/.ide/goland` | `Go - Current Problem` |

**打开的是 `idea` / `pycharm` / `goland` 文件夹，不是里面的 `.idea`、`.iml`、`current_problem.xml`，也不是仓库根。** 隐藏目录不好找时，将命令打印的路径粘贴到 Open 对话框。不要再新建空项目。

如果窗口标题是 `current_problem.xml`，左侧也只显示这个文件，说明打开了单文件；关闭它，按上表重新打开。JetBrains 区分[打开文件与打开项目目录](https://www.jetbrains.com/help/idea/opening-files-from-command-line.html)，运行 XML 不能替代项目。

## 3. 在 IDE 中选择 SDK 和当前题目

### SDK / 解释器（每个项目只选一次）

| IDE | 在哪里设置 | 选择什么 |
|:---|:---|:---|
| IDEA | File → Project Structure → Project → SDK | JDK 17；模块使用 Project SDK |
| PyCharm | Settings → Python Interpreter → Add Interpreter → Add Local Interpreter | 选择已有环境 `<当前仓库>/.venv/bin/python`，不新建环境 |
| GoLand | Settings → Go → GOROOT | `go.mod` 要求的 Go 版本 |

`.venv` 由 `make init` 创建；如果还没有，先完成仓库初始化。复制了仓库时，解释器也要选**当前副本**的路径，不能沿用原仓库的 SDK 路径。初始化命令不猜测本机 SDK 名称。

### 选择题目（首次运行和以后换题都在这里操作）

初始 Go/Python 配置里的 `__PROBLEM_DIR__` 表示**尚未选题**，Java 也未加入任何题目源码根。不要立即点 Run；不会默认选第 94 题，也不会拿 `0000.Template` 的跳过结果当通过。没有题目时先 `make new ID=1`。

以下以已有的第 94 题为例；换题时换成实际目录。配置名称保持 `Current Problem`，不要点 `+` 新建配置。

**GoLand**：Run → Edit Configurations → `Go - Current Problem`。

- Test kind 选 **Package**，Package path 填 `github.com/nickmyb/algorithm-notes/leetcode/0094.Binary-Tree-Inorder-Traversal`（仓库 module 路径以 `go.mod` 为准）。
- Pattern 保持空白，Working directory 保持当前仓库根。若手工编辑 XML，`package` 和 `directory` 两处目标一起修改。

**PyCharm**：Run → Edit Configurations → `Python - Current Problem`。

- Target 选 **Script path / script**，通过文件选择器选择当前题目的 `solution_test.py`；不要使用 `solution_test.test_xxx` 这样的模块名。
- Working directory 保持当前仓库根。`structures/python` 已由命令标为 Sources Root，不用再手工设置。

**IDEA**：在目录树中将 `leetcode/0094.Binary-Tree-Inorder-Traversal` 右键 → Mark Directory as → **Sources Root**。

- 换题时先取消旧题的 Sources Root，再标记新题；同一时刻只能有一道题，不能标记整个 `leetcode`。
- `structures/java` 已标记，`structures/java/test` 已排除，保持不动。Run → Edit Configurations 中，Main class 保持 `SolutionTest`，模块保持 `algorithm-notes-java`。
- 执行 **Build → Rebuild Project**，再运行 `Java - Current Problem`。

保存后从运行下拉框选择对应的 `Current Problem`，不要在旧 Run 窗口点 Rerun。确认输出的路径及用例属于当前题目；`SKIPPED`、没有测试或目标仍含占位符，都不等于题解通过。

以后换题只重复上面“选择题目”的操作：不运行 `make ide`，不改 SDK，不改名称，不修改固定模板。

## 固定模板与本机文件

| IDE | 固定运行模板（提交 Git） | 可编辑的当前配置（仅本机） |
|:---|:---|:---|
| GoLand | [Go 模板](./goland/current_problem.xml.tpl) | `.ide/goland/.idea/runConfigurations/current_problem.xml` |
| PyCharm | [Python 模板](./pycharm/current_problem.xml.tpl) | `.ide/pycharm/.idea/runConfigurations/current_problem.xml` |
| IDEA | [Java 模板](./idea/current_problem.xml.tpl) | `.ide/idea/.idea/runConfigurations/current_problem.xml` |

`make ide` 使用这些固定模板创建当前配置，并补齐项目文件，不需要手工复制、改后缀或创建模块。`make new` 只创建题解，不改 IDE 配置。

`.ide/` 已被 Git 忽略。保留 IDE 生成的 `.iml`、`misc.xml`、`modules.xml`、`workspace.xml` 等必要文件；它们不是多余的运行配置。直接编辑磁盘上的 XML 前必须先关闭项目，避免 IDE 用内存里的设置覆盖它。`$PROJECT_DIR$` / `$MODULE_DIR$` 是 IDE 的相对路径宏，不要手工替换成个人绝对路径。

## 已有配置或需要恢复

已有可用项目时直接打开，不必执行 `make ide`。如果 `.ide` 已删除，重新执行即可。

若配置残缺需要重建，**先关闭对应 IDE 项目**。以下命令仅备份并重建 IDEA，保留其他 IDE 和 `.venv`：

```sh
mkdir -p .ide-backups
ide_backup_dir=$(mktemp -d .ide-backups/idea.XXXXXX)
mv .ide/idea "$ide_backup_dir/project"
make ide IDE=idea
```

旧项目完整保留在 `$ide_backup_dir/project`（即 `.ide-backups/idea.*/project`）；需要恢复时先关闭项目。不要打开备份作为当前项目，也不要为解决 SDK 冲突反复重建 `.venv`。

三个 IDE 的配置相互独立，但 Content Root 都指向同一份源码。`.ide`、`.ide-backups`、`.idea`、`.venv`、`out` 和缓存均已排除；IDEA 输出为仓库根 `out/idea`。旧根目录 `.idea`、旧 `structures/java/structures.iml` 不作为项目入口。编译副本和过滤原因见[根 README](../README.md#ide-配置)。

## 协作致谢

感谢 ChatGPT（Codex）参与 IDE 配置冲突排查、独立项目初始化、固定运行模板、回归检查及本说明的整理；也感谢 Claude（Claude Code）对仓库初始多语言工具链的建设。「固定模板 + 一份当前题目配置」方案经仓库作者讨论确认，后续由作者维护。
