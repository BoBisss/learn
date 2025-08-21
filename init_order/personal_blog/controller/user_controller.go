package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"personal_blog/global"
	"personal_blog/model"
	"personal_blog/utils"
)

func Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	password, err := utils.HashPassWord(user.PassWord)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	user.PassWord = password

	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "user created success")
}

func Login(ctx *gin.Context) {
	var receive struct {
		Username string
		Password string
	}
	var user model.User
	//转换数据
	if err := ctx.ShouldBindJSON(&receive); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	//查询数据库
	if err := global.DB.Where("user_name = ?", receive.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "The username does not exist"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}
	//密码校验
	if !utils.ValidPassWord(receive.Password, user.PassWord) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Wrong password"})
		return
	}
	//生成token
	token, err := utils.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
