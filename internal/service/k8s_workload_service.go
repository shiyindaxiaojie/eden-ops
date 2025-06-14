package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"fmt"
	"strings"
	"time"
)

// K8sWorkloadService Kubernetes工作负载服务接口
type K8sWorkloadService interface {
	Create(workload *model.K8sWorkload) error
	Update(workload *model.K8sWorkload) error
	Delete(id int64) error
	Get(id int64) (*model.K8sWorkload, error)
	List(configID int64, page, pageSize int) ([]model.K8sWorkload, int64, error)
	ListWithFilter(page, pageSize int, name, namespace, workloadType, status, replicas, sortBy, sortOrder string, startTime, endTime *string, configId *int64) ([]*model.K8sWorkloadResponse, int64, error)
	ListByConfigID(configID int64) ([]model.K8sWorkload, error)
	DeleteByConfigID(configID int64) error
	SyncWorkloads(configID int64, workloads []model.K8sWorkload) error
}

// k8sWorkloadService Kubernetes工作负载服务实现
type k8sWorkloadService struct {
	repo        repository.K8sWorkloadRepository
	historyRepo repository.K8sWorkloadHistoryRepository
}

// NewK8sWorkloadService 创建Kubernetes工作负载服务
func NewK8sWorkloadService(repo repository.K8sWorkloadRepository, historyRepo repository.K8sWorkloadHistoryRepository) K8sWorkloadService {
	return &k8sWorkloadService{
		repo:        repo,
		historyRepo: historyRepo,
	}
}

// Create 创建工作负载
func (s *k8sWorkloadService) Create(workload *model.K8sWorkload) error {
	return s.repo.Create(workload)
}

// Update 更新工作负载
func (s *k8sWorkloadService) Update(workload *model.K8sWorkload) error {
	return s.repo.Update(workload)
}

// Delete 删除工作负载
func (s *k8sWorkloadService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Get 获取工作负载
func (s *k8sWorkloadService) Get(id int64) (*model.K8sWorkload, error) {
	return s.repo.Get(id)
}

// List 获取工作负载列表
func (s *k8sWorkloadService) List(configID int64, page, pageSize int) ([]model.K8sWorkload, int64, error) {
	total, workloads, err := s.repo.List(configID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return workloads, total, nil
}

// ListWithFilter 获取工作负载列表（支持筛选）
func (s *k8sWorkloadService) ListWithFilter(page, pageSize int, name, namespace, workloadType, status, replicas, sortBy, sortOrder string, startTime, endTime *string, configId *int64) ([]*model.K8sWorkloadResponse, int64, error) {
	total, workloads, err := s.repo.ListWithFilter(page, pageSize, name, namespace, workloadType, status, replicas, sortBy, sortOrder, startTime, endTime, configId)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应结构
	result := make([]*model.K8sWorkloadResponse, len(workloads))
	for i := range workloads {
		result[i] = workloads[i].ToResponse()
	}

	return result, total, nil
}

// ListByConfigID 根据配置ID获取所有工作负载
func (s *k8sWorkloadService) ListByConfigID(configID int64) ([]model.K8sWorkload, error) {
	return s.repo.ListByConfigID(configID)
}

// DeleteByConfigID 根据配置ID删除所有工作负载
func (s *k8sWorkloadService) DeleteByConfigID(configID int64) error {
	return s.repo.DeleteByConfigID(configID)
}

// SyncWorkloads 同步工作负载数据（支持历史表归档）
func (s *k8sWorkloadService) SyncWorkloads(configID int64, workloads []model.K8sWorkload) error {
	// 顺序执行，每个步骤使用独立的小事务，避免长事务

	// 1. 先归档不存在的工作负载到历史表（独立小事务）
	if err := s.retryOnLockTimeout(func() error {
		return s.historyRepo.ArchiveWorkloadsNotInList(configID, workloads, model.ArchiveReasonSyncCleanup)
	}); err != nil {
		return fmt.Errorf("failed to archive workloads: %v", err)
	}

	// 2. 删除已归档的工作负载（独立小事务）
	if err := s.retryOnLockTimeout(func() error {
		return s.repo.DeleteNotInList(configID, workloads)
	}); err != nil {
		return fmt.Errorf("failed to delete old workloads: %v", err)
	}

	// 3. 批量创建或更新新的工作负载（独立小事务）
	if len(workloads) > 0 {
		if err := s.retryOnLockTimeout(func() error {
			return s.repo.BatchCreateOrUpdate(workloads)
		}); err != nil {
			return fmt.Errorf("failed to create/update workloads: %v", err)
		}
	}

	return nil
}

// retryOnLockTimeout 在锁等待超时时重试
func (s *k8sWorkloadService) retryOnLockTimeout(fn func() error) error {
	maxRetries := 3
	baseDelay := 1000 // 1秒

	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		// 检查是否是锁等待超时错误
		if strings.Contains(err.Error(), "Lock wait timeout exceeded") {
			if i < maxRetries-1 { // 不是最后一次重试
				delay := baseDelay * (i + 1) // 递增延迟
				time.Sleep(time.Duration(delay) * time.Millisecond)
				continue
			}
		}

		// 非锁超时错误或已达到最大重试次数
		return err
	}

	return fmt.Errorf("max retries exceeded")
}
