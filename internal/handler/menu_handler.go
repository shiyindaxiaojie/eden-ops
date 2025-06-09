package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	menuService service.MenuService
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// MenuRequest 菜单请求
type MenuRequest struct {
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	ParentID   uint   `json:"parentId"`
	Status     string `json:"status"`
	Permission string `json:"permission"`
	Type       string `json:"type" binding:"required"` // M:目录 C:菜单 F:按钮
}

// Create 创建菜单
func (h *MenuHandler) Create(c *gin.Context) {
	var req MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, err)
		return
	}

	// 转换菜单类型
	menuType := 1 // 默认为菜单
	switch req.Type {
	case "M": // 目录
		menuType = 0
	case "C": // 菜单
		menuType = 1
	case "F": // 按钮
		menuType = 2
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	menu := &model.Menu{
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		ParentID:  req.ParentID,
		Status:    status,
		Type:      menuType,
	}

	err := h.menuService.Create(menu)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, menu)
}

// Update 更新菜单
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取原菜单信息
	menu, err := h.menuService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "菜单不存在")
		return
	}

	// 转换菜单类型
	menuType := 1 // 默认为菜单
	switch req.Type {
	case "M": // 目录
		menuType = 0
	case "C": // 菜单
		menuType = 1
	case "F": // 按钮
		menuType = 2
	}

	// 转换状态
	status := 1 // 默认启用
	if req.Status == "0" {
		status = 0
	}

	// 更新菜单信息
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.ParentID = req.ParentID
	menu.Status = status
	menu.Type = menuType

	err = h.menuService.Update(menu)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, menu)
}

// Delete 删除菜单
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	err = h.menuService.Delete(uint(id))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// Get 获取菜单详情
func (h *MenuHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	menu, err := h.menuService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "菜单不存在")
		return
	}

	response.Success(c, menu)
}

// List 获取菜单列表
func (h *MenuHandler) List(c *gin.Context) {
	menus, err := h.menuService.List()
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, menus)
}

// ListByRoleID 根据角色ID获取菜单列表
func (h *MenuHandler) ListByRoleID(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
	if err != nil {
		response.Failed(c, errors.New("无效的角色ID"))
		return
	}

	menus, err := h.menuService.ListByRoleID(uint(roleID))
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, menus)
}
