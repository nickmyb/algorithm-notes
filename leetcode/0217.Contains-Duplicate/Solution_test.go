package leetcode

import (
	"reflect"
	"testing"
)

func TestContainsDuplicate(t *testing.T) {
	qs := []struct {
		name string
		in   []int
		want bool
	}{
		{"Example 1", []int{1, 2, 3, 1}, true},
		{"Example 2", []int{1, 2, 3, 4}, false},
		{"Example 3", []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := containsDuplicate(q.in)
			if !reflect.DeepEqual(got, q.want) {
				t.Fatalf("containsDuplicate(%v) = %v, want %v", q.in, got, q.want)
			}
		})
	}
}
