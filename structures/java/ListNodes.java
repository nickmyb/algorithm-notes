import java.util.ArrayList;
import java.util.List;

/**
 * 单链表的测试辅助方法，对应 structures/go/ListNode.go 里的同名函数：
 *
 * <pre>
 * Go                     Java
 * Ints2List          ->  ListNodes.of
 * List2Ints          ->  ListNodes.toList
 * Ints2ListWithCycle ->  ListNodes.withCycle
 * GetNodeWith        ->  ListNodes.find
 * </pre>
 */
public final class ListNodes {

    /**
     * toList 的遍历上限。链表成环时不加限制会死循环；这里照搬 Go 版的做法，
     * 超出就抛异常，把"测试挂住"变成"测试报错"。
     */
    private static final int LIMIT = 100;

    private ListNodes() {}

    /** 按数组建链表，空数组返回 null。 */
    public static ListNode of(int... vals) {
        ListNode dummy = new ListNode();
        ListNode cur = dummy;
        for (int v : vals) {
            cur.next = new ListNode(v);
            cur = cur.next;
        }
        return dummy.next;
    }

    /** 链表展开成列表。遇到环会抛异常而不是死循环。 */
    public static List<Integer> toList(ListNode head) {
        List<Integer> out = new ArrayList<>();
        int times = 0;
        while (head != null) {
            if (++times > LIMIT) {
                throw new IllegalStateException(
                        "链条深度超过 " + LIMIT + "，可能有环。检查题解，或放宽 ListNodes.LIMIT。");
            }
            out.add(head.val);
            head = head.next;
        }
        return out;
    }

    /**
     * 建一个尾部指回第 pos 个节点的链表（头节点下标为 0），pos 为 -1 表示无环。
     * 用于第 141、142 题这类判环的题目。
     */
    public static ListNode withCycle(int[] vals, int pos) {
        ListNode head = of(vals);
        if (pos == -1 || head == null) {
            return head;
        }

        ListNode entry = head;
        for (int i = 0; i < pos; i++) {
            entry = entry.next;
        }
        ListNode tail = entry;
        while (tail.next != null) {
            tail = tail.next;
        }
        tail.next = entry;
        return head;
    }

    /** 找出第一个 val 等于 target 的节点，找不到返回 null。 */
    public static ListNode find(ListNode head, int target) {
        while (head != null && head.val != target) {
            head = head.next;
        }
        return head;
    }
}
