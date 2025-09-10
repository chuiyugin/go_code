package server

import (
	_ "embed"
	"net/http"

	v1 "review-service/api/review/v1" // ← 按你的模块名改
	"review-service/internal/conf"    // ← 按你的模块名改
	"review-service/internal/service" // ← 按你的模块名改

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

/*
目录结构（关键是 openapi.yaml 必须在当前包的子目录）：
internal/
  server/
    http.go
    docs/
      openapi.yaml
*/
//go:embed docs/openapi.yaml
var openapiSpec []byte

// 最小 Swagger UI（CDN）
const swaggerHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>API Docs</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css" />
  <style>html,body,#swagger-ui{height:100%;margin:0}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',            // 同源读取 spec
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, reviewer *service.ReviewService, logger log.Logger) *kratoshttp.Server {
	opts := []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := kratoshttp.NewServer(opts...)

	// 提供 openapi.yaml（浏览器可直接预览；Swagger UI 亦可解析）
	srv.Handle("/openapi.yaml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // 或 "application/yaml"
		_, _ = w.Write(openapiSpec)
	}))

	// Swagger UI：/q 与 /q/ 都能打开
	srv.Handle("/doc", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerHTML))
	}))
	srv.Handle("/doc/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerHTML))
	}))

	// 健康检查（可选）
	srv.HandleFunc("/q/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 业务路由
	v1.RegisterReviewHTTPServer(srv, reviewer)
	return srv
}
