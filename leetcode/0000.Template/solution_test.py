"""题解测试骨架。

参数里的 solution 是仓库根目录 conftest.py 提供的 fixture：它按路径加载同目录的
solution.py，模块名带上了题目目录名。所以每道题的文件都可以叫 solution.py，
不会像普通 import 那样在 sys.modules 里互相顶替。

题解本体是 class Solution（和 LeetCode 的模板一致，便于原样复制），所以要先
实例化：solution.Solution()。
"""

import pytest

# 一行一个用例，三列分别是用例名、输入、期望输出。
#
# 第一列会通过下面的 ids= 成为 pytest 的用例 id，-v 输出里显示 [Example 1]
# 而不是默认的 [vals0-want0]（那个完全看不出是哪个用例）。
#
# 用例照着该题 README 里 ## 题目 那节的 Example 转录，题目给几个就写几个。
CASES = [
    # ("Example 1", [2, 7, 11, 15], [0, 1]),
]


@pytest.mark.parametrize("name,in_,want", CASES, ids=[c[0] for c in CASES])
def test_solve(solution, name, in_, want):
    pytest.skip("骨架还没有题解，写完后删掉这一行")

    got = solution.Solution().solve(in_)
    assert got == want, f"{name}: solve({in_}) = {got}, want {want}"


# 输入是树或链表时，用 structures/python 里的辅助建结构：
#
#     from tree_node import build_tree
#     got = solution.Solution().maxDepth(build_tree(vals))
#
# 一题写了多种解法时，把方法名也参数化，同一组用例跑全部实现，
# 顺带保证几种解法结果一致：
#
#     @pytest.mark.parametrize("impl", ["solve", "solveBruteForce"])
#     @pytest.mark.parametrize("name,in_,want", CASES, ids=[c[0] for c in CASES])
#     def test_solve(solution, impl, name, in_, want):
#         got = getattr(solution.Solution(), impl)(in_)
