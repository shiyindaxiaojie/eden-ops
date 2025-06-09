package model

import (
	"time"

	"gorm.io/gorm"
)

// ServerStatus 服务器状态
type ServerStatus string

const (
	// ServerStatusEnabled 启用
	ServerStatusEnabled ServerStatus = "enabled"
	// ServerStatusDisabled 禁用
	ServerStatusDisabled ServerStatus = "disabled"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:服务器名称"`
	Host        string         `json:"host" gorm:"type:varchar(100);not null;comment:主机地址"`
	Port        int            `json:"port" gorm:"not null;comment:端口"`
	Username    string         `json:"username" gorm:"type:varchar(100);not null;comment:用户名"`
	Password    string         `json:"-" gorm:"type:varchar(100);comment:密码"` // 不返回密码
	PrivateKey  string         `json:"-" gorm:"type:text;comment:私钥"`         // 不返回私钥
	Description string         `json:"description" gorm:"type:varchar(255);comment:服务器描述"`
	Status      ServerStatus   `json:"status" gorm:"type:varchar(20);not null;default:enabled;comment:状态(enabled/disabled)"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for ServerConfig
func (ServerConfig) TableName() string {
	return "infra_server_config"
}
