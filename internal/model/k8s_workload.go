package model

import (
	"time"
)

// K8sWorkload Kubernetes工作负载
type K8sWorkload struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	ConfigID      int64      `json:"config_id" gorm:"not null"`
	Name          string     `json:"name" gorm:"size:200;not null"`
	Namespace     string     `json:"namespace" gorm:"size:100;not null"`
	Kind          string     `json:"kind" gorm:"size:20;not null"`
	Replicas      int        `json:"replicas" gorm:"default:0"`
	ReadyReplicas int        `json:"ready_replicas" gorm:"default:0"`
	Status        string     `json:"status" gorm:"size:20"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
	Config        *K8sConfig `json:"config,omitempty" gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;-"`
}

// TableName 表名
func (K8sWorkload) TableName() string {
	return "infra_k8s_workload"
}
