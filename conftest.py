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
import sys
from pathlib import Path

import pytest

SOLUTION_FILENAME = "solution.py"

# 共享数据结构（TreeNode、ListNode 等）放在 structures/python，加进 sys.path 后
# 题解里直接 `from tree_node import TreeNode, build_tree` 即可，不用每道树的题重抄一遍。
# 对应 Go 的 structures/go 包和 Java 的 structures/java 目录。
#
# 这里只加这一个路径，题解模块本身仍然由下面的 fixture 按文件路径加载——
# 那是为了绕开几百个同名 solution.py 的模块名冲突，不能改成普通 import。
#
# 用 append 而不是 insert(0)：插到最前面的话，这个目录会排在标准库之前，
# 哪天在里面放个 queue.py 或 heapq.py 就会把整个测试进程的标准库遮蔽掉。
# 追加到末尾则是标准库优先，重名时最多是自己的模块取不到，立刻就能发现。
_STRUCTURES = Path(__file__).parent / "structures" / "python"
if _STRUCTURES.is_dir() and str(_STRUCTURES) not in sys.path:
    sys.path.append(str(_STRUCTURES))


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
