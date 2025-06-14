package model

import (
	"time"
)

// K8sWorkloadHistory Workload历史记录模型
type K8sWorkloadHistory struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalID    uint      `gorm:"not null;index" json:"original_id"`                    // 原始工作负载ID
	ConfigID      int64     `gorm:"not null;index" json:"config_id"`                     // 集群配置ID
	Name          string    `gorm:"size:200;not null" json:"name"`                       // 工作负载名称
	Namespace     string    `gorm:"size:100;not null;index" json:"namespace"`            // 命名空间
	Kind          string    `gorm:"size:20;not null;index" json:"kind"`                  // 工作负载类型
	Replicas      int       `gorm:"default:0" json:"replicas"`                           // 副本数
	ReadyReplicas int       `gorm:"default:0" json:"ready_replicas"`                     // 就绪副本数
	Status        *string   `gorm:"size:20" json:"status"`                               // 状态
	Labels        *string   `gorm:"type:text" json:"labels"`                             // 标签(JSON格式)
	Selector      *string   `gorm:"type:text" json:"selector"`                           // 选择器(JSON格式)
	Images        *string   `gorm:"type:text" json:"images"`                             // 容器镜像列表(JSON格式)
	CPURequest    *string   `gorm:"size:20" json:"cpu_request"`                          // CPU请求
	CPULimit      *string   `gorm:"size:20" json:"cpu_limit"`                            // CPU限制
	MemoryRequest *string   `gorm:"size:20" json:"memory_request"`                       // 内存请求
	MemoryLimit   *string   `gorm:"size:20" json:"memory_limit"`                         // 内存限制
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`                          // 原始创建时间
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`                          // 原始更新时间
	DeletedAt     *time.Time `json:"deleted_at"`                                         // 原始删除时间
	ArchivedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"archived_at"` // 归档时间
	ArchiveReason string    `gorm:"size:100;default:'sync_cleanup';index" json:"archive_reason"` // 归档原因
}

// TableName 指定表名
func (K8sWorkloadHistory) TableName() string {
	return "infra_k8s_workload_history"
}
