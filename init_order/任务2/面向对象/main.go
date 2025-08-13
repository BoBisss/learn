package main

import "fmt"

type Shape interface {
	Perimeter()
	Area()
}

type Rectangle struct {
}

type Circle struct {
}

func (rectangle *Rectangle) Perimeter() {
	fmt.Println("Rectangle实现了Perimeter方法")
}
func (rectangle *Rectangle) Area() {
	fmt.Println("Rectangle实现了Area方法")
}

func (circle *Circle) Perimeter() {
	fmt.Println("Circle实现了Perimeter方法")
}
func (circle *Circle) Area() {
	fmt.Println("Circle实现了Area方法")
}

type Person struct {
	Name string
	Age  int8
}
type Employee struct {
	EmployeeID string
	Person
}

func (e *Employee) PrintInfo() {
	fmt.Println("当前员工的员工ID是：", e.EmployeeID)
	fmt.Println("当前员工的名字是：", e.Name)
	fmt.Println("当前员工的年龄是：", e.Age)
}

func main() {
	/*
		题目 ：定义一个 Shape 接口，包含 Area() 和 Area() 两个方法。
		然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
		在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
		考察点 ：接口的定义与实现、面向对象编程风格。
	*/
	var rectangle Shape = &Rectangle{}
	var circle Shape = &Circle{}
	rectangle.Perimeter()
	rectangle.Area()
	circle.Perimeter()
	circle.Area()

	/*
		题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
		再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
		为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
		考察点 ：组合的使用、方法接收者。
	*/
	employee := Employee{
		EmployeeID: "12874665",
		Person: Person{
			Name: "张三",
			Age:  25,
		},
	}
	employee.PrintInfo()

}
