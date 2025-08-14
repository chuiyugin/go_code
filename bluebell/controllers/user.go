package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamSignUP)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invaild param", zap.Error(err)) // 日志
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.Re_password) == 0 || p.Password != p.Re_password {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invaild param") // 日志
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	fmt.Println(p)
	// 2.业务处理
	logic.SignUp(p)
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
