package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// ServerConfigRepository handles database operations for server configurations
type ServerConfigRepository struct {
	db *gorm.DB
}

// NewServerConfigRepository creates a new ServerConfigRepository instance
func NewServerConfigRepository(db *gorm.DB) *ServerConfigRepository {
	return &ServerConfigRepository{
		db: db,
	}
}

// List 获取服务器配置列表
func (r *ServerConfigRepository) List(page, pageSize int, name string) (int64, []model.ServerConfig, error) {
	var total int64
	var configs []model.ServerConfig

	query := r.db.Model(&model.ServerConfig{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&configs).Error
	if err != nil {
		return 0, nil, err
	}

	return total, configs, nil
}

// Create 创建服务器配置
func (r *ServerConfigRepository) Create(config *model.ServerConfig) error {
	return r.db.Create(config).Error
}

// Get 获取服务器配置详情
func (r *ServerConfigRepository) Get(id uint) (*model.ServerConfig, error) {
	var config model.ServerConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Update 更新服务器配置
func (r *ServerConfigRepository) Update(config *model.ServerConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除服务器配置
func (r *ServerConfigRepository) Delete(id uint) error {
	return r.db.Delete(&model.ServerConfig{}, id).Error
}
