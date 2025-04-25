package controller

import "database/sql"

// Controller 控制器结构体，用于处理应用程序的业务逻辑。
type Controller struct {
	db *sql.DB
}

// NewController 创建一个新的 Controller 实例。
// 参数:
//   db (*sql.DB): 数据库连接对象，用于控制器与数据库进行交互。
// 返回值:
//   *Controller: 返回一个指向新创建的 Controller 实例的指针。
func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}
