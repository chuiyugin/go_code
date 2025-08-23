package main

import (
	"client/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接rpc server
	conn, err := grpc.NewClient("127.0.0.1:8973", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc.NewClient failed, err:%v", err)
		return
	}
	defer conn.Close()
	// 创建 rpc client 端
	client := proto.NewCalcServiceClient(conn)
	// 发起rpc调用
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Add(ctx, &proto.AddRequest{X: 10, Y: 100})
	if err != nil {
		log.Fatal("client.Add failed, err:%v", err)
		return
	}
	// 打印结果
	log.Printf("result:%v\n", resp.GetResult())
}
