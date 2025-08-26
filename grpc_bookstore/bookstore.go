package main

import (
	"context"
	"fmt"
	"grpc_bookstore/pb"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// bookstore grpc 服务

const (
	defaultCursor       = "0" // 默认游标
	defaultPageSize int = 2   // 默认每页显示书本数量
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

func (s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 参数 check
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invaild shelf id")
	}
	var (
		cursor       = defaultCursor
		pageSize int = defaultPageSize
	)
	// if in.GetPageToken() == "" {
	// 	// 没有分页 token 默认第一页
	// } else {
	if len(in.GetPageToken()) > 0 {
		// 有分页的话就解析分页数据
		pageInfo := Token(in.GetPageToken()).Decode()
		// 再判断解析结果是否有效
		if pageInfo.InVaild() {
			return nil, status.Error(codes.InvalidArgument, "invaild page_token")
		}
		cursor = pageInfo.NextID
		pageSize = int(pageInfo.PageSize)
	}
	// 查询数据库
	bookList, err := s.bs.GetBookListByShelfId(ctx, in.GetShelf(), cursor, pageSize+1)
	if err != nil {
		fmt.Printf("GetBookListByShelfId failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 如果查询出来的结果比 pageSize 大，那么就说明有下一页
	var (
		haveNextPage  bool
		nextPageToken string
		realSize      int = len(bookList)
	)
	// 当查询数据库的结果数大于 pageSize
	if len(bookList) > pageSize {
		haveNextPage = true // 有下一页
		realSize = pageSize // 下面格式化数据没必要把所有查询结果都返回，只需要返回 pageSize 个数据即可
	}
	// 封装返回的数据
	// 将 book --> []*pb.Books
	res := make([]*pb.Book, 0, realSize)
	//for _, b := range bookList {
	for i := 0; i < realSize; i++ {
		res = append(res, &pb.Book{
			Id:     bookList[i].ID,
			Author: bookList[i].Author,
			Title:  bookList[i].Title,
		})
	}
	// 如果有下一页，就要生成下一页的 page_toke
	if haveNextPage {
		nextPageInfo := Page{
			NextID:        strconv.FormatInt(res[realSize-1].Id, 10), // 最后一个返回结果的 id (strconv.FormatInt用于将整数转换为字符串)
			NextTimeAtUTC: time.Now().Unix(),
			PageSize:      int64(pageSize),
		}
		nextPageToken = string(nextPageInfo.Encode())
	}
	return &pb.ListBooksResponse{Books: res, NextPageToken: nextPageToken}, nil
}
