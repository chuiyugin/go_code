package controllers

type ResCode int64

const (
	CodeSucess          ResCode = 1000 + iota // 1000, 类型是 ResCode
	CodeInvalidParam                          // 1001, 类型是 ResCode
	CodeUserExist                             // 1002, 类型是 ResCode
	CodeUserNotExist                          // 1003, 类型是 ResCode
	CodeInvalidPassword                       // 1004, 类型是 ResCode
	CodeServerBusy                            // 1005, 类型是 ResCode

	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSucess:          "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
