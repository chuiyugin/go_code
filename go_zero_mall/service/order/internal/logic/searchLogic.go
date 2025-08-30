package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"order/internal/svc"
	"order/internal/types"
	"order/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// todo: add your logic here and delete this line
	// 1.根据请求参数中的订单号查询数据库找到订单记录
	order_id, err := strconv.ParseUint(req.OrderID, 10, 64)
	fmt.Printf("---> order_id:%v\n", order_id)
	if err != nil {
		logx.Errorw("order_id strconv.ParseUint", logx.Field("err", err))
		return nil, err
	}
	od, err := l.svcCtx.OrdersModel.FindOneByOrderId(l.ctx, order_id)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("订单查询不存在")
	}
	if err != nil {
		logx.Errorw("OrdersModel.FindOne failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	// 2.根据订单记录中的 user_id 去查询用户数据（通过RPC调用user服务）
	userResp, err := l.svcCtx.UserRPC.GetUser(l.ctx, &userclient.GetUserRequest{UserID: 1756458293})
	if err != nil {
		logx.Errorw("UserRPC.GetUser failed", logx.Field("err", err))
		return nil, err
	}
	// 3.拼接返回结果（这个接口的数据是由多个服务组成的）
	return &types.SearchResponse{
		OrderID:  strconv.FormatUint(od.OrderId, 10),
		Status:   int(od.Status),
		Username: userResp.Username,
	}, nil
}
