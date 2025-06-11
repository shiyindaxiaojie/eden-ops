package model

import (
	"time"

	"gorm.io/gorm"
)

// CloudProvider 云厂商模型
type CloudProvider struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:云厂商名称"`
	Code        string         `json:"code" gorm:"size:50;not null;unique;comment:云厂商代码"`
	Description string         `json:"description" gorm:"size:500;comment:描述"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 表名
func (CloudProvider) TableName() string {
	return "infra_cloud_provider"
}
