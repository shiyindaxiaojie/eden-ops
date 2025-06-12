package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// K8sPodService K8s Pod服务接口
type K8sPodService interface {
	Create(pod *model.K8sPod) error
	Update(pod *model.K8sPod) error
	Delete(id int64) error
	Get(id int64) (*model.K8sPod, error)
	List(configID int64, page, pageSize int) ([]model.K8sPod, int64, error)
	ListWithFilter(page, pageSize int, name, namespace, workloadName, status, sortBy, sortOrder string, startTime, endTime *string, configId *int64) ([]*model.K8sPodResponse, int64, error)
	ListByConfigID(configID int64) ([]model.K8sPod, error)
	DeleteByConfigID(configID int64) error
	SyncPods(configID int64, pods []model.K8sPod) error
}

// k8sPodService K8s Pod服务实现
type k8sPodService struct {
	repo repository.K8sPodRepository
}

// NewK8sPodService 创建K8s Pod服务
func NewK8sPodService(repo repository.K8sPodRepository) K8sPodService {
	return &k8sPodService{
		repo: repo,
	}
}

// Create 创建Pod
func (s *k8sPodService) Create(pod *model.K8sPod) error {
	return s.repo.Create(pod)
}

// Update 更新Pod
func (s *k8sPodService) Update(pod *model.K8sPod) error {
	return s.repo.Update(pod)
}

// Delete 删除Pod
func (s *k8sPodService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Get 获取Pod
func (s *k8sPodService) Get(id int64) (*model.K8sPod, error) {
	return s.repo.Get(id)
}

// List 获取Pod列表
func (s *k8sPodService) List(configID int64, page, pageSize int) ([]model.K8sPod, int64, error) {
	total, pods, err := s.repo.List(configID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return pods, total, nil
}

// ListWithFilter 获取Pod列表（支持筛选）
func (s *k8sPodService) ListWithFilter(page, pageSize int, name, namespace, workloadName, status, sortBy, sortOrder string, startTime, endTime *string, configId *int64) ([]*model.K8sPodResponse, int64, error) {
	total, pods, err := s.repo.ListWithFilter(page, pageSize, name, namespace, workloadName, status, sortBy, sortOrder, startTime, endTime, configId)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应结构
	result := make([]*model.K8sPodResponse, len(pods))
	for i := range pods {
		result[i] = pods[i].ToResponse()
	}

	return result, total, nil
}

// ListByConfigID 根据配置ID获取所有Pod
func (s *k8sPodService) ListByConfigID(configID int64) ([]model.K8sPod, error) {
	return s.repo.ListByConfigID(configID)
}

// DeleteByConfigID 根据配置ID删除所有Pod
func (s *k8sPodService) DeleteByConfigID(configID int64) error {
	return s.repo.DeleteByConfigID(configID)
}

// SyncPods 同步Pod数据
func (s *k8sPodService) SyncPods(configID int64, pods []model.K8sPod) error {
	// 先删除该集群的所有Pod
	if err := s.repo.DeleteByConfigID(configID); err != nil {
		return err
	}
	
	// 批量创建新的Pod
	if len(pods) > 0 {
		return s.repo.BatchCreate(pods)
	}
	
	return nil
}
