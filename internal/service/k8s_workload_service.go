package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// K8sWorkloadService Kubernetes工作负载服务接口
type K8sWorkloadService interface {
	Create(workload *model.K8sWorkload) error
	Update(workload *model.K8sWorkload) error
	Delete(id int64) error
	Get(id int64) (*model.K8sWorkload, error)
	List(configID int64, page, pageSize int) ([]model.K8sWorkload, int64, error)
	ListByConfigID(configID int64) ([]model.K8sWorkload, error)
	DeleteByConfigID(configID int64) error
	SyncWorkloads(configID int64, workloads []model.K8sWorkload) error
}

// k8sWorkloadService Kubernetes工作负载服务实现
type k8sWorkloadService struct {
	repo repository.K8sWorkloadRepository
}

// NewK8sWorkloadService 创建Kubernetes工作负载服务
func NewK8sWorkloadService(repo repository.K8sWorkloadRepository) K8sWorkloadService {
	return &k8sWorkloadService{
		repo: repo,
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

// ListByConfigID 根据配置ID获取所有工作负载
func (s *k8sWorkloadService) ListByConfigID(configID int64) ([]model.K8sWorkload, error) {
	return s.repo.ListByConfigID(configID)
}

// DeleteByConfigID 根据配置ID删除所有工作负载
func (s *k8sWorkloadService) DeleteByConfigID(configID int64) error {
	return s.repo.DeleteByConfigID(configID)
}

// SyncWorkloads 同步工作负载数据
func (s *k8sWorkloadService) SyncWorkloads(configID int64, workloads []model.K8sWorkload) error {
	// 先删除该集群的所有工作负载
	if err := s.repo.DeleteByConfigID(configID); err != nil {
		return err
	}
	
	// 批量创建新的工作负载
	if len(workloads) > 0 {
		return s.repo.BatchCreate(workloads)
	}
	
	return nil
}
