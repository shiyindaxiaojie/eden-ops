package handler

import (
	"eden-ops/internal/repository"
	"eden-ops/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// K8sWorkloadHistoryHandler Workload历史数据处理器
type K8sWorkloadHistoryHandler struct {
	workloadHistoryRepo repository.K8sWorkloadHistoryRepository
}

// NewK8sWorkloadHistoryHandler 创建Workload历史数据处理器
func NewK8sWorkloadHistoryHandler(workloadHistoryRepo repository.K8sWorkloadHistoryRepository) *K8sWorkloadHistoryHandler {
	return &K8sWorkloadHistoryHandler{
		workloadHistoryRepo: workloadHistoryRepo,
	}
}

// GetWorkloadHistory 获取Workload历史记录
func (h *K8sWorkloadHistoryHandler) GetWorkloadHistory(c *gin.Context) {
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

	histories, total, err := h.workloadHistoryRepo.GetWorkloadHistory(configID, page, pageSize, startTime, endTime)
	if err != nil {
		logger.Error("获取Workload历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workload history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     histories,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CleanupWorkloadHistory 清理Workload历史数据
func (h *K8sWorkloadHistoryHandler) CleanupWorkloadHistory(c *gin.Context) {
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

	err = h.workloadHistoryRepo.CleanupWorkloadHistory(beforeDate)
	if err != nil {
		logger.Error("清理Workload历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup workload history"})
		return
	}

	logger.Info("手动清理Workload历史数据成功，清理 %s 之前的数据", beforeDate.Format("2006-01-02"))
	c.JSON(http.StatusOK, gin.H{"message": "Workload history cleanup completed successfully"})
}

// GetWorkloadHistoryStatistics 获取Workload历史数据统计
func (h *K8sWorkloadHistoryHandler) GetWorkloadHistoryStatistics(c *gin.Context) {
	configIDStr := c.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	// 获取Workload历史数据统计
	workloadHistories, workloadTotal, err := h.workloadHistoryRepo.GetWorkloadHistory(configID, 1, 1, nil, nil)
	if err != nil {
		logger.Error("获取Workload历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workload history statistics"})
		return
	}

	statistics := gin.H{
		"workloadHistoryCount": workloadTotal,
	}

	// 获取最新的归档时间
	if len(workloadHistories) > 0 {
		statistics["lastArchivedAt"] = workloadHistories[0].ArchivedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": statistics,
	})
}
