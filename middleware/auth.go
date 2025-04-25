package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// TokenExpiration 定义令牌有效期为24小时（秒）
	TokenExpiration = 86400
	// QueryTimeout 定义数据库查询超时时间
	QueryTimeout = 5 * time.Second
)

func Auth(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		authorization := c.GetHeader("Authorization")
		fmt.Println("authorization:", authorization)
		if authorization == "" {
			sendUnauthorized(c)
			return
		}
		// 拆分 token 并检查格式
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			sendUnauthorized(c)
			return
		}
		// 5秒超时
		ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
		// 释放资源
		defer cancel()
		// 从缓存中读出userid
		if tokenKey, exists := c.Get(parts[1]); exists {
			// 断言并检查是否为非空切片
			if key, ok := tokenKey.([]int64); ok && len(key) == 2 {
				// 检查令牌是否过期
				if key[1] >= time.Now().Unix() {
					c.Set("userid", key[0])
					c.Next()
					return
				}
			}
		}
		// 用户id
		var userid int64
		err := db.QueryRowContext(
			ctx,
			"SELECT `id` FROM `users` WHERE `token` = ?",
			parts[1],
		).Scan(&userid)
		// 错误处理
		if err != nil {
			sendUnauthorized(c)
			return
		}
		// 写入令牌，值是userid 和 令牌有效时间
		c.Set(parts[1], []int64{userid, time.Now().Unix() + TokenExpiration})
		c.Set("userid", userid)
		c.Next()
	}
}

// sendUnauthorized 向客户端发送401未授权错误响应。
// 这个函数在令牌验证失败时被调用，向客户端返回固定格式的JSON数据，
// 包含状态码和过期提示信息，告知用户需要重新登录。
func sendUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": "令牌已过期，请重新登录...",
	})
}
