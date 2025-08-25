package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"server/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	proto.UnimplementedCalcServiceServer
}

func (s server) Add(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	sum := in.GetX() + in.GetY()
	return &proto.AddResponse{Result: int64(sum)}, nil
}

func main() {
	l, err := net.Listen("tcp", ":8973")
	if err != nil {
		log.Fatal("net.Listen failed, err:%v", err)
		return
	}
	s := grpc.NewServer()
	// 注册服务
	proto.RegisterCalcServiceServer(s, &server{})
	// 启动该服务
	// err = s.Serve(l)
	// if err != nil {
	// 	log.Fatal("s.Serve failed, err:%v", err)
	// 	return
	// }

	// 开启goroutine启动grpc server
	go func() {
		log.Fatalln(s.Serve(l))
	}()

	// 以下是grpc gateway 新加入的内容
	// 创建一个连接到我们刚刚启动的 gRPC 服务器的客户端连接
	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.NewClient(
		"127.0.0.1:8973",
		grpc.WithBlock(), // 阻塞直到连接成功
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// 注册 Calc 路由
	err = proto.RegisterCalcServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// 定义HTTP server配置
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	// 8090端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe()) // 启动HTTP服务

}
