package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func GetPostById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var post model.Post
	if err := global.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"err": "data id empty"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, post)

}

func CreatePost(ctx *gin.Context) {
	var post model.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	uid, _ := ctx.Get("id")
	fmt.Println(reflect.TypeOf(uid))
	post.UserID = uid.(uint)

	if err := global.DB.Create(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "create success"})
}

func DeletePost(ctx *gin.Context) {
	var receive struct {
		Id uint
	}
	if err := ctx.ShouldBindJSON(&receive); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var userId uint
	if err := global.DB.Model(&model.Post{}).Select("user_id").Where("id = ?", receive.Id).Scan(&userId).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if value, _ := ctx.Get("id"); value != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "You do not have permission to manipulate this data"})
		return
	}

	if err := global.DB.Delete(&model.Post{}, receive.Id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "delete success"})

}

func UpdatePost(ctx *gin.Context) {
	var post model.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if post.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "id cannot be empty"})
		return
	}
	var userId uint
	if err := global.DB.Model(&model.Post{}).Select("user_id").Where("id = ?", post.ID).Scan(&userId).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if value, _ := ctx.Get("id"); value != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "You do not have permission to manipulate this data"})
		return
	}
	if err := global.DB.Where("id = ?", post.ID).Updates(model.Post{
		Title:   post.Title,
		Content: post.Content,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "update success"})

}
