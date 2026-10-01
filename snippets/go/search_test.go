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
