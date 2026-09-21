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

# 单题模式用 -v：那时你在反复改一道题，想看到每个 Example 的名字和通过与否。
# 全量模式用 -q：几十上百个用例逐条刷屏会把真正要看的结论淹掉（make init 里
# 尤其明显）。失败详情两种模式都会打印，安静的只是通过的那些。
if [ $# -gt 0 ]; then
    "$PYTEST" -v "${test_dirs[@]}"
else
    "$PYTEST" -q "${test_dirs[@]}"
fi
rc=$?

# 5 = no tests collected。只有在指定了目录时才当作"这题没有 Python 题解"放过；
# 全量模式下一个测试都收不到说明环境或配置坏了，照常报错。
if [ "$rc" -eq 5 ] && [ $# -gt 0 ]; then
    echo "跳过 Python：指定的题目没有 Python 题解"
    TEST_STATUS=SKIP
    exit 0
fi

exit "$rc"
