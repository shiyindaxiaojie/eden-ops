package handler

import (
	"eden-ops/internal/repository"
	"eden-ops/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// K8sHistoryHandler K8s历史数据处理器（统一接口）
type K8sHistoryHandler struct {
	podHistoryHandler      *K8sPodHistoryHandler
	nodeHistoryHandler     *K8sNodeHistoryHandler
	workloadHistoryHandler *K8sWorkloadHistoryHandler
	podHistoryRepo         repository.K8sPodHistoryRepository
	nodeHistoryRepo        repository.K8sNodeHistoryRepository
	workloadHistoryRepo    repository.K8sWorkloadHistoryRepository
}

// NewK8sHistoryHandler 创建K8s历史数据处理器
func NewK8sHistoryHandler(
	podHistoryRepo repository.K8sPodHistoryRepository,
	nodeHistoryRepo repository.K8sNodeHistoryRepository,
	workloadHistoryRepo repository.K8sWorkloadHistoryRepository) *K8sHistoryHandler {
	return &K8sHistoryHandler{
		podHistoryHandler:      NewK8sPodHistoryHandler(podHistoryRepo),
		nodeHistoryHandler:     NewK8sNodeHistoryHandler(nodeHistoryRepo),
		workloadHistoryHandler: NewK8sWorkloadHistoryHandler(workloadHistoryRepo),
		podHistoryRepo:         podHistoryRepo,
		nodeHistoryRepo:        nodeHistoryRepo,
		workloadHistoryRepo:    workloadHistoryRepo,
	}
}

// GetPodHistory 获取Pod历史记录
func (h *K8sHistoryHandler) GetPodHistory(c *gin.Context) {
	h.podHistoryHandler.GetPodHistory(c)
}

// GetNodeHistory 获取Node历史记录
func (h *K8sHistoryHandler) GetNodeHistory(c *gin.Context) {
	h.nodeHistoryHandler.GetNodeHistory(c)
}

// GetWorkloadHistory 获取Workload历史记录
func (h *K8sHistoryHandler) GetWorkloadHistory(c *gin.Context) {
	h.workloadHistoryHandler.GetWorkloadHistory(c)
}

// CleanupHistory 手动清理历史数据
func (h *K8sHistoryHandler) CleanupHistory(c *gin.Context) {
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

	// 清理Pod历史
	if err := h.podHistoryRepo.CleanupPodHistory(beforeDate); err != nil {
		logger.Error("清理Pod历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup pod history"})
		return
	}

	// 清理Node历史
	if err := h.nodeHistoryRepo.CleanupNodeHistory(beforeDate); err != nil {
		logger.Error("清理Node历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup node history"})
		return
	}

	// 清理Workload历史
	if err := h.workloadHistoryRepo.CleanupWorkloadHistory(beforeDate); err != nil {
		logger.Error("清理Workload历史数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup workload history"})
		return
	}

	logger.Info("手动清理历史数据成功，清理 %s 之前的数据", beforeDate.Format("2006-01-02"))
	c.JSON(http.StatusOK, gin.H{"message": "History cleanup completed successfully"})
}

// GetHistoryStatistics 获取历史数据统计
func (h *K8sHistoryHandler) GetHistoryStatistics(c *gin.Context) {
	configIDStr := c.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	// 获取各类型历史数据的统计
	podTotal, err := h.podHistoryRepo.CountPodHistory(configID)
	if err != nil {
		logger.Error("获取Pod历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pod history statistics"})
		return
	}

	nodeTotal, err := h.nodeHistoryRepo.CountNodeHistory(configID)
	if err != nil {
		logger.Error("获取Node历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get node history statistics"})
		return
	}

	workloadTotal, err := h.workloadHistoryRepo.CountWorkloadHistory(configID)
	if err != nil {
		logger.Error("获取Workload历史统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workload history statistics"})
		return
	}

	// 获取最新的归档时间
	var lastArchivedAt *time.Time
	
	// 获取Pod最新归档时间
	podHistories, _, err := h.podHistoryRepo.GetPodHistory(configID, 1, 1, nil, nil)
	if err == nil && len(podHistories) > 0 {
		lastArchivedAt = &podHistories[0].ArchivedAt
	}
	
	// 获取Node最新归档时间
	nodeHistories, _, err := h.nodeHistoryRepo.GetNodeHistory(configID, 1, 1, nil, nil)
	if err == nil && len(nodeHistories) > 0 && (lastArchivedAt == nil || nodeHistories[0].ArchivedAt.After(*lastArchivedAt)) {
		lastArchivedAt = &nodeHistories[0].ArchivedAt
	}
	
	// 获取Workload最新归档时间
	workloadHistories, _, err := h.workloadHistoryRepo.GetWorkloadHistory(configID, 1, 1, nil, nil)
	if err == nil && len(workloadHistories) > 0 && (lastArchivedAt == nil || workloadHistories[0].ArchivedAt.After(*lastArchivedAt)) {
		lastArchivedAt = &workloadHistories[0].ArchivedAt
	}

	statistics := gin.H{
		"podHistoryCount":      podTotal,
		"nodeHistoryCount":     nodeTotal,
		"workloadHistoryCount": workloadTotal,
		"totalHistoryCount":    podTotal + nodeTotal + workloadTotal,
	}

	if lastArchivedAt != nil {
		statistics["lastArchivedAt"] = lastArchivedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": statistics,
	})
}
