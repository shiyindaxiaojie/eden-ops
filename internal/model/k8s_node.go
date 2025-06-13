package model

import (
	"encoding/json"
	"time"
)

// K8sNode Kubernetes节点模型
type K8sNode struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID           int64     `gorm:"type:bigint;not null;index" json:"configId"`
	Name               string    `gorm:"type:varchar(255);not null" json:"name"`
	InternalIP         string    `gorm:"type:varchar(45)" json:"internalIP"`
	ExternalIP         string    `gorm:"type:varchar(45)" json:"externalIP"`
	Hostname           string    `gorm:"type:varchar(255)" json:"hostname"`
	OSImage            string    `gorm:"type:varchar(255)" json:"osImage"`
	KernelVersion      string    `gorm:"type:varchar(100)" json:"kernelVersion"`
	ContainerRuntime   string    `gorm:"type:varchar(100)" json:"containerRuntime"`
	KubeletVersion     string    `gorm:"type:varchar(50)" json:"kubeletVersion"`
	KubeProxyVersion   string    `gorm:"type:varchar(50)" json:"kubeProxyVersion"`
	CPUCapacity        string    `gorm:"type:varchar(20)" json:"cpuCapacity"`
	MemoryCapacity     string    `gorm:"type:varchar(20)" json:"memoryCapacity"`
	PodsCapacity       string    `gorm:"type:varchar(20)" json:"podsCapacity"`
	CPUAllocatable     string    `gorm:"type:varchar(20)" json:"cpuAllocatable"`
	MemoryAllocatable  string    `gorm:"type:varchar(20)" json:"memoryAllocatable"`
	PodsAllocatable    string    `gorm:"type:varchar(20)" json:"podsAllocatable"`
	CPUUsage           string    `gorm:"type:varchar(20)" json:"cpuUsage"`
	MemoryUsage        string    `gorm:"type:varchar(20)" json:"memoryUsage"`
	PodsUsage          int       `gorm:"type:int;default:0" json:"podsUsage"`
	Labels             string    `gorm:"type:text" json:"labels"`
	Annotations        string    `gorm:"type:text" json:"annotations"`
	Taints             string    `gorm:"type:text" json:"taints"`
	Conditions         string    `gorm:"type:text" json:"conditions"`
	Status             string    `gorm:"type:varchar(50);default:'Unknown'" json:"status"`
	Ready              bool      `gorm:"type:tinyint(1);default:0" json:"ready"`
	Schedulable        bool      `gorm:"type:tinyint(1);default:1" json:"schedulable"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 指定表名
func (K8sNode) TableName() string {
	return "infra_k8s_node"
}

// K8sNodeResponse 节点响应模型
type K8sNodeResponse struct {
	ID                 int64                  `json:"id"`
	ConfigID           int64                  `json:"configId"`
	Name               string                 `json:"name"`
	InternalIP         string                 `json:"internalIP"`
	ExternalIP         string                 `json:"externalIP"`
	Hostname           string                 `json:"hostname"`
	OSImage            string                 `json:"osImage"`
	KernelVersion      string                 `json:"kernelVersion"`
	ContainerRuntime   string                 `json:"containerRuntime"`
	KubeletVersion     string                 `json:"kubeletVersion"`
	KubeProxyVersion   string                 `json:"kubeProxyVersion"`
	CPUCapacity        string                 `json:"cpuCapacity"`
	MemoryCapacity     string                 `json:"memoryCapacity"`
	PodsCapacity       string                 `json:"podsCapacity"`
	CPUAllocatable     string                 `json:"cpuAllocatable"`
	MemoryAllocatable  string                 `json:"memoryAllocatable"`
	PodsAllocatable    string                 `json:"podsAllocatable"`
	CPUUsage           string                 `json:"cpuUsage"`
	MemoryUsage        string                 `json:"memoryUsage"`
	PodsUsage          int                    `json:"podsUsage"`
	Labels             map[string]interface{} `json:"labels"`
	Annotations        map[string]interface{} `json:"annotations"`
	Taints             []interface{}          `json:"taints"`
	Conditions         []interface{}          `json:"conditions"`
	Status             string                 `json:"status"`
	Ready              bool                   `json:"ready"`
	Schedulable        bool                   `json:"schedulable"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
}

// ToResponse 转换为响应模型
func (n *K8sNode) ToResponse() *K8sNodeResponse {
	resp := &K8sNodeResponse{
		ID:                n.ID,
		ConfigID:          n.ConfigID,
		Name:              n.Name,
		InternalIP:        n.InternalIP,
		ExternalIP:        n.ExternalIP,
		Hostname:          n.Hostname,
		OSImage:           n.OSImage,
		KernelVersion:     n.KernelVersion,
		ContainerRuntime:  n.ContainerRuntime,
		KubeletVersion:    n.KubeletVersion,
		KubeProxyVersion:  n.KubeProxyVersion,
		CPUCapacity:       n.CPUCapacity,
		MemoryCapacity:    n.MemoryCapacity,
		PodsCapacity:      n.PodsCapacity,
		CPUAllocatable:    n.CPUAllocatable,
		MemoryAllocatable: n.MemoryAllocatable,
		PodsAllocatable:   n.PodsAllocatable,
		CPUUsage:          n.CPUUsage,
		MemoryUsage:       n.MemoryUsage,
		PodsUsage:         n.PodsUsage,
		Status:            n.Status,
		Ready:             n.Ready,
		Schedulable:       n.Schedulable,
		CreatedAt:         n.CreatedAt,
		UpdatedAt:         n.UpdatedAt,
	}

	// 解析JSON字段
	if n.Labels != "" {
		json.Unmarshal([]byte(n.Labels), &resp.Labels)
	}
	if n.Annotations != "" {
		json.Unmarshal([]byte(n.Annotations), &resp.Annotations)
	}
	if n.Taints != "" {
		json.Unmarshal([]byte(n.Taints), &resp.Taints)
	}
	if n.Conditions != "" {
		json.Unmarshal([]byte(n.Conditions), &resp.Conditions)
	}

	return resp
}
