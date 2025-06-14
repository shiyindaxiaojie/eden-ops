package model

import (
	"time"
)

// CloudAccount 云账号模型
type CloudAccount struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"size:50;not null;uniqueIndex"`
	ProviderID  *int64     `json:"providerId" gorm:"column:provider_id;comment:云厂商ID"`
	AccessKey   string     `json:"accessKey" gorm:"column:access_key;size:100;not null"`
	SecretKey   string     `json:"secretKey" gorm:"column:secret_key;size:100;not null"`
	Description string     `json:"description" gorm:"size:200"`
	Status      int        `json:"status" gorm:"default:1"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"index"`

	// 关联字段
	Provider     *CloudProvider `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	ProviderName string         `json:"providerName,omitempty" gorm:"-"`
}

// TableName 表名
func (CloudAccount) TableName() string {
	return "infra_cloud_account"
}
