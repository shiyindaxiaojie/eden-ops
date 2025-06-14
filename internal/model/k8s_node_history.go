package model

import (
	"time"
)

// K8sNodeHistory Node历史记录模型
type K8sNodeHistory struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalID          uint      `gorm:"not null;index" json:"original_id"`                    // 原始节点ID
	ConfigID            int64     `gorm:"not null;index" json:"config_id"`                     // K8s配置ID
	Name                string    `gorm:"size:255;not null;index" json:"name"`                 // 节点名称
	InternalIP          *string   `gorm:"size:45;index" json:"internal_ip"`                    // 内部IP
	ExternalIP          *string   `gorm:"size:45" json:"external_ip"`                          // 外部IP
	Hostname            *string   `gorm:"size:255" json:"hostname"`                            // 主机名
	OSImage             *string   `gorm:"size:255" json:"os_image"`                            // 操作系统镜像
	KernelVersion       *string   `gorm:"size:100" json:"kernel_version"`                      // 内核版本
	ContainerRuntime    *string   `gorm:"size:100" json:"container_runtime"`                   // 容器运行时
	KubeletVersion      *string   `gorm:"size:50" json:"kubelet_version"`                      // Kubelet版本
	KubeProxyVersion    *string   `gorm:"size:50" json:"kube_proxy_version"`                   // KubeProxy版本
	CPUCapacity         *string   `gorm:"size:20" json:"cpu_capacity"`                         // CPU容量
	MemoryCapacity      *string   `gorm:"size:20" json:"memory_capacity"`                      // 内存容量
	PodsCapacity        *string   `gorm:"size:20" json:"pods_capacity"`                        // Pod容量
	CPUAllocatable      *string   `gorm:"size:20" json:"cpu_allocatable"`                      // CPU可分配
	MemoryAllocatable   *string   `gorm:"size:20" json:"memory_allocatable"`                   // 内存可分配
	PodsAllocatable     *string   `gorm:"size:20" json:"pods_allocatable"`                     // Pod可分配
	CPUUsage            *string   `gorm:"size:20" json:"cpu_usage"`                            // CPU使用量
	MemoryUsage         *string   `gorm:"size:20" json:"memory_usage"`                         // 内存使用量
	PodsUsage           int       `gorm:"default:0" json:"pods_usage"`                         // Pod使用量
	Labels              *string   `gorm:"type:text" json:"labels"`                             // 标签(JSON格式)
	Annotations         *string   `gorm:"type:text" json:"annotations"`                        // 注解(JSON格式)
	Taints              *string   `gorm:"type:text" json:"taints"`                             // 污点(JSON格式)
	Conditions          *string   `gorm:"type:text" json:"conditions"`                         // 状态条件(JSON格式)
	Status              string    `gorm:"size:50;default:'Unknown';index" json:"status"`       // 节点状态
	Ready               bool      `gorm:"default:false" json:"ready"`                          // 是否就绪
	Schedulable         bool      `gorm:"default:true" json:"schedulable"`                     // 是否可调度
	CreatedAt           time.Time `gorm:"not null" json:"created_at"`                          // 原始创建时间
	UpdatedAt           time.Time `gorm:"not null" json:"updated_at"`                          // 原始更新时间
	DeletedAt           *time.Time `json:"deleted_at"`                                         // 原始删除时间
	ArchivedAt          time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"archived_at"` // 归档时间
	ArchiveReason       string    `gorm:"size:100;default:'sync_cleanup';index" json:"archive_reason"` // 归档原因
}

// TableName 指定表名
func (K8sNodeHistory) TableName() string {
	return "infra_k8s_node_history"
}
