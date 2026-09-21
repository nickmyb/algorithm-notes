"""测试入口的回归检查：针对已复现的假成功、错题号和过期产物。"""

import os
from pathlib import Path
import shutil
import subprocess

import pytest

ROOT = Path(__file__).resolve().parents[1]


@pytest.fixture
def repo(tmp_path):
    root = tmp_path / "repo"
    root.mkdir()
    (root / "scripts").mkdir()
    for name in (
        "Makefile", "gotest.sh", "pytest.sh", "javatest.sh",
        "scripts/test.sh", "scripts/test-common.sh", "scripts/problem-dir.sh",
        "scripts/find-tool.sh",
    ):
        shutil.copy2(ROOT / name, root / name)
    for name in ("0000.Template", "0094.Example", "0104.Example"):
        (root / "leetcode" / name).mkdir(parents=True)
    return root


def run(repo, script, *args, env=None, cwd=None):
    return subprocess.run(
        ["bash", str(repo / script), *map(str, args)],
        cwd=cwd or repo,
        env={**os.environ, **(env or {})},
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=30,
    )


def executable(repo, name, body):
    path = repo / name
    path.write_text("#!/usr/bin/env bash\n" + body)
    path.chmod(0o755)
    return str(path)


@pytest.mark.parametrize("number", ["94", "0094", "00094", "0104"])
def test_decimal_problem_id(repo, number):
    result = run(repo, "scripts/problem-dir.sh", number)
    assert result.returncode == 0, result.stdout
    assert result.stdout.strip() == f"leetcode/{int(number):04d}.Example"


@pytest.mark.parametrize("number", ["", "missing", "9999"])
def test_bad_id_never_falls_back_to_all_tests(repo, number):
    result = subprocess.run(
        ["make", "-n", "test", f"ID={number}"], cwd=repo,
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=10,
    )
    assert result.returncode != 0
    assert "scripts/test.sh" not in result.stdout


def test_missing_go_is_failure(repo):
    result = run(repo, "gotest.sh", env={"GO": str(repo / "missing-go")})
    assert result.returncode != 0
    assert "\n  Go: FAIL\n" in result.stdout


def test_partial_go_list_failure_is_not_swallowed(repo):
    go = executable(repo, "fake-go", 'echo example/leetcode/0094.Example\nexit 2\n')
    result = run(repo, "gotest.sh", env={"GO": go})
    assert result.returncode != 0
    assert "\n  Go: FAIL\n" in result.stdout


def test_go_filters_ide_copies_and_runs_tool_tests(repo):
    go = executable(repo, "fake-go", '''
if [ "$1" = list ]; then
    printf '%s\n' example/leetcode/0094.Example example/leetcode/0094.Example/out/production/copy
else
    printf '%s\n' "$@"
fi
''')
    result = run(repo, "gotest.sh", env={"GO": go}, cwd=repo / "scripts")
    assert result.returncode == 0, result.stdout
    assert "example/leetcode/0094.Example" in result.stdout
    assert "/out/" not in result.stdout
    assert "./ctl/..." in result.stdout
    assert "./structures/..." in result.stdout


@pytest.mark.parametrize("script", ["gotest.sh", "pytest.sh", "javatest.sh"])
def test_missing_directory_is_failure(repo, script):
    result = run(repo, script, "leetcode/does-not-exist")
    assert result.returncode != 0
    assert ": FAIL\n" in result.stdout


@pytest.mark.parametrize("script", ["gotest.sh", "pytest.sh", "javatest.sh"])
def test_unwritten_language_is_skipped(repo, script):
    pytest_cmd = executable(repo, "fake-pytest", "exit 5\n")
    result = run(repo, script, repo / "leetcode/0094.Example", env={"PYTEST": pytest_cmd},
                 cwd=repo / "scripts")
    assert result.returncode == 0, result.stdout
    assert ": SKIP\n" in result.stdout


