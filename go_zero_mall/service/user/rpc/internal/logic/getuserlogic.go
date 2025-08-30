package logic

import (
	"context"
	"errors"

	"user/rpc/internal/svc"
	"user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 方法
func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserReply, error) {
	// todo: add your logic here and delete this line
	// 根据 UserID 查询数据库返回用户信息
	one, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserID)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errors.New("无效的UserID")
	}
	if err != nil {
		logx.Errorw(
			"use.rpc.getUser FindOneByUsername failed",
			logx.Field("err", err),
		)
		return nil, errors.New("查询失败")
	}
	return &user.GetUserReply{
		UserID:   one.UserId,
		Username: one.Username,
		Gender:   one.Gender,
	}, nil
}
