#!/usr/bin/env bash
#
# 逐个题目目录编译并运行 Java 题解测试。
#
#   bash ./javatest.sh                                  全部题目
#   bash ./javatest.sh leetcode/0094.Binary-Tree-...    只测指定目录
#
# 为什么要逐目录：每道题的题解类都叫 Solution（放在 default package），
# 一次性 `javac leetcode/*/*.java` 会直接报 "duplicate class: Solution"。
# 按目录编译到各自的 out/ 子目录就不会冲突。
#
# 测试没有引入 JUnit，就是一个 main：断言不成立抛 AssertionError，
# 进程非零退出，这里据此判失败。
set -eu

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT="$ROOT/out/java"
cd "$ROOT"
TEST_LANGUAGE=Java
source ./scripts/test-common.sh

# 题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
ROOTS=(leetcode)

# 目标 Java 版本。--release 同时约束语言特性和可用 API，所以哪怕本地装的是
# JDK 21，写了 Java 21 的 switch 模式匹配也会在本地就报错，而不是推上去才被 CI 拦下。
# 这里是版本的唯一来源，CI 只要装的 JDK 不低于它即可。
JAVA_RELEASE="${JAVA_RELEASE:-17}"

# 共享数据结构（TreeNode、ListNode 等），编译每个题目目录时一并加进源文件列表，
# 这样题解里直接用 TreeNode 就行，不用每道树的题重抄一遍。
# 对应 Go 的 structures/go 包。因为是逐目录分别编译，同一份源文件在各目录的
# 产物里各有一份，不会互相冲突。
SHARED_DIR="$ROOT/structures/java"
shared=()
if [ -d "$SHARED_DIR" ]; then
    shared=("$SHARED_DIR"/*.java)
    # 目录存在但为空时，glob 不展开，会留下一个字面量路径，这里清掉
    [ -e "${shared[0]}" ] || shared=()
fi

# 要测哪些目录：给了参数就只测这些，否则遍历所有题解根目录
dirs=()
if [ $# -gt 0 ]; then
    dirs=("${test_dirs[@]}")
else
    for root in "${ROOTS[@]}"; do
        for dir in "$ROOT/$root"/*/; do
            [ -d "$dir" ] || continue
            dirs+=("${dir%/}")
        done
    done
fi

# 每次用独立的编译目录，单题模式也不能复用上次的 SolutionTest.class。
# 否则删掉或改名后的测试仍会运行，甚至与另一次并行测试串在一起。
mkdir -p "$OUT"
build_dir="$(mktemp -d "$OUT/run.XXXXXX")"
finish_java() {
    local rc=$?
    rm -rf -- "$build_dir"
    test_summary "$rc"
}
trap finish_java EXIT

total=0
failed=0
failed_dirs=()

for dir in "${dirs[@]}"; do
    [ -d "$dir" ] || continue
    # shellcheck disable=SC2206
    sources=("$dir"/*.java)
    # 那道题没写 Java 题解不算失败，跳过
    [ -e "${sources[0]}" ] || continue

    name="${dir#"$ROOT"/}"
    classes="$build_dir/$total"
    mkdir -p "$classes"
    total=$((total + 1))

    if ! javac --release "$JAVA_RELEASE" -encoding UTF-8 -d "$classes" \
        "${sources[@]}" ${shared[@]+"${shared[@]}"}; then
        echo "COMPILE FAIL  $name"
        failed=$((failed + 1))
        failed_dirs+=("$name")
        continue
    fi

    # 没有测试类的目录只验证能编译通过
    if [ ! -f "$classes/SolutionTest.class" ]; then
        echo "COMPILE OK    $name (无 SolutionTest)"
        continue
    fi

    if java -cp "$classes" SolutionTest; then
        echo "PASS          $name"
    else
        echo "TEST FAIL     $name"
        failed=$((failed + 1))
        failed_dirs+=("$name")
    fi
done

# 共享结构自己的测试。这些辅助方法是所有树/链表题的地基，写错了会让题解测试给出
# 假结果，所以必须跟着跑。只在全量模式下跑——单题模式的目的是快速迭代一道题。
if [ $# -eq 0 ] && [ -f "$SHARED_DIR/test/StructuresTest.java" ]; then
    total=$((total + 1))
    classes="$build_dir/structures"
    mkdir -p "$classes"
    if javac --release "$JAVA_RELEASE" -encoding UTF-8 -d "$classes" \
        ${shared[@]+"${shared[@]}"} "$SHARED_DIR/test/StructuresTest.java" \
        && java -cp "$classes" StructuresTest; then
        echo "PASS          structures/java"
    else
        echo "TEST FAIL     structures/java"
        failed=$((failed + 1))
        failed_dirs+=("structures/java")
    fi
fi

if [ "$total" -eq 0 ] && [ "$failed" -eq 0 ]; then
    echo "跳过 Java：指定的题目没有 Java 题解"
    TEST_STATUS=SKIP
    exit 0
fi

echo "Java: $total 个测试/编译目录，$failed 个失败"
if [ "$failed" -gt 0 ]; then
    printf '  %s\n' "${failed_dirs[@]}"
    exit 1
fi