@pytest.mark.parametrize("script", ["gotest.sh", "javatest.sh"])
def test_template_is_skipped_by_default(repo, script):
    """0000.Template 是 ctl new 的复制源，属于工具链，日常 make test 不该测它。

    它每次只贡献一条 SKIP 和几行噪音，却会出现在每道题的测试输出里。
    """
    (repo / "leetcode/0000.Template/Solution.go").write_text("package leetcode\n")
    (repo / "leetcode/0094.Example/Solution.go").write_text("package leetcode\n")
    go = executable(repo, "fake-go", 'if [ "$1" = list ]; then\n'
                    'printf "%s\\n" example/leetcode/0000.Template example/leetcode/0094.Example\n'
                    'else printf "%s\\n" "$@"; fi\n')
    result = run(repo, script, env={"GO": go, "TEST_TEMPLATE": "0"})
    assert result.returncode == 0, result.stdout
    assert "0000.Template" not in result.stdout


@pytest.mark.parametrize("script", ["gotest.sh", "javatest.sh"])
def test_template_runs_when_requested(repo, script):
    """make init 和 CI 设 TEST_TEMPLATE=1：骨架坏了每次 make new 都产出坏目录，
    不能完全不验。"""
    (repo / "leetcode/0000.Template/Solution.go").write_text("package leetcode\n")
    go = executable(repo, "fake-go", 'if [ "$1" = list ]; then\n'
                    'printf "%s\\n" example/leetcode/0000.Template\n'
                    'else printf "%s\\n" "$@"; fi\n')
    result = run(repo, script, env={"GO": go, "TEST_TEMPLATE": "1"})
    assert result.returncode == 0, result.stdout
    if script == "gotest.sh":
        assert "0000.Template" in result.stdout


def test_template_runs_when_named_explicitly(repo):
    """显式 make test ID=0 时始终尊重用户的选择，不受默认排除影响。"""
    (repo / "leetcode/0000.Template/Solution.go").write_text("package leetcode\n")
    go = executable(repo, "fake-go", 'printf "%s\\n" "$@"\n')
    result = run(repo, "gotest.sh", repo / "leetcode/0000.Template",
                 env={"GO": go, "TEST_TEMPLATE": "0"})
    assert result.returncode == 0, result.stdout
    assert "0000.Template" in result.stdout


def test_all_runs_remaining_languages_after_failure(repo):
    pytest_cmd = executable(repo, "fake-pytest", "echo PYTHON_RAN\nexit 5\n")
    (repo / "leetcode/0094.Example/Solution.go").write_text("package leetcode\n")
    result = run(repo, "scripts/test.sh", "leetcode/0094.Example",
                 env={"GO": str(repo / "missing-go"), "PYTEST": pytest_cmd})
    assert result.returncode != 0
    assert "\n  Go: FAIL\n" in result.stdout
    assert "PYTHON_RAN" in result.stdout
    assert "\n  Java: SKIP\n" in result.stdout
    assert "  All: FAIL\n" in result.stdout
    # 总结论用更重的横线框住，和单门语言的结论区分开
    assert result.stdout.rstrip().endswith("═" * 56)


@pytest.mark.skipif(not shutil.which("javac") or not shutil.which("java"), reason="需要 JDK")
def test_java_does_not_run_deleted_test_class(repo):
    problem = repo / "leetcode/0094.Example"
    (problem / "Solution.java").write_text("class Solution {}\n")
    test = problem / "SolutionTest.java"
    test.write_text('public class SolutionTest { public static void main(String[] args) {'
                    'System.out.println("OLD_TEST_RAN"); }}\n')
    first = run(repo, "javatest.sh", "leetcode/0094.Example")
    assert first.returncode == 0, first.stdout
    assert "OLD_TEST_RAN" in first.stdout
    test.unlink()
    second = run(repo, "javatest.sh", problem, cwd=repo / "scripts")
    assert second.returncode == 0, second.stdout
    assert "OLD_TEST_RAN" not in second.stdout
    assert "COMPILE OK" in second.stdout


@pytest.mark.skipif(not shutil.which("javac") or not shutil.which("java"), reason="需要 JDK")
def test_java_failure_propagates(repo):
    problem = repo / "leetcode/0094.Example"
    (problem / "SolutionTest.java").write_text(
        'public class SolutionTest { public static void main(String[] args) {'
        'throw new AssertionError("expected failure"); }}\n')
    result = run(repo, "javatest.sh", problem)
    assert result.returncode != 0
    assert "\n  Java: FAIL\n" in result.stdout


