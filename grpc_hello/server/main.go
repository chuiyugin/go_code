package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"server/pb"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
) // proto 文件定义中的option

// grpc server
type server struct {
	pb.UnimplementedGreeterServer                // 用于未实现的结构体，避免报错
	mu                            sync.Mutex     //count的并发锁
	count                         map[string]int // 存储每个name调用sayhello的次数（注意map要初始化）
}

// SayHello 是需要实现的方法
// 这个方法是对外提供的服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 理由defer 在发送完相应数据后，发送trailer
	defer func() {
		trailer := metadata.Pairs(
			"timestamp", strconv.Itoa(int(time.Now().Unix())),
		)
		grpc.SetTrailer(ctx, trailer)
		return
	}()
	// 在执行业务逻辑之前check metadata中是否包含token
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Printf("md:%#v ok:%#v\n", md, ok)
	if !ok { // 没有token，拒绝接受
		return nil, status.Error(codes.Unauthenticated, "无效请求")
	}
	vl := md.Get("token")
	if len(vl) < 1 || vl[0] != "api-test-yugin" {
		return nil, status.Error(codes.Unauthenticated, "无效请求")
	}
	// 限制同一个名字的调用次数
	s.mu.Lock()             //加锁（保证并发安全）
	s.count[in.GetName()]++ // 记录name的请求次数
	s.mu.Unlock()
	if s.count[in.GetName()] > 1 {
		// 返回请求次数限制的错误
		st := status.New(codes.ResourceExhausted, "request limit")
		// 添加错误的详情信息
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{
					{
						Subject:     fmt.Sprintf("name:%s", in.Name),
						Description: "每个name只能调用一次SayHello",
					},
				},
			},
		)
		if err != nil { // st.WithDetails 执行失败，返回原来的st
			return nil, st.Err()
		}
		return nil, ds.Err() // 返回带details的error
	}
	// 正常执行
	reply := "hello " + in.GetName()
	// 发送数据前发送header
	header := metadata.New(map[string]string{
		"location": "GuangZhou",
	})
	grpc.SendHeader(ctx, header)
	return &pb.HelloResponse{Reply: reply}, nil
}

// LotsOfReplies 返回使用多种语言打招呼
func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}

	for _, word := range words {
		data := &pb.HelloResponse{
			Reply: word + in.GetName(),
		}
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "你好："
	for {
		// 接收客户端发来的流式数据
		res, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			return stream.SendAndClose(&pb.HelloResponse{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += res.GetName()
	}
}

// BidiHello 双向流式打招呼
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := magic(in.GetName()) // 对收到的数据做些处理

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

func main() {
	// 启动服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen, :err:%v\n", err)
		return
	}
	// 加载证书信息
	creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		fmt.Printf("credentials.NewServerTLSFromFile failed, err:%v\n", err)
		return
	}
	s := grpc.NewServer(grpc.Creds(creds)) // 创建rpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		fmt.Printf("failed to listen, :err:%v\n", err)
		return
	}
}
