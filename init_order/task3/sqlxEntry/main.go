package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
	题目1：使用SQL扩展库进行查询
	假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	要求 ：
	编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

	题目2：实现类型安全映射
	假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	要求 ：
	定义一个 Book 结构体，包含与 books 表对应的字段。
	编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     int
}

type Book struct {
	ID       int
	Title    string
	Authored string `db:"author"`
	Price    int
}

var db *sqlx.DB

func main() {
	initDB()
	createTable()
	//insertData()
	//条件查询数据
	selectByCondition()
	//查询书籍
	selectBooks()

}

func selectBooks() {
	//defer db.Close()
	var books []Book
	if err := db.Select(&books, "select * from books where price > ?", 50); err != nil {
		panic(err)
	}
	fmt.Println(books)
}

func selectByCondition() {
	//defer db.Close()

	depart := "技术部"
	//查询技术部员工
	var list []Employee
	if err := db.Select(&list, "select * from employees where department = ?", depart); err != nil {
		panic(err)
	}
	fmt.Println(list)

	//查询最高工资员工
	var employee Employee
	if err := db.Get(&employee, "select * from employees order by salary desc limit 1"); err != nil {
		panic(err)
	}
	fmt.Println(employee)
}

func insertData() {
	sql := `insert into employees(name, department, salary) values (?,?,?),(?,?,?),(?,?,?),(?,?,?)`
	db.MustExec(sql, "张三", "技术部", 3000, "李四", "技术部", 5000, "王五", "管理部", 6000, "赵六", "技术部", 7000)
	fmt.Println("插入成功")
}

func createTable() {
	sql := `create table if not exists employees(
		id  bigint auto_increment primary key,
		name varchar(10),
		department varchar(15),
		salary bigint
	)`
	db.MustExec(sql)
	fmt.Println("创建成功")
}

func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}
