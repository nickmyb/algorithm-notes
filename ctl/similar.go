package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// 题解 README 里「### 同类题」一节的解析与检查，格式和规则见 AGENTS.md 约束 12。
//
// 规则要求每一对同类题两边都写、前置只标在后做的一侧。这里只报告问题，
// 不改写题解 README：那是作者的笔记，补哪一行、说明怎么写由作者决定。

const similarHeading = "### 同类题"

var (
	// 「- 前置：[49. Group Anagrams](../0049.Group-Anagrams/)：说明」，前置和说明都可选
	similarLine = regexp.MustCompile(`^- (前置：)?\[[^\]]+\]\(\.\./([^/()]+)/?\)`)
	// 题解 README 的一级标题「# [49. Group Anagrams](https://...)」
	problemTitle = regexp.MustCompile(`^# \[([^\]]+)\]\(`)
	// 只认 NNNN.标题 形式的题目目录，0000.Template 另外跳过
	problemDirName = regexp.MustCompile(`^\d{4}\.`)
)

type similarLink struct {
	prereq bool
	line   int
}

type similarNotes struct {
	title      string // 「49. Group Anagrams」，报错时用来给出该补的那一行
	hasSection bool
	links      map[string]similarLink // 键是目标题目目录名
}

// parseSimilar 解析一份题解 README 的同类题小节，返回格式不符的行。
// 代码块里的内容不算（骨架 README 的格式示例就在代码块里）。
func parseSimilar(readme string) (similarNotes, []string) {
	notes := similarNotes{links: map[string]similarLink{}}
	var errs []string
	inFence, inSection := false, false
	scanner := bufio.NewScanner(strings.NewReader(readme))
	for n := 1; scanner.Scan(); n++ {
		line := strings.TrimRight(scanner.Text(), " \t")
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if notes.title == "" {
			if m := problemTitle.FindStringSubmatch(line); m != nil {
				notes.title = m[1]
			}
		}
		if strings.HasPrefix(line, "#") {
			inSection = line == similarHeading
			notes.hasSection = notes.hasSection || inSection
			continue
		}
		if !inSection || line == "" {
			continue
		}
		m := similarLine.FindStringSubmatch(line)
		if m == nil {
			errs = append(errs, fmt.Sprintf("%d: 不符合同类题格式，应为「- [题号. 标题](../题目目录/)」，建议先做的加「前置：」: %s", n, line))
			continue
		}
		notes.links[m[2]] = similarLink{prereq: m[1] != "", line: n}
	}
	return notes, errs
}

// checkSimilar 检查 root 下所有题目的同类题：链接指向的目录存在、两边都写了、
// 前置不成环（两道题互标、多道题绕一圈）、没有能由传递推出的冗余前置。
func checkSimilar(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return []string{err.Error()}
	}
	// 报错里的路径从题解根目录名开始，如 leetcode/0049.Group-Anagrams/README.md
	base := filepath.Base(filepath.Clean(root))
	readmePath := func(dir string) string { return base + "/" + dir + "/README.md" }

	all := map[string]similarNotes{}
	var dirs, errs []string
	for _, e := range entries {
		if !e.IsDir() || !problemDirName.MatchString(e.Name()) || e.Name() == "0000.Template" {
			continue
		}
		body, err := os.ReadFile(filepath.Join(root, e.Name(), "README.md"))
		if err != nil {
			continue // 没有 README 的目录不参与；被别的题链接到时在下面报出
		}
		notes, formatErrs := parseSimilar(string(body))
		for _, fe := range formatErrs {
			errs = append(errs, readmePath(e.Name())+":"+fe)
		}
		all[e.Name()] = notes
		dirs = append(dirs, e.Name())
	}
	sort.Strings(dirs)

	for _, a := range dirs {
		targets := make([]string, 0, len(all[a].links))
		for b := range all[a].links {
			targets = append(targets, b)
		}
		sort.Strings(targets)
		for _, b := range targets {
			link := all[a].links[b]
			at := fmt.Sprintf("%s:%d", readmePath(a), link.line)
			if b == a {
				errs = append(errs, at+": 同类题链接到了本题自己")
				continue
			}
			other, ok := all[b]
			if !ok {
				if _, err := os.Stat(filepath.Join(root, b)); err != nil {
					errs = append(errs, fmt.Sprintf("%s: 链接的题目目录 %s 不存在", at, b))
				} else {
					errs = append(errs, fmt.Sprintf("%s: %s 没有 README.md，无法核对另一侧", at, b))
				}
				continue
			}
			back, ok := other.links[a]
			if !ok {
				msg := fmt.Sprintf("%s: 列了 %s，但 %s 没有列回来。同类题两边都写，应在它的「%s」下补一行：\n  - [%s](../%s/)",
					at, b, readmePath(b), similarHeading, all[a].title, a)
				if !other.hasSection {
					msg += fmt.Sprintf("\n  （%s 还没有「%s」小节，先在「## 解题思路」下加上）", b, similarHeading)
				}
				errs = append(errs, msg)
				continue
			}
			// 两侧都标前置时只报一次
			if link.prereq && back.prereq && a < b {
				errs = append(errs, fmt.Sprintf("%s 和 %s:%d 互相标了前置，这是成环。前置只标在后做的一侧，另一侧写成普通同类题",
					at, readmePath(b), back.line))
			}
		}
	}
	errs = append(errs, redundantPrereqs(all, dirs, readmePath)...)
	errs = append(errs, prereqCycles(all, dirs, readmePath)...)
	sort.Strings(errs)
	return errs
}

