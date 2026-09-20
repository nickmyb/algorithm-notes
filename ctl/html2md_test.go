package main

import "testing"

func TestHTMLToMarkdown(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "空输入",
			in:   "   \n ",
			want: "",
		},
		{
			name: "段落和行内代码",
			in:   `<p>Given <code>nums</code> and <code>target</code>.</p>`,
			want: "Given `nums` and `target`.",
		},
		{
			name: "加粗与斜体",
			in:   `<p><strong>Note:</strong> return <em>indices</em>.</p>`,
			want: "**Note:** return *indices*.",
		},
		{
			name: "示例代码块原样保留",
			in:   "<p>Example:</p><pre><strong>Input:</strong> nums = [2,7]\n<strong>Output:</strong> [0,1]</pre>",
			want: "Example:\n\n```\nInput: nums = [2,7]\nOutput: [0,1]\n```",
		},
		{
			name: "代码块里的实体要反转义且不被当成行内标记",
			in:   `<pre>1 &lt;= n &lt;= 10^4 &amp;&amp; a &gt; b</pre>`,
			want: "```\n1 <= n <= 10^4 && a > b\n```",
		},
		{
			name: "列表",
			in:   `<ul><li>first</li><li>second</li></ul>`,
			want: "- first\n- second",
		},
		{
			name: "上标转成 ^",
			in:   `<p>1 &lt;= n &lt;= 10<sup>5</sup></p>`,
			want: "1 <= n <= 10^5",
		},
		{
			name: "链接和图片",
			in:   `<p>see <a href="https://x.com/a">here</a></p><img src="https://x.com/i.png" />`,
			want: "see [here](https://x.com/a)\n\n![](https://x.com/i.png)",
		},
		{
			name: "nbsp 换成普通空格",
			in:   `<p>a&nbsp;b</p>`,
			want: "a b",
		},
		{
			// LeetCode 实际返回的就是这个形态：尾随空白是实体而不是真空格
			name: "强调标记内侧的 nbsp 也要挪出去",
			in:   `<p><strong>Follow-up:&nbsp;</strong>Can you do better?</p>`,
			want: "**Follow-up:** Can you do better?",
		},
		{
			name: "只含 nbsp 的段落不留下空强调",
			in:   `<p>&nbsp;</p><p>next</p>`,
			want: "next",
		},
		{
			// CommonMark 会把 "**进阶：**你" 整段渲染成字面星号，补个空格才生效
			name: "中文标点收尾的强调后面补空格",
			in:   `<p><strong>进阶：</strong>你可以想出更优解吗？</p>`,
			want: "**进阶：** 你可以想出更优解吗？",
		},
		{
			name: "英文标点收尾不受影响",
			in:   `<p><strong>Note:</strong>see below</p>`,
			want: "**Note:**see below",
		},
		{
			name: "未覆盖的标签只剥标签不丢字",
			in:   `<p>a <mark>b</mark> c</p>`,
			want: "a b c",
		},
		{
			// LeetCode 常见 <strong>Follow-up: </strong>，空格留在 ** 里面
			// Markdown 就不认这个强调了，星号会原样显示
			name: "强调标记内侧的空白要挪到外面",
			in:   `<p><strong>Follow-up: </strong>Can you do better?</p>`,
			want: "**Follow-up:** Can you do better?",
		},
		{
			// 中文题面里 <strong>…空格</strong><em>…</em> 会叠成 ***，渲染全乱
			name: "相邻强调不粘成三星号",
			in:   `<p><strong>和为目标值 </strong><em><code>target</code></em> 的整数</p>`,
			want: "**和为目标值** *`target`* 的整数",
		},
		{
			name: "全是空白的强调标签原样返回",
			in:   `<p>a<strong> </strong>b</p>`,
			want: "a b",
		},
		{
			name: "相邻列表项之间不留空行",
			in:   `<ul><li>a</li></ul><p></p><ul><li>b</li></ul>`,
			want: "- a\n- b",
		},
		{
			name: "列表前的正文仍然隔一个空行",
			in:   `<p>Constraints:</p><ul><li>a</li><li>b</li></ul>`,
			want: "Constraints:\n\n- a\n- b",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := htmlToMarkdown(c.in); got != c.want {
				t.Errorf("htmlToMarkdown(%q)\n got = %q\nwant = %q", c.in, got, c.want)
			}
		})
	}
}
