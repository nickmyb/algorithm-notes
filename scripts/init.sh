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
JAVAC="${JAVAC:-javac}"
VENV="$ROOT/.venv"
LOCAL_MK="$ROOT/local.mk"

source "$ROOT/scripts/find-tool.sh"

# 项目验证过的版本下限。Go 的下限跟着 go.mod 的 go 指令走（改 go.mod 这里自动跟上），
# Java 的由 javatest.sh 的 --release 管。三者都写在 README 的「环境要求」里。
GO_MIN="$(awk '/^go /{print $2; exit}' "$ROOT/go.mod")"
PYTHON_MIN=3.12
# JDK 要能编译出 --release 指定的字节码版本，所以下限就是 javatest.sh 里的那个数
JAVA_MIN="$(sed -n 's/.*JAVA_RELEASE:-\([0-9][0-9]*\).*/\1/p' "$ROOT/javatest.sh" | head -1)"

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

# 解析一门语言的解释器/编译器。两条路径，区别很重要：
#
#   显式指定（make init GO=... 或 local.mk 里记着的）→ 只验这一个。
#       不合用就报错退出，绝不背着用户换一个：他既然写了路径，就是要用那个。
#   没指定（还是默认的 go / python3 / javac）→ 才去 PATH 和常见安装位置找。
#
# 约定用退出码回传结果，不用全局变量：resolve 总是在 $(...) 里调用，
# 子 shell 里的赋值传不回来（踩过）。
#   0 = 可用，stdout 是真实路径
#   2 = 找到了但版本不够，stdout 是实际版本
#   1 = 找不到，stdout 为空
resolve() {
    local kind="$1" min="$2" given="$3" default="$4"
    if [ "$given" = "$default" ]; then
        find_tool "$kind" "$min" "$default"
    else
        check_tool "$kind" "$min" "$given"
    fi
}

# 显式指定失败时的统一收尾。$3 是 stdout（版本号或空），$4 是退出码。
die_explicit() {
    local var="$1" label="$2" detail="$3" rc="$4"
    if [ "$rc" = 2 ]; then
        printf '\033[31m[error] 指定的 %s 版本是 %s，低于要求的 %s\033[0m\n' \
            "$label" "$detail" "$5"
    else
        printf '\033[31m[error] 指定的 %s 找不到或跑不起来\033[0m\n' "$label"
    fi
    echo "  来源是命令行的 $var=，或 local.mk 里记着的路径"
    echo "  换一个路径重跑，或删掉 local.mk 让 make init 重新探测"
    exit 1
}

step "检查工具链"

if resolved_go="$(resolve go "$GO_MIN" "$GO" go)"; then
    GO="$resolved_go"
    printf '  go      %s  (%s)\n' "$(go_version "$GO")" "$GO"
else
    rc=$?
    [ "$GO" = go ] || die_explicit GO Go "$resolved_go" "$rc" "$GO_MIN"
    # 没显式指定，说明 PATH 和常见安装位置都没有满足下限的 Go，只能让用户去装
    printf '\033[31m[error] 没找到 go %s 或更高版本。ctl 工具是 Go 写的，这一项必需。\033[0m\n' "$GO_MIN"
    echo "  已经找过：PATH、/usr/local/go、/opt/go、/snap/bin、~/go/go*、~/sdk/go*、/usr/lib/go-*"
    echo "  去装：    https://go.dev/dl/ （需要 $GO_MIN 或更高）"
    echo "  装在别处：make init GO=/path/to/go/bin/go"
    exit 1
fi

has_python=1
if resolved_python="$(resolve python "$PYTHON_MIN" "$PYTHON" python3)"; then
    PYTHON="$resolved_python"
    printf '  python  %s  (%s)\n' "$(python_version "$PYTHON")" "$PYTHON"
else
    rc=$?
    # Python 缺了平时只警告，但用户自己指定的不能当没看见
    [ "$PYTHON" = python3 ] || die_explicit PYTHON Python "$resolved_python" "$rc" "$PYTHON_MIN"
    has_python=0
    warn "没找到 Python $PYTHON_MIN 或更高版本（PATH、/usr/bin、pyenv 都找过了），跳过 Python 题解；装好后重跑 make init，或 make init PYTHON=/path/to/python3"
fi

has_java=1
JAVA=java
if resolved_javac="$(resolve java "$JAVA_MIN" "$JAVAC" javac)"; then
    JAVAC="$resolved_javac"
    # java 必须和 javac 同一个 JDK，否则编译产物跑起来会报 class file has wrong version。
    # 自动探测时 find-tool.sh 的候选过滤已经保证了这点，显式指定时没走到，这里补一道。
    if [ ! -x "${JAVAC%/javac}/java" ]; then
        printf '\033[31m[error] %s 旁边没有配套的 java\033[0m\n' "$JAVAC"
        echo "  javac 和 java 必须来自同一个 JDK，否则编译产物跑起来会报 class file has wrong version"
        exit 1
    fi
    JAVA="${JAVAC%/javac}/java"
    printf '  javac   %s  (%s)\n' "$(java_version "$JAVAC")" "$JAVAC"
else
    rc=$?
    [ "$JAVAC" = javac ] || die_explicit JAVAC JDK "$resolved_javac" "$rc" "$JAVA_MIN"
    has_java=0
    warn "没找到 JDK $JAVA_MIN 或更高版本（PATH、\$JAVA_HOME、/usr/lib/jvm、sdkman 都找过了），跳过 Java 题解；装好后重跑 make init，或 make init JAVAC=/path/to/javac"
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

# 这三处带 TEST_TEMPLATE=1：初始化要连 0000.Template 一起验，确认骨架本身
# 能编译能跑。日常 make test 不测它，免得每次都多出一条 SKIP 和几行噪音。
step "跑 Go 题解测试"
TEST_TEMPLATE=1 GO="$GO" bash ./gotest.sh || die "Go 测试失败"

if [ "$has_python" = 1 ]; then
    step "跑 Python 题解测试"
    TEST_TEMPLATE=1 PYTEST="$VENV/bin/pytest" bash ./pytest.sh || die "Python 测试失败"
fi

if [ "$has_java" = 1 ]; then
    step "跑 Java 题解测试"
    TEST_TEMPLATE=1 JAVAC="$JAVAC" JAVA="$JAVA" bash ./javatest.sh || die "Java 测试失败"
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
