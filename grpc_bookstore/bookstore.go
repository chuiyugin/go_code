package main

import (
	"context"
	"grpc_bookstore/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedBookstoreServer
	bs *bookstore
}

// ListShelves 列出所有书架的RPC方法
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	// 调用 orm 操作那些方法
	sl, err := s.bs.ListShelfs(ctx)
	if err == gorm.ErrEmptySlice { // 没有数据
		return &pb.ListShelvesResponse{}, nil
	}
	if err != nil { // 查询数据库失败
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回数据 (弄成 []*pb.Shelf 的格式)
	nsl := make([]*pb.Shelf, 0, len(sl))
	for _, si := range sl {
		nsl = append(nsl, &pb.Shelf{
			Id:    si.ID,
			Theme: si.Theme,
			Size:  si.Size,
		})
	}
	return &pb.ListShelvesResponse{Shelves: nsl}, nil
}

// CreateShelf 创建书架
func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	// 参数检查
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invaild theme")
	}
	// 准备数据
	data := Shelf{
		Theme: in.GetShelf().GetTheme(),
		Size:  in.GetShelf().GetSize(),
	}
	// 调用orm的创建函数
	ns, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Shelf{
		Id:    ns.ID,
		Theme: ns.Theme,
		Size:  ns.Size,
	}, nil
}

// GetShelf 根据id获取书架
func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	// 参数检查
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invaild shelf id")
	}
	// 查询数据库
	si, err := s.bs.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	return &pb.Shelf{
		Id:    si.ID,
		Theme: si.Theme,
		Size:  si.Size,
	}, nil
}

// DeleteShelf 根据id删除书架
func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	// 参数检查
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invaild shelf id")
	}
	err := s.bs.DeleteShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil
}
