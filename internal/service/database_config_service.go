package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// DatabaseConfigService 数据库配置服务接口
type DatabaseConfigService interface {
	Create(config *model.DatabaseConfig) error
	Update(config *model.DatabaseConfig) error
	Delete(id uint) error
	Get(id uint) (*model.DatabaseConfig, error)
	List(page, pageSize int, name string) (int64, []model.DatabaseConfig, error)
	TestConnection(config *model.DatabaseConfig) error
}

// databaseConfigService 数据库配置服务实现
type databaseConfigService struct {
	repo *repository.DatabaseConfigRepository
}

// NewDatabaseConfigService 创建数据库配置服务
func NewDatabaseConfigService(repo *repository.DatabaseConfigRepository) DatabaseConfigService {
	return &databaseConfigService{repo: repo}
}

// List 获取数据库配置列表
func (s *databaseConfigService) List(page, pageSize int, name string) (int64, []model.DatabaseConfig, error) {
	return s.repo.List(page, pageSize, name)
}

// Create 创建数据库配置
func (s *databaseConfigService) Create(config *model.DatabaseConfig) error {
	return s.repo.Create(config)
}

// Get 获取数据库配置详情
func (s *databaseConfigService) Get(id uint) (*model.DatabaseConfig, error) {
	return s.repo.Get(id)
}

// Update 更新数据库配置
func (s *databaseConfigService) Update(config *model.DatabaseConfig) error {
	return s.repo.Update(config)
}

// Delete 删除数据库配置
func (s *databaseConfigService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// TestConnection 测试数据库连接
func (s *databaseConfigService) TestConnection(config *model.DatabaseConfig) error {
	// TODO: 实现数据库连接测试逻辑
	return nil
}
