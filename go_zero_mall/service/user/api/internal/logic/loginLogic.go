package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func passwordMd5(password []byte) string {

	h := md5.New()
	h.Write(password) // 密码计算md5
	h.Write(secret)
	PasswordStr := hex.EncodeToString(h.Sum(nil))
	return PasswordStr
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	// 实现登录功能
	// 1. 处理用户发来的请求，拿到用户名和密码
	// 2. 判断输入的用户名和密码跟数据库中的是否一致
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err == sqlx.ErrNotFound {
		return &types.LoginResponse{
			Message: "用户名不存在",
		}, nil
	}
	if err != nil {
		logx.Errorw("UserModel.FindOneByUsername failed", logx.Field("err", err))
		return &types.LoginResponse{
			Message: "内部错误",
		}, errors.New("内部错误")
	}
	if u.Password != passwordMd5([]byte(req.Password)) {
		return &types.LoginResponse{
			Message: "用户名或密码错误",
		}, nil
	}
	// 3. 如果结果一致就登录成功，否则就登录失败
	return &types.LoginResponse{
		Message: "登录成功！",
	}, nil
}
