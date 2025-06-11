package handler

import (
	"eden-ops/internal/pkg/middleware"
	"eden-ops/internal/service"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/logger"
	"eden-ops/pkg/response"

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
		logger.Error("登录请求参数错误: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	logger.Info("用户登录请求: username=%s", req.Username)

	// 调用服务层登录方法
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		logger.Error("用户登录失败: username=%s, error=%v", req.Username, err)
		response.Unauthorized(c, err.Error())
		return
	}

	// 获取用户信息
	claims, err := h.jwtAuth.ParseToken(token)
	if err != nil {
		logger.Error("解析token失败: username=%s, error=%v", req.Username, err)
		response.Failed(c, err)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserInfo(claims.UserID)
	if err != nil {
		logger.Error("获取用户信息失败: username=%s, userID=%d, error=%v", req.Username, claims.UserID, err)
		response.Failed(c, err)
		return
	}

	// 清除敏感信息
	user.Password = ""

	logger.Info("用户登录成功: username=%s, userID=%d", req.Username, user.ID)
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
