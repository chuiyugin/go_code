package logic

import (
	"context"
	"errors"
	"fmt"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	// todo: add your logic here and delete this line
	// 获取并打印JWT鉴权的数据
	fmt.Printf("JWT UserId:%v\n", l.ctx.Value("UserId"))
	// 1. 拿到请求参数
	// 2. 根据用户id查数据库
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserID)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("细节查询不存在")
	}
	if err != nil {
		logx.Errorw("UserModel.FindOneByUserId failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	// 3. 格式化数据（数据库里存的数据和前端要求的字段不太一致）
	// 4. 返回响应
	return &types.DetailResponse{
		Username: u.Username,
		Gender:   int(u.Gender),
	}, nil
}
