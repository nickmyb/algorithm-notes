package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nickmyb/algorithm-notes/ctl/util"
)

// TestSimilarProblemsSymmetric 是真正对仓库生效的检查，make similar 和全量测试都跑它。
func TestSimilarProblemsSymmetric(t *testing.T) {
	for _, e := range checkSimilar(util.SolutionsDir) {
		t.Error(e)
	}
}

// ---------- 检查器自身的测试，在临时目录里造样本 ----------

func similarReadme(title, section string) string {
	return "# [" + title + "](https://leetcode.cn/problems/x/)\n\n## 解题思路\n\n思路\n\n" +
		section + "\n## 复杂度\n\n- 时间复杂度：O(n)\n"
}

func similarRepo(t *testing.T, readmes map[string]string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "leetcode")
	for dir, body := range readmes {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			t.Fatal(err)
		}
		if body == "" {
			continue // 只建目录，不写 README
		}
		if err := os.WriteFile(filepath.Join(root, dir, "README.md"), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func wantSimilarErrors(t *testing.T, errs []string, substrs ...string) {
	t.Helper()
	if len(errs) != len(substrs) {
		t.Fatalf("want %d errors, got %d:\n%s", len(substrs), len(errs), strings.Join(errs, "\n"))
	}
	for i, s := range substrs {
		if !strings.Contains(errs[i], s) {
			t.Errorf("error %d = %q, want it to contain %q", i, errs[i], s)
		}
	}
}

func TestSimilarAcceptsBothSides(t *testing.T) {
	root := similarRepo(t, map[string]string{
		// 前置只标在后做的一侧，两侧说明不同、冒号后带空格都可以
		"0049.Group-Anagrams": similarReadme("49. Group Anagrams",
			"### 同类题\n\n- 前置：[242. Valid Anagram](../0242.Valid-Anagram/)： 如何检查字母异位词\n"),
		"0242.Valid-Anagram": similarReadme("242. Valid Anagram",
			"### 同类题\n\n- [49. Group Anagrams](../0049.Group-Anagrams/)\n"),
		// 没有同类题小节的题不受影响
		"0001.Two-Sum": similarReadme("1. Two Sum", ""),
	})
	wantSimilarErrors(t, checkSimilar(root))
}

func TestSimilarReportsMissingOtherSide(t *testing.T) {
	root := similarRepo(t, map[string]string{
		"0049.Group-Anagrams": similarReadme("49. Group Anagrams",
			"### 同类题\n\n- 前置：[242. Valid Anagram](../0242.Valid-Anagram/)\n"),
		"0242.Valid-Anagram": similarReadme("242. Valid Anagram", ""),
	})
	wantSimilarErrors(t, checkSimilar(root),
		// 指出在哪一侧、补哪一行；对方还没有小节时提示先加上
		"leetcode/0049.Group-Anagrams/README.md:9: 列了 0242.Valid-Anagram，但 leetcode/0242.Valid-Anagram/README.md 没有列回来")
	errs := checkSimilar(root)
	for _, s := range []string{"- [49. Group Anagrams](../0049.Group-Anagrams/)", "还没有「### 同类题」小节"} {
		if !strings.Contains(errs[0], s) {
			t.Errorf("报错里应有 %q:\n%s", s, errs[0])
		}
	}
}

func TestSimilarReportsBrokenLinks(t *testing.T) {
	root := similarRepo(t, map[string]string{
		"0049.Group-Anagrams": similarReadme("49. Group Anagrams",
			"### 同类题\n\n- [18. 4Sum](../0018.4Sum/)\n- [242. Valid Anagram](../0242.Valid-Anagram/)\n- [49. Group Anagrams](../0049.Group-Anagrams/)\n"),
		"0242.Valid-Anagram": "", // 目录在，README 不在
	})
	// 报错按字符串排序，行号 10、11、9 依次对应 242、自己、18
	wantSimilarErrors(t, checkSimilar(root),
		"0242.Valid-Anagram 没有 README.md",
		"链接到了本题自己",
		"链接的题目目录 0018.4Sum 不存在")
}

func TestSimilarReportsMutualPrereq(t *testing.T) {
	root := similarRepo(t, map[string]string{
		"0049.Group-Anagrams": similarReadme("49. Group Anagrams",
			"### 同类题\n\n- 前置：[242. Valid Anagram](../0242.Valid-Anagram/)\n"),
		"0242.Valid-Anagram": similarReadme("242. Valid Anagram",
			"### 同类题\n\n- 前置：[49. Group Anagrams](../0049.Group-Anagrams/)\n"),
	})
	// 同一对只报一次
	wantSimilarErrors(t, checkSimilar(root), "互相标了前置，这是成环")
}

func TestSimilarReportsBadFormat(t *testing.T) {
	root := similarRepo(t, map[string]string{
		"0049.Group-Anagrams": similarReadme("49. Group Anagrams",
			"### 同类题\n\n* [242. Valid Anagram](../0242.Valid-Anagram/)\n- 前置: [1. Two Sum](../0001.Two-Sum/)\n242 也是同类题\n"),
	})
	// 格式不符的行会漏过对称性检查，所以单独报出：列表符号、半角冒号、非列表的文字
	wantSimilarErrors(t, checkSimilar(root),
		"README.md:10: 不符合同类题格式",
		"README.md:11: 不符合同类题格式",
		"README.md:9: 不符合同类题格式")
}

func TestSimilarIgnoresCodeBlocksOtherSectionsAndTemplate(t *testing.T) {
	example := "```markdown\n### 同类题\n\n- [18. 4Sum](../0018.4Sum/)\n```\n"
	root := similarRepo(t, map[string]string{
		// 骨架 README 的格式示例在代码块里，链到的题并不存在
		"0000.Template": similarReadme("0. 骨架", "### 同类题\n\n- [18. 4Sum](../0018.4Sum/)\n"),
		"0001.Two-Sum":  similarReadme("1. Two Sum", example),
		// 同类题小节在下一个标题处结束，复杂度下的列表项不算
		"0015.3Sum": similarReadme("15. 3Sum", "### 同类题\n"),
	})
	wantSimilarErrors(t, checkSimilar(root))
}

func TestSimilarReportsRedundantPrereq(t *testing.T) {
	root := similarRepo(t, map[string]string{
		"0704.Binary-Search": similarReadme("704. Binary Search",
			"### 同类题\n\n- [35. Search Insert Position](../0035.Search-Insert-Position/)\n- [34. Find First](../0034.Find-First/)\n"),
		"0035.Search-Insert-Position": similarReadme("35. Search Insert Position",
			"### 同类题\n\n- 前置：[704. Binary Search](../0704.Binary-Search/)\n- [34. Find First](../0034.Find-First/)\n"),
		// 704 → 35 → 34 已经推出 704 在 34 之前，这里再标前置 704 是冗余边
		"0034.Find-First": similarReadme("34. Find First",
			"### 同类题\n\n- 前置：[35. Search Insert Position](../0035.Search-Insert-Position/)\n- 前置：[704. Binary Search](../0704.Binary-Search/)\n"),
	})
	// 只报冗余的那一条，并给出推出它的路径；直接前置 35 不报
	wantSimilarErrors(t, checkSimilar(root),
		"leetcode/0034.Find-First/README.md:10: 前置 0704.Binary-Search 可由 0704.Binary-Search → 0035.Search-Insert-Position → 0034.Find-First 推出")
}

func TestSimilarAcceptsIndependentPrereqs(t *testing.T) {
	// 多个前置之间没有先后（与关系），一道题也可以是多道题的前置：都不是冗余
	root := similarRepo(t, map[string]string{
		"0001.Two-Sum": similarReadme("1. Two Sum",
			"### 同类题\n\n- [15. 3Sum](../0015.3Sum/)\n- [18. 4Sum](../0018.4Sum/)\n"),
		"0167.Two-Sum-II": similarReadme("167. Two Sum II",
			"### 同类题\n\n- [15. 3Sum](../0015.3Sum/)\n"),
		"0015.3Sum": similarReadme("15. 3Sum",
			"### 同类题\n\n- 前置：[1. Two Sum](../0001.Two-Sum/)\n- 前置：[167. Two Sum II](../0167.Two-Sum-II/)\n"),
		"0018.4Sum": similarReadme("18. 4Sum",
			"### 同类题\n\n- 前置：[1. Two Sum](../0001.Two-Sum/)\n"),
	})
	wantSimilarErrors(t, checkSimilar(root))
}

func TestSimilarRedundancyIgnoresLongCycle(t *testing.T) {
	// 1 → 2 → 3 → 1 是三道题的环，3 另有前置 4。从 1 出发经过 3 本身才能回到 3，
	// 不算冗余；环本身暂不检查，所以这里不报任何错
	root := similarRepo(t, map[string]string{
		"0001.A": similarReadme("1. A", "### 同类题\n\n- 前置：[3. C](../0003.C/)\n- [2. B](../0002.B/)\n"),
		"0002.B": similarReadme("2. B", "### 同类题\n\n- 前置：[1. A](../0001.A/)\n- [3. C](../0003.C/)\n"),
		"0003.C": similarReadme("3. C", "### 同类题\n\n- 前置：[2. B](../0002.B/)\n- 前置：[4. D](../0004.D/)\n- [1. A](../0001.A/)\n"),
		"0004.D": similarReadme("4. D", "### 同类题\n\n- [3. C](../0003.C/)\n"),
	})
	wantSimilarErrors(t, checkSimilar(root))
}
