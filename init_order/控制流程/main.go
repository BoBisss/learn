package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 找出只出现一次的数字
	arr := []int{1, 2, 2, 1, 5, 8, 9, 9, 8}
	fmt.Println(findOnlyNumber(arr))
	//判断是否是回文数
	fmt.Println(isPalindrome(12321))
	fmt.Println(isPalindrome(123421))

}

// 找出只出现一次的数字
func findOnlyNumber(nums []int) int {
	countMap := make(map[int]int, 9)
	for _, value := range nums {
		countMap[value]++
	}
	for key, value := range countMap {
		if value == 1 {
			return key
		}
	}
	return -1
}

func isPalindrome(x int) bool {
	num := strconv.Itoa(x)
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		if num[i] != num[j] {
			return false
		}
	}
	return true
}
