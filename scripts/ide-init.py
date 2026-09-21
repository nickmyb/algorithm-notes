#!/usr/bin/env python3
"""一次性创建独立 JetBrains 项目；不选题、不注册 SDK、不覆盖已有配置。"""

import argparse
from pathlib import Path
import re
import sys
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[1]
CONFIGS = {
    "goland": ("Go", "go", "WEB_MODULE"),
    "pycharm": ("Python", "python", "PYTHON_MODULE"),
    "idea": ("Java", "java", "JAVA_MODULE"),
}
EXCLUDES = (".ide", ".ide-backups", ".idea", ".venv", "out", "ctl/.cache")


def xml_text(root):
    ET.indent(root, space="  ")
    return ET.tostring(root, encoding="unicode", xml_declaration=True) + "\n"


def project_component(name, **attributes):
    root = ET.Element("project", version="4")
    component = ET.SubElement(root, "component", name=name, **attributes)
    return root, component


def project_files(ide):
    """先在内存中准备完整项目，校验失败时不留下半套配置。"""
    language, suffix, module_type = CONFIGS[ide]
    module_name = f"algorithm-notes-{suffix}"
    module_file = f"{module_name}.iml"
    files = {".idea/.name": f"{ROOT.name} {language}\n"}

    project, manager = project_component("ProjectModuleManager")
    modules = ET.SubElement(manager, "modules")
    ET.SubElement(modules, "module", {
        "fileurl": f"file://$PROJECT_DIR$/{module_file}",
        "filepath": f"$PROJECT_DIR$/{module_file}",
    })
    files[".idea/modules.xml"] = xml_text(project)

    module = ET.Element("module", type=module_type, version="4")
    if ide == "goland":
        ET.SubElement(module, "component", name="Go", enabled="true")
    roots = ET.SubElement(module, "component", name="NewModuleRootManager")
    if ide == "idea":
        roots.set("inherit-compiler-output", "true")
        ET.SubElement(roots, "exclude-output")
    content = ET.SubElement(roots, "content", url="file://$MODULE_DIR$/../..")
    if ide != "goland":
        ET.SubElement(content, "sourceFolder", {
            "url": f"file://$MODULE_DIR$/../../structures/{suffix}",
            "isTestSource": "false",
        })
    for path in EXCLUDES + (("structures/java/test",) if ide == "idea" else ()):
        ET.SubElement(content, "excludeFolder", url=f"file://$MODULE_DIR$/../../{path}")
    for pattern in ("out", "__pycache__", ".pytest_cache"):
        ET.SubElement(content, "excludePattern", pattern=pattern)
    if ide != "goland":
        ET.SubElement(roots, "orderEntry", type="inheritedJdk")
    ET.SubElement(roots, "orderEntry", type="sourceFolder", forTests="false")
    if ide == "pycharm":
        runner = ET.SubElement(module, "component", name="TestRunnerService")
        ET.SubElement(runner, "option", name="PROJECT_TEST_RUNNER", value="pytest")
    files[module_file] = xml_text(module)

    project, manager = project_component("ProjectRootManager", version="2")
    if ide == "idea":
        # Java 版本仍只在 javatest.sh 声明，不在 IDE 初始化器里另写一份。
        match = re.search(r'^JAVA_RELEASE="\$\{JAVA_RELEASE:-(\d+)\}"$',
                          (ROOT / "javatest.sh").read_text(encoding="utf-8"), re.M)
        if not match:
            raise ValueError("无法从 javatest.sh 读取 JAVA_RELEASE")
        manager.set("languageLevel", f"JDK_{match[1]}")
        ET.SubElement(manager, "output", url="file://$PROJECT_DIR$/../../out/idea")
    # SDK 名称是机器本地状态，首次在 IDE 中选择后才由 IDE 写入。
    files[".idea/misc.xml"] = xml_text(project)

    project, mappings = project_component("VcsDirectoryMappings")
    ET.SubElement(mappings, "mapping", directory="$PROJECT_DIR$/../..", vcs="Git")
    files[".idea/vcs.xml"] = xml_text(project)

    template = ROOT / "ide-templates" / ide / "current_problem.xml.tpl"
    run_config = ET.fromstring(template.read_text(encoding="utf-8"))
    if ide == "goland":
        match = re.search(r"^module\s+(\S+)", (ROOT / "go.mod").read_text(encoding="utf-8"), re.M)
        if not match:
            raise ValueError("无法从 go.mod 读取 module 路径")
        for element in run_config.iter():
            for key, value in element.attrib.items():
                element.set(key, value.replace("__GO_MODULE__", match[1].strip('"')))
    # Go/Python 故意保留不存在的 __PROBLEM_DIR__，不得留空后退化成全量测试。
    # Java 初始只接入共享结构；选择一个题目 Sources Root 后才有 SolutionTest。
    files[".idea/runConfigurations/current_problem.xml"] = xml_text(run_config)
    return files


