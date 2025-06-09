package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// ServerConfigHandler 服务器配置处理器
type ServerConfigHandler struct {
	serverConfigService service.ServerConfigService
	logger              *logrus.Logger
}

// NewServerConfigHandler 创建服务器配置处理器
func NewServerConfigHandler(serverConfigService service.ServerConfigService, logger *logrus.Logger) *ServerConfigHandler {
	return &ServerConfigHandler{
		serverConfigService: serverConfigService,
		logger:              logger,
	}
}

// List 获取服务器配置列表
func (h *ServerConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")

	total, configs, err := h.serverConfigService.List(page, pageSize, name)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, configs, total)
}

// Get 获取服务器配置
func (h *ServerConfigHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务器配置ID")
		return
	}

	config, err := h.serverConfigService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "服务器配置不存在")
		return
	}

	response.Success(c, config)
}

// Create 创建服务器配置
func (h *ServerConfigHandler) Create(c *gin.Context) {
	var config model.ServerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.serverConfigService.Create(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, config)
}

// Update 更新服务器配置
func (h *ServerConfigHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的服务器配置ID")
		return
	}

	var config model.ServerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	config.ID = uint(id)
	if err := h.serverConfigService.Update(&config); err != nil {
		h.logger.Error("更新服务器配置失败", zap.Error(err))
		response.InternalServerError(c, "更新服务器配置失败")
		return
	}

	response.Success(c, config)
}

// Delete 删除服务器配置
func (h *ServerConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的服务器配置ID")
		return
	}

	if err := h.serverConfigService.Delete(uint(id)); err != nil {
		h.logger.Error("删除服务器配置失败", zap.Error(err))
		response.InternalServerError(c, "删除服务器配置失败")
		return
	}

	response.Success(c, nil)
}

// TestConnection 测试连接
func (h *ServerConfigHandler) TestConnection(c *gin.Context) {
	var config model.ServerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	if err := h.serverConfigService.TestConnection(&config); err != nil {
		h.logger.Error("测试服务器连接失败", zap.Error(err))
		response.InternalServerError(c, "测试服务器连接失败")
		return
	}

	response.Success(c, nil)
}
