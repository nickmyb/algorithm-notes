package snippets

import "testing"

// 每个用例都注明它能抓到的错误写法，说不出来的不加。

func TestCountLetters(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want map[byte]int // 只列非零项，其余下标必须为 0
	}{
		// 重复字母：抓把 += 1 写成 = 1
		{"重复字母", "anagram", map[byte]int{'a': 3, 'n': 1, 'g': 1, 'r': 1, 'm': 1}},
		// 两端字母：抓下标偏移（如 s-'a'+1），'z' 会越界、'a' 会落到下标 1
		{"两端字母", "az", map[byte]int{'a': 1, 'z': 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countLetters(tt.str)
			for i := range got {
				if want := tt.want[byte('a'+i)]; got[i] != want {
					t.Errorf("countLetters(%q)['%c'] = %d, want %d", tt.str, 'a'+i, got[i], want)
				}
			}
		})
	}
}

// 调用方依赖的性质：互为变位词结果相等，否则不等。
func TestCountLettersAsAnagramKey(t *testing.T) {
	if countLetters("listen") != countLetters("silent") {
		t.Error("互为变位词的计数应相等")
	}
	// 字母集合相同、个数不同：抓只记「出现过没有」的写法
	if countLetters("aab") == countLetters("abb") {
		t.Error("字母个数不同的计数不应相等")
	}
}

// 文档注释写明非小写字母会 panic，这里核实这句话属实。
func TestCountLettersPanicsOnNonLowercase(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("大写字母应越界 panic")
		}
	}()
	countLetters("A")
}

func TestAnagrammatize(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		// 重复字符：抓排序时去重
		{"重复字符", "banana", "aaabnn"},
		// 多字节字符：抓按 []byte 排序，那样会把 UTF-8 编码拆碎，得到非法字符串
		{"多字节字符", "文中", "中文"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anagrammatize(tt.str); got != tt.want {
				t.Errorf("anagrammatize(%q) = %q, want %q", tt.str, got, tt.want)
			}
		})
	}
}

// 调用方依赖的性质：互为变位词结果相等，否则不等。
func TestAnagrammatizeAsAnagramKey(t *testing.T) {
	if anagrammatize("eat") != anagrammatize("tea") {
		t.Error("互为变位词的规范形式应相等")
	}
	if anagrammatize("aab") == anagrammatize("abb") {
		t.Error("字符个数不同的规范形式不应相等")
	}
}