def initialize(ides):
    ide_root = ROOT / ".ide"
    if ide_root.is_symlink() or (ide_root.exists() and not ide_root.is_dir()):
        raise ValueError(".ide 不是普通目录；为保护现有文件，停止初始化")
    for ide in ides:
        target = ide_root / ide
        if target.exists() or target.is_symlink():
            raise ValueError(
                f"{target} 已存在，未覆盖任何配置。\n"
                "已有项目直接打开，换题在 IDE 中编辑 Current Problem。\n"
                "只初始化其他 IDE 时指定 IDE=idea/pycharm/goland；"
                "需要重建时先关闭项目并备份该目录，见 ide-templates/README.md。"
            )
    plans = {ide: project_files(ide) for ide in ides}
    ide_root.mkdir(exist_ok=True)
    for ide, files in plans.items():
        target = ide_root / ide
        target.mkdir()  # 不使用 exist_ok，防止检查之后出现的目录被覆盖。
        for name, body in files.items():
            path = target / name
            path.parent.mkdir(parents=True, exist_ok=True)
            with path.open("x", encoding="utf-8") as stream:
                stream.write(body)
        print(f"{CONFIGS[ide][0]} 项目已创建。用 IDE 的 File → Open 打开文件夹：\n  {target}")

    print("\n不要打开 .idea、.iml 或 current_problem.xml 文件；请打开上面列出的文件夹。")
    print("项目文件不含 SDK 和题目，下面两步必须在 IDE 里手工完成一次。")
    for ide in ides:
        print()
        for line in next_steps(ide):
            print(line)
    print("\n详细步骤：ide-templates/README.md。换题只改运行配置，无须再次执行 make ide。")


def example_problem_dir():
    """挑一个真实的题目目录当示例，没有题解时退回一个明显是占位的名字。

    提示里给具体路径而不是 __PROBLEM_DIR__ 本身，是因为占位符不告诉用户该填
    "leetcode/0094.Binary-Tree-Inorder-Traversal" 还是绝对路径还是包名，只能猜。
    """
    problems = sorted(
        p.name for p in (ROOT / "leetcode").glob("[0-9][0-9][0-9][0-9].*")
        if p.is_dir() and p.name != "0000.Template"
    )
    return f"leetcode/{problems[0]}" if problems else "leetcode/0001.Two-Sum"


def go_module_path():
    match = re.search(r"^module\s+(\S+)$", (ROOT / "go.mod").read_text(encoding="utf-8"), re.M)
    return match[1] if match else "github.com/you/your-repo"


def next_steps(ide):
    problem = example_problem_dir()
    if ide == "goland":
        return [
            "GoLand：",
            "  1. Settings → Go → GOROOT 选一个不低于 go.mod 要求的 Go SDK",
            "  2. Run → Edit Configurations → Go - Current Problem：",
            "     Test kind 选 Package，Pattern 留空，Package path 填当前题目，例如",
            f"       {go_module_path()}/{problem}",
            f"     Directory 同步改成 {ROOT / problem}",
        ]
    if ide == "pycharm":
        return [
            "PyCharm：",
            "  1. Settings → Project → Python Interpreter → Add → Existing，选",
            f"       {ROOT / '.venv' / 'bin' / 'python'}",
            "     （这个虚拟环境由 make init 创建，没有就先跑一次 make init）",
            "  2. Run → Edit Configurations → Python - Current Problem：",
            "     Target 选 Script path，填当前题目的测试文件，例如",
            f"       {ROOT / problem / 'solution_test.py'}",
            "     注意选 Script path 而不是同名的 Module name。",
        ]
    return [
        "IntelliJ IDEA：",
        "  1. File → Project Structure → Project → SDK 选 JDK 17（语言级别已按 javatest.sh 设好）",
        "  2. 在项目树里右键当前这一道题的目录 → Mark Directory as → Sources Root，例如",
        f"       {problem}",
        "     structures/java 已经标好了，不用动；标完 Build → Rebuild Project。",
        "  换题时先右键上一题 → Unmark as Sources Root，再标新的那道：同一时刻只能有",
        "  一个题目目录是 Sources Root，多于一个会撞 Duplicate class Solution。",
        "  IDEA 有时会自己把含 .java 的题目目录标成 Sources Root，留意别同时标上多道。",
    ]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--ide", choices=["all", *CONFIGS], default="all")
    args = parser.parse_args()
    ides = list(CONFIGS) if args.ide == "all" else [args.ide]
    try:
        initialize(ides)
    except (OSError, ValueError, ET.ParseError) as error:
        print(f"IDE 初始化失败：{error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
