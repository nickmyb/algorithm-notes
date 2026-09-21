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
    print("尚未选择 SDK/解释器和题目。首次打开后先完成这两项，再运行 Current Problem。")
    print("Go/Python：在 Edit Configurations 中替换 __PROBLEM_DIR__ 目标。")
    print("Java：将当前一道题标为 Sources Root，然后 Rebuild Project。")
    print("详细步骤：ide-templates/README.md。换题无须再次执行 make ide。")


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
