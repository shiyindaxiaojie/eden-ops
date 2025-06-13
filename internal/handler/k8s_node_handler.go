package handler

import (
	"eden-ops/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// K8sNodeHandler 节点处理器
type K8sNodeHandler struct {
	nodeService service.K8sNodeService
}

// NewK8sNodeHandler 创建节点处理器
func NewK8sNodeHandler(nodeService service.K8sNodeService) *K8sNodeHandler {
	return &K8sNodeHandler{
		nodeService: nodeService,
	}
}

// List 获取节点列表
func (h *K8sNodeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	configID, _ := strconv.ParseInt(c.Query("configId"), 10, 64)
	name := c.Query("name")
	internalIP := c.Query("internalIP")
	status := c.Query("status")

	var ready *bool
	if readyStr := c.Query("ready"); readyStr != "" {
		if r, err := strconv.ParseBool(readyStr); err == nil {
			ready = &r
		}
	}

	nodes, total, err := h.nodeService.List(page, pageSize, configID, name, internalIP, status, ready)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取节点列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":     nodes,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetByID 根据ID获取节点详情
func (h *K8sNodeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的节点ID",
		})
		return
	}

	node, err := h.nodeService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "节点不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    node.ToResponse(),
	})
}

// Delete 删除节点
func (h *K8sNodeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的节点ID",
		})
		return
	}

	err = h.nodeService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除节点失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
