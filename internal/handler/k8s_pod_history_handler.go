package handler

import (
	"eden-ops/internal/repository"
	"eden-ops/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// K8sPodHistoryHandler Pod历史数据处理器
type K8sPodHistoryHandler struct {
	podHistoryRepo repository.K8sPodHistoryRepository
}

// NewK8sPodHistoryHandler 创建Pod历史数据处理器
func NewK8sPodHistoryHandler(podHistoryRepo repository.K8sPodHistoryRepository) *K8sPodHistoryHandler {
	return &K8sPodHistoryHandler{
		podHistoryRepo: podHistoryRepo,
	}
}

// GetPodHistory 获取Pod历史记录
func (h *K8sPodHistoryHandler) GetPodHistory(c *gin.Context) {
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

	histories, total, err := h.podHistoryRepo.GetPodHistory(configID, page, pageSize, startTime, endTime)
	if err != nil {
		logger.Error("获取Pod历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pod history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     histories,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CleanupPodHistory 清理Pod历史数据
func (h *K8sPodHistoryHandler) CleanupPodHistory(c *gin.Context) {
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

	err = h.podHistoryRepo.CleanupPodHistory(beforeDate)
	if err != nil {
		logger.Error("清理Pod历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup pod history"})
		return
	}

	logger.Info("手动清理Pod历史数据成功，清理 %s 之前的数据", beforeDate.Format("2006-01-02"))
	c.JSON(http.StatusOK, gin.H{"message": "Pod history cleanup completed successfully"})
}

// GetPodHistoryStatistics 获取Pod历史数据统计
func (h *K8sPodHistoryHandler) GetPodHistoryStatistics(c *gin.Context) {
	configIDStr := c.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	// 获取Pod历史数据统计
	podHistories, podTotal, err := h.podHistoryRepo.GetPodHistory(configID, 1, 1, nil, nil)
	if err != nil {
		logger.Error("获取Pod历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pod history statistics"})
		return
	}

	statistics := gin.H{
		"podHistoryCount": podTotal,
	}

	// 获取最新的归档时间
	if len(podHistories) > 0 {
		statistics["lastArchivedAt"] = podHistories[0].ArchivedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": statistics,
	})
}
