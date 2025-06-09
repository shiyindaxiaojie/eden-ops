package model

import (
	"time"

	"gorm.io/gorm"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:数据库名称"`
	Type        string         `json:"type" gorm:"type:varchar(20);not null;comment:数据库类型(mysql/postgresql/mongodb)"`
	Host        string         `json:"host" gorm:"type:varchar(100);not null;comment:主机地址"`
	Port        int           `json:"port" gorm:"not null;comment:端口"`
	Username    string         `json:"username" gorm:"type:varchar(100);not null;comment:用户名"`
	Password    string         `json:"password" gorm:"type:varchar(100);not null;comment:密码"`
	Database    string         `json:"database" gorm:"type:varchar(100);not null;comment:数据库名"`
	Description string         `json:"description" gorm:"type:varchar(255);comment:数据库描述"`
	Status      string         `json:"status" gorm:"type:varchar(20);not null;default:enabled;comment:状态(enabled/disabled)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Provider *CloudProvider `json:"provider,omitempty" gorm:"foreignKey:ProviderID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE;-"`
}

// TableName specifies the table name for DatabaseConfig
func (DatabaseConfig) TableName() string {
	return "infra_database_config"
}
