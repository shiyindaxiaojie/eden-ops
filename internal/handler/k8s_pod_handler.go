package handler

import (
	"eden-ops/internal/service"
	"eden-ops/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// K8sPodHandler K8s Pod处理器
type K8sPodHandler struct {
	podService service.K8sPodService
}

// NewK8sPodHandler 创建K8s Pod处理器
func NewK8sPodHandler(podService service.K8sPodService) *K8sPodHandler {
	return &K8sPodHandler{podService: podService}
}

// List 获取Pod列表
func (h *K8sPodHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	namespace := c.Query("namespace")
	workloadName := c.Query("workloadName")
	status := c.Query("status")
	instanceIP := c.Query("instanceIP")
	sortBy := c.Query("sortBy")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var configId *int64
	if configIdStr := c.Query("configId"); configIdStr != "" {
		if id, err := strconv.ParseInt(configIdStr, 10, 64); err == nil {
			configId = &id
		}
	}

	// 处理时间参数
	var startTime, endTime *string
	if startTimeStr != "" {
		startTime = &startTimeStr
	}
	if endTimeStr != "" {
		endTime = &endTimeStr
	}

	pods, total, err := h.podService.ListWithFilter(page, pageSize, name, namespace, workloadName, status, instanceIP, sortBy, sortOrder, startTime, endTime, configId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.PageSuccess(c, pods, total)
}

// Get 获取Pod详情
func (h *K8sPodHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的Pod ID")
		return
	}

	pod, err := h.podService.Get(id)
	if err != nil {
		response.Failed(c, err)
		return
	}

	if pod == nil {
		response.NotFound(c, "Pod不存在")
		return
	}

	response.Success(c, pod)
}