# ---------- 工具链探测 ----------
#
# 这组测试盯的是同一个卡点：go / python3 装了但没进 PATH。它是新手跑 make init
# 时最常见的失败，而 README 承诺"三条命令开始刷题"，探测不到就等于承诺落空。


def fake_go(path, version):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(f'#!/usr/bin/env bash\necho "go version go{version} linux/amd64"\n')
    path.chmod(0o755)
    return path


def find(repo, home, kind, minimum, name, path="/usr/bin:/bin", **extra):
    """在受控的 HOME / PATH / JAVA_HOME 下跑 find_tool，返回选中的路径（找不到则为空串）。"""
    result = subprocess.run(
        ["bash", "-c", f'source "$1"; find_tool {kind} {minimum} {name}',
         "_", str(repo / "scripts/find-tool.sh")],
        env={"HOME": str(home), "PATH": path, **extra},
        text=True, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, timeout=30,
    )
    return result.stdout.strip()


def test_finds_go_installed_outside_path(repo, tmp_path):
    """~/go/goX/bin 是官方下载器和手动解压最常见的落点，都不会自动进 PATH。"""
    home = tmp_path / "home"
    expected = fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    assert find(repo, home, "go", "1.26.4", "go") == str(expected)


def test_ignores_go_below_minimum(repo, tmp_path):
    """版本不够的 go 要当作没有：让 make init 直接说清楚，
    好过放过去在 go mod tidy 里炸出一句看不懂的报错。"""
    home = tmp_path / "home"
    fake_go(home / "go/go1.21.0/bin/go", "1.21.0")
    assert find(repo, home, "go", "1.26.4", "go") == ""


def test_picks_newest_when_several_are_installed(repo, tmp_path):
    home = tmp_path / "home"
    fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    fake_go(home / "sdk/go1.9.7/bin/go", "1.9.7")
    # 1.9.7 的字典序比 1.26.4 大，按字符串比会选错
    assert find(repo, home, "go", "1.2", "go").endswith("go1.26.4/bin/go")


def test_path_wins_when_it_satisfies_minimum(repo, tmp_path):
    """用户自己配好的 PATH 优先，哪怕别处装着更新的版本。"""
    home = tmp_path / "home"
    fake_go(home / "go/go1.99.0/bin/go", "1.99.0")
    on_path = fake_go(tmp_path / "bin/go", "1.26.4")
    got = find(repo, home, "go", "1.26.4", "go", path=f"{tmp_path / 'bin'}:/usr/bin:/bin")
    assert got == str(on_path)


def test_local_mk_supplies_go_to_later_commands(repo):
    """探测的结果必须落到 local.mk：否则 make init 能跑通，
    下一条 make new 又回到 Makefile 的 GO ?= go，等于没找。"""
    (repo / "local.mk").write_text("GO := /opt/custom/go\n")
    result = subprocess.run(
        ["make", "-n", "test-go"], cwd=repo,
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=10,
    )
    assert result.returncode == 0, result.stdout
    assert "/opt/custom/go" in result.stdout


def test_command_line_go_overrides_local_mk(repo):
    (repo / "local.mk").write_text("GO := /opt/custom/go\n")
    result = subprocess.run(
        ["make", "-n", "test-go", "GO=/opt/explicit/go"], cwd=repo,
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=10,
    )
    assert result.returncode == 0, result.stdout
    assert "/opt/explicit/go" in result.stdout
    assert "/opt/custom/go" not in result.stdout


def test_ci_can_read_python_minimum():
    """CI 用 sed 从 init.sh 抠 PYTHON_MIN 来选 setup-python 的版本。
    改这个变量名会让 CI 静默拿到空串，错误要等到 setup-python 才暴露。"""
    init = (ROOT / "scripts/init.sh").read_text()
    version = [ln.split("=", 1)[1] for ln in init.splitlines()
               if ln.startswith("PYTHON_MIN=")]
    assert len(version) == 1, "scripts/init.sh 必须有且只有一行 PYTHON_MIN=<版本>"
    assert version[0].count(".") == 1 and all(x.isdigit() for x in version[0].split("."))


