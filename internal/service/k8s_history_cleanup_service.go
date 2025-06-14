package service

import (
	"eden-ops/internal/repository"
	"eden-ops/pkg/logger"
	"fmt"
	"time"
)

// K8sHistoryCleanupConfig 历史表清理配置
type K8sHistoryCleanupConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	CleanupEnabled  bool          `mapstructure:"cleanup_enabled"`
	CleanupDays     int           `mapstructure:"cleanup_days"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	BatchSize       int           `mapstructure:"batch_size"`
}

// K8sHistoryCleanupService K8s历史表清理服务
type K8sHistoryCleanupService struct {
	podHistoryRepo      repository.K8sPodHistoryRepository
	nodeHistoryRepo     repository.K8sNodeHistoryRepository
	workloadHistoryRepo repository.K8sWorkloadHistoryRepository
	config              K8sHistoryCleanupConfig
	stopChan            chan struct{}
}

// NewK8sHistoryCleanupService 创建K8s历史表清理服务
func NewK8sHistoryCleanupService(
	podHistoryRepo repository.K8sPodHistoryRepository,
	nodeHistoryRepo repository.K8sNodeHistoryRepository,
	workloadHistoryRepo repository.K8sWorkloadHistoryRepository,
	config K8sHistoryCleanupConfig) *K8sHistoryCleanupService {
	return &K8sHistoryCleanupService{
		podHistoryRepo:      podHistoryRepo,
		nodeHistoryRepo:     nodeHistoryRepo,
		workloadHistoryRepo: workloadHistoryRepo,
		config:              config,
		stopChan:            make(chan struct{}),
	}
}

// Start 启动清理服务
func (s *K8sHistoryCleanupService) Start() {
	if !s.config.Enabled || !s.config.CleanupEnabled {
		logger.Info("K8s历史表清理服务已禁用")
		return
	}

	logger.Info("启动K8s历史表清理服务，清理间隔: %v，保留天数: %d", s.config.CleanupInterval, s.config.CleanupDays)

	ticker := time.NewTicker(s.config.CleanupInterval)
	defer ticker.Stop()

	// 立即执行一次清理
	s.cleanup()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopChan:
			logger.Info("K8s历史表清理服务已停止")
			return
		}
	}
}

// Stop 停止清理服务
func (s *K8sHistoryCleanupService) Stop() {
	close(s.stopChan)
}

// cleanup 执行清理操作
func (s *K8sHistoryCleanupService) cleanup() {
	beforeDate := time.Now().AddDate(0, 0, -s.config.CleanupDays)

	logger.Info("开始清理K8s历史数据，清理 %s 之前的数据", beforeDate.Format("2006-01-02 15:04:05"))

	// 清理Pod历史
	if err := s.podHistoryRepo.CleanupPodHistory(beforeDate); err != nil {
		logger.Error("清理Pod历史数据失败: %v", err)
	} else {
		logger.Info("Pod历史数据清理完成")
	}

	// 清理Node历史
	if err := s.nodeHistoryRepo.CleanupNodeHistory(beforeDate); err != nil {
		logger.Error("清理Node历史数据失败: %v", err)
	} else {
		logger.Info("Node历史数据清理完成")
	}

	// 清理Workload历史
	if err := s.workloadHistoryRepo.CleanupWorkloadHistory(beforeDate); err != nil {
		logger.Error("清理Workload历史数据失败: %v", err)
	} else {
		logger.Info("Workload历史数据清理完成")
	}

	logger.Info("K8s历史数据清理完成")
}

// GetHistoryStatistics 获取历史数据统计
func (s *K8sHistoryCleanupService) GetHistoryStatistics() (*K8sHistoryStatistics, error) {
	// 这里可以添加统计逻辑，比如各个历史表的记录数量、占用空间等
	return &K8sHistoryStatistics{
		PodHistoryCount:      0, // 实际实现时需要查询数据库
		NodeHistoryCount:     0,
		WorkloadHistoryCount: 0,
		LastCleanupTime:      time.Now(), // 实际实现时需要记录上次清理时间
	}, nil
}

// K8sHistoryStatistics K8s历史数据统计
type K8sHistoryStatistics struct {
	PodHistoryCount      int64     `json:"pod_history_count"`
	NodeHistoryCount     int64     `json:"node_history_count"`
	WorkloadHistoryCount int64     `json:"workload_history_count"`
	LastCleanupTime      time.Time `json:"last_cleanup_time"`
}

// ManualCleanup 手动清理历史数据
func (s *K8sHistoryCleanupService) ManualCleanup(beforeDate time.Time) error {
	logger.Info("开始手动清理K8s历史数据，清理 %s 之前的数据", beforeDate.Format("2006-01-02 15:04:05"))

	// 清理Pod历史
	if err := s.podHistoryRepo.CleanupPodHistory(beforeDate); err != nil {
		return fmt.Errorf("清理Pod历史数据失败: %v", err)
	}

	// 清理Node历史
	if err := s.nodeHistoryRepo.CleanupNodeHistory(beforeDate); err != nil {
		return fmt.Errorf("清理Node历史数据失败: %v", err)
	}

	// 清理Workload历史
	if err := s.workloadHistoryRepo.CleanupWorkloadHistory(beforeDate); err != nil {
		return fmt.Errorf("清理Workload历史数据失败: %v", err)
	}

	logger.Info("手动清理K8s历史数据完成")
	return nil
}
