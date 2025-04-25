// middleware/cors.go
package middleware

import "github.com/gin-gonic/gin"

// CORS 是一个 Gin 框架的中间件函数，用于处理跨域资源共享（CORS）请求。
// 它通过设置 HTTP 响应头来允许跨域访问，并支持预检请求（OPTIONS 方法）。
// 该函数无需参数，返回值为 gin.HandlerFunc 类型，可以直接用于 Gin 路由中间件链。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的跨域来源为任意域名（"*"），表示所有域名都可以访问当前服务。
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许跨域请求中携带身份验证信息（如 Cookie）。
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 指定客户端可以使用的请求头字段。
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 指定客户端可以使用的 HTTP 请求方法。
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 如果请求方法是 OPTIONS，则认为这是浏览器发起的预检请求（Preflight Request）。
		// 直接返回 204 状态码，表示预检请求成功，但不返回任何内容。
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// 继续执行后续的中间件或路由处理函数。
		c.Next()
	}
}
