package leetcode

// ===== 本地接线区 =====

import "math/rand"

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====

// region @snippets/go/sort.go

// Quick 按 cmp 的顺序对切片做原地快速排序。cmp(a, b) 分别用负数、零、正数表示
// a 小于、等于、大于 b；只看符号，不要求返回 -1、0、1。
// 必须通过构造函数设置 cmp，零值 Quick 不能排序多于一个元素。
type Quick[T any] struct {
	arr []T
	cmp func(T, T) int
}

// NewQuick 创建一个使用 cmp 的空排序器；cmp 不能为 nil。
func NewQuick[T any](cmp func(T, T) int) *Quick[T] {
	return &Quick[T]{
		cmp: cmp,
		arr: make([]T, 0),
	}
}

// NewQuickCap 创建一个空排序器，预留 maxN 个元素的容量。
// maxN 只是容量提示，Append 超过它时仍会扩容。
func NewQuickCap[T any](cmp func(T, T) int, maxN int) *Quick[T] {
	return &Quick[T]{
		cmp: cmp,
		arr: make([]T, 0, maxN),
	}
}

// NewQuickFrom 直接使用 x 创建排序器，不复制切片。O(1)。
// Sort 会原地打乱并排序 x；Append 也可能写入 x 的底层数组。
func NewQuickFrom[T any](cmp func(T, T) int, x []T) *Quick[T] {
	return &Quick[T]{arr: x, cmp: cmp}
}

// Append 向待排序切片末尾添加一个元素。
func (q *Quick[T]) Append(x T) {
	q.arr = append(q.arr, x)
}

// Arr 返回待排序切片本身，修改其中的元素也会修改 Quick 中的数据。
func (q *Quick[T]) Arr() []T {
	return q.arr
}

// partition 以 arr[lo] 为枢轴划分闭区间 arr[lo..hi]，返回枢轴的最终位置。
// 左侧元素不大于枢轴，右侧元素不小于枢轴；两端扫描分别在 hi 和 lo 停住。
func (q *Quick[T]) partition(arr []T, lo, hi int) int {
	pivot := arr[lo]

	i, j := lo, hi+1
	for {
		for {
			i++
			if q.cmp(pivot, arr[i]) <= 0 {
				break
			}
			if i == hi {
				break
			}
		}
		for {
			j--
			if q.cmp(arr[j], pivot) <= 0 {
				break
			}
			if j == lo {
				break
			}
		}

		if j <= i {
			break
		}

		arr[i], arr[j] = arr[j], arr[i]

	}

	arr[lo], arr[j] = arr[j], arr[lo]

	return j
}

// quickSortRange 递归排序闭区间 arr[lo..hi]。
func (q *Quick[T]) quickSortRange(arr []T, lo, hi int) {
	if hi <= lo {
		return
	}

	j := q.partition(arr, lo, hi)
	q.quickSortRange(arr, lo, j-1)
	q.quickSortRange(arr, j+1, hi)
}

// Sort 先随机打乱当前切片，再原地排序；排序不稳定。
// 平均时间 O(n log n)，最坏 O(n²)；递归栈平均 O(log n)，最坏 O(n)。
func (q *Quick[T]) Sort() {
	rand.Shuffle(len(q.arr), func(i, j int) {
		q.arr[i], q.arr[j] = q.arr[j], q.arr[i]
	})
	q.quickSortRange(q.arr, 0, len(q.arr)-1)
}

// endregion @snippets/go/sort.go

func merge(intervals [][]int) [][]int {
	quick := NewQuickFrom(func(i, j []int) int {
		if i[0] == j[0] {
			return i[1] - j[1]
		}
		return i[0] - j[0]
	}, intervals)
	quick.Sort()

	sortedIntervals := quick.Arr()

	var ret [][]int

	ret = append(ret, sortedIntervals[0])

	index := 1
	for index < len(intervals) {
		interval := sortedIntervals[index]
		lo, hi := interval[0], interval[1]

		length := len(ret)
		if lo <= ret[length-1][1] {
			ret[length-1][1] = max(ret[length-1][1], hi)
		} else {
			ret = append(ret, interval)
		}

		index += 1
	}

	return ret
}

func mergeRange(intervals [][]int) [][]int {
	m := make(map[int]int)
	var ret [][]int

	for _, interval := range intervals {
		l, h := interval[0], interval[1]
		low, hi := l, h
		isNew := true

		if len(m) == 0 {
			m[l] = h
			continue
		}

		for k, v := range m {
			// 和interval相关的range都统一删除,插入新range

			if k <= low && hi <= v {
				// 被包含在已有区间,保持现在的范围结果即可
				isNew = false
				break
			} else if low < k && v < hi {
				// 包含已有区间
			} else if k <= low && low <= v {
				if v <= hi {
					low = k
				}
			} else if k <= hi && hi <= v {
				if low <= k {
					hi = v
				}
			} else {
				// 无交叉新范围
				// 没有旧范围需要删除
				continue
			}

			delete(m, k)
		}

		if isNew {
			m[low] = hi
		}
	}

	for k, v := range m {
		ret = append(ret, []int{k, v})
	}
	return ret
}
