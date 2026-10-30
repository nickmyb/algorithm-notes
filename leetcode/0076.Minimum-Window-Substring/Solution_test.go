package leetcode

import "testing"

func TestMinWindow(t *testing.T) {
	qs := []struct {
		name string
		s    string
		t    string
		want string
	}{
		{"Example 1", "ADOBECODEBANC", "ABC", "BANC"},
		{"Example 2", "a", "a", "a"},
		{"Example 3", "a", "aa", ""},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := minWindow(q.s, q.t)
			if got != q.want {
				t.Fatalf("minWindow(%q, %q) = %q, want %q", q.s, q.t, got, q.want)
			}
		})
	}
}
