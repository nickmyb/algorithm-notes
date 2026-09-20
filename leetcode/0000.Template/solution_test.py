"""题解测试骨架。

参数里的 solution 是仓库根目录 conftest.py 提供的 fixture：它按路径加载同目录的
solution.py，模块名带上了题目目录名。所以每道题的文件都可以叫 solution.py，
不会像普通 import 那样在 sys.modules 里互相顶替。

题解是 class Solution（和 LeetCode 的模板一致，便于原样复制），所以要先实例化：
solution.Solution()。
"""

import pytest

# 一题有多种解法时，把方法名都列进来，同一组用例跑全部实现，
# 顺带保证它们结果一致。
IMPLS = ["solve"]


@pytest.mark.parametrize("impl", IMPLS)
def test_solve(solution, impl):
    pytest.skip("骨架还没有题解，写完后删掉这一行")

    solve = getattr(solution.Solution(), impl)
    assert solve([2, 7, 11, 15]) == [0, 1]
