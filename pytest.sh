#!/usr/bin/env bash
#
# 跑 Python 题解测试。
#
#   bash ./pytest.sh                                  全部题解（testpaths 在 pytest.ini 里）
#   bash ./pytest.sh leetcode/0094.Binary-Tree-...    只测指定目录
#
# 存在的理由是最后那段退出码处理：pytest 在「一个测试都没收集到」时返回 5。
# 单题模式下那道题没写 Python 题解是正常的，不该让 make 变红。
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTEST="${PYTEST:-$ROOT/.venv/bin/pytest}"

# 切到仓库根再干活。pytest 的 rootdir 和 pytest.ini 的 testpaths 都跟着 cwd 走，
# 从子目录启动会静默缩小测试范围（只跑当前目录那几个），看起来还是绿的。
cd "$ROOT"
TEST_LANGUAGE=Python
source ./scripts/test-common.sh

if ! command -v "$PYTEST" >/dev/null 2>&1; then
    echo "找不到 pytest（$PYTEST），先跑 make init" >&2
    exit 1
fi

# --no-header 去掉 "platform linux -- Python 3.12.3 ... rootdir ... configfile" 那几行，
# 它们每次都一样，夹在三门语言的输出之间纯属噪音。
if [ $# -gt 0 ]; then
    # 单题模式：只有这一道题，逐条显示每个 Example 的名字和通过与否。
    "$PYTEST" -v --no-header "${test_dirs[@]}"
    rc=$?
else
    # 全量模式分两段跑，因为两类测试想看的东西不一样：
    #   题解测试   每题就几个 Example，逐条显示才知道哪道题的哪个用例挂了
    #   工具链测试 structures 和 scripts 加起来几十个，只要一行结论
    # 合在一起用 -v 会把题解淹掉（make init 里刷过 59 行），全用 -q 又看不见题解用例。
    "$PYTEST" -v --no-header leetcode
    rc=$?
    "$PYTEST" -q --no-header structures/python scripts
    tooling_rc=$?
    [ "$rc" -eq 0 ] && rc=$tooling_rc
fi

# 5 = no tests collected。只有在指定了目录时才当作"这题没有 Python 题解"放过；
# 全量模式下一个测试都收不到说明环境或配置坏了，照常报错。
if [ "$rc" -eq 5 ] && [ $# -gt 0 ]; then
    echo "跳过 Python：指定的题目没有 Python 题解"
    TEST_STATUS=SKIP
    exit 0
fi

exit "$rc"
