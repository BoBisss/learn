package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"personal_blog/global"
	"personal_blog/model"
	"personal_blog/utils"
)

func Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		global.Logger.Warn("Invalid user registration request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	password, err := utils.HashPassWord(user.PassWord)
	if err != nil {
		global.Logger.Error("Password hashing failed",
			zap.Error(err),
			zap.String("username", user.UserName),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	user.PassWord = password

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("User creation failed",
			zap.Error(err),
			zap.String("username", user.UserName),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	global.Logger.Info("User registered successfully",
		zap.String("username", user.UserName),
		zap.String("email", user.Email),
	)

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
		global.Logger.Warn("Invalid login request",
			zap.Error(err),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	global.Logger.Info("User login attempt",
		zap.String("username", receive.Username),
		zap.String("client_ip", ctx.ClientIP()),
	)

	//查询数据库
	if err := global.DB.Where("user_name = ?", receive.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("Login failed - user not found",
				zap.String("username", receive.Username),
				zap.String("client_ip", ctx.ClientIP()),
			)
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "The username does not exist"})
			return
		} else {
			global.Logger.Error("Database error during login",
				zap.Error(err),
				zap.String("username", receive.Username),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}
	//密码校验
	if !utils.ValidPassWord(receive.Password, user.PassWord) {
		global.Logger.Warn("Login failed - invalid password",
			zap.String("username", receive.Username),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Wrong password"})
		return
	}
	//生成token
	token, err := utils.GenerateToken(user)
	if err != nil {
		global.Logger.Error("Token generation failed",
			zap.Error(err),
			zap.String("username", receive.Username),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to generate token"})
		return
	}

	global.Logger.Info("User login successful",
		zap.String("username", receive.Username),
		zap.String("client_ip", ctx.ClientIP()),
	)

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
