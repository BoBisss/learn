package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `binding:"required" gorm:"not null"`
	UserID  uint
	PostID  uint
}
