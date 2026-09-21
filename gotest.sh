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
set -eu

# 切到仓库根再干活。下面用的是 ./leetcode/... 这类相对路径，而 IDE 的 Run
# Configuration、编辑器终端常常在别的目录启动，不切的话会直接报
# "pattern ./leetcode/...: no such file or directory"。
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_LANGUAGE=Go
source ./scripts/test-common.sh
GO="${GO:-go}"

# 题解根目录。以后新增 lcp/ 这类同级目录，在这里加一项即可。
ROOTS=(./leetcode/...)

if [ $# -gt 0 ]; then
    # 单题模式：挑出确实有 Go 文件的目录。那道题没写 Go 题解不算失败，
    # 直接跳过——否则 go test 会报 "no Go files" 让整条命令变红。
    targets=()
    for dir in "${test_dirs[@]}"; do
        if compgen -G "$dir/*.go" > /dev/null; then
            targets+=("$dir")
        fi
    done

    if [ ${#targets[@]} -eq 0 ]; then
        echo "跳过 Go：指定的题目没有 Go 题解"
        TEST_STATUS=SKIP
        exit 0
    fi

    "$GO" test -v "${targets[@]}"
    exit
fi

# 排除 IDE 的编译输出目录。IDEA 把题目目录当项目打开时会在里面建 out/，并且把
# 所有非 .java 文件当资源复制进去——包括 Solution.go 和 Solution_test.go 的旧副本。
# go 的 ./... 会把那儿当成一个真实的包编译测试：旧题解配旧测试，永远是绿的，
# 但它早就不是你在写的那份了。这种"测试说谎"比测试失败危险得多。
# 先单独取结果：进程替换里的 go list 失败不会触发外层 set -e。
# 否则 go 缺失、依赖损坏都会被误报成「没有包」，整次测试仍然成功。
package_list="$("$GO" list "${ROOTS[@]}")"
pkgs=()
while IFS= read -r package; do
    case "$package" in
        '' | */out/*) continue ;;
    esac
    pkgs+=("$package")
done <<< "$package_list"

if [ ${#pkgs[@]} -gt 0 ]; then
    "$GO" test -v -covermode=atomic -coverprofile=coverage.txt "${pkgs[@]}"
fi

# 工具链自己的测试（ctl 的 HTML→Markdown 转换、structures 的数据结构）不进覆盖率统计，
# 但必须跑，否则改坏了 ctl 要到下次生成 README 才发现。
"$GO" test ./ctl/... ./structures/...
