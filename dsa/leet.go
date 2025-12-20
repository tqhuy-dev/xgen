package dsa

import "math"

// LeetCode238ProductOfArrayExceptSelf [1,2,3,4] -> [24,12,8,6]
func LeetCode238ProductOfArrayExceptSelf(array []int) []int {
	res := make([]int, len(array))
	prefix := 1
	for i := range array {
		res[i] = prefix
		prefix *= array[i]
	}
	suffix := 1
	for i := len(array) - 1; i >= 0; i-- {
		res[i] *= suffix
		suffix *= array[i]
	}
	return res
}

// LeetCode1925CountSquareSumTriples 1 <= n <= 100
func LeetCode1925CountSquareSumTriples(n int) int {
	var count int
	for i := 1; i <= n; i++ {
		a2 := i * i
		for j := i + 1; j <= n; j++ {
			c2 := j*j + a2
			c := int(math.Sqrt(float64(c2)))
			if c <= n && c*c == c2 {
				count += 2
			}
		}
	}
	return count
}

func LeetCode3583CountSpecialTriplets(nums []int) int {
	const MAXV = 100001
	const MOD = 1000000007
	var memo, duplets [MAXV]int

	res := 0
	for _, number := range nums {
		res = (res + duplets[number]) % MOD
		twon := number * 2
		if twon < MAXV {
			duplets[twon] = (duplets[twon] + memo[twon]) % MOD
		}

		memo[number]++
	}

	return res
}

func LeetCode239SlidingWindowMaximum(nums []int, k int) []int {
	result := make([]int, len(nums)-k+1)

	deque := make([]int, 0)
	for i, num := range nums {
		if len(deque) > 0 && deque[0] == i-k {
			deque = deque[1:]
		}
		for len(deque) > 0 && nums[deque[len(deque)-1]] < num {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		if i >= k-1 {
			result[i-k+1] = nums[deque[0]]
		}
	}
	return result
}

func LeetCode643MaximumAverageSubarrayI(nums []int, k int) float64 {
	sum := 0
	for i := 0; i < k; i++ {
		sum += nums[i]
	}
	maxSum := sum
	for i := k; i < len(nums); i++ {
		sum += nums[i] - nums[i-k]
		maxSum = max(maxSum, sum)
	}
	return float64(maxSum) / float64(k)
}

func LeetCode713SubarrayProductLessThanK(nums []int, k int) int {
	var res int
	left := 0
	value := 1
	for i, num := range nums {
		value *= num
		for value >= k && left < i {
			value /= nums[left]
			left++
		}
		if value < k {
			res += i - left + 1
		}
	}
	return res
}

func LeetCode1004MaxConsecutiveOnesIII(nums []int, k int) int {
	var maxLen int
	zeroCount := 0
	left := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] == 0 {
			zeroCount++
		}
		if zeroCount > k {
			for zeroCount > k {
				if nums[left] == 0 {
					zeroCount--
				}
				left++
			}
		}

		value := i - left + 1
		if value > maxLen {
			maxLen = value
		}
	}
	return maxLen
}
