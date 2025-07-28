package main

// main.go 117.72.167.29:9090/

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Todo Model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func main() {
	// 连接mysql数据库
	db, err := gorm.Open("mysql", "yugin:1234@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	// gin框架模板文件引用的静态文件路径
	r.Static("/static", "static")
	// gin框架模板文件路径
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// v1
	v1Group := r.Group("V1"){
		// 待办事项
		// 添加
		v1Group.POST("/todo",func(c *gin.Context){

		})
		// 查看所有待办事项
		v1Group.GET("/todo",func(c *gin.Context){

		})
		// 查看一个待办事项
		v1Group.GET("/todo/:id",func(c *gin.Context){

		})
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id",func(c *gin.Context){

		})
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id",func(c *gin.Context){

		})
	}

	r.Run(":9090")
}
