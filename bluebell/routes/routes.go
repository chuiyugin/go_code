package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	// 使用中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务逻辑
	r.POST("/signup", controllers.SignHandler)

	// 登录业务逻辑
	r.POST("/login", controllers.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录用户，判断请求头中是否有 有效的JWT token
		c.String(http.StatusOK, "pong")
	})
	return r
}
