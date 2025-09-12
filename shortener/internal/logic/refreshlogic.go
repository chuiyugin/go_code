package logic

import (
	"context"
	"errors"
	"time"

	"shortener/internal/svc"
	"shortener/internal/types"

	"shortener/internal/pkg/auth"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshLogic) Refresh(req *types.RefreshRequest) (resp *types.RefreshResponse, err error) {
	// 1) 解析并校验 RefreshToken（签名/exp）
	claims, err := auth.ParseJwt(req.RefreshToken, l.svcCtx.Config.Auth.RefreshSecret)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}
	// typ 检查
	if t := auth.ClaimString(claims, "typ"); t != "refresh" {
		return nil, errors.New("不是刷新令牌")
	}
	userID := auth.ClaimInt64(claims, "UserId")
	if userID == 0 {
		return nil, errors.New("刷新令牌缺少用户信息")
	}
	oldJTI := auth.ClaimString(claims, "jti")
	if oldJTI == "" {
		return nil, errors.New("刷新令牌缺少 jti")
	}

	// 2) 校验白名单
	allowed, err := auth.RTIsAllowed(l.svcCtx.Rds, userID, oldJTI)
	if err != nil {
		l.Errorf("RTIsAllowed error: %v", err)
		return nil, errors.New("内部错误")
	}
	if !allowed {
		return nil, errors.New("刷新令牌已失效")
	}

	now := time.Now().Unix()

	// 3) 生成新的 AT
	at, err := auth.GetJwtTokenWithClaims(
		l.svcCtx.Config.Auth.AccessSecret,
		now,
		l.svcCtx.Config.Auth.AccessExpire,
		userID,
		"access",
		"",
	)
	if err != nil {
		l.Errorf("getJwtToken(access) error: %v", err)
		return nil, errors.New("内部错误")
	}
	atExp := int(now + l.svcCtx.Config.Auth.AccessExpire)

	// 4) 轮换新的 RT（先 Allow 新，再撤销旧）
	newJTI := uuid.NewString()
	rt, err := auth.GetJwtTokenWithClaims(
		l.svcCtx.Config.Auth.RefreshSecret,
		now,
		l.svcCtx.Config.Auth.RefreshExpire,
		userID,
		"refresh",
		newJTI,
	)
	if err != nil {
		l.Errorf("getJwtToken(refresh) error: %v", err)
		return nil, errors.New("内部错误")
	}
	rtExp := int(now + l.svcCtx.Config.Auth.RefreshExpire)

	if err := auth.RTAllow(l.svcCtx.Rds, userID, newJTI, l.svcCtx.Config.Auth.RefreshExpire); err != nil {
		l.Errorf("RTAllow error: %v", err)
		return nil, errors.New("内部错误")
	}
	// 撤销旧 RT（幂等）
	if err := auth.RTRevoke(l.svcCtx.Rds, userID, oldJTI); err != nil {
		l.Errorf("RTRevoke warning: %v", err)
	}

	return &types.RefreshResponse{
		AccessToken:   at,
		AccessExpire:  atExp,
		RefreshToken:  rt,
		RefreshExpire: rtExp,
	}, nil
}
