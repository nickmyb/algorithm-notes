package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// LeetCode 题目描述里实际会出现的标签就这么一小撮，没必要引一个通用 HTML 解析库。
// 没覆盖到的标签会被直接剥掉，文字内容保留，不会丢信息。
var (
	reCRLF = regexp.MustCompile(`\r\n?`)
	// &nbsp; 得在处理标签之前就归一成普通空格。它是不换行空格的实体写法，
	// 而实体反转义在整个流程的末尾才做——到那时 <strong>Follow-up:&nbsp;</strong>
	// 早就被替换成了 "**Follow-up:&nbsp;**"，空白已经进了强调标记内侧，修不回来。
	reNbsp     = regexp.MustCompile(`(?i)&nbsp;|&#160;|&#xa0;`)
	rePre      = regexp.MustCompile(`(?is)<pre[^>]*>(.*?)</pre>`)
	reCode     = regexp.MustCompile(`(?is)<code[^>]*>(.*?)</code>`)
	reSup      = regexp.MustCompile(`(?is)<sup[^>]*>(.*?)</sup>`)
	reSub      = regexp.MustCompile(`(?is)<sub[^>]*>(.*?)</sub>`)
	reStrong   = regexp.MustCompile(`(?is)<(?:strong|b)[^>]*>(.*?)</(?:strong|b)>`)
	reEm       = regexp.MustCompile(`(?is)<(?:em|i)[^>]*>(.*?)</(?:em|i)>`)
	reImg      = regexp.MustCompile(`(?is)<img[^>]*\bsrc="([^"]*)"[^>]*>`)
	reLink     = regexp.MustCompile(`(?is)<a[^>]*\bhref="([^"]*)"[^>]*>(.*?)</a>`)
	reLi       = regexp.MustCompile(`(?is)<li[^>]*>(.*?)</li>`)
	reBr       = regexp.MustCompile(`(?is)<br\s*/?>`)
	reBlockTag = regexp.MustCompile(`(?is)</?(?:p|div|ul|ol|blockquote|section|h[1-6])[^>]*>`)
	reAnyTag   = regexp.MustCompile(`(?is)<[^>]+>`)
	reBlankRun = regexp.MustCompile(`\n{3,}`)
	reTrailWS  = regexp.MustCompile(`[ \t]+\n`)
	// 匹配「以中文标点收尾的强调 + 紧跟着的非空白字符」，如 "**进阶：**你"。
	// 详见 padCJKEmphasis。
	reCJKEmphasisClose = regexp.MustCompile(`([：，。；！？、）】」』])(\*{1,2})([^\s*])`)
)

// preSentinel 是抽出 <pre> 内容后留下的占位符。用 \x00 包起来是因为
// 题目描述正文里不可能出现 NUL 字节，不会跟真实内容撞上。
const preSentinel = "\x00PRE%d\x00"

