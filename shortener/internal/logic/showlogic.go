package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Err404 = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// todo: add your logic here and delete this line
	// 查看短链接，输入 yugin.cn/f --> 重定向到真实链接
	// req.ShortUrl = f
	// 1.0 布隆过滤器（不存在的短链接直接返回404，不需要后续处理）
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("show.Filter.Exists() failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, Err404
	}
	// 不存在
	if !exist {
		return nil, Err404
	}
	fmt.Println("开始查询缓存和DB...")
	// 1 根据短链接查询原始的长链接（在查询数据前增加了缓存层）
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			logx.Errorw("show.ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, Err404
		}
		return nil, Err404
	}
	// 2 返回查询到的长链接响应，在调用handler层调用重定向响应
	return &types.ShowResponse{
		LongUrl: u.Lurl.String,
	}, nil
}
