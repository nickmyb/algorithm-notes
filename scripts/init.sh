#!/usr/bin/env bash
#
# 把一个新 clone 的仓库带到能开始刷题的状态：
#   1. 检查工具链（Go 必需，Python / Java 缺了只警告，不写那门语言就不影响）
#   2. 整理 Go 依赖并确认 ctl 能编译
#   3. 建 Python 虚拟环境并装 pytest
#   4. 跑一遍三门语言的测试，确认链路通
#   5. 生成仓库根 README.md
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

GO="${GO:-go}"
PYTHON="${PYTHON:-python3}"
VENV="$ROOT/.venv"

# 项目验证过的版本下限。Go 的下限由 go.mod 的 go 指令管，Java 的由 javatest.sh
# 的 --release 管，这里只需要管 Python。三者都写在 README 的「环境要求」里。
PYTHON_MIN_MAJOR=3
PYTHON_MIN_MINOR=12

warnings=()

step() { printf '\n\033[36m==> %s\033[0m\n' "$1"; }
warn() {
    printf '\033[33m[warn] %s\033[0m\n' "$1"
    warnings+=("$1")
}
die() {
    printf '\033[31m[error] %s\033[0m\n' "$1"
    exit 1
}

step "检查工具链"
if ! command -v "$GO" >/dev/null 2>&1; then
    # 装了 Go 但没进 PATH 是最常见的情况（比如装在 ~/go/go1.26.4/bin 却没导出），
    # 所以除了"去装"还要给出"已经装了怎么办"的两条出路。
    printf '\033[31m[error] 找不到 go。ctl 工具是 Go 写的，这一项必需。\033[0m\n'
    echo "  没装：      https://go.dev/dl/ ，版本见 go.mod（当前 $(grep -m1 '^go ' "$ROOT/go.mod" | awk '{print $2}')）"
    echo "  装了没进 PATH：export PATH=/path/to/go/bin:\$PATH 后重跑，或直接指定："
    echo "                 make init GO=/path/to/go/bin/go"
    exit 1
fi
echo "  go      $($GO version | awk '{print $3}')"

has_python=1
if command -v "$PYTHON" >/dev/null 2>&1; then
    python_version="$($PYTHON --version 2>&1 | awk '{print $2}')"
    echo "  python  $python_version"
    if ! "$PYTHON" -c "import sys; sys.exit(0 if sys.version_info >= ($PYTHON_MIN_MAJOR, $PYTHON_MIN_MINOR) else 1)"; then
        die "Python 需要 ${PYTHON_MIN_MAJOR}.${PYTHON_MIN_MINOR} 或更高，当前是 $python_version"
    fi
else
    has_python=0
    warn "找不到 $PYTHON，跳过 Python 题解的环境准备（装了但不在 PATH 时用 make init PYTHON=/path/to/python3）"
fi

has_java=1
if command -v javac >/dev/null 2>&1 && command -v java >/dev/null 2>&1; then
    # 实际的语言版本由 javatest.sh 的 --release 约束，这里只报一下本机装的是哪个
    echo "  javac   $(javac -version 2>&1 | awk '{print $2}')"
else
    has_java=0
    warn "找不到 javac/java，跳过 Java 题解的测试"
fi

step "整理 Go 依赖"
$GO mod tidy || die "go mod tidy 失败"

step "编译 ctl"
(cd ctl && $GO build -o /dev/null .) || die "ctl 编译失败"
echo "  ok"

if [ "$has_python" = 1 ]; then
    step "准备 Python 虚拟环境（$VENV）"
    if [ ! -x "$VENV/bin/pytest" ]; then
        $PYTHON -m venv "$VENV" || die "创建虚拟环境失败"
        "$VENV/bin/pip" install --quiet --upgrade pip
        "$VENV/bin/pip" install --quiet -r requirements-dev.txt || die "安装 Python 依赖失败"
    fi
    echo "  $("$VENV/bin/pytest" --version)"
fi

step "跑 Go 题解测试"
bash ./gotest.sh || die "Go 测试失败"

if [ "$has_python" = 1 ]; then
    step "跑 Python 题解测试"
    PYTEST="$VENV/bin/pytest" bash ./pytest.sh || die "Python 测试失败"
fi

if [ "$has_java" = 1 ]; then
    step "跑 Java 题解测试"
    bash ./javatest.sh || die "Java 测试失败"
fi

step "生成 README.md"
# 需要联网拉 LeetCode 题库；拉不到就跳过，不影响本地开发
if (cd ctl && $GO run . build readme); then
    echo "  ok"
else
    warn "生成 README 失败（多半是网络不通），联网后跑 make readme 补上"
fi

printf '\n\033[32m初始化完成\033[0m\n'
if [ ${#warnings[@]} -gt 0 ]; then
    printf '有 %s 条提醒：\n' "${#warnings[@]}"
    printf '  - %s\n' "${warnings[@]}"
fi

cat <<'EOF'

下一步：
  make new ID=1        新开一题（会去 LeetCode 查标题并建好目录）
  make test ID=1       只测这一道；不带 ID 则跑三门语言的全部测试
  make readme          刷新 README 的题目表格
  make help            看全部命令

用 JetBrains IDE 的话，先跑一次 make ide 创建独立项目（只在 IDE 里写题才需要，
用编辑器 + make test 的话可以跳过）。之后还要在 IDE 里选 SDK/解释器并指定当前题目，
步骤见 ide-templates/README.md。

可选：想让 README 的「个人数据」表格有内容，按 ctl/README.md 配好 ctl/config.toml。
EOF
