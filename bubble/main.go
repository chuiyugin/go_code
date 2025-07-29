package main

// main.go 117.72.167.29:9090/

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Todo Model
type Todo struct {
	ID     int    `json:"id" gorm:"primary_key;autoIncrement`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	// 连接mysql数据库
	DB, err = gorm.Open("mysql", "yugin:1234@(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	// 模型绑定
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	// gin框架模板文件引用的静态文件路径
	r.Static("/static", "static")
	// gin框架模板文件路径
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// v1
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 前端界面填写待办事项 点击提交 会发送请求到这里
			// 1.从请求中把数据拿出来
			var todo Todo
			if err := c.ShouldBindJSON(&todo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// 2.存入数据库并返回响应
			err = DB.Create(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todo)
			}

		})
		// 查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			// 查询todos表里的所有数据
			var todolist []Todo
			err = DB.Find(&todolist).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todolist)
			}
		})
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			err = DB.Where("id=?", id).First(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			if err := c.ShouldBindJSON(&todo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
				return
			}
			err = DB.Where("id=?", id).Delete(Todo{}).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": "处理成功！"})
			}
		})
	}

	r.Run(":9090")
}
