package snippets

import "math/rand"

func quickSort(arr []int, lo, hi int) {
	if hi <= lo {
		return
	}

	rand.Shuffle(hi-lo+1, func(i, j int) {
		arr[lo+i], arr[lo+j] = arr[lo+j], arr[lo+i]
	})
	quickSortRange(arr, lo, hi)
}

func quickSortRange(arr []int, lo, hi int) {
	if hi <= lo {
		return
	}

	j := partition(arr, lo, hi)
	quickSortRange(arr, lo, j-1)
	quickSortRange(arr, j+1, hi)
}

func partition(arr []int, lo, hi int) int {
	pivot := arr[lo]

	i, j := lo, hi+1
	for {
		for {
			i++
			if pivot <= arr[i] {
				break
			}
			if i == hi {
				break
			}
		}
		for {
			j--
			if arr[j] <= pivot {
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
