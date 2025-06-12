package model

import (
	"time"

	"gorm.io/gorm"
)

// Migration 数据库迁移记录
type Migration struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	Version       string     `json:"version" gorm:"size:50;not null;uniqueIndex"`
	Description   string     `json:"description" gorm:"size:200"`
	Script        string     `json:"script" gorm:"size:100;not null"`
	Checksum      string     `json:"checksum" gorm:"size:32;not null"`
	InstalledBy   string     `json:"installedBy" gorm:"size:100;not null"`
	InstalledOn   time.Time  `json:"installedOn" gorm:"not null;default:CURRENT_TIMESTAMP"`
	ExecutionTime int        `json:"executionTime" gorm:"not null"`
	Success       bool       `json:"success" gorm:"not null"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt" gorm:"index"`
}

// TableName 表名
func (Migration) TableName() string {
	return "sys_db_version"
}

// AutoMigrate 自动迁移数据库模型
func AutoMigrate(db *gorm.DB) error {
	// 自动迁移模型
	return db.AutoMigrate(
		&User{},
		&Role{},
		&UserRole{},
		&Menu{},
		&RoleMenu{},
		&CloudAccount{},
		&DatabaseConfig{},
		&K8sConfig{},
		&K8sWorkload{},
		&K8sWorkloadNamespace{},
		&ServerConfig{},
		&Migration{},
	)
}
