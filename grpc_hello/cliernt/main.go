package main

import (
	"client/pb"
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "yugin", "通过-name指定server中的name字段")

func main() {
	flag.Parse() // 解析命令行参数
	// 连接server
	conn, err := grpc.NewClient("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc.NewClient failed, err:%v", err)
		return
	}
	defer conn.Close()
	// 创建客户端
	c := pb.NewGreeterClient(conn) // 使用生成的go代码
	// 调用rpc方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatal("c.SayHello failed, err:%v", err)
		return
	}
	// 收到执行结果resp
	log.Printf("resp:%v", resp.GetReply())
}
