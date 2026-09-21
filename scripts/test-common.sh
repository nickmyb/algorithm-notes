#!/usr/bin/env bash
# 由各语言的测试入口 source：统一开始横幅、结束结论和退出状态。
#
# 三门语言的原生输出格式差异很大（go test 的 ok/FAIL、pytest 的进度条、
# javac + 自写 main 的逐条 PASS），混在一起很难一眼看出"跑到哪了、谁挂了"。
# 这里不改各自的原生输出，只在外面套一层统一的框。
TEST_STATUS=PASS

# 终端才上色，重定向到文件或管道时输出纯文本，免得 CI 日志里全是转义序列。
if [ -t 1 ]; then
    _c_head=$'\033[36m'; _c_pass=$'\033[32m'; _c_fail=$'\033[31m'; _c_off=$'\033[0m'
else
    _c_head=''; _c_pass=''; _c_fail=''; _c_off=''
fi
_rule='────────────────────────────────────────────────────────'
_rule_all='════════════════════════════════════════════════════════'

# 总入口（scripts/test.sh）设 TEST_BANNER=0：它不需要开始横幅，
# 最后用更重的框把整次运行的结论和单门语言区分开。
TEST_BANNER=${TEST_BANNER:-1}

test_begin() {
    [ "$TEST_BANNER" = 1 ] || return 0
    printf '\n%s%s\n  %s\n%s%s\n' "$_c_head" "$_rule" "$TEST_LANGUAGE" "$_rule" "$_c_off"
}

test_summary() {
    local rc=${1:-$?}
    [ "$rc" -eq 0 ] || TEST_STATUS=FAIL
    local color=$_c_pass
    [ "$TEST_STATUS" = FAIL ] && color=$_c_fail
    if [ "$TEST_BANNER" = 1 ]; then
        printf '%s  %s: %s%s\n' "$color" "$TEST_LANGUAGE" "$TEST_STATUS" "$_c_off"
    else
        printf '\n%s%s\n  %s: %s\n%s%s\n' "$color" "$_rule_all" \
            "$TEST_LANGUAGE" "$TEST_STATUS" "$_rule_all" "$_c_off"
    fi
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

test_begin
