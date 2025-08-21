package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"personal_blog/global"
	"personal_blog/utils"
)

func MiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			global.Logger.Warn("Missing authorization token",
				zap.String("client_ip", ctx.ClientIP()),
				zap.String("path", ctx.Request.URL.Path),
			)
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token is missing"})
			ctx.Abort()
			return
		}
		data, err := utils.ParseToken(token)
		if err != nil {
			global.Logger.Warn("Invalid token",
				zap.Error(err),
				zap.String("client_ip", ctx.ClientIP()),
				zap.String("token", token),
			)
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("id", uint(data[0].(float64)))
		ctx.Set("username", data[1])
		
		global.Logger.Info("Token validated successfully",
			zap.String("username", data[1].(string)),
			zap.String("path", ctx.Request.URL.Path),
		)
		ctx.Next()
	}
}