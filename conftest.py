"""pytest 的全局配置，给每道题的测试提供 solution fixture。

为什么需要它：Python 的模块名就是文件名，几百道题都叫 solution.py / solution_test.py
的话，普通 import 会在 sys.modules 里互相顶替 —— pytest 默认的 prepend 模式会直接
报 "import file mismatch" 中断收集，换成 importlib 模式虽然能收集，但测试里
`from solution import ...` 又会 ModuleNotFoundError（importlib 模式不把题目目录
加进 sys.path）。

这里的做法是按文件路径加载模块，并用题目目录名当模块名，从根本上避免重名：
每个测试通过 solution fixture 拿到同目录的题解模块，文件名就能在全仓库统一。
"""

import importlib.util
from pathlib import Path

import pytest

SOLUTION_FILENAME = "solution.py"


@pytest.fixture
def solution(request):
    """加载与测试文件同目录的 solution.py，返回模块对象。"""
    path = Path(request.path).parent / SOLUTION_FILENAME
    if not path.exists():
        pytest.fail(f"同目录下没有 {SOLUTION_FILENAME}: {path.parent}")

    # 模块名带上题目目录名（如 sol::0001.Two-Sum），保证全仓库唯一。
    # 目录名里有点和横杠，不是合法标识符，但这里的名字只是 sys.modules 的键，
    # 不会被 import 语句用到，所以没关系。
    spec = importlib.util.spec_from_file_location(f"sol::{path.parent.name}", path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module
