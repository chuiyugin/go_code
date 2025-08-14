package models

// 定义请求参数的结构体

type ParamSignUP struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Re_password string `json:"re_password"`
}
