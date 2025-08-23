package main

import (
	"context"
	"log"
	"net"
	"server/proto"

	"google.golang.org/grpc"
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
	err = s.Serve(l)
	if err != nil {
		log.Fatal("s.Serve failed, err:%v", err)
		return
	}
}
