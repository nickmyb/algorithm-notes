import java.util.Arrays;

/**
 * 题解测试骨架。
 *
 * <p>没有引入 JUnit，测试就是一个 main：断言不成立就抛 AssertionError，
 * 进程以非零码退出，javatest.sh 据此判定失败。
 *
 * <p>一题有多种解法时，在 main 里对每种实现都调一遍 check。
 */
public class SolutionTest {

    public static void main(String[] args) {
        // 骨架还没有题解。写完后把下面这行换成真实断言，例如：
        // Solution s = new Solution();
        // check("twoSum", s.twoSum(new int[] {2, 7, 11, 15}, 9), new int[] {0, 1});
        // check("twoSumBruteForce", s.twoSumBruteForce(new int[] {2, 7, 11, 15}, 9), new int[] {0, 1});
        System.out.println("SKIP 骨架还没有题解");
    }

    static void check(String impl, int[] got, int[] want) {
        if (!Arrays.equals(got, want)) {
            throw new AssertionError(
                    impl + " = " + Arrays.toString(got) + ", want " + Arrays.toString(want));
        }
    }
}
