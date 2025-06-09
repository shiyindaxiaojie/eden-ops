package model

import (
	"time"
)

// CloudAccount 云账号模型
type CloudAccount struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"size:50;not null;uniqueIndex"`
	Provider    string     `json:"provider" gorm:"size:20;not null"`
	SecretID    string     `json:"secretId" gorm:"size:100;not null"`
	SecretKey   string     `json:"secretKey" gorm:"size:100;not null"`
	Region      string     `json:"region" gorm:"size:50"`
	Description string     `json:"description" gorm:"size:255"`
	Status      int        `json:"status" gorm:"default:1"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"index"`
}

// TableName 表名
func (CloudAccount) TableName() string {
	return "infra_cloud_account"
}
