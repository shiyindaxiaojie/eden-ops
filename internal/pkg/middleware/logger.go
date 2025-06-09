package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回统一格式的日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()
		if clientIP == "::1" {
			clientIP = "127.0.0.1"
		}

		// 构建日志格式：时间 | 状态码 | 耗时ms | 客户端IP | 请求方法 "请求路径"
		timestamp := time.Now().Format("2006/01/02 15:04:05.000")
		logMsg := fmt.Sprintf("%s | %3d | %dms | %15s | %-7s %q",
			timestamp,
			statusCode,
			latencyTime.Milliseconds(),
			clientIP,
			reqMethod,
			reqURI,
		)

		// 输出到控制台
		fmt.Println(logMsg)
	}
}