// prereqGraph 收集前置边 B → A（B 建议先做，记在 A 的 README 里）：
// next 是每道题的后继，prereqs 是每道题的直接前置，都已排序。
func prereqGraph(all map[string]similarNotes, dirs []string) (next, prereqs map[string][]string) {
	next, prereqs = map[string][]string{}, map[string][]string{}
	for _, a := range dirs {
		for b, link := range all[a].links {
			if _, ok := all[b]; !ok || b == a || !link.prereq {
				continue // 链接本身的问题在 checkSimilar 里报
			}
			next[b] = append(next[b], a)
			prereqs[a] = append(prereqs[a], b)
		}
	}
	for _, m := range []map[string][]string{next, prereqs} {
		for _, v := range m {
			sort.Strings(v)
		}
	}
	return next, prereqs
}

// prereqCycles 报出多于两道题的前置环（两道题互标前置在 checkSimilar 里报）。
// 对每道题找经过它的最短环，同一组题的环只报一次。
func prereqCycles(all map[string]similarNotes, dirs []string, readmePath func(string) string) []string {
	next, _ := prereqGraph(all, dirs)
	seen := map[string]bool{}
	var errs []string
	for _, s := range dirs {
		var cycle []string // s → … → s，不含末尾重复的 s
		for _, n := range next[s] {
			if path := prereqPath(next, n, s, ""); path != nil && (cycle == nil || len(path) < len(cycle)) {
				cycle = append([]string{s}, path[:len(path)-1]...)
			}
		}
		if len(cycle) <= 2 {
			continue
		}
		members := append([]string(nil), cycle...)
		sort.Strings(members)
		if key := strings.Join(members, " "); !seen[key] {
			seen[key] = true
			// 闭合环的最后一条边记在 s 的 README 里：「前置：cycle 的最后一道」
			errs = append(errs, fmt.Sprintf("%s:%d: 前置成环：%s → %s。前置表示建议先做，环里的题排不出先后，把其中一条改成普通同类题",
				readmePath(s), all[s].links[cycle[len(cycle)-1]].line, strings.Join(cycle, " → "), s))
		}
	}
	return errs
}

// redundantPrereqs 报出能由传递推出的前置：A 写了「前置：B」，而 B 经由 A 的
// 另一个前置 C 已经排在 A 之前（B → … → C → A），那么 B → A 这条边是冗余的。
// 找路径时绕开 A 本身，否则环会被误报成冗余（环由 prereqCycles 报）。
func redundantPrereqs(all map[string]similarNotes, dirs []string, readmePath func(string) string) []string {
	next, prereqs := prereqGraph(all, dirs)
	var errs []string
	for _, a := range dirs {
		ps := prereqs[a]
		for _, b := range ps {
			for _, c := range ps {
				if c == b {
					continue
				}
				if path := prereqPath(next, b, c, a); path != nil {
					errs = append(errs, fmt.Sprintf("%s:%d: 前置 %s 可由 %s → %s 推出，是冗余前置。只写直接前置：去掉这一行的「前置：」，或者两侧一起删掉",
						readmePath(a), all[a].links[b].line, b, strings.Join(path, " → "), a))
					break
				}
			}
		}
	}
	return errs
}

// prereqPath 沿前置边找一条 from 到 to 的路径（BFS，取最短的），不经过 avoid。
// 找不到时返回 nil。
func prereqPath(next map[string][]string, from, to, avoid string) []string {
	parent := map[string]string{from: ""}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == to {
			var path []string
			for n := to; n != ""; n = parent[n] {
				path = append([]string{n}, path...)
			}
			return path
		}
		for _, n := range next[cur] {
			if _, seen := parent[n]; seen || n == avoid {
				continue
			}
			parent[n] = cur
			queue = append(queue, n)
		}
	}
	return nil
}
