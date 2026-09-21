# IDE 运行配置模板

每个 IDE 的项目只初始化一次，之后保留一份可编辑的「当前题目」配置。这里的模板是固定参照，用于首次创建或恢复；不要为了换题而修改模板、重建项目或新增一批运行配置。`make new` 不会创建或改写 IDE 配置。

## 两类文件

| IDE | 固定模板（提交 Git） | 当前题目配置（仅本机） |
|:---|:---|:---|
| GoLand | [goland/current_problem.xml.tpl](./goland/current_problem.xml.tpl) | `.ide/goland/.idea/runConfigurations/current_problem.xml` |
| PyCharm | [pycharm/current_problem.xml.tpl](./pycharm/current_problem.xml.tpl) | `.ide/pycharm/.idea/runConfigurations/current_problem.xml` |
| IDEA | [idea/current_problem.xml.tpl](./idea/current_problem.xml.tpl) | `.ide/idea/.idea/runConfigurations/current_problem.xml` |

模板故意使用 `.xml.tpl` 扩展名，且放在 `.idea` 之外，不出现在运行列表里。不要把 `0000.Template` 当成可运行模板配置：它的跳过状态不能证明当前题目通过。

每个 IDE 的运行下拉框只保留当前题目，名称固定为 `Go - Current Problem`、`Python - Current Problem`、`Java - Current Problem`，不带题号。文件名始终是 `current_problem.xml`。换题只改目标，不改名称；具体题号从目标路径和测试输出确认。若 IDE 保存配置时自动改了文件名，保留它实际使用的那一份即可，不要再补出同内容的第二份。

## 首次创建或恢复

1. 按[根 README 的 IDE 配置](../README.md#ide-配置)创建独立项目；已经配置好的项目不用重建。
2. 关闭对应项目窗口，避免 IDE 将内存中的旧配置重新写回磁盘。
3. 从上表复制该 IDE 的模板到当前题目配置路径，去掉 `.tpl` 后缀；已有可用配置时，不要重复复制。恢复前先备份原文件。
4. 只在复制出来的 Go/Python 文件中替换路径占位符，原模板保持不变；Java 的题目目录在模块中设置。
5. 重新打开 `.ide/goland`、`.ide/pycharm` 或 `.ide/idea`，从下拉框选择当前题目运行。不要在旧 Run 窗口点击 Rerun。

占位符示例：

| 占位符 | 第 94 题的值 | 从哪里取 |
|:---|:---|:---|
| `__PROBLEM_DIR__` | `leetcode/0094.Binary-Tree-Inorder-Traversal` | 题目目录，相对仓库根，不带末尾 `/` |
| `__GO_MODULE__` | `github.com/nickmyb/algorithm-notes` | `go.mod` 第一行 `module` 后的路径，仅 Go 模板需要 |

Java 模板没有题号或路径占位符：题目源码的选择在模块里，见下一节。模板中的 `$PROJECT_DIR$` 是 JetBrains 自带的宏，**不要替换**；这里它指 `.ide/<IDE>`，所以 `$PROJECT_DIR$/../..` 才是仓库根。路径若含 `&`、`"` 等 XML 特殊字符，手工编辑 XML 时需转义，或直接在 IDE 的配置界面里填写。

模板引用的模块名固定为 `algorithm-notes-go`、`algorithm-notes-python`、`algorithm-notes-java`，须与各自项目的模块名一致。SDK/解释器由项目和模块提供，模板不会注册或安装 SDK，也不会写入个人绝对路径。换题不用改 SDK。

## 日常换题：直接修改已有配置

在 **Run → Edit Configurations** 中编辑已有的当前题目，不点 `+` 新建；在 IDE 界面里修改不需要关闭项目。只有直接编辑磁盘上的 XML 时，才先关闭项目。

### GoLand

1. 编辑 `Go - Current Problem`。Test kind 保持 **Package**，Package path 改为新题的完整包路径，例如 `github.com/nickmyb/algorithm-notes/leetcode/0104.Maximum-Depth-of-Binary-Tree`。
2. 测试函数筛选（Pattern）保持空白，让整个包的测试都运行，避免留着 `TestInorderTraversal` 导致新题没有测试可跑。不要只看退出码 0，还要确认实际用例出现。
3. Working directory 保持仓库根。若直接编辑 XML，同时更新 `<directory>` 的题目路径；`root_directory` 和 `filePath` 的仓库根设置不动。

### PyCharm

1. 编辑 `Python - Current Problem`。Target 类型保持 **Script path / script**，路径改为新题的 `solution_test.py`。不要改成 `solution_test.test_xxx` 这样的同名模块目标。
2. Working directory 保持仓库根，解释器仍使用当前模块的 `.venv/bin/python`。换题不会改变 `structures/python` 的 Sources Root。
3. 在控制台确认 `Launching pytest with arguments ...` 指向新题，而不是 `0000.Template` 或 `out/` 里的副本。

只有确实生成了 Python 文件的题目才运行 Python 配置；仅写 Go 的题不需要为它创建 Python 配置。

### IntelliJ IDEA

每题的入口都叫 `SolutionTest`，实际运行哪题由模块的 Sources Root 决定，**不能靠改运行配置名称换题**。

1. 先取消旧题目录的 **Sources Root**，再将新题目录标记为 Sources Root；`structures/java` 的标记保持不变，`structures/java/test` 保持 Excluded。
2. 当前编译范围应只有「新题目录 + `structures/java`」，不能同时包含两道题，也不要把整个 `leetcode` 标成 Sources Root。
3. 运行配置始终使用 `Java - Current Problem`。Main class 保持 `SolutionTest`，模块保持 `algorithm-notes-java`，Working directory 保持仓库根；正常换题不需要改这份运行配置的字段。
4. 执行一次 **Build → Rebuild Project**，再运行当前配置，避免旧题的 `.class` 留在输出目录造成混淆。

若关闭 IDE 后手工修改模块文件，对应的是 `.ide/idea/algorithm-notes-java.iml` 中这一行，只替换题目路径：

```xml
<sourceFolder url="file://$MODULE_DIR$/../../leetcode/0001.Two-Sum" isTestSource="false" />
```

这里的 `0001.Two-Sum` 只是换题示例，需先用 `make new ID=1` 建好对应语言的文件；不要将它和旧题同时保留为源码根。

## 不随题目变化的基础设置

SDK、Content Root、模块名、`structures/` 接入、输出目录和排除项都属于一次性项目设置，不是需要清理的多余配置。`.iml`、`misc.xml`、`modules.xml`、`vcs.xml`、`workspace.xml` 等必要文件应保留。

三个项目均排除 `.ide`、`.ide-backups`、`.venv`、`out`；IDEA 输出仍为根目录的 `out/idea`。旧根目录 `.idea` 和旧 `structures/java/structures.iml` 不再使用。迁移备份放在被 Git 忽略的 `.ide-backups/`，仅用于恢复，不作为项目打开；恢复时先关闭相关项目。

## 协作致谢

感谢 ChatGPT（Codex）参与 IDE 配置冲突排查、独立项目配置、固定运行模板、回归检查及本说明的整理；也感谢 Claude（Claude Code）对仓库初始多语言工具链的建设。本文采用的「固定模板 + 一份当前题目配置」方案经仓库作者讨论确认，后续由作者维护。
