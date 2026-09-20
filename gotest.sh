#!/usr/bin/env bash
#
# 跑 Go 题解测试并生成单一、合法的覆盖率文件 coverage.txt。
#
# Go 1.10+ 支持对多个包一次性 -coverprofile，直接产出单个合法 profile，
# 不要退回成「每个包分别跑再 cat 拼接」——那样产出的文件夹着多行重复的
# "mode: atomic" 头，新版 Codecov 上传器会把它算成 0% 覆盖率。
set -e

# 题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
# ctl/ 和 structures/ 的单测由后面那条 go test ./... 覆盖。
ROOTS=(./leetcode/...)

go test -covermode=atomic -coverprofile=coverage.txt "${ROOTS[@]}"

# 工具链自己的测试（ctl 的 HTML→Markdown 转换、structures 的数据结构）不进覆盖率统计，
# 但必须跑，否则改坏了 ctl 要到下次生成 README 才发现。
go test ./ctl/... ./structures/...
