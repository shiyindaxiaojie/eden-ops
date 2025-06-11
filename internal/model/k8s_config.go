package model

import (
	"time"
)

// K8sConfig Kubernetes配置模型
type K8sConfig struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description"`
	Kubeconfig   string     `gorm:"type:text;not null" json:"kubeconfig"`
	ProviderId   *int64     `gorm:"column:provider_id" json:"providerId"`
	ProviderName string     `gorm:"-" json:"providerName"`
	Status       int        `gorm:"type:tinyint;default:1" json:"status"`
	SyncInterval int        `gorm:"type:int;default:30;comment:同步间隔(秒)" json:"syncInterval"`
	Version      string     `gorm:"type:varchar(20)" json:"version"`
	LastSyncTime *time.Time `json:"last_sync_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
}

// K8sWorkloadInfo Kubernetes工作负载信息
type K8sWorkloadInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Replicas  int    `json:"replicas"`
	Status    string `json:"status"`
}

// TableName specifies the table name for K8sConfig
func (K8sConfig) TableName() string {
	return "infra_k8s_config"
}
