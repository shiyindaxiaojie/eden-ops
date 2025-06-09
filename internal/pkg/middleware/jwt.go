package middleware

import (
	"eden-ops/pkg/auth"
	"eden-ops/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// UserIDKey 用户ID上下文键
	UserIDKey = "user_id"
	// UsernameKey 用户名上下文键
	UsernameKey = "username"
)

// JWT 中间件
func JWT(jwtAuth *auth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未登录或非法访问")
			c.Abort()
			return
		}

		// 检查 token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证信息格式错误")
			c.Abort()
			return
		}

		claims, err := jwtAuth.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 设置用户信息
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)

		c.Next()
	}
}
