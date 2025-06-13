package handler

import (
	"eden-ops/internal/model"
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CloudAccountHandler 云账号处理器
type CloudAccountHandler struct {
	cloudAccountService service.CloudAccountService
}

// NewCloudAccountHandler 创建云账号处理器
func NewCloudAccountHandler(cloudAccountService service.CloudAccountService) *CloudAccountHandler {
	return &CloudAccountHandler{
		cloudAccountService: cloudAccountService,
	}
}

// List 获取云账号列表
func (h *CloudAccountHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")
	providerIDStr := c.DefaultQuery("providerId", "")
	statusStr := c.DefaultQuery("status", "")

	var providerID *int64
	if providerIDStr != "" {
		if pid, err := strconv.ParseInt(providerIDStr, 10, 64); err == nil {
			providerID = &pid
		}
	}

	var status *int
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	accounts, total, err := h.cloudAccountService.ListWithFilter(page, pageSize, name, providerID, status)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, accounts, total)
}

// Get 获取云账号详情
func (h *CloudAccountHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云账号ID")
		return
	}

	account, err := h.cloudAccountService.Get(uint(id))
	if err != nil {
		response.NotFound(c, "云账号不存在")
		return
	}

	response.Success(c, account)
}

// Create 创建云账号
func (h *CloudAccountHandler) Create(c *gin.Context) {
	var account model.CloudAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cloudAccountService.Create(&account); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, account)
}

// Update 更新云账号
func (h *CloudAccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云账号ID")
		return
	}

	var account model.CloudAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	account.ID = int64(id)
	if err := h.cloudAccountService.Update(&account); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, account)
}

// Delete 删除云账号
func (h *CloudAccountHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的云账号ID")
		return
	}

	if err := h.cloudAccountService.Delete(uint(id)); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}

// TestConnection 测试云账号连接
func (h *CloudAccountHandler) TestConnection(c *gin.Context) {
	var account model.CloudAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cloudAccountService.TestConnection(&account); err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, nil)
}