def test_ci_python_version_matches_init():
    workflow = next((ROOT / ".github/workflows").glob("*.yml")).read_text()
    assert "s/^PYTHON_MIN=//p" in workflow, "CI 的解析表达式和 scripts/init.sh 对不上了"


@pytest.fixture
def clean_path(tmp_path):
    """只含 find-tool.sh 所需工具的 PATH，不含 go / python3 / javac。

    探测测试必须在这种 PATH 下跑：否则本机真实的 /usr/bin/javac 会先被
    "PATH 优先" 那条规则命中，测不到候选目录那段逻辑。
    """
    binx = tmp_path / "cleanbin"
    binx.mkdir()
    for name in ("bash", "sh", "awk", "sed", "sort", "head", "rm", "env", "cat", "dirname"):
        src = shutil.which(name)
        assert src, f"本机没有 {name}，测试无法构造干净 PATH"
        (binx / name).symlink_to(src)
    return str(binx)


def fake_jdk(home, version, with_java=True):
    """造一个假 JDK，返回它的 bin 目录。"""
    binx = home / "jdk/bin"
    binx.mkdir(parents=True, exist_ok=True)
    (binx / "javac").write_text(f'#!/usr/bin/env bash\necho "javac {version}" >&2\n')
    (binx / "javac").chmod(0o755)
    if with_java:
        (binx / "java").write_text('#!/usr/bin/env bash\nexit 0\n')
        (binx / "java").chmod(0o755)
    return binx


def test_finds_jdk_via_java_home(repo, tmp_path, clean_path):
    """JAVA_HOME 设了但 bin 没进 PATH 是 JDK 最常见的装法。"""
    home = tmp_path / "home"
    binx = fake_jdk(home, "21.0.1")
    got = find(repo, home, "java", "17", "javac", path=clean_path, JAVA_HOME=str(binx.parent))
    assert got == str(binx / "javac")


def test_ignores_jdk_below_minimum(repo, tmp_path, clean_path):
    """JDK 8 打印 "javac 1.8.0_432"，按字符串比会误判成大于 17。

    断言的是"这个候选没被选中"而不是"什么都没找到"：跑测试的机器上
    /usr/lib/jvm 里可能真装着合规的 JDK，那是正确行为，不该让测试变红。
    """
    home = tmp_path / "home"
    binx = fake_jdk(home, "1.8.0_432")
    got = find(repo, home, "java", "17", "javac", path=clean_path,
               JAVA_HOME=str(binx.parent))
    assert got != str(binx / "javac")


def test_rejects_javac_without_matching_java(repo, tmp_path, clean_path):
    """javac 和 java 必须同一个 JDK，否则跑测试时会撞上 class file has wrong version。
    只有 javac 的目录（比如只装了 jdk-headless 的残缺安装）不能算数。"""
    home = tmp_path / "home"
    binx = fake_jdk(home, "21.0.1", with_java=False)
    got = find(repo, home, "java", "17", "javac", path=clean_path,
               JAVA_HOME=str(binx.parent))
    assert got != str(binx / "javac")


def test_javatest_honours_javac_and_java_overrides(repo):
    """探测到的 JDK 要真能传进 javatest.sh，否则探测了也白搭。"""
    problem = repo / "leetcode/0094.Example"
    (problem / "Solution.java").write_text("class Solution {}\n")
    javac = executable(repo, "fake-javac", 'echo "FAKE_JAVAC $*"\nexit 1\n')
    result = run(repo, "javatest.sh", problem, env={"JAVAC": javac, "JAVA": "/bin/false"})
    assert "FAKE_JAVAC" in result.stdout, result.stdout


def test_local_mk_supplies_jdk_to_later_commands(repo):
    (repo / "local.mk").write_text("JAVAC := /opt/jdk/bin/javac\nJAVA := /opt/jdk/bin/java\n")
    result = subprocess.run(
        ["make", "-n", "test-java"], cwd=repo,
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=10,
    )
    assert result.returncode == 0, result.stdout
    assert "/opt/jdk/bin/javac" in result.stdout
    assert "/opt/jdk/bin/java" in result.stdout


