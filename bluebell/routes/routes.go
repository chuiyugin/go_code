package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"
	"time"

	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	// 使用中间件（令牌桶限流）
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	v1.POST("/signup", controllers.SignHandler)

	// 登录业务路由
	v1.POST("/login", controllers.LoginHandler)

	// 应用JWT认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)

		// 根据时间或分数获取帖子列表
		v1.GET("/posts2", controllers.GetPostListHandler2)

		// 投票
		v1.POST("/vote", controllers.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
