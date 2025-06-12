package handler

import (
	"eden-ops/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// K8sWorkloadNamespaceHandler K8s工作负载命名空间处理器
type K8sWorkloadNamespaceHandler struct {
	repo repository.K8sWorkloadNamespaceRepository
}

// NewK8sWorkloadNamespaceHandler 创建K8s工作负载命名空间处理器
func NewK8sWorkloadNamespaceHandler(repo repository.K8sWorkloadNamespaceRepository) *K8sWorkloadNamespaceHandler {
	return &K8sWorkloadNamespaceHandler{repo: repo}
}

// GetNamespacesByConfigID 根据配置ID获取命名空间列表
func (h *K8sWorkloadNamespaceHandler) GetNamespacesByConfigID(c *gin.Context) {
	configIDStr := c.Query("configId")
	if configIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "configId is required",
			"data":    nil,
		})
		return
	}

	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid configId",
			"data":    nil,
		})
		return
	}

	namespaces, err := h.repo.GetByConfigID(configID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get namespaces",
			"data":    nil,
		})
		return
	}

	// 转换为响应格式
	var result []string
	for _, ns := range namespaces {
		result = append(result, ns.Namespace)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}
