package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const VERSION = "1.0.0"

func (ctrl *Controller) Index(c *gin.Context) {
	// 输出文件头
	c.Header("Content-Type", "text/plain; charset=utf-8")
	// 直接输出text就可以了
	c.String(http.StatusOK, "欢迎使用基于lade.io的说说系统，当前版本号：v"+VERSION)
}
