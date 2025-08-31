package svc

import (
	"order/internal/config"
	"order/internal/interceptor"
	"order/model"
	"order/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	UserRPC     userclient.User
	OrdersModel model.OrdersModel // 加入User表增删改查操作
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		UserRPC:     userclient.NewUser(zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptor.UserInterceptor))), // 初始化user服务的RPC客户端
		OrdersModel: model.NewOrdersModel(sqlConn, c.CacheRedis),
	}
}
