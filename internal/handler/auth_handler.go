package handler

import (
	"eden-ops/internal/pkg/middleware"
	"eden-ops/internal/service"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService service.UserService
	jwtAuth     *auth.JWTAuth
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userService service.UserService, jwtAuth *auth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtAuth:     jwtAuth,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("请求参数绑定失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	log.Printf("接收到登录请求: username=%s", req.Username)

	// 调用服务层登录方法
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("登录失败: %v", err)
		response.Unauthorized(c, err.Error())
		return
	}

	log.Printf("登录成功，解析token")

	// 获取用户信息
	claims, err := h.jwtAuth.ParseToken(token)
	if err != nil {
		log.Printf("解析token失败: %v", err)
		response.Failed(c, err)
		return
	}

	log.Printf("获取用户信息: userID=%d", claims.UserID)

	// 获取用户信息
	user, err := h.userService.GetUserInfo(claims.UserID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		response.Failed(c, err)
		return
	}

	log.Printf("获取用户信息成功，返回响应")

	// 清除敏感信息
	user.Password = ""

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 前端应该清除token，后端不需要做特殊处理
	response.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 调用服务层获取用户信息
	user, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, user)
}
