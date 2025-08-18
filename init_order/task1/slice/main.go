package main

import "fmt"

func main() {

	//删除有序数组中的重复项
	fmt.Println("新长度是：", removeDuplicates([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))

	//合并区间
	fmt.Println(merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	index := 0
	for i, _ := range nums {
		if i == 0 {
			continue
		}
		if nums[index] != nums[i] {
			index++
			nums[index] = nums[i]
		}
	}
	return index + 1
}

// 合并区间
func merge(intervals [][]int) [][]int {
	var res [][]int

	for i, _ := range intervals {
		for j := i + 1; j < len(intervals); j++ {
			if intervals[i][0] > intervals[j][0] {
				intervals[i], intervals[j] = intervals[j], intervals[i]
			}
		}
	}
	for i := 0; i < len(intervals); i++ {
		if i == 0 {
			res = append(res, intervals[i])
		} else {
			if intervals[i][0] <= res[len(res)-1][1] {
				res[len(res)-1][1] = max(res[len(res)-1][1], intervals[i][1])
			} else {
				res = append(res, intervals[i])
			}
		}
	}
	return res

}
