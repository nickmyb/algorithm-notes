import java.io.PrintStream;
import java.util.Arrays;
import java.util.Objects;
import java.util.function.Supplier;

/** Java 题解测试的用例报告器；不依赖 JUnit，IDE 直接运行 main 也有汇总。 */
public final class ExampleTests {
    private final PrintStream out;
    private int passed;
    private int failed;
    private int skipped;

    public ExampleTests() {
        this(System.out);
    }

    // 共享结构测试用独立输出流验证报告，不污染正常用例的输出。
    ExampleTests(PrintStream out) {
        this.out = out;
    }

    /** 延迟求值，以便题解抛异常时也能记录失败，并继续跑剩余用例。 */
    public void check(String name, Supplier<?> solution, Object want) {
        try {
            Object got = solution.get();
            if (!Objects.deepEquals(got, want)) {
                throw new AssertionError("got " + display(got) + ", want " + display(want));
            }
            passed++;
            out.println("    PASS  " + name);
        } catch (AssertionError | RuntimeException error) {
            failed++;
            out.println("    FAIL  " + name);
            error.printStackTrace(out);
        }
    }

    public void skip(String name) {
        skipped++;
        out.println("    SKIP  " + name);
    }

    /** 失败时在汇总之后抛异常，保证 IDE 和脚本都得到非零退出码。 */
    public void finish() {
        String status = failed > 0 ? "FAIL" : passed > 0 ? "PASS" : "SKIP";
        out.printf("  %s: %d passed, %d failed, %d skipped%n",
                status, passed, failed, skipped);
        if (failed > 0) {
            throw new AssertionError(failed + " test(s) failed");
        }
    }

    private static String display(Object value) {
        String wrapped = Arrays.deepToString(new Object[] {value});
        return wrapped.substring(1, wrapped.length() - 1);
    }
}
