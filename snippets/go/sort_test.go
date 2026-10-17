package snippets

import (
	"cmp"
	"slices"
	"testing"
)

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"nil 切片", nil, nil},
		{"空切片", []int{}, []int{}},
		{"单元素", []int{7}, []int{7}},
		{"两个元素", []int{2, 1}, []int{1, 2}},
		{"重复值", []int{3, 1, 3, 2, 1}, []int{1, 1, 2, 3, 3}},
		{"全部相等", []int{4, 4, 4, 4}, []int{4, 4, 4, 4}},
		{"正负混合", []int{0, -2, 3, -1}, []int{-2, -1, 0, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.in)
			q := NewQuickFrom(cmp.Compare[int], in)
			q.Sort()
			if !slices.Equal(q.Arr(), tt.want) {
				t.Errorf("Arr() = %v, want %v", q.Arr(), tt.want)
			}
			if !slices.Equal(in, tt.want) {
				t.Errorf("Sort() changed input to %v, want %v", in, tt.want)
			}
		})
	}
}

func TestQuickAppendAndResort(t *testing.T) {
	q := NewQuickCap(cmp.Compare[int], 2)
	for _, x := range []int{3, 1, 2} {
		q.Append(x)
	}
	q.Sort()
	if want := []int{1, 2, 3}; !slices.Equal(q.Arr(), want) {
		t.Fatalf("first Sort(): got %v, want %v", q.Arr(), want)
	}

	q.Append(0)
	q.Sort()
	if want := []int{0, 1, 2, 3}; !slices.Equal(q.Arr(), want) {
		t.Errorf("second Sort(): got %v, want %v", q.Arr(), want)
	}
}

func TestQuickFromSortsInputInPlace(t *testing.T) {
	in := []int{3, 1, 2}
	q := NewQuickFrom(cmp.Compare[int], in)
	q.Sort()
	if want := []int{1, 2, 3}; !slices.Equal(in, want) {
		t.Fatalf("NewQuickFrom input after Sort() = %v, want %v", in, want)
	}

	in[0] = 99
	if q.Arr()[0] != 99 {
		t.Errorf("input mutation not reflected in Quick: got %v", q.Arr())
	}
	q.Arr()[0] = 10
	if in[0] != 10 {
		t.Errorf("Arr() mutation not reflected in input: in[0] = %d", in[0])
	}
}

func TestQuickCustomComparator(t *testing.T) {
	type item struct {
		name string
		rank int
	}
	in := []item{{"low", 1}, {"high", 3}, {"middle", 2}}
	q := NewQuickFrom(func(a, b item) int { return cmp.Compare(b.rank, a.rank) }, in)
	q.Sort()
	want := []item{{"high", 3}, {"middle", 2}, {"low", 1}}
	if !slices.Equal(q.Arr(), want) {
		t.Errorf("Sort() = %v, want %v", q.Arr(), want)
	}
	if !slices.Equal(in, want) {
		t.Errorf("Sort() changed input to %v, want %v", in, want)
	}
}
