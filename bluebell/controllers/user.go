package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamSignUP)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invaild param", zap.Error(err)) // 日志
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		// 如果不是validator.ValidationErrors类型的错误则直接返回错误信息
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{

			"msg": RemoveTopStruct(errs.Translate(trans)), // 将错误信息翻译成中文并去除结构体的信息
		})
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.Re_password) == 0 || p.Password != p.Re_password {
	// 	// 请求参数有误，直接返回响应
	// 	zap.L().Error("SignUp with invaild param") // 日志
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }
	fmt.Println(p)
	// 2.业务处理
	logic.SignUp(p)
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
