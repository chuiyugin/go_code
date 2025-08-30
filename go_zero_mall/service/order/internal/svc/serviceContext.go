package svc

import (
	"order/internal/config"
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
		UserRPC:     userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		OrdersModel: model.NewOrdersModel(sqlConn, c.CacheRedis),
	}
}
