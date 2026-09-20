package structures

import "testing"

func TestFormatInts(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want string
	}{
		{"空数组", []int{}, "[]"},
		{"纯数字", []int{1, 2, 3}, "[1,2,3]"},
		{"含哨兵", []int{1, NULL, 2, 3}, "[1,null,2,3]"},
		{"负数不受影响", []int{-1, NULL, -100}, "[-1,null,-100]"},
		{"全是哨兵", []int{NULL, NULL}, "[null,null]"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FormatInts(c.in); got != c.want {
				t.Fatalf("FormatInts(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}
