package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"personal_blog/global"
	"personal_blog/model"
	"reflect"
	"strconv"
)

func GetPostList(ctx *gin.Context) {
	var list []model.Post
	if err := global.DB.Find(&list).Error; err != nil {
		global.Logger.Error("Failed to get post list",
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	global.Logger.Info("Post list retrieved successfully",
		zap.Int("count", len(list)),
	)
	ctx.JSON(http.StatusOK, list)
}

func GetPostById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		global.Logger.Warn("Invalid post ID requested",
			zap.String("id", idStr),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var post model.Post
	if err := global.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("Post not found",
				zap.Uint("post_id", uint(id)),
				zap.String("client_ip", ctx.ClientIP()),
			)
			ctx.JSON(http.StatusNotFound, gin.H{"err": "data id empty"})
			return
		} else {
			global.Logger.Error("Database error while getting post",
				zap.Error(err),
				zap.Uint("post_id", uint(id)),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	// 获取用户信息
	userID, _ := ctx.Get("id")
	username, _ := ctx.Get("username")

	global.Logger.Info("Post retrieved successfully",
		zap.Uint("post_id", uint(id)),
		zap.Any("requested_by_user_id", userID),
		zap.String("requested_by_username", username.(string)),
	)
	ctx.JSON(http.StatusOK, post)

}

func CreatePost(ctx *gin.Context) {
	var post model.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		global.Logger.Warn("Invalid post creation request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	uid, _ := ctx.Get("id")
	fmt.Println(reflect.TypeOf(uid))
	post.UserID = uid.(uint)

	if err := global.DB.Create(&post).Error; err != nil {
		global.Logger.Error("Post creation failed",
			zap.Error(err),
			zap.Any("user_id", uid),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	global.Logger.Info("Post created successfully",
		zap.Uint("post_id", post.ID),
		zap.Any("created_by_user_id", uid),
	)

	ctx.JSON(http.StatusOK, gin.H{"msg": "create success"})
}

func DeletePost(ctx *gin.Context) {
	var receive struct {
		Id uint
	}
	if err := ctx.ShouldBindJSON(&receive); err != nil {
		global.Logger.Warn("Invalid post deletion request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var userId uint
	if err := global.DB.Model(&model.Post{}).Select("user_id").Where("id = ?", receive.Id).Scan(&userId).Error; err != nil {
		global.Logger.Error("Failed to get post owner",
			zap.Error(err),
			zap.Uint("post_id", receive.Id),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if value, _ := ctx.Get("id"); value != userId {
		global.Logger.Warn("Unauthorized post deletion attempt",
			zap.Any("user_id", value),
			zap.Uint("post_id", receive.Id),
			zap.Uint("post_owner_id", userId),
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "You do not have permission to manipulate this data"})
		return
	}

	if err := global.DB.Delete(&model.Post{}, receive.Id).Error; err != nil {
		global.Logger.Error("Post deletion failed",
			zap.Error(err),
			zap.Uint("post_id", receive.Id),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// 获取用户信息用于日志记录
	deletedBy, _ := ctx.Get("id")
	global.Logger.Info("Post deleted successfully",
		zap.Uint("post_id", receive.Id),
		zap.Any("deleted_by_user_id", deletedBy),
	)
	ctx.JSON(http.StatusOK, gin.H{"msg": "delete success"})

}

func UpdatePost(ctx *gin.Context) {
	var post model.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		global.Logger.Warn("Invalid post update request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if post.ID == 0 {
		global.Logger.Warn("Post update failed - missing ID",
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "id cannot be empty"})
		return
	}
	var userId uint
	if err := global.DB.Model(&model.Post{}).Select("user_id").Where("id = ?", post.ID).Scan(&userId).Error; err != nil {
		global.Logger.Error("Failed to get post owner for update",
			zap.Error(err),
			zap.Uint("post_id", post.ID),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if value, _ := ctx.Get("id"); value != userId {
		global.Logger.Warn("Unauthorized post update attempt",
			zap.Any("user_id", value),
			zap.Uint("post_id", post.ID),
			zap.Uint("post_owner_id", userId),
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "You do not have permission to manipulate this data"})
		return
	}
	if err := global.DB.Where("id = ?", post.ID).Updates(model.Post{
		Title:   post.Title,
		Content: post.Content,
	}).Error; err != nil {
		global.Logger.Error("Post update failed",
			zap.Error(err),
			zap.Uint("post_id", post.ID),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// 获取用户信息用于日志记录
	updatedBy, _ := ctx.Get("id")
	global.Logger.Info("Post updated successfully",
		zap.Uint("post_id", post.ID),
		zap.Any("updated_by_user_id", updatedBy),
	)
	ctx.JSON(http.StatusOK, gin.H{"msg": "update success"})

}
