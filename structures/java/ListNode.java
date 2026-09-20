/**
 * 单链表节点。
 *
 * <p>字段和构造器照搬 LeetCode Java 编辑器里给的那段注释（Definition for singly-linked list），
 * 一个字不改，题解可以在本地和提交框之间原样复制。
 *
 * <p>建链表、展开成数组这类只有本地测试才需要的辅助方法放在 {@link ListNodes}。
 */
public class ListNode {
    int val;
    ListNode next;

    ListNode() {}

    ListNode(int val) {
        this.val = val;
    }

    ListNode(int val, ListNode next) {
        this.val = val;
        this.next = next;
    }
}
