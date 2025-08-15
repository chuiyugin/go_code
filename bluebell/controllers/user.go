package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

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
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans)))
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
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp with invaild param") // 日志
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 1.获取请求参数以及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invaild param", zap.Error(err)) // 日志
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		// 如果不是validator.ValidationErrors类型的错误则直接返回错误信息
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans)))
		return
	}
	// 2.业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err)) // 日志
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}
