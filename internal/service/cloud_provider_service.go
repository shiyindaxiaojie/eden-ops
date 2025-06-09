package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// CloudProviderService 云服务商服务接口
type CloudProviderService interface {
	Create(provider *model.CloudProvider) error
	Update(provider *model.CloudProvider) error
	Delete(id int64) error
	Get(id int64) (*model.CloudProvider, error)
	List(page, pageSize int, name string) (int64, []model.CloudProvider, error)
}

// cloudProviderService 云服务商服务实现
type cloudProviderService struct {
	repo repository.CloudProviderRepository
}

// NewCloudProviderService 创建云服务商服务实例
func NewCloudProviderService(repo repository.CloudProviderRepository) CloudProviderService {
	return &cloudProviderService{
		repo: repo,
	}
}

// Create 创建云服务商
func (s *cloudProviderService) Create(provider *model.CloudProvider) error {
	return s.repo.Create(provider)
}

// Update 更新云服务商
func (s *cloudProviderService) Update(provider *model.CloudProvider) error {
	return s.repo.Update(provider)
}

// Delete 删除云服务商
func (s *cloudProviderService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Get 获取云服务商详情
func (s *cloudProviderService) Get(id int64) (*model.CloudProvider, error) {
	return s.repo.FindByID(id)
}

// List 获取云服务商列表
func (s *cloudProviderService) List(page, pageSize int, name string) (int64, []model.CloudProvider, error) {
	return s.repo.List(page, pageSize, name)
}
