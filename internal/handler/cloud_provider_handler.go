package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CloudProviderHandler 云厂商处理器
type CloudProviderHandler struct {
	cloudProviderService service.CloudProviderService
}

// NewCloudProviderHandler 创建云厂商处理器
func NewCloudProviderHandler(cloudProviderService service.CloudProviderService) *CloudProviderHandler {
	return &CloudProviderHandler{
		cloudProviderService: cloudProviderService,
	}
}

// List 获取云厂商列表
func (h *CloudProviderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")
	status := c.DefaultQuery("status", "")

	var statusInt *int
	if status != "" {
		s, err := strconv.Atoi(status)
		if err == nil {
			statusInt = &s
		}
	}

	providers, total, err := h.cloudProviderService.List(page, pageSize, name, statusInt)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, providers, total)
}

// Get 获取云厂商详情
func (h *CloudProviderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云厂商ID")
		return
	}

	provider, err := h.cloudProviderService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "云厂商不存在")
		return
	}

	response.Success(c, provider)
}

// Create 创建云厂商
func (h *CloudProviderHandler) Create(c *gin.Context) {
	var provider model.CloudProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cloudProviderService.Create(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, provider)
}

// Update 更新云厂商
func (h *CloudProviderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云厂商ID")
		return
	}

	var provider model.CloudProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider.ID = uint(id)
	if err := h.cloudProviderService.Update(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, provider)
}

// Delete 删除云厂商
func (h *CloudProviderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云厂商ID")
		return
	}

	if err := h.cloudProviderService.Delete(uint(id)); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}
