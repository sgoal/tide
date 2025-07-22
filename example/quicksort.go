package main

import (
	"fmt"
)

// 快速排序函数
func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[len(arr)/2]
	left := []int{}
	right := []int{}
	middle := []int{}
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v > pivot {
			right = append(right, v)
		} else {
			middle = append(middle, v)
		}
	}
	left = quickSort(left)
	right = quickSort(right)
	return append(append(left, middle...), right...)
}

func main() {
	arr := []int{3, 6, 8, 10, 1, 2, 1}
	sortedArr := quickSort(arr)
	fmt.Println("排序后的数组:", sortedArr)
}
