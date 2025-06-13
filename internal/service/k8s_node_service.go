package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// K8sNodeService 节点服务接口
type K8sNodeService interface {
	Create(node *model.K8sNode) error
	Update(node *model.K8sNode) error
	Delete(id int64) error
	GetByID(id int64) (*model.K8sNode, error)
	List(page, pageSize int, configID int64, name, internalIP, status string, ready *bool) ([]*model.K8sNodeResponse, int64, error)
	BatchCreateOrUpdate(nodes []model.K8sNode) error
}

// k8sNodeService 节点服务实现
type k8sNodeService struct {
	repo repository.K8sNodeRepository
}

// NewK8sNodeService 创建节点服务
func NewK8sNodeService(repo repository.K8sNodeRepository) K8sNodeService {
	return &k8sNodeService{repo: repo}
}

// Create 创建节点
func (s *k8sNodeService) Create(node *model.K8sNode) error {
	return s.repo.Create(node)
}

// Update 更新节点
func (s *k8sNodeService) Update(node *model.K8sNode) error {
	return s.repo.Update(node)
}

// Delete 删除节点
func (s *k8sNodeService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID获取节点
func (s *k8sNodeService) GetByID(id int64) (*model.K8sNode, error) {
	return s.repo.GetByID(id)
}

// List 获取节点列表
func (s *k8sNodeService) List(page, pageSize int, configID int64, name, internalIP, status string, ready *bool) ([]*model.K8sNodeResponse, int64, error) {
	total, nodes, err := s.repo.List(page, pageSize, configID, name, internalIP, status, ready)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.K8sNodeResponse
	for i := range nodes {
		result = append(result, nodes[i].ToResponse())
	}

	return result, total, nil
}

// BatchCreateOrUpdate 批量创建或更新节点
func (s *k8sNodeService) BatchCreateOrUpdate(nodes []model.K8sNode) error {
	return s.repo.BatchCreateOrUpdate(nodes)
}
