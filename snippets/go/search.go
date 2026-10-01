package snippets

// binarySearch 二分查找
func binarySearch(key int, a []int) int {
	lo, hi := 0, len(a)-1

	for lo <= hi {
		mid := (lo + hi) / 2

		if key == a[mid] {
			return mid
		} else if key < a[mid] {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return -1
}
