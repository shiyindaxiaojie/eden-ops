package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DatabaseConfigHandler 数据库配置处理器
type DatabaseConfigHandler struct {
	databaseConfigService service.DatabaseConfigService
}

// NewDatabaseConfigHandler 创建数据库配置处理器
func NewDatabaseConfigHandler(databaseConfigService service.DatabaseConfigService) *DatabaseConfigHandler {
	return &DatabaseConfigHandler{
		databaseConfigService: databaseConfigService,
	}
}

// List 获取数据库配置列表
func (h *DatabaseConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")

	total, configs, err := h.databaseConfigService.List(page, pageSize, name)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, configs, total)
}

// Get 获取数据库配置详情
func (h *DatabaseConfigHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的数据库配置ID")
		return
	}

	config, err := h.databaseConfigService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "数据库配置不存在")
		return
	}

	response.Success(c, config)
}

// Create 创建数据库配置
func (h *DatabaseConfigHandler) Create(c *gin.Context) {
	var config model.DatabaseConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.databaseConfigService.Create(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, config)
}

// Update 更新数据库配置
func (h *DatabaseConfigHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的数据库配置ID")
		return
	}

	var config model.DatabaseConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config.ID = uint(id)
	if err := h.databaseConfigService.Update(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, config)
}

// Delete 删除数据库配置
func (h *DatabaseConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的数据库配置ID")
		return
	}

	if err := h.databaseConfigService.Delete(uint(id)); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// TestConnection 测试数据库连接
func (h *DatabaseConfigHandler) TestConnection(c *gin.Context) {
	var config model.DatabaseConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.databaseConfigService.TestConnection(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}
