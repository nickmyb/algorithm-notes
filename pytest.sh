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

if [ ! -x "$PYTEST" ]; then
    echo "找不到 pytest（$PYTEST），先跑 make init" >&2
    exit 1
fi

"$PYTEST" -q "$@"
rc=$?

# 5 = no tests collected。只有在指定了目录时才当作"这题没有 Python 题解"放过；
# 全量模式下一个测试都收不到说明环境或配置坏了，照常报错。
if [ "$rc" -eq 5 ] && [ $# -gt 0 ]; then
    echo "跳过 Python：指定的题目没有 Python 题解"
    exit 0
fi

exit "$rc"
