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
