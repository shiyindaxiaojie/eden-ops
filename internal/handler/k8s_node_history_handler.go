package handler

import (
	"eden-ops/internal/repository"
	"eden-ops/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// K8sNodeHistoryHandler Node历史数据处理器
type K8sNodeHistoryHandler struct {
	nodeHistoryRepo repository.K8sNodeHistoryRepository
}

// NewK8sNodeHistoryHandler 创建Node历史数据处理器
func NewK8sNodeHistoryHandler(nodeHistoryRepo repository.K8sNodeHistoryRepository) *K8sNodeHistoryHandler {
	return &K8sNodeHistoryHandler{
		nodeHistoryRepo: nodeHistoryRepo,
	}
}

// GetNodeHistory 获取Node历史记录
func (h *K8sNodeHistoryHandler) GetNodeHistory(c *gin.Context) {
	configIDStr := c.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析时间范围
	var startTime, endTime *time.Time
	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
			startTime = &t
		}
	}
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
			endTime = &t
		}
	}

	histories, total, err := h.nodeHistoryRepo.GetNodeHistory(configID, page, pageSize, startTime, endTime)
	if err != nil {
		logger.Error("获取Node历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get node history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     histories,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CleanupNodeHistory 清理Node历史数据
func (h *K8sNodeHistoryHandler) CleanupNodeHistory(c *gin.Context) {
	var req struct {
		BeforeDate string `json:"beforeDate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	beforeDate, err := time.Parse("2006-01-02", req.BeforeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	err = h.nodeHistoryRepo.CleanupNodeHistory(beforeDate)
	if err != nil {
		logger.Error("清理Node历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup node history"})
		return
	}

	logger.Info("手动清理Node历史数据成功，清理 %s 之前的数据", beforeDate.Format("2006-01-02"))
	c.JSON(http.StatusOK, gin.H{"message": "Node history cleanup completed successfully"})
}

// GetNodeHistoryStatistics 获取Node历史数据统计
func (h *K8sNodeHistoryHandler) GetNodeHistoryStatistics(c *gin.Context) {
	configIDStr := c.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	// 获取Node历史数据统计
	nodeHistories, nodeTotal, err := h.nodeHistoryRepo.GetNodeHistory(configID, 1, 1, nil, nil)
	if err != nil {
		logger.Error("获取Node历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get node history statistics"})
		return
	}

	statistics := gin.H{
		"nodeHistoryCount": nodeTotal,
	}

	// 获取最新的归档时间
	if len(nodeHistories) > 0 {
		statistics["lastArchivedAt"] = nodeHistories[0].ArchivedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": statistics,
	})
}
