package main

import "fmt"

func main() {
	//定义一个整数，+10将结果输出
	num := 10
	questionOne(&num)
	println("函数执行的结果是：", num)

	//将切片中的元素都乘以2
	arr := []int{1, 2, 5, 8, 4, 6, 9, 7}
	questionTwo(&arr)
	fmt.Println("切片的结果是：", arr)
}

// 题目一
func questionOne(num *int) {
	*num += 10
}

// 题目二
func questionTwo(arr *[]int) {
	for i := range *arr {
		(*arr)[i] *= 2
	}
}
