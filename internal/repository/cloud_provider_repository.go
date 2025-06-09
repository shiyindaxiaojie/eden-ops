package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// CloudProviderRepository 云厂商仓库接口
type CloudProviderRepository interface {
	Create(provider *model.CloudProvider) error
	Update(provider *model.CloudProvider) error
	Delete(id int64) error
	FindByID(id int64) (*model.CloudProvider, error)
	List(page, pageSize int, name string) (int64, []model.CloudProvider, error)
}

// cloudProviderRepository 云厂商仓库实现
type cloudProviderRepository struct {
	db *gorm.DB
}

// NewCloudProviderRepository 创建云厂商仓库实例
func NewCloudProviderRepository(db *gorm.DB) CloudProviderRepository {
	return &cloudProviderRepository{
		db: db,
	}
}

// Create 创建云厂商
func (r *cloudProviderRepository) Create(provider *model.CloudProvider) error {
	return r.db.Create(provider).Error
}

// Update 更新云厂商
func (r *cloudProviderRepository) Update(provider *model.CloudProvider) error {
	return r.db.Save(provider).Error
}

// Delete 删除云厂商
func (r *cloudProviderRepository) Delete(id int64) error {
	return r.db.Delete(&model.CloudProvider{}, id).Error
}

// FindByID 根据ID查找云厂商
func (r *cloudProviderRepository) FindByID(id int64) (*model.CloudProvider, error) {
	var provider model.CloudProvider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// List 获取云厂商列表
func (r *cloudProviderRepository) List(page, pageSize int, name string) (int64, []model.CloudProvider, error) {
	var total int64
	var providers []model.CloudProvider

	query := r.db.Model(&model.CloudProvider{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&providers).Error
	if err != nil {
		return 0, nil, err
	}

	return total, providers, nil
}
