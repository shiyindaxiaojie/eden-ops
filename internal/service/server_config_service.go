package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// ServerConfigService 服务器配置服务接口
type ServerConfigService interface {
	Create(config *model.ServerConfig) error
	Update(config *model.ServerConfig) error
	Delete(id uint) error
	Get(id uint) (*model.ServerConfig, error)
	List(page, pageSize int, name string) (int64, []model.ServerConfig, error)
	TestConnection(config *model.ServerConfig) error
}

// serverConfigService 服务器配置服务实现
type serverConfigService struct {
	repo *repository.ServerConfigRepository
}

// NewServerConfigService 创建服务器配置服务
func NewServerConfigService(repo *repository.ServerConfigRepository) ServerConfigService {
	return &serverConfigService{repo: repo}
}

// List 获取服务器配置列表
func (s *serverConfigService) List(page, pageSize int, name string) (int64, []model.ServerConfig, error) {
	return s.repo.List(page, pageSize, name)
}

// Create 创建服务器配置
func (s *serverConfigService) Create(config *model.ServerConfig) error {
	return s.repo.Create(config)
}

// Get 获取服务器配置详情
func (s *serverConfigService) Get(id uint) (*model.ServerConfig, error) {
	return s.repo.Get(id)
}

// Update 更新服务器配置
func (s *serverConfigService) Update(config *model.ServerConfig) error {
	return s.repo.Update(config)
}

// Delete 删除服务器配置
func (s *serverConfigService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// TestConnection 测试服务器连接
func (s *serverConfigService) TestConnection(config *model.ServerConfig) error {
	// TODO: 实现服务器连接测试逻辑
	return nil
}
