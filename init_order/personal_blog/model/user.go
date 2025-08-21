package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `binding:"required" gorm:"unique;not nul"`
	PassWord string `binding:"required" gorm:"not null"`
	Email    string `binding:"required" gorm:"unique;not null"`
	Contents []Post
	Comments []Comment
}
