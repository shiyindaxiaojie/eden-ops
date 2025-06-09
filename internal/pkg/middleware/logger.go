package middleware

import (
	"eden-ops/pkg/logger"
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

		// 使用全局日志格式记录
		logger.LogHTTPRequest(statusCode, reqMethod, reqURI, latencyTime, clientIP)
	}
}
