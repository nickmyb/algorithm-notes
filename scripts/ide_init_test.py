"""从没有 .ide 的仓库初始化，避免只有运行 XML 却无法打开项目。"""

import json
from pathlib import Path
import shutil
import subprocess
import sys
import xml.etree.ElementTree as ET

import pytest

ROOT = Path(__file__).resolve().parents[1]
CONFIGS = [("goland", "Go", "go"), ("pycharm", "Python", "python"), ("idea", "Java", "java")]


@pytest.fixture
def repo(tmp_path):
    root = tmp_path / "repo with spaces & 中文"
    (root / "scripts").mkdir(parents=True)
    for name in ("Makefile", "go.mod", "javatest.sh", "scripts/ide-init.py"):
        shutil.copy2(ROOT / name, root / name)
    shutil.copytree(ROOT / "ide-templates", root / "ide-templates")
    return root


def run(repo, *args, cwd=None):
    return subprocess.run(
        ["make", "-C", str(repo), "ide", f"PYTHON={sys.executable}", *args],
        cwd=cwd or repo, text=True, stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT, timeout=20,
    )


def snapshot(root):
    return {str(path.relative_to(root)): path.read_bytes()
            for path in root.rglob("*") if path.is_file()}


@pytest.mark.parametrize("ide,language,suffix", CONFIGS)
def test_complete_project_without_problem_or_sdk(repo, ide, language, suffix):
    # 没有 leetcode 目录、虚拟环境或全局 SDK，初始化仍应成功。
    result = run(repo, f"IDE={ide}")
    assert result.returncode == 0, result.stdout
    project = repo / ".ide" / ide
    assert str(project) in result.stdout
    assert "尚未选择 SDK/解释器和题目" in result.stdout
    assert sorted(path.name for path in (repo / ".ide").iterdir()) == [ide]

    modules = ET.parse(project / ".idea/modules.xml")
    registered = modules.findall("./component/modules/module")
    assert len(registered) == 1
    assert registered[0].get("filepath") == f"$PROJECT_DIR$/algorithm-notes-{suffix}.iml"
    module = ET.parse(project / f"algorithm-notes-{suffix}.iml")
    content = module.find("./component[@name='NewModuleRootManager']/content")
    assert content.get("url") == "file://$MODULE_DIR$/../.."
    sources = {node.get("url") for node in content.findall("sourceFolder")}
    assert sources == (set() if ide == "goland" else {f"file://$MODULE_DIR$/../../structures/{suffix}"})
    excluded = {node.get("url") for node in content.findall("excludeFolder")}
    for folder in (".ide", ".ide-backups", ".idea", ".venv", "out", "ctl/.cache"):
        assert f"file://$MODULE_DIR$/../../{folder}" in excluded
    if ide == "idea":
        assert "file://$MODULE_DIR$/../../structures/java/test" in excluded
        misc = ET.parse(project / ".idea/misc.xml")
        assert misc.find("./component/output").get("url") == "file://$PROJECT_DIR$/../../out/idea"
    vcs = ET.parse(project / ".idea/vcs.xml")
    assert vcs.find("./component/mapping").get("directory") == "$PROJECT_DIR$/../.."

    configs = list((project / ".idea/runConfigurations").glob("*.xml"))
    assert len(configs) == 1
    config = ET.parse(configs[0]).getroot().find("configuration")
    assert config.get("name") == f"{language} - Current Problem"
    assert config.find("module").get("name") == f"algorithm-notes-{suffix}"
    if ide == "goland":
        assert config.find("package").get("value").endswith("/__PROBLEM_DIR__")
        assert "__GO_MODULE__" not in configs[0].read_text()
    elif ide == "pycharm":
        target = config.find("option[@name='_new_target']").get("value")
        assert json.loads(target) == "$PROJECT_DIR$/../../__PROBLEM_DIR__/solution_test.py"
        assert json.loads(config.find("option[@name='_new_targetType']").get("value")) == "PATH"
    for name, body in snapshot(project).items():
        text = body.decode("utf-8")
        assert str(repo) not in text  # 移动或复制仓库后仍通过相对宏引用当前仓库。
        assert "0000.Template" not in text
        assert "0094" not in text
        assert "project-jdk-name" not in text
        assert "jdkName=" not in text
        if name.endswith((".xml", ".iml")):
            ET.fromstring(text)


def test_default_initializes_all_from_other_working_directory(repo, tmp_path):
    result = run(repo, cwd=tmp_path)
    assert result.returncode == 0, result.stdout
    assert {path.name for path in (repo / ".ide").iterdir()} == {row[0] for row in CONFIGS}
    assert not (repo / ".idea").exists()
    assert not (repo / "leetcode").exists()
    assert not (repo / ".venv").exists()


def test_existing_project_stops_all_before_writing(repo):
    first = run(repo, "IDE=idea")
    assert first.returncode == 0, first.stdout
    config = repo / ".ide/idea/.idea/misc.xml"
    config.write_text("用户已选择的 SDK，必须保持原样", encoding="utf-8")
    before = snapshot(repo)
    again = run(repo)
    assert again.returncode != 0
    assert "已存在，未覆盖任何配置" in again.stdout
    assert snapshot(repo) == before


def test_incomplete_project_is_not_overwritten(repo):
    project = repo / ".ide/pycharm"
    project.mkdir(parents=True)
    (project / "current_problem.xml").write_text("用户的旧文件")
    before = snapshot(repo)
    result = run(repo, "IDE=pycharm")
    assert result.returncode != 0
    assert snapshot(repo) == before


def test_bad_template_is_detected_before_creating_projects(repo):
    (repo / "ide-templates/idea/current_problem.xml.tpl").write_text("<broken")
    result = run(repo)
    assert result.returncode != 0
    assert not (repo / ".ide").exists()


@pytest.mark.parametrize("target", ["parent", "project"])
def test_symlink_is_not_followed(repo, tmp_path, target):
    outside = tmp_path / "other project"
    outside.mkdir()
    if target == "parent":
        (repo / ".ide").symlink_to(outside, target_is_directory=True)
    else:
        (repo / ".ide").mkdir()
        (repo / ".ide/idea").symlink_to(outside, target_is_directory=True)
    result = run(repo, "IDE=idea")
    assert result.returncode != 0
    assert not list(outside.iterdir())


def test_unknown_ide_creates_nothing(repo):
    result = run(repo, "IDE=unknown")
    assert result.returncode != 0
    assert not (repo / ".ide").exists()


def test_java_release_comes_from_test_script(repo):
    (repo / "javatest.sh").write_text('JAVA_RELEASE="${JAVA_RELEASE:-21}"\n')
    result = run(repo, "IDE=idea")
    assert result.returncode == 0, result.stdout
    misc = ET.parse(repo / ".ide/idea/.idea/misc.xml")
    assert misc.find("component").get("languageLevel") == "JDK_21"


def test_initializer_can_be_called_from_outside_repo(repo, tmp_path):
    result = subprocess.run(
        [sys.executable, str(repo / "scripts/ide-init.py"), "--ide", "pycharm"],
        cwd=tmp_path, text=True, capture_output=True, timeout=20,
    )
    assert result.returncode == 0, result.stderr
    assert (repo / ".ide/pycharm/.idea/modules.xml").is_file()
