package middleware

import (
	"bytes"
	"fmt"
	"net/http"
)

// 定义全局中间件

type bodyCopy struct {
	http.ResponseWriter               // 结构体嵌入的接口类型
	body                *bytes.Buffer // 用于记录请求的响应体信息
}

func newbodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

// Write 构造一个针对bodyCopy的Write函数
func (bc *bodyCopy) Write(buf []byte) (int, error) {
	// 1.在自己的buffer中写入响应体
	bc.body.Write(buf)
	// 2.返回正常响应体的
	return bc.ResponseWriter.Write(buf)
}

// 功能：记录所有请求的响应体信息

// CopyRes 复制请求的响应体
func CopyRes(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 处理请求前
		b := newbodyCopy(w)
		next(b, r) // 实际执行的后续handler处理函数
		// 处理请求后
		fmt.Printf("req:%v, reqs:%s\n", r.URL, b.body.String())

	}
}
