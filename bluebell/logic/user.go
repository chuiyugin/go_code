package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑代码

func SignUp(p *models.ParamSignUP) {
	// 1.判断用户存不存在
	mysql.QueryUserByUsername()
	// 2.生成UID
	snowflake.GenID()

	// 3.保存进数据库
	mysql.InsertUser()
}
