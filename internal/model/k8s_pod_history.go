package model

import (
	"time"
)

// K8sPodHistory Pod历史记录模型
type K8sPodHistory struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalID     uint      `gorm:"not null;index" json:"original_id"`                    // 原始Pod ID
	ConfigID       int64     `gorm:"not null;index" json:"config_id"`                     // K8s配置ID
	WorkloadID     *int64    `gorm:"index" json:"workload_id"`                            // 工作负载ID
	Name           string    `gorm:"size:255;not null" json:"name"`                       // Pod名称
	Namespace      string    `gorm:"size:100;not null;index" json:"namespace"`            // 命名空间
	WorkloadName   *string   `gorm:"size:255" json:"workload_name"`                       // 工作负载名称
	WorkloadKind   *string   `gorm:"size:50" json:"workload_kind"`                        // 工作负载类型
	Status         string    `gorm:"size:50;not null" json:"status"`                      // Pod状态
	Phase          *string   `gorm:"size:50" json:"phase"`                                // Pod阶段
	NodeName       *string   `gorm:"size:255" json:"node_name"`                           // 节点名称
	PodIP          *string   `gorm:"size:45" json:"pod_ip"`                               // Pod IP
	HostIP         *string   `gorm:"size:45" json:"host_ip"`                              // 主机IP
	InstanceIP     *string   `gorm:"size:45" json:"instance_ip"`                          // 实例IP
	CPURequest     *string   `gorm:"size:20" json:"cpu_request"`                          // CPU请求
	CPULimit       *string   `gorm:"size:20" json:"cpu_limit"`                            // CPU限制
	MemoryRequest  *string   `gorm:"size:20" json:"memory_request"`                       // 内存请求
	MemoryLimit    *string   `gorm:"size:20" json:"memory_limit"`                         // 内存限制
	RestartCount   int       `gorm:"default:0" json:"restart_count"`                      // 重启次数
	StartTime      *time.Time `json:"start_time"`                                         // 启动时间
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`                          // 原始创建时间
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`                          // 原始更新时间
	DeletedAt      *time.Time `json:"deleted_at"`                                         // 原始删除时间
	ArchivedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"archived_at"` // 归档时间
	ArchiveReason  string    `gorm:"size:100;default:'sync_cleanup';index" json:"archive_reason"` // 归档原因
}

// TableName 指定表名
func (K8sPodHistory) TableName() string {
	return "infra_k8s_pod_history"
}

// ArchiveReason 归档原因常量
const (
	ArchiveReasonSyncCleanup = "sync_cleanup" // 同步清理
	ArchiveReasonManual      = "manual"       // 手动归档
	ArchiveReasonExpired     = "expired"      // 过期清理
)
