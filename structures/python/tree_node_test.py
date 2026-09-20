"""tree_node 里各辅助方法的测试，对应 structures/go/TreeNode_test.go。

这些辅助是所有树题的地基，写错了会让题解测试给出假结果，所以必须跟着跑。
pytest.ini 的 testpaths 包含了本目录。
"""

import pytest
from tree_node import (
    build_tree,
    equals,
    find,
    inorder,
    postorder,
    preorder,
    tree_to_level_order,
)

# 第 94 题的示例
SMALL = [1, None, 2, 3]
# 带空洞的大树
BIG = [1, 2, 3, 4, 5, None, 8, None, None, 6, 7, 9]


@pytest.mark.parametrize(
    "traverse,want",
    [
        (inorder, [1, 3, 2]),
        (preorder, [1, 2, 3]),
        (postorder, [3, 2, 1]),
    ],
)
def test_traversals(traverse, want):
    assert traverse(build_tree(SMALL)) == want


@pytest.mark.parametrize("vals", [[], [1], SMALL, BIG])
def test_build_level_order_roundtrip(vals):
    assert tree_to_level_order(build_tree(vals)) == vals


def test_big_tree_inorder():
    assert inorder(build_tree(BIG)) == [4, 2, 6, 5, 7, 1, 3, 9, 8]


def test_build_empty():
    assert build_tree([]) is None
    assert build_tree([None]) is None


def test_find():
    big = build_tree(BIG)
    assert find(big, 7).val == 7
    assert find(big, 99) is None
    assert find(None, 1) is None


def test_equals():
    assert equals(build_tree([1, 2]), build_tree([1, 2]))
    assert equals(None, None)
    # 左右子树的位置必须区分开
    assert not equals(build_tree([1, 2]), build_tree([1, None, 2]))
    assert not equals(build_tree([1]), None)
