package errorx

// 自定义的错误

const (
	defaultErrCode = 1001
	RPCErrCode     = 1002
)

// CodeError 自定义错误
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErrorResponse 自定义的错误响应
type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError 返回自定义错误
func NewCodeError(code int, msg string) error {
	return CodeError{
		Code: code,
		Msg:  msg,
	}
}

// NewDefaultCodeError 返回默认自定义错误
func NewDefaultCodeError(msg string) error {
	return CodeError{
		Code: defaultErrCode,
		Msg:  msg,
	}
}

// Error CodeError 实现 error 接口
func (e CodeError) Error() string {
	return e.Msg
}

// Data 返回自定义类型的错误响应
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
