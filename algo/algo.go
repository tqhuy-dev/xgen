package algo

import (
	"sort"
	"strings"
)

// Input : arr = [2,0,2,1,1,0]
// Output : [0,0,1,1,2,2]
func DutchNationFlag(arr []int) []int {
	low, mid, high := 0, 0, len(arr)-1
	for mid <= high {
		switch arr[mid] {
		case 2:
			arr[mid], arr[high] = arr[high], arr[mid]
			high -= 1
		case 0:
			arr[mid], arr[low] = arr[low], arr[mid]
			low += 1
			mid += 1
		case 1:
			mid += 1
		}
	}
	return arr
}

// Input : arr = [2,1,4,5,1,0,3,4,1,2]
// Output : [0 1 1 1 2 2 3 4 4 5]

func CountingSort(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	mapCount := make(map[int]int)
	for _, v := range arr {
		mapCount[v]++
	}
	keys := make([]int, len(mapCount))
	i := 0
	for k := range mapCount {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	newArr := make([]int, 0, len(arr))
	for _, k := range keys {
		for j := 0; j < mapCount[k]; j++ {
			newArr = append(newArr, k)
		}
	}
	return newArr
}

func FindAllIndexes(s, tmp string) []int {
	result := []int{}
	start := 0

	for {
		idx := strings.Index(s[start:], tmp)
		if idx == -1 {
			break
		}

		realIdx := start + idx

		// thêm tất cả index liên tục
		for i := 0; i < len(tmp); i++ {
			result = append(result, realIdx+i)
		}

		// tiếp tục tìm từ sau vị trí vừa match
		start = realIdx + 1
	}

	return result
}

func MinExtraCharacterString(s string, dictionary []string) int {
	n := len(s)
	dp := make([]int, n+1)

	// dp[n] = 0 (base case), tự động có vì slice default = 0

	for i := n - 1; i >= 0; i-- {
		// 1. Assume s[i] is an extra character
		dp[i] = 1 + dp[i+1]

		// 2. Try matching dictionary words
		for _, w := range dictionary {
			wl := len(w)
			if i+wl <= n && s[i:i+wl] == w {
				if dp[i+wl] < dp[i] {
					dp[i] = dp[i+wl]
				}
			}
		}
	}

	return dp[0]
}
