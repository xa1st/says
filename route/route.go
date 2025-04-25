package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/xa1st/blog-func/controller"
	"github.com/xa1st/blog-func/middleware"
)

// SetupRouter 配置并返回一个gin路由器。
// 该函数接受一个实现了Controller接口的控制器实例作为参数，
// 以便在路由中绑定处理函数。
func SetupRouter(ctrl *controller.Controller, db *sql.DB) *gin.Engine {
	// 创建一个默认的gin路由器
	r := gin.Default()
	// 创建一个分组，用于处理/say相关的路由
	say := r.Group("/api/v1")
	// 使用CORS中间件
	say.Use(middleware.CORS())
	// say.Use(middleware.CORS(), middleware.Auth(db))
	{
		// GET请求，用于获取所有memo
		say.GET("/memo", ctrl.List)
		// 用户登陆+状态检测
		say.POST("/auth/status", middleware.Auth(db), ctrl.Auth)
		// POST请求，用于添加memo
		say.POST("/memos", middleware.Auth(db), ctrl.Add)
		// GET请求，用于获取所有memo
		say.GET("/index", ctrl.Index)
	}

	// 返回gin路由器
	return r
}
