package leetcode

import (
	"reflect"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	qs := []struct {
		name string
		in   string
		want bool
	}{
		{"Example 1", "A man, a plan, a canal: Panama", true},
		{"Example 2", "race a car", false},
		{"Example 3", " ", true},
		{"数字不能跳过 数字与字母相差32", "0P", false},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := isPalindrome(q.in)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("isPalindrome(%v) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}
