#!/usr/bin/env bash
# 由测试入口 source；保留各语言原生报告，统一最后一行和退出状态。
TEST_STATUS=PASS

test_summary() {
    local rc=${1:-$?}
    [ "$rc" -eq 0 ] || TEST_STATUS=FAIL
    printf '\n===== %s: %s =====\n' "$TEST_LANGUAGE" "$TEST_STATUS"
}
trap test_summary EXIT

# 调用方先切到仓库根。接受仓库相对路径和绝对路径，拼错目录必须失败。
test_dirs=()
for test_dir in "$@"; do
    if [ ! -d "$test_dir" ]; then
        printf '测试目录不存在：%s\n' "$test_dir" >&2
        exit 1
    fi
    test_dirs+=("$(cd "$test_dir" && pwd)")
done
