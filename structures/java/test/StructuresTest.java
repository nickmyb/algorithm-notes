import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.util.Arrays;
import java.util.List;

/**
 * structures/java 里各辅助方法的测试，对应 structures/go 下的那些 *_test.go。
 *
 * <p>放在 test/ 子目录而不是和被测代码同级：javatest.sh 取共享结构用的是
 * {@code structures/java/*.java} 这个非递归 glob，放同级的话这个测试类会被编进
 * 每一道题的产物里。
 *
 * <p>由 javatest.sh 在全量模式下单独编译运行（单题模式不跑）。
 *
 * <p>这些辅助方法是所有树/链表题的地基，写错了会让题解测试给出假结果——
 * 初版 toLevelOrder 就因为往 ArrayDeque 里塞 null 子节点而 NPE，是这个测试抓到的。
 */
public class StructuresTest {

    public static void main(String[] args) {
        testTreeNodes();
        testListNodes();
        testExampleTests();
        System.out.println("  structures/java 全部通过");
    }

    static void testTreeNodes() {
        // 第 94 题的示例
        TreeNode t = TreeNodes.build(1, null, 2, 3);
        eq("inorder", TreeNodes.inorder(t), Arrays.asList(1, 3, 2));
        eq("preorder", TreeNodes.preorder(t), Arrays.asList(1, 2, 3));
        eq("postorder", TreeNodes.postorder(t), Arrays.asList(3, 2, 1));
        eq("toLevelOrder 往返", TreeNodes.toLevelOrder(t), Arrays.asList(1, null, 2, 3));

        eq("空树", TreeNodes.toLevelOrder(TreeNodes.build()), Arrays.asList());
        eq("单节点", TreeNodes.toLevelOrder(TreeNodes.build(1)), Arrays.asList(1));
        if (TreeNodes.build((Integer) null) != null) {
            throw new AssertionError("首元素为 null 应该建出空树");
        }

        // 带空洞的大树，往返必须一致
        List<Integer> vals = Arrays.asList(1, 2, 3, 4, 5, null, 8, null, null, 6, 7, 9);
        TreeNode big = TreeNodes.build(vals.toArray(new Integer[0]));
        eq("大树 inorder", TreeNodes.inorder(big), Arrays.asList(4, 2, 6, 5, 7, 1, 3, 9, 8));
        eq("大树往返", TreeNodes.toLevelOrder(big), vals);

        if (TreeNodes.find(big, 7) == null || TreeNodes.find(big, 7).val != 7) {
            throw new AssertionError("find 没找到存在的值");
        }
        if (TreeNodes.find(big, 99) != null) {
            throw new AssertionError("find 对不存在的值应返回 null");
        }
        if (TreeNodes.find(null, 1) != null) {
            throw new AssertionError("find(null) 应返回 null");
        }

        if (!TreeNodes.equals(TreeNodes.build(1, 2), TreeNodes.build(1, 2))) {
            throw new AssertionError("equals 对相同的树应为 true");
        }
        if (TreeNodes.equals(TreeNodes.build(1, 2), TreeNodes.build(1, null, 2))) {
            throw new AssertionError("equals 应该区分左右子树");
        }
        if (!TreeNodes.equals(null, null)) {
            throw new AssertionError("equals(null, null) 应为 true");
        }
    }

    static void testListNodes() {
        eq("链表往返", ListNodes.toList(ListNodes.of(1, 2, 3)), Arrays.asList(1, 2, 3));
        eq("空链表", ListNodes.toList(ListNodes.of()), Arrays.asList());

        // 超过旧上限 100 的正常链表必须可展开；重复值不代表环。
        eq("101 个重复值节点", ListNodes.toList(ListNodes.of(new int[101])),
                java.util.Collections.nCopies(101, 0));

        if (ListNodes.find(ListNodes.of(1, 2, 3), 2).val != 2) {
            throw new AssertionError("find 没找到存在的值");
        }
        if (ListNodes.find(ListNodes.of(1, 2, 3), 9) != null) {
            throw new AssertionError("find 对不存在的值应返回 null");
        }

        // 成环时 toList 必须抛异常而不是死循环
        ListNode cyclic = ListNodes.withCycle(new int[] {3, 2, 0, -4}, 1);
        try {
            ListNodes.toList(cyclic);
            throw new AssertionError("成环链表 toList 应该抛异常");
        } catch (IllegalStateException expected) {
            // 预期行为
        }

        if (ListNodes.withCycle(new int[] {1, 2}, -1).next.next != null) {
            throw new AssertionError("pos = -1 不应该成环");
        }
        if (ListNodes.withCycle(new int[] {}, 0) != null) {
            throw new AssertionError("空数组应返回 null");
        }
    }

    static void testExampleTests() {
        ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        ExampleTests passing = new ExampleTests(new PrintStream(buffer));
        passing.check("数组按值比较", () -> new int[] {1, 2}, new int[] {1, 2});
        passing.check("列表", () -> Arrays.asList(1, 2), Arrays.asList(1, 2));
        passing.check("整数", () -> 3, 3);
        passing.check("null", () -> null, null);
        passing.finish();
        if (!buffer.toString().contains("PASS: 4 passed, 0 failed, 0 skipped")) {
            throw new AssertionError("成功汇总错误: " + buffer);
        }

        buffer.reset();
        ExampleTests failing = new ExampleTests(new PrintStream(buffer));
        failing.check("值不同", () -> new int[] {1}, new int[] {2});
        failing.check("题解抛异常", () -> { throw new IllegalStateException("boom"); }, 0);
        failing.check("失败之后继续跑", () -> 1, 1);
        boolean rejected = false;
        try {
            failing.finish();
        } catch (AssertionError expected) {
            rejected = true;
        }
        if (!rejected || !buffer.toString().contains("FAIL: 1 passed, 2 failed, 0 skipped")) {
            throw new AssertionError("失败不能被报成通过: " + buffer);
        }

        buffer.reset();
        ExampleTests skeleton = new ExampleTests(new PrintStream(buffer));
        skeleton.skip("骨架");
        skeleton.finish();
        if (!buffer.toString().contains("SKIP: 0 passed, 0 failed, 1 skipped")) {
            throw new AssertionError("骨架不能被报成用例通过: " + buffer);
        }
    }

    static void eq(String name, List<Integer> got, List<Integer> want) {
        if (!got.equals(want)) {
            throw new AssertionError(name + " = " + got + ", want " + want);
        }
    }
}
