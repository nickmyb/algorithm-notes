package leetcode

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	qs := []struct {
		name string
		in   string
		want bool
	}{
		{"Example 1", "()", true},
		{"Example 2", "()[]{}", true},
		{"Example 3", "(]", false},
		{"Example 4", "([])", true},
		{"Example 5", "([)]", false},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := isValid(q.in)
			if got != q.want {
				t.Fatalf("isValid(%q) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}
