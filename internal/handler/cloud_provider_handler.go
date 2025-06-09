package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/pkg/response"
	"eden-ops/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CloudProviderHandler 云服务商处理器
type CloudProviderHandler struct {
	cloudProviderService service.CloudProviderService
}

// NewCloudProviderHandler 创建云服务商处理器
func NewCloudProviderHandler(cloudProviderService service.CloudProviderService) *CloudProviderHandler {
	return &CloudProviderHandler{
		cloudProviderService: cloudProviderService,
	}
}

// List 获取云服务商列表
func (h *CloudProviderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")

	total, providers, err := h.cloudProviderService.List(page, pageSize, name)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, providers, total)
}

// Get 获取云服务商详情
func (h *CloudProviderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, err)
		return
	}

	provider, err := h.cloudProviderService.Get(id)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, provider)
}

// Create 创建云服务商
func (h *CloudProviderHandler) Create(c *gin.Context) {
	var provider model.CloudProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	if err := h.cloudProviderService.Create(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// Update 更新云服务商
func (h *CloudProviderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, err)
		return
	}

	var provider model.CloudProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	provider.ID = uint(id)
	if err := h.cloudProviderService.Update(&provider); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// Delete 删除云服务商
func (h *CloudProviderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, err)
		return
	}

	if err := h.cloudProviderService.Delete(id); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}
