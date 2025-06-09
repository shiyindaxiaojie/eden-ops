package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

type CloudAccountRepository struct {
	db *gorm.DB
}

func NewCloudAccountRepository(db *gorm.DB) *CloudAccountRepository {
	return &CloudAccountRepository{
		db: db,
	}
}

// List 获取云账号列表
func (r *CloudAccountRepository) List(page, pageSize int, name string) (int64, []model.CloudAccount, error) {
	var total int64
	var accounts []model.CloudAccount

	query := r.db.Model(&model.CloudAccount{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&accounts).Error
	if err != nil {
		return 0, nil, err
	}

	return total, accounts, nil
}

// Create 创建云账号
func (r *CloudAccountRepository) Create(account *model.CloudAccount) error {
	return r.db.Create(account).Error
}

// Get 获取云账号详情
func (r *CloudAccountRepository) Get(id uint) (*model.CloudAccount, error) {
	var account model.CloudAccount
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Update 更新云账号
func (r *CloudAccountRepository) Update(account *model.CloudAccount) error {
	return r.db.Save(account).Error
}

// Delete 删除云账号
func (r *CloudAccountRepository) Delete(id uint) error {
	return r.db.Delete(&model.CloudAccount{}, id).Error
}
