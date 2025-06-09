package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// CloudAccountService 云账号服务接口
type CloudAccountService interface {
	Create(account *model.CloudAccount) error
	Update(account *model.CloudAccount) error
	Delete(id uint) error
	Get(id uint) (*model.CloudAccount, error)
	List(page, pageSize int, name string) ([]*model.CloudAccount, int64, error)
	TestConnection(account *model.CloudAccount) error
}

// cloudAccountService 云账号服务实现
type cloudAccountService struct {
	repo *repository.CloudAccountRepository
}

// NewCloudAccountService 创建云账号服务
func NewCloudAccountService(repo *repository.CloudAccountRepository) CloudAccountService {
	return &cloudAccountService{
		repo: repo,
	}
}

// List 获取云账号列表
func (s *cloudAccountService) List(page, pageSize int, name string) ([]*model.CloudAccount, int64, error) {
	total, accounts, err := s.repo.List(page, pageSize, name)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.CloudAccount
	for i := range accounts {
		result = append(result, &accounts[i])
	}

	return result, total, nil
}

// Create 创建云账号
func (s *cloudAccountService) Create(account *model.CloudAccount) error {
	return s.repo.Create(account)
}

// Get 获取云账号详情
func (s *cloudAccountService) Get(id uint) (*model.CloudAccount, error) {
	return s.repo.Get(id)
}

// Update 更新云账号
func (s *cloudAccountService) Update(account *model.CloudAccount) error {
	return s.repo.Update(account)
}

// Delete 删除云账号
func (s *cloudAccountService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// TestConnection 测试云账号连接
func (s *cloudAccountService) TestConnection(account *model.CloudAccount) error {
	// TODO: 实现云账号连接测试逻辑
	return nil
}
