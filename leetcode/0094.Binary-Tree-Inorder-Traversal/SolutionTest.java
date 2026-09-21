import java.util.Arrays;

public class SolutionTest {

    public static void main(String[] args) {
        ExampleTests tests = new ExampleTests();
        Solution s = new Solution();

        // 题目英文版给的 4 个示例，原样转录。
        // 中文版只有 3 个，缺了 Example 2 那棵大树——官方翻译滞后，以英文版为准。
        tests.check("Example 1",
                () -> s.inorderTraversal(TreeNodes.build(1, null, 2, 3)),
                Arrays.asList(1, 3, 2));
        tests.check("Example 2",
                () -> s.inorderTraversal(TreeNodes.build(1, 2, 3, 4, 5, null, 8, null, null, 6, 7, 9)),
                Arrays.asList(4, 2, 6, 5, 7, 1, 3, 9, 8));
        tests.check("Example 3 空树", () -> s.inorderTraversal(TreeNodes.build()), Arrays.asList());
        tests.check("Example 4 单节点", () -> s.inorderTraversal(TreeNodes.build(1)), Arrays.asList(1));
        tests.finish();
    }
}