def test_init_can_read_java_minimum():
    """init.sh 用 sed 从 javatest.sh 抠 JAVA_RELEASE 当 JDK 版本下限，
    抠不到会 die。改 javatest.sh 里那行的写法要同步改 init.sh。"""
    import re
    init = (ROOT / "scripts/init.sh").read_text()
    expr = re.search(r"JAVA_MIN=\"\$\(sed -n '([^']+)'", init)
    assert expr, "scripts/init.sh 里找不到 JAVA_MIN 的解析表达式"
    got = subprocess.run(["sed", "-n", expr.group(1), str(ROOT / "javatest.sh")],
                         text=True, stdout=subprocess.PIPE, timeout=10).stdout.split()
    assert got and got[0].isdigit(), "没能从 javatest.sh 解析出 JAVA_RELEASE"


@pytest.fixture
def init_repo(tmp_path):
    """够 init.sh 跑到 local.mk 那一步的最小仓库。

    再往后它会在"编译 ctl"失败（没有 ctl/），这正好——local.mk 是在那之前
    写的，用这个失败点就能验证顺序，不用给 init.sh 加测试专用的开关。
    """
    root = tmp_path / "initrepo"
    (root / "scripts").mkdir(parents=True)
    for name in ("scripts/init.sh", "scripts/find-tool.sh", "javatest.sh"):
        shutil.copy2(ROOT / name, root / name)
    (root / "go.mod").write_text("module example\n\ngo 1.26.4\n")
    return root


def run_init_sh(init_repo, home, clean_path, **env):
    return subprocess.run(
        ["bash", str(init_repo / "scripts/init.sh")], cwd=init_repo,
        env={"HOME": str(home), "PATH": clean_path, **env},
        text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=60,
    )


# ---------- 显式指定优先于探测 ----------


def test_explicit_tool_is_never_silently_replaced(init_repo, tmp_path, clean_path):
    """指定了就得用指定的那个。

    早先的实现在指定的版本不够时会静默回退去搜索，挑一个别的用还不吭声——
    用户写了路径就是要用那个，换掉必须说。

    两个细节都是踩出来的：

    跑在隔离的临时仓库里，不能用真实仓库当 cwd。这条测试原先直接对着
    ROOT 跑 init.sh——只要样本版本意外通过了检查，init.sh 就会一路往下
    执行，把真实仓库的 local.mk 写成一个 pytest 临时路径（发生过）。

    样本版本用 1.11.0 这种早于 Go modules 的号，不写成"比当前下限低一点"：
    下限是会调的，写得太近，哪天调到那个数，这条测试就从"验证拒绝"
    悄悄变成"验证接受"。
    """
    home = tmp_path / "home"
    old_go = fake_go(tmp_path / "old/go", "1.11.0")
    result = run_init_sh(init_repo, home, clean_path, GO=str(old_go))
    assert result.returncode != 0, result.stdout
    assert "1.11.0" in result.stdout and "低于" in result.stdout
    # 关键：没有换成别的 go 接着往下跑
    assert "整理 Go 依赖" not in result.stdout


def test_explicit_missing_tool_reports_instead_of_searching(init_repo, tmp_path, clean_path):
    home = tmp_path / "home"
    fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    result = run_init_sh(init_repo, home, clean_path, JAVAC="/nope/javac")
    assert result.returncode != 0, result.stdout
    assert "找不到或跑不起来" in result.stdout


def test_default_still_auto_detects(repo, tmp_path):
    """没显式指定时，探测照常工作——上一条不能把自动探测一起关掉。"""
    home = tmp_path / "home"
    expected = fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    assert find(repo, home, "go", "1.26.4", "go") == str(expected)


# ---------- init.sh 真的跑一遍 ----------
#
# 上面那些 local.mk 的测试是手工造出 local.mk 再看 Makefile 认不认，
# 漏掉了"init.sh 到底写没写"这一环。实测漏过一次：一轮重构把写入步骤
# 连同 export 一起删掉了，init 照样报成功，下一条 make new 才炸。


