#!/usr/bin/env bash
#
# 逐个题目目录编译并运行 Java 题解测试。
#
# 为什么要逐目录：每道题的题解类都叫 Solution（放在 default package），
# 一次性 `javac leetcode/*/*.java` 会直接报 "duplicate class: Solution"。
# 按目录编译到各自的 out/ 子目录就不会冲突。
#
# 测试没有引入 JUnit，就是一个 main：断言不成立抛 AssertionError，
# 进程非零退出，这里据此判失败。
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT="$ROOT/out/java"

# 题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
ROOTS=(leetcode)

# 目标 Java 版本。--release 同时约束语言特性和可用 API，所以哪怕本地装的是
# JDK 21，写了 record 或 switch 模式匹配也会在本地就报错，而不是推上去才被 CI 拦下。
# 这里是版本的唯一来源，CI 只要装的 JDK 不低于它即可。
JAVA_RELEASE="${JAVA_RELEASE:-17}"

rm -rf "$OUT"

total=0
failed=0
failed_dirs=()

for root in "${ROOTS[@]}"; do
    for dir in "$ROOT/$root"/*/; do
        [ -d "$dir" ] || continue
        # shellcheck disable=SC2206
        sources=("$dir"*.java)
        [ -e "${sources[0]}" ] || continue

        name="$(basename "$dir")"
        classes="$OUT/$root/$name"
        mkdir -p "$classes"
        total=$((total + 1))

        if ! javac --release "$JAVA_RELEASE" -encoding UTF-8 -d "$classes" "${sources[@]}"; then
            echo "COMPILE FAIL  $root/$name"
            failed=$((failed + 1))
            failed_dirs+=("$root/$name")
            continue
        fi

        # 没有测试类的目录只验证能编译通过
        if [ ! -f "$classes/SolutionTest.class" ]; then
            echo "COMPILE OK    $root/$name (无 SolutionTest)"
            continue
        fi

        if java -cp "$classes" SolutionTest; then
            echo "PASS          $root/$name"
        else
            echo "TEST FAIL     $root/$name"
            failed=$((failed + 1))
            failed_dirs+=("$root/$name")
        fi
    done
done

echo "===== Java: $total 个目录，$failed 个失败 ====="
if [ "$failed" -gt 0 ]; then
    printf '  %s\n' "${failed_dirs[@]}"
    exit 1
fi
