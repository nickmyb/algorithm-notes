import java.util.Arrays;
import java.util.List;

/**
 * 题解测试骨架。
 *
 * <p>没有引入 JUnit，测试就是一个 main：check 逐条打印通过与否，断言不成立就抛
 * AssertionError，进程以非零码退出，javatest.sh 据此判定失败。
 *
 * <p>用例照着该题 README 里 ## 题目 那节的 Example 转录，题目给几个就写几个。
 * 第一个参数写用例名（Example 几），和 Go 的 t.Run、pytest 的用例 id 对齐，
 * 三门语言的测试输出看起来是一回事。
 */
public class SolutionTest {

    public static void main(String[] args) {
        // 骨架还没有题解。写完后删掉这行，换成真实断言：
        //
        //     Solution s = new Solution();
        //     check("Example 1", s.solve(new int[] {2, 7, 11, 15}), Arrays.asList(0, 1));
        //
        // 输入是树或链表时用 structures/java 里的辅助建结构：
        //
        //     check("Example 1", s.maxDepth(TreeNodes.build(3, 9, 20, null, null, 15, 7)), 3);
        //
        // 一题写了多种解法时，对每种实现都调一遍 check，用例名带上实现名区分：
        //
        //     check("Example 1/solve", s.solve(...), want);
        //     check("Example 1/solveBruteForce", s.solveBruteForce(...), want);
        System.out.println("SKIP 骨架还没有题解");
    }

    /** 断言一个用例并打印结果。失败时抛 AssertionError，让进程非零退出。 */
    static void check(String name, List<Integer> got, List<Integer> want) {
        if (!got.equals(want)) {
            throw new AssertionError(name + ": got " + got + ", want " + want);
        }
        System.out.println("    PASS  " + name);
    }

    /** 返回值是数组时用这个重载。 */
    static void check(String name, int[] got, int[] want) {
        if (!Arrays.equals(got, want)) {
            throw new AssertionError(
                    name + ": got " + Arrays.toString(got) + ", want " + Arrays.toString(want));
        }
        System.out.println("    PASS  " + name);
    }

    /** 返回值是单个整数时用这个重载。 */
    static void check(String name, int got, int want) {
        if (got != want) {
            throw new AssertionError(name + ": got " + got + ", want " + want);
        }
        System.out.println("    PASS  " + name);
    }
}
