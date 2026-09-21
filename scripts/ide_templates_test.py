"""固定运行模板的回归检查：避免换题后仍指向旧题或同名 Python 模块。"""

import json
from pathlib import Path
import re
import xml.etree.ElementTree as ET

import pytest

ROOT = Path(__file__).resolve().parents[1]
TEMPLATES = ROOT / "ide-templates"
GO_MODULE = re.search(r"^module\s+(\S+)", (ROOT / "go.mod").read_text(), re.M)[1]
CONFIGS = [
    ("goland", "Go", "go", "GoTestRunConfiguration", "Go Test"),
    ("pycharm", "Python", "python", "tests", "py.test"),
    ("idea", "Java", "java", "Application", "Application"),
]


@pytest.mark.parametrize("ide,language,module,kind,factory", CONFIGS)
@pytest.mark.parametrize("problem", [
    "leetcode/0094.Binary-Tree-Inorder-Traversal",
    "leetcode/0001.Two-Sum",
])
def test_current_problem_template(ide, language, module, kind, factory, problem):
    template = (TEMPLATES / ide / "current_problem.xml.tpl").read_text()
    values = {"__PROBLEM_DIR__": problem, "__GO_MODULE__": GO_MODULE}
    for placeholder, value in values.items():
        template = template.replace(placeholder, value)
    assert not re.search(r"__[A-Z_]+__", template)
    root = ET.fromstring(template)
    assert root.attrib["name"] == "ProjectRunConfigurationManager"
    assert len(root.findall("configuration")) == 1
    config = root.find("configuration")
    assert config.attrib["name"] == f"{language} - Current Problem"
    assert config.attrib["type"] == kind
    assert config.attrib["factoryName"] == factory
    assert config.find("module").attrib["name"] == f"algorithm-notes-{module}"
    options = {option.attrib["name"]: option.attrib["value"] for option in config.findall("option")}
    if ide == "goland":
        assert config.find("kind").attrib["value"] == "PACKAGE"
        assert config.find("package").attrib["value"] == f"{GO_MODULE}/{problem}"
        assert config.find("directory").attrib["value"] == f"$PROJECT_DIR$/../../{problem}"
        assert config.find("working_directory").attrib["value"] == "$PROJECT_DIR$/../.."
        assert config.find("pattern").attrib["value"] == ""
    else:
        assert options["WORKING_DIRECTORY"] == "$PROJECT_DIR$/../.."
        if ide == "pycharm":
            assert options["IS_MODULE_SDK"] == "true"
            assert options["SDK_HOME"] == ""
            assert json.loads(options["_new_targetType"]) == "PATH"
            assert json.loads(options["_new_target"]) == f"$PROJECT_DIR$/../../{problem}/solution_test.py"
        else:
            assert options["MAIN_CLASS_NAME"] == "SolutionTest"
            assert config.find("method/option[@name='Make']").attrib["enabled"] == "true"


def test_templates_cannot_be_loaded_as_run_configurations():
    paths = list(TEMPLATES.rglob("*.tpl"))
    assert len(paths) == 3
    assert not list(TEMPLATES.rglob("*.xml"))
    for path in paths:
        template = path.read_text()
        assert ".idea" not in path.parts
        assert "__ID__" not in template
        assert "0000.Template" not in template
        assert "0094" not in template
        assert "/home/" not in template
        assert "/Users/" not in template
