package logic

import (
	"context"

	"go_zero_demo_2/internal/svc"
	"go_zero_demo_2/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Go_zero_demo_2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGo_zero_demo_2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Go_zero_demo_2Logic {
	return &Go_zero_demo_2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Go_zero_demo_2Logic) Go_zero_demo_2(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
