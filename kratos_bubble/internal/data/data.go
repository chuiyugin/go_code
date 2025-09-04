package data

import (
	"kratos_bubble/internal/biz"
	"kratos_bubble/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewTodoRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: db,
	}, cleanup, nil
}

// 符合控制反转的写法
func NewDB(c *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&biz.Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
