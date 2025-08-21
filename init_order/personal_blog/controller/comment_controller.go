package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal_blog/global"
	"personal_blog/model"
)

func AddComment(ctx *gin.Context) {
	var receive struct {
		Content string `binding:"required"`
		PostID  uint   `binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&receive); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var comment = model.Comment{Content: receive.Content, PostID: receive.PostID}
	userId, _ := ctx.Get("id")
	comment.UserID = userId.(uint)
	if err := global.DB.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"mgs": "Reviews are successful"})
}

func GetCommentByPost(ctx *gin.Context) {
	var PostId struct {
		PostID uint `binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&PostId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var comments []model.Comment
	if err := global.DB.Where("post_id = ?", PostId.PostID).Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}
