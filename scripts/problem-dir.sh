#!/usr/bin/env bash
#
# 把题号解析成题目目录，如 94 -> leetcode/0094.Binary-Tree-Inorder-Traversal
#
# 目录路径（相对仓库根）打到 stdout，错误信息打到 stderr，供 Makefile 用
# $(shell ...) 取值。找不到或匹配到多个都算失败。
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

id="${1:-}"
if [ -z "$id" ]; then
    echo "用法: $0 <题号>" >&2
    exit 1
fi
case "$id" in
    '' | *[!0-9]*)
        echo "题号必须是数字，收到 '$id'" >&2
        exit 1
        ;;
esac

# 题号是十进制文本。printf %d 会把 0094 当八进制，报错后甚至得到 0000。
# 用字符串补齐宽度，同时避免超长数字的整数溢出。
id="${id#"${id%%[!0]*}"}"
id="${id:-0}"
padded="$(printf '%4s' "$id")"
padded="${padded// /0}"

# 题解根目录，和 gotest.sh / javatest.sh / pytest.ini 保持一致。
# 以后新增 lcp/ 这类同级目录，在这里加一项。
ROOTS=(leetcode)

matches=()
for root in "${ROOTS[@]}"; do
    for dir in "$root/$padded".*/; do
        [ -d "$dir" ] || continue
        matches+=("${dir%/}")
    done
done

case "${#matches[@]}" in
    0)
        echo "找不到题号 $id（期望 ${ROOTS[0]}/$padded.*），先跑 make new ID=$id" >&2
        exit 1
        ;;
    1)
        echo "${matches[0]}"
        ;;
    *)
        echo "题号 $id 匹配到多个目录，需要人工处理：" >&2
        printf '  %s\n' "${matches[@]}" >&2
        exit 1
        ;;
esac