def test_init_actually_writes_local_mk(init_repo, tmp_path, clean_path):
    """探测到 PATH 外的 go 之后，必须把它记下来。

    只在 init.sh 里解析出路径是不够的：那样 make init 报成功，
    下一条 make new 回到 Makefile 的 GO ?= go，直接 go: not found。
    """
    home = tmp_path / "home"
    go = fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    result = run_init_sh(init_repo, home, clean_path)

    local_mk = init_repo / "local.mk"
    assert local_mk.exists(), f"init.sh 没有生成 local.mk\n{result.stdout}"
    assert f"GO := {go}" in local_mk.read_text()


def test_init_omits_tools_already_on_path(init_repo, tmp_path, clean_path):
    """在 PATH 里的就别记进 local.mk。

    记一份没用的绝对路径，将来换版本反而会被这条陈旧记录钉住——而且按
    "显式指定优先"的规则，陈旧路径会被当成用户点名，直接报错而不是重新探测。

    只断言 GO 那一行不在：跑测试的机器上 python3 / JDK 装在 PATH 之外是
    正常的，它们被探测到并记下来是正确行为。
    """
    home = tmp_path / "home"
    binx = tmp_path / "onpath"
    binx.mkdir()
    fake_go(tmp_path / "src/go", "1.26.4").replace(binx / "go")
    run_init_sh(init_repo, home, f"{binx}:{clean_path}")
    local_mk = init_repo / "local.mk"
    text = local_mk.read_text() if local_mk.exists() else ""
    assert "GO :=" not in text, text


def test_tool_path_with_spaces_survives(init_repo, tmp_path, clean_path):
    """工具装在带空格的路径下（macOS 的 Application Support、Windows 挂载盘常见）。

    探测阶段本来就过得去，炸在后面没加引号的 $GO mod tidy 上——路径被拆成两半。
    """
    home = tmp_path / "home"
    go = fake_go(tmp_path / "my tools/go/bin/go", "1.26.4")
    result = run_init_sh(init_repo, home, clean_path, GO=str(go))

    assert "go mod tidy 失败" not in result.stdout, result.stdout
    assert (init_repo / "local.mk").exists(), result.stdout
    assert f"GO := {go}" in (init_repo / "local.mk").read_text()


def test_init_stops_when_version_floor_cannot_be_parsed(init_repo, tmp_path, clean_path):
    """go.mod 解析不出版本下限时必须停，不能带着空下限往下跑。

    空下限会让版本比较失去意义，而报错写成"没找到 go  或更高版本"——
    中间空着的正是版本号，看着像工具坏了而不是仓库坏了。
    """
    home = tmp_path / "home"
    fake_go(home / "go/go1.26.4/bin/go", "1.26.4")
    (init_repo / "go.mod").write_text("module example\n")
    result = run_init_sh(init_repo, home, clean_path)
    assert result.returncode != 0
    assert "没能从 go.mod 解析出" in result.stdout, result.stdout


def test_empty_repo_is_not_a_failure(repo, tmp_path):
    """一道题都没有时，全量 make test 不能报错。

    这是从 init 标签起步的第一现场：make init 成功，紧接着 make test
    就该是绿的。pytest 在"没收集到用例"时返回 5，题解那一段收到 5 是
    正常状态（还没写题），工具链那一段收到 5 才是真坏了。
    """
    calls = tmp_path / "calls"
    pytest_cmd = executable(
        repo, "fake-pytest",
        f'echo "$@" >> {calls}\n'
        # 第一次调用（题解）返回 5，第二次（工具链）返回 0
        f'[ "$(wc -l < {calls})" = 1 ] && exit 5\nexit 0\n')
    result = run(repo, "pytest.sh", env={"PYTEST": pytest_cmd, "TEST_TEMPLATE": "0"})
    assert result.returncode == 0, result.stdout
    assert "还没有 Python 题解" in result.stdout
    assert "structures/python" in calls.read_text(), "工具链那一段仍然要跑"


def test_broken_tooling_collection_still_fails(repo, tmp_path):
    """反过来：工具链那一段收不到用例说明环境或配置坏了，必须报错。"""
    pytest_cmd = executable(repo, "fake-pytest", "exit 5\n")
    result = run(repo, "pytest.sh", env={"PYTEST": pytest_cmd, "TEST_TEMPLATE": "0"})
    assert result.returncode != 0, result.stdout
