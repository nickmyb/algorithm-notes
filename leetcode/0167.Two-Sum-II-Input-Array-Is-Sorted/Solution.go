package leetcode

// ===== 本地接线区 =====

// ===== 以下是题解本体，必须和提交到 LeetCode 的代码一字不差 =====
func twoSum(numbers []int, target int) []int {
	i := 0
	j := len(numbers) - 1

	for i < j {
		sum := numbers[i] + numbers[j]

		if sum == target {
			return []int{i + 1, j + 1}
		}
		if sum < target {
			i += 1
		}
		if sum > target {
			j -= 1
		}
	}

	return nil
}

func twoSumTraversal(numbers []int, target int) []int {
	length := len(numbers)
	i := 0

	for i < length-1 {
		j := i + 1
		lastSum := numbers[i] + numbers[j]

		for j < length {
			sum := numbers[i] + numbers[j]

			if sum == target {
				return []int{i + 1, j + 1}
			}
			if lastSum < target && sum > target {
				break
			}

			lastSum = sum

			j++
		}

		i++
	}

	return nil
}
