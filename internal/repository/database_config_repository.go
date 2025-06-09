package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// DatabaseConfigRepository handles database operations for database configurations
type DatabaseConfigRepository struct {
	db *gorm.DB
}

// NewDatabaseConfigRepository creates a new DatabaseConfigRepository instance
func NewDatabaseConfigRepository(db *gorm.DB) *DatabaseConfigRepository {
	return &DatabaseConfigRepository{
		db: db,
	}
}

// List 获取数据库配置列表
func (r *DatabaseConfigRepository) List(page, pageSize int, name string) (int64, []model.DatabaseConfig, error) {
	var total int64
	var configs []model.DatabaseConfig

	query := r.db.Model(&model.DatabaseConfig{})
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

// Create 创建数据库配置
func (r *DatabaseConfigRepository) Create(config *model.DatabaseConfig) error {
	return r.db.Create(config).Error
}

// Get 获取数据库配置详情
func (r *DatabaseConfigRepository) Get(id uint) (*model.DatabaseConfig, error) {
	var config model.DatabaseConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Update 更新数据库配置
func (r *DatabaseConfigRepository) Update(config *model.DatabaseConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除数据库配置
func (r *DatabaseConfigRepository) Delete(id uint) error {
	return r.db.Delete(&model.DatabaseConfig{}, id).Error
}
