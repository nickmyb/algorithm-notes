"""list_node 里各辅助方法的测试，对应 structures/go/ListNode_test.go。"""

import pytest
from list_node import build_list, build_list_with_cycle, find, list_to_values


@pytest.mark.parametrize("vals", [[], [1], [1, 2, 3]])
def test_build_to_values_roundtrip(vals):
    assert list_to_values(build_list(vals)) == vals


def test_build_empty():
    assert build_list([]) is None


def test_long_list_with_repeated_values():
    # 抓住旧的 100 节点上限；同时防止把值相等误当成访问了同一个节点。
    vals = [0] * 101
    assert list_to_values(build_list(vals)) == vals


def test_find():
    head = build_list([1, 2, 3])
    assert find(head, 2).val == 2
    assert find(head, 9) is None
    assert find(None, 1) is None


def test_cycle_raises_instead_of_hanging():
    # 第 142 题的示例：尾节点指回下标 1
    cyclic = build_list_with_cycle([3, 2, 0, -4], 1)
    with pytest.raises(RuntimeError):
        list_to_values(cyclic)


def test_no_cycle_when_pos_is_minus_one():
    head = build_list_with_cycle([1, 2], -1)
    assert list_to_values(head) == [1, 2]


def test_cycle_on_empty_list():
    assert build_list_with_cycle([], 0) is None
