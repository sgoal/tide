package main

import (
	"fmt"
)

// 冒泡排序函数
func bubbleSort(arr []int) []int {
	// 拷贝一个数组避免修改原数组
	sortedArr := make([]int, len(arr))
	copy(sortedArr, arr)

	n := len(sortedArr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if sortedArr[j] > sortedArr[j+1] {
				sortedArr[j], sortedArr[j+1] = sortedArr[j+1], sortedArr[j]
			}
		}
	}
	return sortedArr
}

func main() {
	arr := []int{5, 2, 9, 1, 5, 6}
	sortedArr := bubbleSort(arr)
	fmt.Println("排序后的数组:", sortedArr)
}
