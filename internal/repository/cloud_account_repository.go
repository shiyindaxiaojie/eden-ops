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

// ListWithFilter 获取云账号列表（支持过滤）
func (r *CloudAccountRepository) ListWithFilter(page, pageSize int, name string, providerID *int64, status *int) (int64, []model.CloudAccount, error) {
	var total int64
	var accounts []model.CloudAccount

	// 构建基础查询，明确指定表名避免歧义
	query := r.db.Model(&model.CloudAccount{}).Table("infra_cloud_account")
	if name != "" {
		query = query.Where("infra_cloud_account.name LIKE ?", "%"+name+"%")
	}
	if providerID != nil {
		query = query.Where("infra_cloud_account.provider_id = ?", *providerID)
	}
	if status != nil {
		query = query.Where("infra_cloud_account.status = ?", *status)
	}

	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize

	// 使用原生SQL查询以获取云厂商名称
	type AccountWithProvider struct {
		model.CloudAccount
		ProviderName string `json:"providerName"`
	}

	var accountsWithProvider []AccountWithProvider
	err = query.Select("infra_cloud_account.*, COALESCE(infra_cloud_provider.name, '') as provider_name").
		Joins("LEFT JOIN infra_cloud_provider ON infra_cloud_account.provider_id = infra_cloud_provider.id").
		Offset(offset).Limit(pageSize).Scan(&accountsWithProvider).Error
	if err != nil {
		return 0, nil, err
	}

	// 转换为CloudAccount类型并设置ProviderName
	for _, item := range accountsWithProvider {
		item.CloudAccount.ProviderName = item.ProviderName
		accounts = append(accounts, item.CloudAccount)
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
