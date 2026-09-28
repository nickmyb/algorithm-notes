package leetcode

import (
	"reflect"
	"testing"
)

func TestIsAnagram(t *testing.T) {
	qs := []struct {
		name string
		s    string
		t    string
		want bool
	}{
		{"Example 1", "anagram", "nagaram", true},
		{"Example 2", "rat", "car", false},
	}

	impls := map[string]func(string, string) bool{
		"isAnagram":       isAnagram,
		"isAnagramBySort": isAnagramBySort,
	}

	for _, q := range qs {
		for implName, impl := range impls {
			t.Run(q.name+"/"+implName, func(t *testing.T) {
				got := impl(q.s, q.t)

				if !reflect.DeepEqual(got, q.want) {
					t.Fatalf("%v(%v, %v) = %v, want %v", implName, q.s, q.t, got, q.want)
				}
			})
		}
	}
}
