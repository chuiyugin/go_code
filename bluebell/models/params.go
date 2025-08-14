package models

// 定义请求参数的结构体

type ParamSignUP struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Re_password string `json:"re_password" binding:"required,eqfield=Password"`
}
