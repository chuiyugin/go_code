package server

import (
	"context"
	"fmt"
	v1 "kratos_bubble/api/bubble/v1"
	"kratos_bubble/internal/conf"
	"kratos_bubble/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// 自定义中间件
func Middleware1() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 逻辑执行之前
			fmt.Println("自定义中间件执行handler执行之前")
			defer func() {
				fmt.Println("自定义中间件执行handler执行之后")
			}()
			return handler(ctx, req) // 执行目标handler
		}
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, todo *service.TodoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			Middleware1(), // 注册上自定义中间件
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 替换默认的HTTP响应编码器
	opts = append(opts, http.ResponseEncoder(responseEncoder))

	// 替换默认的错误响应编码器
	opts = append(opts, http.ErrorEncoder(errorEncoder))

	srv := http.NewServer(opts...)
	v1.RegisterTodoHTTPServer(srv, todo)
	return srv
}
