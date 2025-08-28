package svc

import (
	"api/internal/config"
	"api/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UsersModel // 加入User表增删改查操作
}

func NewServiceContext(c config.Config) *ServiceContext {
	// UserModel 接口类型
	// 调用构造函数得到 *model.defaultUserModel
	// NewUserModel(conn sqlx.SqlConn)
	// 需要 sqlx.SqlConn 的数据库链接
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUsersModel(sqlConn),
	}
}
