package logic

import (
	"context"
	"fmt"

	"api/internal/svc"
	"api/internal/types"
	"api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	// todo: add your logic here and delete this line
	fmt.Printf("req:%#v\n", req)

	// 把用户注册信息保存到数据库
	user := &model.Users{
		UserId:   111,
		Username: req.Username,
		Password: req.Password,
		Gender:   1,
	}
	if _, err := l.svcCtx.UserModel.Insert(context.Background(), user); err != nil {
		return nil, err
	}
	return &types.SignupResponse{Message: "success!"}, nil
}
