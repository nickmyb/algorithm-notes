R"""94. Binary Tree Inorder Traversal 的测试。

题解本体是 class Solution（和 LeetCode 模板一致），所以要先实例化。
"""

import pytest
from tree_node import build_tree

# 题目英文版给的 4 个示例，输入直接照抄题面里的层序数组。
# 中文版只有 3 个，缺了 Example 2 那棵大树——官方翻译滞后，以英文版为准。
#
# 每行第一列是用例名，它会成为 pytest 的用例 id（见文件末尾的 ids=），
# 这样 -v 输出里显示的是 [Example 1] 而不是默认的 [vals0-want0]。
CASES = [
    ("Example 1", [1, None, 2, 3], [1, 3, 2]),
    ("Example 2", [1, 2, 3, 4, 5, None, 8, None, None, 6, 7, 9], [4, 2, 6, 5, 7, 1, 3, 9, 8]),
    ("Example 3 空树", [], []),
    ("Example 4 单节点", [1], [1]),
]


@pytest.mark.parametrize(
    "name,vals,want", CASES, ids=[c[0] for c in CASES]
)
def test_inorder_traversal(solution, name, vals, want):
    got = solution.Solution().inorderTraversal(build_tree(vals))
    assert got == want, f"{name}: inorderTraversal({vals}) = {got}, want {want}"
