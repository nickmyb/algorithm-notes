package snippets

import "testing"

func TestCountAscii(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want map[byte]int // 只列非零项，其余下标必须为 0
	}{
		// 空输入：验证没有字符时所有计数都为 0。
		{"空字符串", "", nil},
		// 大小写与重复字符：抓大小写合并和只记录是否出现的写法。
		{"大小写及重复字符", "AaAaZz", map[byte]int{'A': 2, 'a': 2, 'Z': 1, 'z': 1}},
		// 非字母：抓保留旧版减 'A' 偏移、忽略数字或空白的写法。
		{"数字标点与空白", "001! !\t\n", map[byte]int{'0': 2, '1': 1, '!': 2, ' ': 1, '\t': 1, '\n': 1}},
		// ASCII 两端：抓下标偏移和遗漏控制字符的写法。
		{"ASCII 两端", "\x00\x7f\x00", map[byte]int{0: 2, 127: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countAscii(tt.str)
			for i, count := range got {
				if want := tt.want[byte(i)]; count != want {
					t.Errorf("countAscii(%q)[0x%02x] = %d, want %d", tt.str, i, count, want)
				}
			}
		})
	}
}

// ASCII 上界之外的第一个字节：验证文档中的非 ASCII 输入契约。
func TestCountAsciiPanicsOnNonASCII(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("非 ASCII 字节应越界 panic")
		}
	}()
	countAscii("\x80")
}
