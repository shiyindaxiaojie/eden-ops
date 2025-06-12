package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// K8sConfigHandler Kubernetes配置处理器
type K8sConfigHandler struct {
	k8sConfigService service.K8sConfigService
}

// NewK8sConfigHandler 创建Kubernetes配置处理器
func NewK8sConfigHandler(k8sConfigService service.K8sConfigService) *K8sConfigHandler {
	return &K8sConfigHandler{
		k8sConfigService: k8sConfigService,
	}
}

// List 获取Kubernetes配置列表
func (h *K8sConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	var providerId *int64
	if providerIdStr := c.Query("providerId"); providerIdStr != "" {
		if p, err := strconv.ParseInt(providerIdStr, 10, 64); err == nil {
			providerId = &p
		}
	}

	configs, total, err := h.k8sConfigService.List(page, pageSize, name, status, providerId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, configs, total)
}

// Get 获取Kubernetes配置详情
func (h *K8sConfigHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的Kubernetes配置ID")
		return
	}

	config, err := h.k8sConfigService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "Kubernetes配置不存在")
		return
	}

	response.Success(c, config)
}

// Create 创建Kubernetes配置
func (h *K8sConfigHandler) Create(c *gin.Context) {
	var config model.K8sConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.k8sConfigService.CreateWithClusterInfo(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, config)
}

// Update 更新Kubernetes配置
func (h *K8sConfigHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的Kubernetes配置ID")
		return
	}

	var config model.K8sConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config.ID = int64(id)
	if err := h.k8sConfigService.UpdateWithClusterInfo(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, config)
}

// Delete 删除Kubernetes配置
func (h *K8sConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的Kubernetes配置ID")
		return
	}

	if err := h.k8sConfigService.Delete(uint(id)); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// TestConnection 测试Kubernetes连接
func (h *K8sConfigHandler) TestConnection(c *gin.Context) {
	var config model.K8sConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.k8sConfigService.TestConnection(&config); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// GetNamespaces 获取命名空间列表
func (h *K8sConfigHandler) GetNamespaces(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, err)
		return
	}

	namespaces, err := h.k8sConfigService.GetNamespaces(id)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, namespaces)
}
