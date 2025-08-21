package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		global.Logger.Warn("Invalid comment creation request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var comment = model.Comment{Content: receive.Content, PostID: receive.PostID}
	userId, _ := ctx.Get("id")
	comment.UserID = userId.(uint)
	if err := global.DB.Create(&comment).Error; err != nil {
		global.Logger.Error("Comment creation failed",
			zap.Error(err),
			zap.Any("user_id", userId),
			zap.Uint("post_id", receive.PostID),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	global.Logger.Info("Comment created successfully",
		zap.Uint("comment_id", comment.ID),
		zap.Any("user_id", userId),
		zap.Uint("post_id", receive.PostID),
	)
	ctx.JSON(http.StatusOK, gin.H{"mgs": "Reviews are successful"})
}

func GetCommentByPost(ctx *gin.Context) {
	var PostId struct {
		PostID uint `binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&PostId); err != nil {
		global.Logger.Warn("Invalid comment query request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var comments []model.Comment
	if err := global.DB.Where("post_id = ?", PostId.PostID).Find(&comments).Error; err != nil {
		global.Logger.Error("Failed to get comments by post",
			zap.Error(err),
			zap.Uint("post_id", PostId.PostID),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	global.Logger.Info("Comments retrieved successfully",
		zap.Uint("post_id", PostId.PostID),
		zap.Int("count", len(comments)),
	)
	ctx.JSON(http.StatusOK, comments)
}
