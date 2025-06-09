package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sConfigService Kubernetes配置服务接口
type K8sConfigService interface {
	Create(config *model.K8sConfig) error
	Update(config *model.K8sConfig) error
	Delete(id uint) error
	Get(id uint) (*model.K8sConfig, error)
	List(page, pageSize int) ([]*model.K8sConfig, int64, error)
	TestConnection(config *model.K8sConfig) error
	SyncCluster(id int64) error
	GetNamespaces(id int64) ([]string, error)
}

// k8sConfigService Kubernetes配置服务实现
type k8sConfigService struct {
	repo repository.K8sConfigRepository
}

// NewK8sConfigService 创建Kubernetes配置服务
func NewK8sConfigService(repo repository.K8sConfigRepository) K8sConfigService {
	return &k8sConfigService{
		repo: repo,
	}
}

// Create 创建Kubernetes配置
func (s *k8sConfigService) Create(config *model.K8sConfig) error {
	return s.repo.Create(config)
}

// Update 更新Kubernetes配置
func (s *k8sConfigService) Update(config *model.K8sConfig) error {
	return s.repo.Update(config)
}

// Delete 删除Kubernetes配置
func (s *k8sConfigService) Delete(id uint) error {
	return s.repo.Delete(int64(id))
}

// Get 获取Kubernetes配置
func (s *k8sConfigService) Get(id uint) (*model.K8sConfig, error) {
	return s.repo.Get(int64(id))
}

// List 获取Kubernetes配置列表
func (s *k8sConfigService) List(page, pageSize int) ([]*model.K8sConfig, int64, error) {
	total, configs, err := s.repo.List(page, pageSize, "")
	if err != nil {
		return nil, 0, err
	}

	var result []*model.K8sConfig
	for i := range configs {
		result = append(result, &configs[i])
	}

	return result, total, nil
}

// TestConnection 测试Kubernetes连接
func (s *k8sConfigService) TestConnection(config *model.K8sConfig) error {
	k8sConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config.Kubeconfig))
	if err != nil {
		return fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get cluster version: %v", err)
	}

	return nil
}

// SyncCluster 同步集群信息
func (s *k8sConfigService) SyncCluster(id int64) error {
	// TODO: 实现集群同步逻辑
	return nil
}

// GetNamespaces 获取命名空间列表
func (s *k8sConfigService) GetNamespaces(id int64) ([]string, error) {
	// 简化实现，返回默认命名空间
	return []string{"default", "kube-system", "kube-public"}, nil
}
