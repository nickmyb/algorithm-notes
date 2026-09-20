#!/usr/bin/env bash
#
# 跑 Go 题解测试。
#
#   bash ./gotest.sh                                  全部题解 + 覆盖率 + 工具链自测
#   bash ./gotest.sh leetcode/0094.Binary-Tree-...    只测指定目录，不生成覆盖率
#
# 全量模式下用一次 -coverprofile 覆盖所有包，直接产出单个合法 profile。
# 不要退回成「每个包分别跑再 cat 拼接」——那样产出的文件夹着多行重复的
# "mode: atomic" 头，新版 Codecov 上传器会把它算成 0% 覆盖率。
set -e

# 题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
ROOTS=(./leetcode/...)

if [ $# -gt 0 ]; then
    # 单题模式：挑出确实有 Go 文件的目录。那道题没写 Go 题解不算失败，
    # 直接跳过——否则 go test 会报 "no Go files" 让整条命令变红。
    targets=()
    for dir in "$@"; do
        if compgen -G "$dir/*.go" > /dev/null; then
            targets+=("./${dir#./}")
        fi
    done

    if [ ${#targets[@]} -eq 0 ]; then
        echo "跳过 Go：指定的题目没有 Go 题解"
        exit 0
    fi

    go test "${targets[@]}"
    exit
fi

go test -covermode=atomic -coverprofile=coverage.txt "${ROOTS[@]}"

# 工具链自己的测试（ctl 的 HTML→Markdown 转换、structures 的数据结构）不进覆盖率统计，
# 但必须跑，否则改坏了 ctl 要到下次生成 README 才发现。
go test ./ctl/... ./structures/...
