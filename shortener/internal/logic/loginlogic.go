package logic

import (
	"errors"
	"time"

	"context"

	"shortener/internal/svc"
	"shortener/internal/types"

	"shortener/internal/pkg/auth" // 上面两个文件所在包

	"github.com/google/uuid"
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

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	u, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err == sqlx.ErrNotFound {
		return &types.LoginResponse{Message: "用户名不存在"}, nil
	}
	if err != nil {
		l.Errorf("FindOneByUsername error: %v", err)
		return nil, errors.New("内部错误")
	}
	if u.Password != passwordMd5([]byte(req.Password)) {
		return &types.LoginResponse{Message: "用户名或密码错误"}, nil
	}

	now := time.Now().Unix()

	// 1) Access Token
	at, err := auth.GetJwtTokenWithClaims(
		l.svcCtx.Config.Auth.AccessSecret,
		now,
		l.svcCtx.Config.Auth.AccessExpire,
		u.UserId,
		"access",
		"",
	)
	if err != nil {
		l.Errorf("getJwtToken(access) error: %v", err)
		return nil, errors.New("内部错误")
	}
	atExp := int(now + l.svcCtx.Config.Auth.AccessExpire)

	// 2) Refresh Token（带 jti）
	jti := uuid.NewString()
	rt, err := auth.GetJwtTokenWithClaims(
		l.svcCtx.Config.Auth.RefreshSecret,
		now,
		l.svcCtx.Config.Auth.RefreshExpire,
		u.UserId,
		"refresh",
		jti,
	)
	if err != nil {
		l.Errorf("getJwtToken(refresh) error: %v", err)
		return nil, errors.New("内部错误")
	}
	rtExp := int(now + l.svcCtx.Config.Auth.RefreshExpire)

	// 3) 将 RT 写入白名单
	if err := auth.RTAllow(l.svcCtx.Rds, u.UserId, jti, l.svcCtx.Config.Auth.RefreshExpire); err != nil {
		l.Errorf("RTAllow error: %v", err)
		return nil, errors.New("内部错误")
	}

	return &types.LoginResponse{
		Message:       "登录成功！",
		AccessToken:   at,
		AccessExpire:  atExp,
		RefreshToken:  rt,    // 新增字段
		RefreshExpire: rtExp, // 新增字段
		RefreshAfter:  int(now + l.svcCtx.Config.Auth.AccessExpire/2),
	}, nil
}
