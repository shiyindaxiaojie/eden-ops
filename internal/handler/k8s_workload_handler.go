package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"eden-ops/internal/service"
	"eden-ops/internal/pkg/response"
)

// K8sWorkloadHandler Kubernetes工作负载处理器
type K8sWorkloadHandler struct {
	workloadService service.K8sWorkloadService
}

// NewK8sWorkloadHandler 创建Kubernetes工作负载处理器
func NewK8sWorkloadHandler(workloadService service.K8sWorkloadService) *K8sWorkloadHandler {
	return &K8sWorkloadHandler{
		workloadService: workloadService,
	}
}

// List 获取工作负载列表
func (h *K8sWorkloadHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	namespace := c.Query("namespace")
	workloadType := c.Query("workloadType")

	var configId *int64
	if configIdStr := c.Query("configId"); configIdStr != "" {
		if id, err := strconv.ParseInt(configIdStr, 10, 64); err == nil {
			configId = &id
		}
	}

	workloads, total, err := h.workloadService.ListWithFilter(page, pageSize, name, namespace, workloadType, configId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, workloads, total)
}

// Get 获取工作负载详情
func (h *K8sWorkloadHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的工作负载ID")
		return
	}

	workload, err := h.workloadService.Get(id)
	if err != nil {
		response.Failed(c, err)
		return
	}

	if workload == nil {
		response.NotFound(c, "工作负载不存在")
		return
	}

	response.Success(c, workload)
}
