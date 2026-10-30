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

	impls := map[string]func(string, string) string{
		"minWindow":           minWindow,
		"minWindowScanCounts": minWindowScanCounts,
	}

	for _, q := range qs {
		for implName, impl := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := impl(q.s, q.t)
				if got != q.want {
					t.Fatalf("%s(%q, %q) = %q, want %q", implName, q.s, q.t, got, q.want)
				}
			})
		}
	}
}
