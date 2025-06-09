package model

import (
	"time"

	"gorm.io/gorm"
)

// CloudProvider 云厂商
type CloudProvider struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:厂商名称"`
	Code        string         `json:"code" gorm:"type:varchar(50);not null;comment:厂商代码"`
	Description string         `json:"description" gorm:"type:varchar(255);comment:厂商描述"`
	Status      string         `json:"status" gorm:"type:varchar(20);not null;default:enabled;comment:状态(enabled/disabled)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 表名
func (CloudProvider) TableName() string {
	return "infra_cloud_provider"
}
