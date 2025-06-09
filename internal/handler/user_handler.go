package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
	jwtAuth     *auth.JWTAuth
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService, jwtAuth *auth.JWTAuth) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtAuth:     jwtAuth,
	}
}

// UserRequest 用户请求
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// 获取用户信息
	claims, err := h.jwtAuth.ParseToken(token)
	if err != nil {
		response.Failed(c, err)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserInfo(claims.UserID)
	if err != nil {
		response.Failed(c, err)
		return
	}

	// 清除敏感信息
	user.Password = ""

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   status,
	}

	err := h.userService.Create(user)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, user)
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	user := &model.User{
		ID:       uint(id),
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   status,
	}

	err = h.userService.Update(user)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// Get 获取用户详情
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	users, total, err := h.userService.List(page, pageSize)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, users, total)
}

// GetRoles 获取用户角色
func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	roles, err := h.userService.GetUserRoles(uint(id))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, roles)
}

// AssignRoles 分配用户角色
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var roleIDs []uint
	if err := c.ShouldBindJSON(&roleIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err = h.userService.AssignRoles(uint(id), roleIDs)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		response.Failed(c, err)
		return
	}

	// 清除敏感信息
	user.Password = ""

	response.Success(c, user)
}
