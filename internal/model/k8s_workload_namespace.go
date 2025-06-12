package model

import (
	"time"
)

// K8sWorkloadNamespace K8s工作负载命名空间模型
type K8sWorkloadNamespace struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	ConfigID      int64      `json:"config_id" gorm:"not null;index"`
	Namespace     string     `json:"namespace" gorm:"size:255;not null;index"`
	WorkloadCount int        `json:"workload_count" gorm:"default:0"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
	Config        *K8sConfig `json:"config,omitempty" gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;-"`
}

// TableName 表名
func (K8sWorkloadNamespace) TableName() string {
	return "infra_k8s_workload_namespace"
}

// K8sWorkloadNamespaceResponse 命名空间响应结构
type K8sWorkloadNamespaceResponse struct {
	ID            int64     `json:"id"`
	ConfigID      int64     `json:"config_id"`
	Namespace     string    `json:"namespace"`
	WorkloadCount int       `json:"workload_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (n *K8sWorkloadNamespace) ToResponse() *K8sWorkloadNamespaceResponse {
	return &K8sWorkloadNamespaceResponse{
		ID:            n.ID,
		ConfigID:      n.ConfigID,
		Namespace:     n.Namespace,
		WorkloadCount: n.WorkloadCount,
		CreatedAt:     n.CreatedAt,
		UpdatedAt:     n.UpdatedAt,
	}
}
