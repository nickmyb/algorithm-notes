"""本地接线区。

LeetCode 的 Python 编辑器预置了 TreeNode 和 typing.List，本地要显式 import。
"""

from typing import List

from tree_node import TreeNode


# ===== 以下是题解本体，与提交到 LeetCode 的代码一字不差 =====


class Solution:
    def inorderTraversal(self, root: TreeNode | None) -> List[int]:
        return self.inorder(root, [])

    def inorder(self, root: TreeNode | None, input: List[int]) -> List[int]:
        if root == None:
            return input

        output = self.inorder(root.left, input)
        output.append(root.val)
        output = self.inorder(root.right, output)

        return output