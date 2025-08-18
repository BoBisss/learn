package main

import (
	"fmt"
)

func main() {

	//是否是有效括号
	fmt.Println(isValid("()[]{}"))
	// 找出最长公共前缀
	prefix := longestCommonPrefix([]string{"flower", "flow", "flight"})
	fmt.Println(prefix)
}

// 是否是有效括号
func isValid(s string) bool {
	stack := make([]rune, 0)
	for _, v := range s {
		switch v {
		case '(', '[', '{':
			stack = append(stack, v)
		case ')':
			if stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		case ']':
			if stack[len(stack)-1] == '[' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		case '}':
			if stack[len(stack)-1] == '{' {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		}
	}
	return true
}

// 找出最长公共前缀
func longestCommonPrefix(strs []string) string {
	num := 0
	strMap := make(map[int]string)
	for i, s := range strs {
		for j, v := range s {
			if i == 0 {
				strMap[j] = string(v)
				continue
			}
			if strMap[j] != string(v) {
				num = j
				break
			}
		}
	}
	return strs[0][:num]
}
