"""本地接线区。

LeetCode 的 Python 编辑器预置了节点类型和 typing，本地要显式 import，写在下面：

    from typing import List, Optional
    from tree_node import TreeNode, build_tree
    from list_node import ListNode

conftest.py 已经把 structures/python 加进了 sys.path，直接 import 即可。
"""

# ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====


class Solution:
    """题解骨架。

    类名和方法签名保持 LeetCode 给的形态（class Solution + 带 self 的方法），
    这样题解可以在本地和提交框之间原样来回复制。把 solve 改成本题要求的方法名，
    例如 Two Sum 是 def twoSum(self, nums, target)。

    同一题的多种解法写成本类里的多个方法，入口方法保持 LeetCode 给的签名名，
    其余解法加后缀说明算法：twoSum / twoSumBruteForce / twoSumTwoPointers。
    """

    def solve(self, nums):
        return []
