import java.util.Arrays;

/** 用例只转录该题 README 英文版的 Example；IDE 直接运行 main 即可。 */
public class SolutionTest {
    public static void main(String[] args) {
        ExampleTests tests = new ExampleTests();
        // 写完题解后删掉 skip，换成真实用例。finish 必须保留。
        //
        // Solution s = new Solution();
        // tests.check("Example 1", () -> s.twoSum(new int[] {2, 7, 11, 15}, 9),
        //         new int[] {0, 1});
        //
        // 树或链表用共享辅助建结构；返回 List、数组或整数都用同一个 check：
        // tests.check("Example 1", () -> s.inorderTraversal(TreeNodes.build(1, null, 2, 3)),
        //         Arrays.asList(1, 3, 2));
        // tests.check("Example 1", () -> s.maxDepth(TreeNodes.build(3, 9, 20, null, null, 15, 7)), 3);
        // 多种解法共用一组用例，用例名带上实现名区分。
        tests.skip("骨架还没有题解");
        tests.finish();
    }
}
