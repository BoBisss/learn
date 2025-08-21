package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"personal_blog/global"
	"time"
)

// RequestLogger 记录请求日志的中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 获取用户信息（如果已登录）
		userID, userExists := c.Get("id")
		username, _ := c.Get("username")

		// 记录访问日志
		if userExists {
			global.Logger.Info("Authenticated request",
				zap.String("client_ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Any("user_id", userID),
				zap.String("username", username.(string)),
				zap.Duration("process_time", endTime.Sub(startTime)),
			)
		} else {
			global.Logger.Info("Anonymous request",
				zap.String("client_ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("process_time", endTime.Sub(startTime)),
			)
		}
	}
}
