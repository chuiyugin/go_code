package main

import (
	"bufio"
	"client/pb"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "yugin", "通过-name指定server中的name字段")

func main() {
	flag.Parse() // 解析命令行参数
	// 连接server
	// 加载证书
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		fmt.Printf("credentials.NewClientTLSFromFile failed, err:%v\n", err)
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:8972", grpc.WithTransportCredentials(creds))
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
	// 普通的rpc调用
	// 携带元数据
	md := metadata.Pairs(
		"token", "api-test-yugin",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 声明两个变量
	var header, trailer metadata.MD
	resp, err := c.SayHello(
		ctx,
		&pb.HelloRequest{Name: *name},
		grpc.Header(&header),
		grpc.Trailer(&trailer))
	if err != nil {
		// 收到带detail的error
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("QuotaFailure:%v\n", info)
			default:
				fmt.Printf("unexpected type\n", info)
			}
		}
		log.Fatal("c.SayHello failed, err:%v\n", err)
		return
	}
	// 拿到数据之前可以获得header
	log.Printf("header:%v\n", header)
	// 收到执行结果resp
	log.Printf("resp:%v\n", resp.GetReply())
	// 拿到数据之前可以获得trailer
	log.Printf("trailer:%v\n", trailer)

	// 流式rpc服务端调用
	// callLotsOfReplies(c)
	// runLotsOfGreeting(c)
	// runBidiHello(c)
}

func callLotsOfReplies(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Println(err)
		return
	}
	// 依次从流式响应中读取返回的响应数据
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("stream.Recv failed, err:%v\n", err)
			return
		}
		log.Printf("recv: %v\n", res.GetReply())
	}
}

func runLotsOfGreeting(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端流式RPC
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed, err: %v", err)
	}
	names := []string{"yugin", "scordingbig", "不傲"}
	for _, name := range names {
		// 发送流式数据
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
		}
	}
	// 流式发送结束后关闭流
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed: %v", err)
	}
	log.Printf("got reply: %v", res.GetReply())
}

func runBidiHello(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// 双向流模式
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed, err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI：%s\n", in.GetReply())
		}
	}()
	// 从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') // 读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
