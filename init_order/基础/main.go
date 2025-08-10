package main

import "fmt"

func main() {

	//两数之和
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
}

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, v := range nums {
		if index, exist := numMap[target-v]; exist {
			return []int{index, i}
		}
		numMap[v] = i
	}
	return nil
}
