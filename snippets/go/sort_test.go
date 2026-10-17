package snippets

import (
	"slices"
	"testing"
)

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		lo   int
		hi   int
		want []int
	}{
		// 空切片：检查入口在分区和打乱前返回。
		{"空切片", []int{}, 0, -1, []int{}},
		// 单元素：检查长度为一时不进入分区。
		{"单元素", []int{7}, 0, 0, []int{7}},
		// 基准最小：检查右侧扫描能在 lo 停下。
		{"基准最小", []int{1, 2}, 0, 1, []int{1, 2}},
		// 重复值：检查等于基准的元素不会导致划分停滞或错排。
		{"重复值", []int{2, 1, 2, 3}, 0, 3, []int{1, 2, 2, 3}},
		// 正负混合：检查比较直接按整数大小进行，不依赖非负输入。
		{"正负混合", []int{0, -2, 3, -1}, 0, 3, []int{-2, -1, 0, 3}},
		// 子区间：检查排序不会改动 lo、hi 以外的元素。
		{"子区间", []int{9, 2, 3, 1, 8}, 1, 3, []int{9, 1, 2, 3, 8}},
		// 基准最大：检查左侧扫描不会访问 hi+1。
		{"基准最大", []int{2, 1}, 0, 1, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Clone(tt.in)
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("quickSort(%v, %d, %d) panic: %v", tt.in, tt.lo, tt.hi, p)
				}
			}()
			quickSort(got, tt.lo, tt.hi)
			if !slices.Equal(got, tt.want) {
				t.Errorf("quickSort(%v, %d, %d) = %v, want %v", tt.in, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}
