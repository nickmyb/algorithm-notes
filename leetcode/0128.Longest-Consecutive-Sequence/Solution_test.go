package leetcode

import "testing"

func TestLongestConsecutive(t *testing.T) {
	qs := []struct {
		name string
		nums []int
		want int
	}{
		{
			name: "Example 1",
			nums: []int{100, 4, 200, 1, 3, 2},
			want: 4,
		},
		{
			name: "Example 2",
			nums: []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1},
			want: 9,
		},
		{
			name: "Example 3",
			nums: []int{1, 0, 1, 2},
			want: 3,
		},
		{
			name: "Empty input",
			nums: []int{},
			want: 0,
		},
		{
			name: "Negative values",
			nums: []int{-2, -1, 0, 1, 5},
			want: 4,
		},
	}

	for _, q := range qs {
		t.Run(q.name, func(t *testing.T) {
			got := longestConsecutive(q.nums)
			if got != q.want {
				t.Fatalf("longestConsecutive(%v) = %d, want %d", q.nums, got, q.want)
			}
		})
	}
}
