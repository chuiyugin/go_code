package main

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	// db, err := gorm.Open(sqlite.Open("/home/yugin/Code/go_code/grpc_bookstore/test.db"), &gorm.Config{})
	dsn := "root:1234@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("1failed to connect database")
	}
	// 迁移
	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

const (
	defaultShelfSize = 100
)

// 定义Model

// Shelf 书架
type Shelf struct {
	ID       int64 `gorm:"primaryKey"`
	Theme    string
	Size     int64
	CreateAt time.Time
	UpdateAt time.Time
}

// Book 图书
type Book struct {
	ID       int64 `gorm:"primaryKey"`
	Author   string
	Title    string
	ShelfID  int64
	CreateAt time.Time
	UpdateAt time.Time
}

// 数据库操作
type bookstore struct {
	db *gorm.DB
}

// CreateShelf 创建书架
func (b *bookstore) CreateShelf(ctx context.Context, data Shelf) (*Shelf, error) {
	if len(data.Theme) <= 0 {
		return nil, errors.New("invaild theme")
	}
	size := data.Size
	if size <= 0 {
		size = defaultShelfSize
	}
	v := Shelf{
		Theme:    data.Theme,
		Size:     size,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

// GetShelf 获取书架
func (b *bookstore) GetShelf(ctx context.Context, id int64) (*Shelf, error) {
	v := Shelf{}
	err := b.db.WithContext(ctx).First(&v, id).Error
	return &v, err
}

// ListShelfs 书架列表
func (b *bookstore) ListShelfs(ctx context.Context) ([]*Shelf, error) {
	var vl []*Shelf
	err := b.db.WithContext(ctx).Find(&vl).Error
	return vl, err
}

// DeleteShelf 书架列表
func (b *bookstore) DeleteShelf(ctx context.Context, id int64) error {
	err := b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
	return err
}
