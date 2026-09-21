#!/usr/bin/env bash
#
# 在 PATH 和常见安装位置里找满足版本下限的 go / python3。
#
# 存在的理由：装了但没进 PATH 是新手最常见的卡点，而各家装的位置都不一样——
# 官方 tarball 在 /usr/local/go/bin，apt 的 golang-1.x 在 /usr/lib/go-1.x/bin，
# `go install golang.org/dl/goX` 在 ~/sdk/goX/bin，pyenv 在 ~/.pyenv/versions/*/bin，
# 没有一个会自动进 PATH。与其让 make init 第一步就退出让人去查文档，不如替他找一遍。
#
# 规则：PATH 里的优先（尊重用户自己配好的），它不满足下限时才翻其他位置，
# 翻到多个就取版本最高的。

# $1 >= $2
version_ge() {
    [ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]
}

go_version()     { "$1" version 2>/dev/null | awk '{print $3}' | sed 's/^go//'; }
# javac -version 打印 "javac 17.0.20"；JDK 8 打印 "javac 1.8.0_432"，sort -V 能正确排在 17 之前
java_version()   { "$1" -version 2>&1 | awk '{print $2; exit}'; }
python_version() { "$1" -c 'import sys; print("%d.%d.%d" % sys.version_info[:3])' 2>/dev/null; }

# 打印所有存在的候选，一行一个。
# 必须显式开 nullglob：默认行为下不匹配的 glob 会原样留下（bash）甚至直接报错
# 中断整个函数（zsh），两种都会让后面的候选被漏掉。
_nullglob_on()  { _saved_nullglob="$(shopt -p nullglob)"; shopt -s nullglob; }
_nullglob_off() { eval "$_saved_nullglob"; }

go_candidates() {
    local p
    _nullglob_on
    for p in /usr/local/go/bin/go /opt/go/bin/go /snap/bin/go \
             /opt/homebrew/bin/go /usr/local/bin/go \
             "$HOME"/go/go*/bin/go "$HOME"/sdk/go*/bin/go /usr/lib/go-*/bin/go; do
        [ -x "$p" ] && printf '%s\n' "$p"
    done
    _nullglob_off
}

python_candidates() {
    local p minor
    _nullglob_on
    # 新的排前面：装了 3.13 又留着系统 3.10 时优先用新的
    for minor in 14 13 12; do
        p="$(command -v "python3.$minor" 2>/dev/null)" && [ -n "$p" ] && printf '%s\n' "$p"
    done
    for p in /usr/bin/python3 /usr/local/bin/python3 /opt/homebrew/bin/python3 \
             "$HOME"/.pyenv/versions/*/bin/python3; do
        [ -x "$p" ] && printf '%s\n' "$p"
    done
    _nullglob_off
}

# java 探测的是 javac，但 javatest.sh 编译完还要用 java 跑。两者必须来自同一个 JDK，
# 否则会撞上 "class file has wrong version"——所以这里只认 java 就在 javac 旁边的候选，
# init.sh 再从选中的 javac 推出配套的 java，不对 java 单独探测一次。
java_candidates() {
    local p
    _nullglob_on
    for p in ${JAVA_HOME:+"$JAVA_HOME/bin/javac"} \
             /usr/lib/jvm/*/bin/javac /usr/java/*/bin/javac \
             /opt/homebrew/opt/openjdk*/bin/javac \
             /Library/Java/JavaVirtualMachines/*/Contents/Home/bin/javac \
             "$HOME"/.sdkman/candidates/java/*/bin/javac; do
        [ -x "$p" ] && [ -x "${p%/javac}/java" ] && printf '%s\n' "$p"
    done
    _nullglob_off
}

# check_tool <go|python|java> <版本下限> <名字或路径>
# 只验这一个，不做任何搜索。用户显式指定时走这条：指定了就得用指定的那个，
# 不合用要能说清是"找不到"还是"版本不够"，不能背着他换一个。
#   0 = 可用，打印真实路径
#   1 = 找不到 / 跑不起来
#   2 = 找到了但版本低于下限（打印实际版本）
check_tool() {
    local kind="$1" min="$2" given="$3" p v
    p="$(command -v "$given" 2>/dev/null)" || return 1
    [ -n "$p" ] || return 1
    v="$("${kind}_version" "$p")"
    [ -n "$v" ] || return 1
    if ! version_ge "$v" "$min"; then
        printf '%s\n' "$v"
        return 2
    fi
    printf '%s\n' "$p"
}

# find_tool <go|python|java> <版本下限> <PATH 里的名字>
# 没有显式指定时才用：PATH 里的优先，不合用再翻常见安装位置。
# 找到打印路径返回 0，找不到返回 1。
find_tool() {
    local kind="$1" min="$2" name="$3" p v best="" best_v=""

    p="$(check_tool "$kind" "$min" "$name")" && { printf '%s\n' "$p"; return 0; }

    while IFS= read -r p; do
        [ -n "$p" ] || continue
        v="$("${kind}_version" "$p")"
        [ -n "$v" ] && version_ge "$v" "$min" || continue
        if [ -z "$best_v" ] || version_ge "$v" "$best_v"; then
            best="$p"; best_v="$v"
        fi
    done < <("${kind}_candidates")

    [ -n "$best" ] || return 1
    printf '%s\n' "$best"
}
