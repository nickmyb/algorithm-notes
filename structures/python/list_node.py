"""单链表节点和测试辅助，对应 structures/go/ListNode.go。

ListNode 的字段照搬 LeetCode Python 编辑器里给的那段注释（Definition for singly-linked
list），一个字不改。下面的函数是只有本地测试才需要的辅助：

    Go                     Python
    Ints2List          ->  build_list
    List2Ints          ->  list_to_values
    Ints2ListWithCycle ->  build_list_with_cycle
    GetNodeWith        ->  find

conftest.py 已经把本目录加进 sys.path，题解里直接 import 即可：

    from list_node import ListNode, build_list
"""

# list_to_values 的遍历上限。链表成环时不加限制会死循环；这里照搬 Go 版的做法，
# 超出就抛异常，把"测试挂住"变成"测试报错"。
LIMIT = 100


class ListNode:
    """LeetCode 的单链表节点定义。"""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def build_list(vals):
    """按数组建链表，空数组返回 None。"""
    dummy = ListNode()
    cur = dummy
    for v in vals:
        cur.next = ListNode(v)
        cur = cur.next
    return dummy.next


def list_to_values(head):
    """链表展开成列表。遇到环会抛异常而不是死循环。"""
    out = []
    times = 0
    while head is not None:
        times += 1
        if times > LIMIT:
            raise RuntimeError(
                f"链条深度超过 {LIMIT}，可能有环。检查题解，或放宽 list_node.LIMIT。"
            )
        out.append(head.val)
        head = head.next
    return out


def build_list_with_cycle(vals, pos):
    """建一个尾部指回第 pos 个节点的链表（头节点下标为 0），pos 为 -1 表示无环。

    用于第 141、142 题这类判环的题目。
    """
    head = build_list(vals)
    if pos == -1 or head is None:
        return head

    entry = head
    for _ in range(pos):
        entry = entry.next
    tail = entry
    while tail.next is not None:
        tail = tail.next
    tail.next = entry
    return head


def find(head, target):
    """找出第一个 val 等于 target 的节点，找不到返回 None。"""
    while head is not None and head.val != target:
        head = head.next
    return head
