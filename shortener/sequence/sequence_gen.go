package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 使用 REPLACE INTO 语句往 sequence 表插入一条数据，并且取出主键id作为号码

const sqlReplaceIntoStub = `REPLACE INTO sequence (stub) VALUES ('a');`

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取下一个号
func (m *MySQL) Next() (seq uint64, err error) {
	// 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw("conn.Prepaer failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	defer stmt.Close()

	// 执行预编译
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}

	// 获取刚插入的ID
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	return uint64(lid), nil
}