// htmlToMarkdown 把 LeetCode 返回的题目描述 HTML 转成 Markdown。
//
// 空输入返回空串（leetcode.com 的 translatedContent 就是空的）。
func htmlToMarkdown(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	s = reCRLF.ReplaceAllString(s, "\n")
	s = reNbsp.ReplaceAllString(s, " ")

	// <pre> 里是示例的输入输出，要原样保留。先抽出来占位，
	// 免得后面的标签替换和实体反转义动到代码块内容。
	var blocks []string
	s = rePre.ReplaceAllStringFunc(s, func(m string) string {
		inner := rePre.FindStringSubmatch(m)[1]
		inner = reBr.ReplaceAllString(inner, "\n")
		inner = reAnyTag.ReplaceAllString(inner, "")
		inner = normalizeText(html.UnescapeString(inner))
		blocks = append(blocks, strings.Trim(inner, "\n"))
		return fmt.Sprintf("\n"+preSentinel+"\n", len(blocks)-1)
	})

	// 行内标签
	s = reCode.ReplaceAllString(s, "`$1`")
	s = reSup.ReplaceAllString(s, "^$1")
	s = reSub.ReplaceAllString(s, "_$1")
	s = replaceEmphasis(s, reStrong, "**")
	s = replaceEmphasis(s, reEm, "*")
	s = reImg.ReplaceAllString(s, "\n![]($1)\n")
	s = reLink.ReplaceAllString(s, "[$2]($1)")

	// 列表项。LeetCode 的题目描述里基本没有嵌套列表，按一层处理。
	s = reLi.ReplaceAllString(s, "\n- $1")

	// 块级标签换成换行，剩下的标签直接剥掉
	s = reBr.ReplaceAllString(s, "\n")
	s = reBlockTag.ReplaceAllString(s, "\n")
	s = reAnyTag.ReplaceAllString(s, "")

	s = normalizeText(html.UnescapeString(s))

	// 把代码块填回去
	for i, block := range blocks {
		s = strings.ReplaceAll(s,
			fmt.Sprintf(preSentinel, i),
			"```\n"+block+"\n```")
	}

	s = reTrailWS.ReplaceAllString(s, "\n")
	s = reBlankRun.ReplaceAllString(s, "\n\n")
	s = tightenLists(s)
	s = padCJKEmphasis(s)
	return strings.TrimSpace(s)
}

// padCJKEmphasis 在以中文标点收尾的强调后面补一个空格。
//
// CommonMark 规定：闭合标记前面是标点、后面又不是空白或标点时，这个标记不算
// "右侧贴合"，整段强调都不生效，星号会原样显示。中文题目里的
// "**进阶：**你可以…" 正好踩中——实测 CommonMark 渲染出来就是字面的
// "**进阶：**你可以…"，而补一个空格后 "**进阶：** 你可以…" 正常加粗。
//
// 英文不受影响，因为英文标点后本来就跟着空格。
func padCJKEmphasis(s string) string {
	return reCJKEmphasisClose.ReplaceAllString(s, "$1$2 $3")
}

// emphasisPad 用来把强调标签内容两端的空白摘出来。
const emphasisPad = " \t\n"

// replaceEmphasis 把 <strong>/<em> 换成 Markdown 标记，并把标签内容两端的空白挪到标记外面。
//
// LeetCode 的题目描述里常见 <strong>Follow-up: </strong> 这种尾随空格，照搬会得到
// "**Follow-up: **"。Markdown 的强调标记内侧不能贴空白，这样写不会渲染成加粗，
// 星号会原样显示出来；中文题目里还会和后面紧跟的 <em> 叠成 "***"，渲染更乱。
func replaceEmphasis(s string, re *regexp.Regexp, marker string) string {
	return re.ReplaceAllStringFunc(s, func(m string) string {
		inner := re.FindStringSubmatch(m)[1]
		trimmed := strings.Trim(inner, emphasisPad)
		if trimmed == "" {
			return inner
		}
		lead := inner[:len(inner)-len(strings.TrimLeft(inner, emphasisPad))]
		trail := inner[len(strings.TrimRight(inner, emphasisPad)):]
		return lead + marker + trimmed + marker + trail
	})
}

// tightenLists 去掉相邻列表项之间的空行。
//
// 块级标签转换会留下多余换行，空行折叠后每个列表项之间都夹一个空行，
// Markdown 会当成"松散列表"渲染，每项都套一层 <p>，行距明显偏大。
func tightenLists(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		isBlank := strings.TrimSpace(line) == ""
		prevIsItem := len(out) > 0 && strings.HasPrefix(out[len(out)-1], "- ")
		nextIsItem := i+1 < len(lines) && strings.HasPrefix(lines[i+1], "- ")
		if isBlank && prevIsItem && nextIsItem {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// normalizeText 收拾反转义之后的文本：&nbsp; 会变成 U+00A0，
// 在 Markdown 里看不出区别却会让对齐和搜索出问题，换成普通空格。
func normalizeText(s string) string {
	return strings.ReplaceAll(s, " ", " ")
}
