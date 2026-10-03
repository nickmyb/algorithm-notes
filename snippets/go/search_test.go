package snippets

import "testing"

// 每个用例都注明它能抓到的错误写法，说不出来的不加。

func TestBinarySearch(t *testing.T) {
	a := []int{1, 3, 5, 7, 9}
	tests := []struct {
		name string
		key  int
		want int
	}{
		// 首元素：抓两个分支写反（key 偏小却收缩左界），会越过目标
		{"首元素", 1, 0},
		// 末元素：抓循环条件写成 lo < hi，lo == hi 时目标还没比较就退出
		{"末元素", 9, 4},
		// 不存在：抓返回 lo（插入位置，那是第 35 题）而不是 -1；
		// 也抓 hi = mid / lo = mid，区间缩不到空会死循环
		{"不存在", 4, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearch(tt.key, a); got != tt.want {
				t.Errorf("binarySearch(%d, %v) = %d, want %d", tt.key, a, got, tt.want)
			}
		})
	}
}

func TestLowerBoundClosed(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		// 有重复：抓把 < 写成 <=，那样返回第一个 > target 的位置（upper bound，这里会得 5）；
		// 也抓循环条件写成 lo < hi，lo == hi 时那一格没染色就退出（这里会得 4）
		{"有重复", []int{5, 7, 7, 8, 8, 10}, 8, 3},
		// 不存在：抓返回最后一个红（hi）而不是第一个蓝（hi + 1）；也抓找不到时返回 -1
		{"不存在", []int{5, 7, 7, 8, 8, 10}, 6, 1},
		// 比所有元素都大：抓 hi 初值写成 len(nums)（混用了左闭右开），mid 会取到 len(nums) 越界
		{"比所有元素都大", []int{5, 7, 7, 8, 8, 10}, 11, 6},
		// 空数组：抓循环后再读 nums[mid] 决定 mid / mid + 1 的写法（第 35 题 searchInsert），空数组会越界；
		// 第 34 题允许 nums 为空
		{"空数组", []int{}, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerBoundClosed(tt.nums, tt.target); got != tt.want {
				t.Errorf("lowerBoundClosed(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
