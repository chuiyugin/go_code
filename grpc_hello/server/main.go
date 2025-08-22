package main

import (
	"context"
	"fmt"
	"net"
	"server/pb"

	"google.golang.org/grpc"
) // proto 文件定义中的option

// grpc server
type server struct {
	pb.UnimplementedGreeterServer // 用于未实现的结构体，避免报错
}

// SayHello 是需要实现的方法
// 这个方法是对外提供的服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "hello " + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	// 启动服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen, :err:%v\n", err)
		return
	}
	s := grpc.NewServer() // 创建rpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		fmt.Printf("failed to listen, :err:%v\n", err)
		return
	}
}
