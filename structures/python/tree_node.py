"""二叉树节点和测试辅助，对应 structures/go/TreeNode.go。

TreeNode 的字段照搬 LeetCode Python 编辑器里给的那段注释（Definition for a binary
tree node），一个字不改，题解可以在本地和提交框之间原样复制。下面的函数是只有本地
测试才需要的辅助，LeetCode 上没有：

    Go                  Python
    Ints2TreeNode   ->  build_tree
    Tree2ints       ->  tree_to_level_order
    Tree2Preorder   ->  preorder
    Tree2Inorder    ->  inorder
    Tree2Postorder  ->  postorder
    GetTargetNode   ->  find
    Equal           ->  equals

Go 版用哨兵 NULL = -1<<63 表示空节点（[]int 装不了 nil），Python 有真正的 None，
所以调用写出来和题目给的输入完全一致：题目是 [1,null,2,3]，代码就是
build_tree([1, None, 2, 3])。

conftest.py 已经把本目录加进 sys.path，题解里直接 import 即可：

    from tree_node import TreeNode, build_tree
"""

from collections import deque


class TreeNode:
    """LeetCode 的二叉树节点定义。"""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def build_tree(vals):
    """按 LeetCode 的层序数组建树，None 表示空节点。"""
    if not vals or vals[0] is None:
        return None

    root = TreeNode(vals[0])
    queue = deque([root])
    i = 1
    while queue and i < len(vals):
        node = queue.popleft()
        # 数组只为非空节点列出孩子，所以每出队一个节点就依次消费两个位置
        if i < len(vals):
            v = vals[i]
            i += 1
            if v is not None:
                node.left = TreeNode(v)
                queue.append(node.left)
        if i < len(vals):
            v = vals[i]
            i += 1
            if v is not None:
                node.right = TreeNode(v)
                queue.append(node.right)
    return root


def tree_to_level_order(root):
    """层序展开，空节点记为 None，末尾的 None 裁掉，便于和题目给的数组直接比对。"""
    if root is None:
        return []

    out = []
    queue = deque([root])
    while queue:
        node = queue.popleft()
        if node is None:
            out.append(None)
            continue
        out.append(node.val)
        queue.append(node.left)
        queue.append(node.right)

    while out and out[-1] is None:
        out.pop()
    return out


def preorder(root):
    """前序遍历。"""
    if root is None:
        return []
    return [root.val] + preorder(root.left) + preorder(root.right)


def inorder(root):
    """中序遍历。"""
    if root is None:
        return []
    return inorder(root.left) + [root.val] + inorder(root.right)


def postorder(root):
    """后序遍历。"""
    if root is None:
        return []
    return postorder(root.left) + postorder(root.right) + [root.val]


def find(root, target):
    """找出第一个 val 等于 target 的节点，找不到返回 None。层序查找，和 Go 版一致。"""
    if root is None:
        return None

    queue = deque([root])
    while queue:
        node = queue.popleft()
        if node is None:
            continue
        if node.val == target:
            return node
        queue.append(node.left)
        queue.append(node.right)
    return None


def equals(a, b):
    """判断两棵树结构和取值是否完全相同。"""
    if a is None or b is None:
        return a is b
    return a.val == b.val and equals(a.left, b.left) and equals(a.right, b.right)
