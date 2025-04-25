package main

import (
	"log"
	"os"

	"github.com/xa1st/blog-func/controller"
	"github.com/xa1st/blog-func/database"
	"github.com/xa1st/blog-func/route"
)

func main() {
	// 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("[ERROR]数据库初始化失败: %v", err)
	}
	// 关闭数据库连接
	defer db.Close()

	// 初始化控制器
	ctrl := controller.NewController(db)

	// 初始化路由
	r := route.SetupRouter(ctrl, db)

	// 从环境变量中获取端口号，如果没有设置则使用默认值3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 启动服务器并监听在8080端口
	r.Run(":" + port)
}
