package snippets

import "math/rand"

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
