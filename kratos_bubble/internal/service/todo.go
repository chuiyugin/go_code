package service

import (
	"context"
	"errors"

	pb "kratos_bubble/api/bubble/v1"
	v1 "kratos_bubble/api/bubble/v1"
	"kratos_bubble/internal/biz"
)

type TodoService struct {
	pb.UnimplementedTodoServer

	uc *biz.TodoUsecase
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{
		uc: uc,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	// 请求入口
	// 1 参数校验
	if len(req.Title) == 0 {
		return &pb.CreateTodoReply{}, errors.New("无效的title")
	}
	// 调用业务逻辑
	data, err := s.uc.CreateTodo(ctx, &biz.Todo{Title: req.Title})
	if err != nil {
		return nil, err
	}
	// 返回响应
	return &pb.CreateTodoReply{
		Id:     data.ID,
		Title:  data.Title,
		Status: data.Status,
	}, nil
}
func (s *TodoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoReply, error) {
	return &pb.UpdateTodoReply{}, nil
}
func (s *TodoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoReply, error) {
	return &pb.DeleteTodoReply{}, nil
}
func (s *TodoService) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	// 1、参数处理
	if req.Id <= 0 {
		return &pb.GetTodoReply{}, errors.New("无效的id")
	}
	// 2、调用biz层业务逻辑
	ret, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		// return nil, err
		return nil, v1.ErrorTodoNotFound("id:%v todo is not found", req.Id)
	}
	// 3、按格式返回响应
	return &pb.GetTodoReply{
		Id:     ret.ID,
		Title:  ret.Title,
		Status: ret.Status,
	}, nil
}
func (s *TodoService) ListTodo(ctx context.Context, req *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	return &pb.ListTodoReply{}, nil
}
