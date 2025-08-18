package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
	题目1：模型定义
	假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	要求 ：
	使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
	编写Go代码，使用Gorm创建这些模型对应的数据库表。

	题目2：关联查询
	基于上述博客系统的模型定义。
	要求 ：
	编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	编写Go代码，使用Gorm查询评论数量最多的文章信息。

	题目3：钩子函数
	继续使用博客系统的模型。
	要求 ：
	为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

type User struct {
	ID        int
	Name      string
	Post      []Post
	PostCount int
}
type Post struct {
	ID            int
	Title         string
	UserID        int
	Comment       []Comment
	CommentStatus string
}
type Comment struct {
	ID      int
	Content string `gorm:"varchar(255)"`
	PostID  int
}

var db *gorm.DB

func main() {
	getConnect()
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	//插入数据
	/*user := []User{
		{
			Name: "张三",
			Post: []Post{
				{Title: "文章1", Comment: []Comment{
					{Content: "评论1"},
					{Content: "评论2"},
					{Content: "评论3"},
				}}, {Title: "文章2", Comment: []Comment{
					{Content: "评论4"},
					{Content: "评论5"},
					{Content: "评论6"},
					{Content: "评论7"},
				},
				},
			},
		}, {
			Name: "李四",
			Post: []Post{
				{Title: "文章3", Comment: []Comment{
					{Content: "评论8"},
				}}, {Title: "文章4", Comment: []Comment{
					{Content: "评论9"},
					{Content: "评论10"},
				},
				},
			},
		},
	}
	db.Create(&user)*/

	user1 := User{}
	//题目2：关联查询 ：使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	db.Preload("Post.Comment").Where("name = ?", "张三").Find(&user1)
	fmt.Println(user1)
	//题目2：使用Gorm查询评论数量最多的文章信息。
	/*
			SELECT
			p.*,
			u.NAME author,
			t.count comment_count
		FROM
			posts AS p
			JOIN ( SELECT post_id, COUNT(*) count FROM `comments` GROUP BY `post_id` ORDER BY count DESC LIMIT 1 ) t ON t.post_id = p.id
			JOIN users u ON u.id = p.user_id
		ORDER BY
			`p`.`title`
			LIMIT 1
	*/
	type Result struct {
		Title        string
		Author       string
		CommentCount int
	}
	var result Result
	query1 := db.Select("post_id, COUNT(*) count").Table("comments").Group("post_id").Order("count DESC").Limit(1)
	db.Table("posts as p").Select("p.*,u.name author,t.count comment_count").Joins("join (?) t ON t.post_id = p.id join users u ON u.id = p.user_id", query1).First(&result)
	fmt.Println(result)

	//测试comment钩子
	db.Where("post_id = ?", 3).Delete(&Comment{PostID: 3})

}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	userID := p.UserID

	return tx.Table("users").Where("id = ?", userID).UpdateColumn("post_count", gorm.Expr("post_count+1")).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	postId := c.PostID
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", postId).Count(&count)
	if count == 0 {
		result := tx.Table("posts").Where("id = ?", postId).UpdateColumn("comment_status", "无评论")
		return result.Error
	}
	return nil
}

func getConnect() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return err
}
