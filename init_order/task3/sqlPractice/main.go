package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
	题目1：基本CRUD操作
	假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	要求 ：
	编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

	题目2：事务语句
	假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	要求 ：
	编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

type Student struct {
	ID    int
	Name  string `gorm:"type:varchar(10)"`
	Age   uint8
	Grade string `gorm:"type:varchar(10)"`
}
type Account struct {
	ID      int
	Balance int
}
type Transaction struct {
	ID            int
	FromAccountID int
	ToAccountID   int
	Amount        int
}

func main() {

	//基本CRUD操作
	questitonOne()
	//事务语句
	questionTwo()
}

func questitonOne() {
	db := getDB()
	//创建表结构
	db.AutoMigrate(&Student{})
	student := Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	}
	db.Create(&student)

	var students []Student
	db.Debug().Where("age > ?", 18).Find(&students)
	fmt.Println(students)

	db.Where("name = ?", "张三").Updates(Student{Grade: "四年级"})

	db.Create(&Student{
		Name:  "李四",
		Age:   13,
		Grade: "二年级",
	})

	db.Where("age < ?", 15).Delete(&Student{})
}

func questionTwo() {
	db := getDB()
	//创建账户信息
	db.AutoMigrate(&Account{}, &Transaction{})
	db.Create([]*Account{{Balance: 100},
		{Balance: 200},
		{Balance: 50},
	})

	err := db.Transaction(func(tx *gorm.DB) error {
		//查询数据
		from := Account{}
		db.First(&from, 2)
		to := Account{}
		db.First(&to, 1)
		amount := 100

		//保存记录
		tx.Create(&Transaction{FromAccountID: from.ID,
			ToAccountID: to.ID,
			Amount:      amount})

		return tx.Transaction(func(tx *gorm.DB) error {
			if (from.Balance - amount) < 0 {
				return errors.New("账户余额不足，无法转账！！！")
			}
			//转账过程
			if err := tx.Model(&from).Where("id = ?", from.ID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
				return err
			}
			if err := tx.Model(&to).Where("id = ?", to.ID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("转账成功")
	}
}

func getDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}
