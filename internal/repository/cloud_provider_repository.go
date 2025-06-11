package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// CloudProviderRepository 云厂商仓库接口
type CloudProviderRepository interface {
	List(page, pageSize int, name string, status *int) ([]model.CloudProvider, int64, error)
	Get(id uint) (*model.CloudProvider, error)
	Create(provider *model.CloudProvider) error
	Update(provider *model.CloudProvider) error
	Delete(id uint) error
}

// cloudProviderRepository 云厂商仓库实现
type cloudProviderRepository struct {
	db *gorm.DB
}

// NewCloudProviderRepository 创建云厂商仓库
func NewCloudProviderRepository(db *gorm.DB) CloudProviderRepository {
	return &cloudProviderRepository{
		db: db,
	}
}

// List 获取云厂商列表
func (r *cloudProviderRepository) List(page, pageSize int, name string, status *int) ([]model.CloudProvider, int64, error) {
	var providers []model.CloudProvider
	var total int64

	query := r.db.Model(&model.CloudProvider{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&providers).Error
	if err != nil {
		return nil, 0, err
	}

	return providers, total, nil
}

// Get 获取云厂商详情
func (r *cloudProviderRepository) Get(id uint) (*model.CloudProvider, error) {
	var provider model.CloudProvider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
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
func (r *cloudProviderRepository) Delete(id uint) error {
	return r.db.Delete(&model.CloudProvider{}, id).Error
}
