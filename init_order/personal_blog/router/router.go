package router

import (
	"github.com/gin-gonic/gin"
	"personal_blog/controller"
	"personal_blog/middleware"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	
	// 添加访问日志中间件
	r.Use(middleware.RequestLogger())
	
	auth := r.Group("auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
	}
	api := r.Group("api")
	api.Use(middleware.MiddleWare())
	{
		api.GET("posts", controller.GetPostList)
		api.GET("posts/getById", controller.GetPostById)
		api.POST("posts", controller.CreatePost)
		api.DELETE("posts", controller.DeletePost)
		api.PATCH("posts", controller.UpdatePost)
		api.POST("comments", controller.AddComment)
		api.GET("comments", controller.GetCommentByPost)
	}

	return r
}