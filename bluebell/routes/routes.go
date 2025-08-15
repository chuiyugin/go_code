package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
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

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}
