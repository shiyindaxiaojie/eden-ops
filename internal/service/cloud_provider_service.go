package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// CloudProviderService 云厂商服务接口
type CloudProviderService interface {
	List(page, pageSize int, name string, status *int) ([]model.CloudProvider, int64, error)
	Get(id uint) (*model.CloudProvider, error)
	Create(provider *model.CloudProvider) error
	Update(provider *model.CloudProvider) error
	Delete(id uint) error
}

// cloudProviderService 云厂商服务实现
type cloudProviderService struct {
	repo repository.CloudProviderRepository
}

// NewCloudProviderService 创建云厂商服务
func NewCloudProviderService(repo repository.CloudProviderRepository) CloudProviderService {
	return &cloudProviderService{
		repo: repo,
	}
}

// List 获取云厂商列表
func (s *cloudProviderService) List(page, pageSize int, name string, status *int) ([]model.CloudProvider, int64, error) {
	return s.repo.List(page, pageSize, name, status)
}

// Get 获取云厂商详情
func (s *cloudProviderService) Get(id uint) (*model.CloudProvider, error) {
	return s.repo.Get(id)
}

// Create 创建云厂商
func (s *cloudProviderService) Create(provider *model.CloudProvider) error {
	return s.repo.Create(provider)
}

// Update 更新云厂商
func (s *cloudProviderService) Update(provider *model.CloudProvider) error {
	return s.repo.Update(provider)
}

// Delete 删除云厂商
func (s *cloudProviderService) Delete(id uint) error {
	return s.repo.Delete(id)
}
