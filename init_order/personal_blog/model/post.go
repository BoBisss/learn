package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string `binding:"required" gorm:"not null"`
	Content  string `binding:"required" gorm:"not null"`
	UserID   uint
	Comments []Comment
}
