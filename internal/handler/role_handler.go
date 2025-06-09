package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// RoleRequest 角色请求
type RoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Status string `json:"status"`
	Remark string `json:"remark"`
}

// Create 创建角色
func (h *RoleHandler) Create(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	role := &model.Role{
		Name:   req.Name,
		Code:   req.Code,
		Status: status,
		Remark: req.Remark,
	}

	err := h.roleService.Create(role)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, role)
}

// Update 更新角色
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取原角色信息
	role, err := h.roleService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	// 更新角色信息
	role.Name = req.Name
	role.Code = req.Code
	role.Status = status
	role.Remark = req.Remark

	err = h.roleService.Update(role)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, role)
}

// Delete 删除角色
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	err = h.roleService.Delete(uint(id))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// Get 获取角色详情
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

// List 获取角色列表
func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	roles, total, err := h.roleService.List(page, pageSize)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  roles,
		"total": total,
	})
}

// AssignMenus 分配菜单
func (h *RoleHandler) AssignMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	// 检查角色是否存在
	_, err = h.roleService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	var req struct {
		MenuIDs []uint `json:"menuIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 分配菜单
	err = h.roleService.AssignMenus(uint(id), req.MenuIDs)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// AssignUserRoles 分配用户角色
func (h *RoleHandler) AssignUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"roleIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 分配角色
	err = h.roleService.AssignUserRoles(uint(userID), req.RoleIDs)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}
