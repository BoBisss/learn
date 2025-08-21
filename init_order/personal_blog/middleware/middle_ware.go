package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal_blog/utils"
)

func MiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token is missing"})
			ctx.Abort()
			return
		}
		data, err := utils.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("id", uint(data[0].(float64)))
		ctx.Set("username", data[1])
		ctx.Next()
	}
}
